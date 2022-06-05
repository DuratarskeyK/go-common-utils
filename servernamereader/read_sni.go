package servernamereader

import (
	"net"
	"time"

	"github.com/duratarskeyk/go-common-utils/idlenet"
)

func ReadSNI(conn net.Conn, timeout time.Duration) (string, []byte, error) {
	header := make([]byte, tlsHeaderLen)
	_, err := idlenet.ReadWithTimeout(conn, timeout, header)
	if err != nil {
		return "", nil, err
	}

	if (header[0]&0x80) > 0 && header[2] == 1 {
		return "", header, ErrSSL2
	}

	tlsContentType := header[0]
	if tlsContentType != tlsHandshakeContentType {
		return "", header, ErrNotTLSHandshake
	}

	tlsVersionMajor := header[1]
	tlsVersionMinor := header[2]

	if tlsVersionMajor < 3 {
		return "", header, ErrBadTLSVersion
	}

	packetLength := (uint(header[3]) << 8) | (uint(header[4]))
	data := make([]byte, packetLength+tlsHeaderLen)
	for i := 0; i < tlsHeaderLen; i++ {
		data[i] = header[i]
	}
	n, err := idlenet.ReadWithTimeout(conn, timeout, data[tlsHeaderLen:])
	dataLen := uint(n) + tlsHeaderLen
	if err != nil {
		return "", nil, err
	}
	var pos uint = tlsHeaderLen

	if data[pos] != tlsHandshakeTypeClientHello {
		return "", data, ErrNotTLSHello
	}

	pos += 38

	if pos+1 > dataLen {
		return "", data, ErrNoSNI
	}
	len := uint(data[pos])
	pos += 1 + len

	if pos+2 > dataLen {
		return "", data, ErrNoSNI
	}
	len = (uint(data[pos]) << 8) | uint(data[pos+1])
	pos += 2 + len

	if pos+1 > dataLen {
		return "", data, ErrNoSNI
	}
	len = uint(data[pos])
	pos += 1 + len

	if pos == dataLen && tlsVersionMajor == 3 && tlsVersionMinor == 0 {
		return "", data, ErrSSL3NoExtensions
	}

	if pos+2 > dataLen {
		return "", data, ErrNoSNI
	}
	len = (uint(data[pos]) << 8) | uint(data[pos+1])
	pos += 2

	if pos+len > dataLen {
		return "", data, ErrNoSNI
	}

	hostname, err := parseExtensions(data[pos:], len)
	if err != nil {
		return "", data, err
	}

	return hostname, data, nil
}

func parseExtensions(data []byte, dataLen uint) (string, error) {
	var pos, len uint

	for pos+4 <= dataLen {
		len = (uint(data[pos+2]) << 8) | uint(data[pos+3])

		if data[pos] == 0x0 && data[pos+1] == 0x0 {
			if pos+4+len > dataLen {
				return "", ErrNoSNI
			}

			return parseSNE(data[pos+4:], len)
		}
		pos += 4 + len
	}

	return "", ErrNoSNI
}

func parseSNE(data []byte, dataLen uint) (string, error) {
	var pos, len uint
	pos = 2

	for pos+3 < dataLen {
		len = (uint(data[pos+1]) << 8) | uint(data[pos+2])
		if pos+3+len > dataLen {
			return "", ErrNoSNI
		}

		if data[pos] == 0x0 {
			return string(data[pos+3 : pos+3+len]), nil
		}
	}

	return "", ErrNoSNI
}
