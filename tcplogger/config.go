package tcplogger

import (
	"encoding/json"
	"errors"
	"time"
)

type TCPLoggerConfig struct {
	Addr         string
	IdleInterval int64
	SendBuf      int
	BufLines     int
	ConnCount    int
	WriteTimeout time.Duration
}

type tcpLoggerConfigJSON struct {
	Addr         string        `json:"addr"`
	IdleInterval int64         `json:"idle_interval"`
	SendBuf      int           `json:"send_buf"`
	BufLines     int           `json:"buf_lines"`
	ConnCount    int           `json:"conn_count"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

var defaultValues = &TCPLoggerConfig{
	IdleInterval: 900,    // 15 minutes
	SendBuf:      131072, // 128 kilobytes
	BufLines:     1000000,
	ConnCount:    2,
	WriteTimeout: 30 * time.Second,
}

var ErrNoAddrField = errors.New("no tcp logger address specified")

func (c *TCPLoggerConfig) UnmarshalJSON(data []byte) error {
	conf := tcpLoggerConfigJSON{
		Addr:         "",
		IdleInterval: -1,
		SendBuf:      -1,
		BufLines:     -1,
		ConnCount:    -1,
		WriteTimeout: -1,
	}
	if err := json.Unmarshal(data, &conf); err != nil {
		return err
	}

	c.Addr = conf.Addr
	if c.Addr == "" {
		return ErrNoAddrField
	}
	c.IdleInterval = conf.IdleInterval
	if c.IdleInterval == -1 {
		c.IdleInterval = defaultValues.IdleInterval
	}
	c.SendBuf = conf.SendBuf
	if c.SendBuf == -1 {
		c.SendBuf = defaultValues.SendBuf
	}
	c.BufLines = conf.BufLines
	if c.BufLines == -1 {
		c.BufLines = defaultValues.BufLines
	}
	c.ConnCount = conf.ConnCount
	if c.ConnCount == -1 {
		c.ConnCount = defaultValues.ConnCount
	}
	c.WriteTimeout = conf.WriteTimeout
	if c.WriteTimeout == -1 {
		c.WriteTimeout = defaultValues.WriteTimeout
	} else {
		c.WriteTimeout *= time.Second
	}

	return nil
}
