package keepalive

import (
	"net"
	"syscall"
)

// SetKeepAlive sets keepalive on the underlying conn socket with given idle, interval and count
func SetKeepAlive(conn *net.TCPConn, idle, interval, count int) error {
	raw, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	var keepAliveErr error
	err = raw.Control(func(fd uintptr) {
		keepAliveErr = realSetKeepAlive(int(fd), idle, interval, count)
	})
	if err != nil {
		return err
	}

	return keepAliveErr
}

func realSetKeepAlive(fdInt, idle, interval, count int) error {
	err := syscall.SetsockoptInt(fdInt, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1)
	if err != nil {
		return err
	}
	err = syscall.SetsockoptInt(fdInt, syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, count)
	if err != nil {
		return err
	}
	err = syscall.SetsockoptInt(fdInt, syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE, idle)
	if err != nil {
		return err
	}
	err = syscall.SetsockoptInt(fdInt, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, interval)
	if err != nil {
		return err
	}
	return nil
}
