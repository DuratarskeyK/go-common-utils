package tcplogger

import (
	"container/list"
	"time"
)

type TCPLogger struct {
	Config *TCPLoggerConfig
	Input  chan string

	receiverStopped chan bool

	output chan string
	conns  []*conn
}

func (l *TCPLogger) Start() {
	l.output = make(chan string)
	l.receiverStopped = make(chan bool)

	l.conns = make([]*conn, l.Config.ConnCount)
	for i := 0; i < l.Config.ConnCount; i++ {
		l.conns[i] = &conn{
			config:       l.Config,
			input:        l.output,
			lastSendTime: time.Now().Unix(),
			stopped:      make(chan bool),
		}
	}

	go l.Receiver()

	for _, conn := range l.conns {
		go conn.start()
	}
}

func (l *TCPLogger) WaitForShutdown() {
	<-l.receiverStopped
	for _, conn := range l.conns {
		<-conn.stopped
	}
}

func (l *TCPLogger) Receiver() {
	var buf list.List
	outCh := func() chan string {
		if buf.Len() == 0 {
			return nil
		}
		return l.output
	}
	curVal := func() string {
		if buf.Len() == 0 {
			return ""
		}
		return buf.Front().Value.(string)
	}
	for buf.Len() > 0 || l.Input != nil {
		select {
		case v, ok := <-l.Input:
			if !ok {
				l.Input = nil
			} else {
				buf.PushBack(v)
				if buf.Len() > l.Config.BufLines {
					buf.Remove(buf.Front())
				}
			}
		case outCh() <- curVal():
			buf.Remove(buf.Front())
		}
	}
	close(l.output)
	close(l.receiverStopped)
}
