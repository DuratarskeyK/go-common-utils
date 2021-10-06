package timeoutread

import (
	"errors"
	"net"
	"os"
	"time"
)

var nilTime time.Time

func ReadWithTimeout(conn net.Conn, timeout time.Duration, p []byte) (int, error) {
	total := len(p)
	pos := 0
	defer conn.SetReadDeadline(nilTime)
	for {
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err := conn.Read(p[pos:])
		total -= n
		pos += n
		if errors.Is(err, os.ErrDeadlineExceeded) {
			if n == 0 {
				return pos, err
			}
			err = nil
		}
		if err != nil || total == 0 {
			return pos, err
		}
	}
}
