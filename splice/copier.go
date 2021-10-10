package splice

import (
	"errors"
	"net"
	"os"
	"time"
)

// Copier makes a 2 way data transfer between Conn1 and Conn2
type Copier struct {
	Done         <-chan struct{}
	doneInternal chan struct{}

	bytesTransferred [2]int64

	timeoutDuration time.Duration

	once        chan struct{}
	returnError error
}

func (bc *Copier) Start(conn1 *net.TCPConn, conn2 *net.TCPConn, idleTimeout uint) (int64, int64, error) {
	bc.bytesTransferred[0], bc.bytesTransferred[1] = 0, 0

	var timeoutCnt uint
	timeoutCnt, bc.timeoutDuration = getCountAndDuration(idleTimeout)

	bc.doneInternal = make(chan struct{})
	bc.returnError = nil

	bc.once = make(chan struct{}, 1)
	bc.once <- struct{}{}

	conn1StatusChan := make(chan bool)
	conn2StatusChan := make(chan bool)

	go bc.srcToDst(conn1, conn2, conn1StatusChan, 0)
	go bc.srcToDst(conn2, conn1, conn2StatusChan, 1)

	var conn1Cnt, conn2Cnt uint
loop:
	for {
		select {
		case <-bc.doneInternal:
			break loop
		case <-bc.Done:
			bc.stop(nil)
			break loop
		case idle, ok := <-conn1StatusChan:
			if !ok {
				break loop
			} else if idle {
				conn1Cnt++
			} else {
				conn1Cnt = 0
			}
		case idle, ok := <-conn2StatusChan:
			if !ok {
				break loop
			} else if idle {
				conn2Cnt++
			} else {
				conn2Cnt = 0
			}
		}
		if conn1Cnt >= timeoutCnt && conn2Cnt >= timeoutCnt {
			bc.stop(os.ErrDeadlineExceeded)
			break loop
		}
	}

	conn1.Close()
	conn2.Close()

	for {
		_, ok1 := <-conn1StatusChan
		_, ok2 := <-conn2StatusChan
		if !(ok1 || ok2) {
			break
		}
	}

	return bc.bytesTransferred[0], bc.bytesTransferred[1], bc.returnError
}

func (bc *Copier) srcToDst(src *net.TCPConn, dst *net.TCPConn, connStatusChan chan<- bool, pos int) {
loop:
	for {
		src.SetReadDeadline(time.Now().Add(bc.timeoutDuration))
		n, err := dst.ReadFrom(src)
		bc.bytesTransferred[pos] += n
		if err == nil || !errors.Is(err, os.ErrDeadlineExceeded) {
			bc.stop(err)
			break
		}

		select {
		case connStatusChan <- n == 0:
		case <-bc.doneInternal:
			break loop
		}
	}

	close(connStatusChan)
}

func (bc *Copier) stop(err error) {
	select {
	case <-bc.once:
		bc.returnError = err
		close(bc.doneInternal)
	default:
	}
}
