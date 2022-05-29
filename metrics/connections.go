package metrics

import (
	"bytes"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type ConnectionsCount struct {
	PackageID string        `json:"id"`
	Ports     map[uint]uint `json:"ports"`
}

type ConnectionsCountGetter interface {
	GetConnectionsCount() []*ConnectionsCount
}

type ConnectionsReporter struct {
	APIKey             string
	APIConnectionsAddr string

	ReportInterval uint

	InfoGetter ConnectionsCountGetter
	Logger     *zap.Logger

	started bool
	stopCh  chan struct{}
}

type connectionsTotal struct {
	Time    int64               `json:"time"`
	Metrics []*ConnectionsCount `json:"metrics"`
}

func (c *ConnectionsReporter) Start() {
	if c.started {
		return
	}
	c.started = true
	c.stopCh = make(chan struct{})

	sleepTime := time.Duration(c.ReportInterval) * time.Minute
	ticker := time.NewTicker(sleepTime)
loop:
	for {
		select {
		case <-ticker.C:
			data := c.InfoGetter.GetConnectionsCount()
			go c.sendConnections(data)
		case <-c.stopCh:
			c.stopCh <- struct{}{}
			break loop
		}
	}
}

func (c *ConnectionsReporter) Stop() {
	c.stopCh <- struct{}{}
	<-c.stopCh
}

func (c *ConnectionsReporter) sendConnections(data []*ConnectionsCount) {
	tmp := &connectionsTotal{
		Time:    time.Now().Unix(),
		Metrics: data,
	}
	jsonBytes, err := json.Marshal(tmp)
	if err != nil {
		c.Logger.Error("Connections count marshalling error", zap.Error(err))
		return
	}

	reader := bytes.NewReader(jsonBytes)
	sendData(c.APIConnectionsAddr, c.APIKey, "connections", reader, c.Logger)
}
