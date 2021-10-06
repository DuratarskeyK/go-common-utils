package splice

import (
	"errors"
	"net"
	"os"
	"time"
)

// Copier makes a 2 way data transfer between Conn1 and Conn2
type Copier struct {
	GlobalStop   chan bool
	Conn1ToConn2 int64
	Conn2ToConn1 int64

	stopChan    chan bool
	once        chan bool
	returnError error
}

const (
	connIdle uint = iota
	connTraffic
	connEOF
)

// Start initiates copying and blocks
func (bc *Copier) Start(conn1 *net.TCPConn, conn2 *net.TCPConn, timeoutSec uint) error {
	timeoutCnt, timeoutDuration := getCountAndDuration(timeoutSec)

	bc.stopChan = make(chan bool)
	bc.once = make(chan bool, 1)
	bc.once <- true
	bc.returnError = nil

	timeoutChan1 := make(chan uint)
	timeoutChan2 := make(chan uint)

	go bc.copySrcToDst(conn1, conn2, timeoutDuration, timeoutChan1, &bc.Conn1ToConn2)
	go bc.copySrcToDst(conn2, conn1, timeoutDuration, timeoutChan2, &bc.Conn2ToConn1)

	var conn1Cnt, conn2Cnt uint
	var conn1EOF, conn2EOF bool
loop:
	for {
		if conn1EOF && conn2EOF {
			break
		}
		select {
		case <-bc.stopChan:
			break loop
		case <-bc.GlobalStop:
			break loop
		case b := <-timeoutChan1:
			switch b {
			case connIdle:
				conn1Cnt++
				if (conn1EOF || conn1Cnt >= timeoutCnt) && (conn2EOF || conn2Cnt >= timeoutCnt) {
					bc.stop(os.ErrDeadlineExceeded)
					break loop
				}
			case connTraffic:
				conn1Cnt = 0
			case connEOF:
				conn1EOF = true
				conn1Cnt = 0
			}
		case b := <-timeoutChan2:
			switch b {
			case connIdle:
				conn2Cnt++
				if (conn1EOF || conn1Cnt >= timeoutCnt) && (conn2EOF || conn2Cnt >= timeoutCnt) {
					bc.stop(os.ErrDeadlineExceeded)
					break loop
				}
			case connTraffic:
				conn2Cnt = 0
			case connEOF:
				conn2EOF = true
				conn2Cnt = 0
			}
		}
	}

	return bc.returnError
}

func (bc *Copier) copySrcToDst(
	src *net.TCPConn,
	dst *net.TCPConn,
	timeoutDuration time.Duration,
	timeoutChan chan<- uint,
	byteCount *int64,
) {
loop:
	for {
		select {
		case <-bc.stopChan:
			break loop
		case <-bc.GlobalStop:
			break loop
		default:
			src.SetReadDeadline(time.Now().Add(timeoutDuration))
			n, err := dst.ReadFrom(src)
			*byteCount += n
			if err == nil {
				select {
				case <-bc.stopChan:
				case <-bc.GlobalStop:
				case timeoutChan <- connEOF:
				}
				break loop
			}

			if !errors.Is(err, os.ErrDeadlineExceeded) {
				bc.stop(err)
				break loop
			}

			snd := connTraffic
			if n == 0 {
				snd = connIdle
			}

			select {
			case <-bc.GlobalStop:
				break loop
			case <-bc.stopChan:
				break loop
			case timeoutChan <- snd:
			}
		}
	}

	// We close dst write and src read, because src has eoffed or errored,
	// so no more data will be read from src and written to dst
	src.CloseRead()
	dst.CloseWrite()
}

func (bc *Copier) stop(err error) {
	select {
	case <-bc.once:
		bc.returnError = err
		close(bc.stopChan)
	default:
	}
}
