package tcplogger

import "time"

type TCPLoggerConfig struct {
	Addr         string        `json:"addr"`
	IdleInterval int64         `json:"idle_interval"`
	SendBuf      int           `json:"send_buf"`
	BufLines     int           `json:"buf_lines"`
	ConnCount    int           `json:"conn_count"`
	WriteTimeout time.Duration `json:"write_timeout"`
}
