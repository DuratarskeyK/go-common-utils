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

var defaultValues = &TCPLoggerConfig{
	IdleInterval: 900,    // 15 minutes
	SendBuf:      131072, // 128 kilobytes
	BufLines:     1000000,
	ConnCount:    2,
	WriteTimeout: 30 * time.Second,
}

var ErrNoAddrField = errors.New("no tcp logger address specified")

func (c *TCPLoggerConfig) UnmarshalJSON(data []byte) error {
	var conf map[string]interface{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return err
	}

	var ok bool
	c.Addr, ok = conf["addr"].(string)
	if !ok {
		return ErrNoAddrField
	}
	c.IdleInterval, ok = conf["idle_interval"].(int64)
	if !ok {
		c.IdleInterval = defaultValues.IdleInterval
	}
	c.SendBuf, ok = conf["send_buf"].(int)
	if !ok {
		c.SendBuf = defaultValues.SendBuf
	}
	c.BufLines, ok = conf["buf_lines"].(int)
	if !ok {
		c.BufLines = defaultValues.BufLines
	}
	c.ConnCount, ok = conf["conn_count"].(int)
	if !ok {
		c.ConnCount = defaultValues.ConnCount
	}
	c.WriteTimeout, ok = conf["write_timeout"].(time.Duration)
	if !ok {
		c.WriteTimeout = defaultValues.WriteTimeout
	} else {
		c.WriteTimeout *= time.Second
	}

	return nil
}
