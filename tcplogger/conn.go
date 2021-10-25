package tcplogger

import (
	"bytes"
	"errors"
	"net"
	"os"
	"time"
)

type conn struct {
	config *TCPLoggerConfig
	conn   *net.TCPConn

	input        chan string
	lastSendTime int64
	stopped      chan bool
}

func (lc *conn) establishConnection() {
	delay := 1
	for {
		conn, err := net.DialTimeout("tcp4", lc.config.Addr, 5*time.Second)
		if err == nil {
			lc.conn = conn.(*net.TCPConn)
			break
		}
		time.Sleep(time.Duration(delay) * time.Second)
		delay *= 2
		if delay > 16 {
			delay = 16
		}
	}
}

func (lc *conn) send(data []byte) {
	if lc.conn == nil {
		lc.establishConnection()
	}

	origLen := len(data)
	pos := 0

	var err error
	for {
		lc.conn.SetWriteDeadline(time.Now().Add(lc.config.WriteTimeout))
		n, err := lc.conn.Write(data[pos:])
		pos += n
		if pos == origLen {
			break
		}
		if errors.Is(err, os.ErrDeadlineExceeded) {
			continue
		}
		if err != nil {
			lc.conn.Close()
			lc.conn = nil
			lc.establishConnection()
			for pos != 0 && data[pos] != '\n' {
				pos--
			}
			if data[pos] == '\n' {
				pos++
			}
		}
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		lc.conn.Close()
		lc.conn = nil
	}
	lc.lastSendTime = time.Now().Unix()
}

func (lc *conn) start() {
	defer close(lc.stopped)

	var buf bytes.Buffer
	buf.Grow(lc.config.SendBuf)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case v, ok := <-lc.input:
			if !ok {
				if buf.Len() > 0 {
					lc.send(buf.Bytes())
				}
				return
			}
			buf.WriteString(v)
			if buf.Len() >= lc.config.SendBuf {
				lc.send(buf.Bytes())
				buf.Reset()
			}
		case <-ticker.C:
			if time.Now().Unix()-lc.lastSendTime >= lc.config.IdleInterval {
				if buf.Len() > 0 {
					lc.send(buf.Bytes())
					buf.Reset()
				} else {
					lc.lastSendTime = time.Now().Unix()
				}
			}
		}
	}
}
