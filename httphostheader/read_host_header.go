package httphostheader

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var ErrNoHostHeader = errors.New("no host header present in the initial request")

func ReadHostHeader(conn net.Conn, buffer *bytes.Buffer, timeout time.Duration, limit int64) (string, error) {
	r := bufio.NewReader(io.LimitReader(conn, limit))

	var hostname string
	for {
		var line string
		for {
			conn.SetReadDeadline(time.Now().Add(timeout))
			str, err := r.ReadString('\n')
			line += str
			if err == nil {
				break
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if len(str) == 0 {
					return "", err
				}
				err = nil
			}
			if err != nil {
				return "", err
			}
		}
		buffer.WriteString(line)
		split := strings.SplitN(line, ":", 2)
		if len(split) == 2 {
			header := strings.TrimSpace(split[0])
			if strings.ToLower(header) == "host" {
				hostname = strings.TrimSpace(split[1])
				if pos := strings.IndexByte(hostname, ':'); pos != -1 {
					hostname = hostname[:pos]
				}
				break
			}
		}
	}

	if n := r.Buffered(); n > 0 {
		b, _ := r.Peek(n)
		buffer.Write(b)
	}

	return hostname, nil
}
