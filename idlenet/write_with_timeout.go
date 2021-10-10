package idlenet

import (
	"errors"
	"net"
	"os"
	"time"
)

func WriteWithTimeout(conn net.Conn, timeout time.Duration, p []byte) (int, error) {
	total := len(p)
	pos := 0
	defer conn.SetWriteDeadline(zeroTime)
	for {
		conn.SetWriteDeadline(time.Now().Add(timeout))
		n, err := conn.Write(p[pos:])
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
