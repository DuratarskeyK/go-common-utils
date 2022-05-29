package stats

import (
	"bytes"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// Data stores metrics for a package, like number of requests
type PackageMetrics struct {
	PackageID       string `json:"id"`
	Hits            uint   `json:"hits"`
	DownloadedBytes int64  `json:"downloaded_bytes"`
	UploadedBytes   int64  `json:"uploaded_bytes"`
}

type requestsTotal struct {
	Time    int64             `json:"time"`
	Metrics []*PackageMetrics `json:"metrics"`
}

// Requests represents a hits counter for each
// backconnect package
type RequestsReporter struct {
	HitsAddr string
	APIKey   string

	ReportInterval uint

	Logger *zap.Logger

	CounterCh chan *PackageMetrics

	counters map[string]*PackageMetrics
	started  bool
	closeCh  chan struct{}
}

func (hc *RequestsReporter) Start() {
	if hc.started {
		return
	}
	hc.started = true
	hc.counters = make(map[string]*PackageMetrics)
	hc.closeCh = make(chan struct{})

	sleepTime := time.Duration(hc.ReportInterval) * time.Minute
	ticker := time.NewTicker(sleepTime)
	for {
		select {
		case v := <-hc.CounterCh:
			packageID := v.PackageID
			s := hc.counters[packageID]
			if s == nil {
				hc.counters[packageID] = v
			} else {
				s.DownloadedBytes += v.DownloadedBytes
				s.UploadedBytes += v.UploadedBytes
				s.Hits++
			}
		case <-ticker.C:
			if len(hc.counters) == 0 {
				continue
			}
			go hc.sendRequests(hc.counters)
			hc.counters = make(map[string]*PackageMetrics)
		case <-hc.closeCh:
			hc.closeCh <- struct{}{}
		}
	}
}

func (hc *RequestsReporter) Stop() {
	hc.closeCh <- struct{}{}
	<-hc.closeCh
}

func (hc *RequestsReporter) sendRequests(data map[string]*PackageMetrics) {
	tmp := &requestsTotal{
		Time:    time.Now().Unix(),
		Metrics: make([]*PackageMetrics, 0, len(data)),
	}
	for _, val := range data {
		tmp.Metrics = append(tmp.Metrics, val)
	}
	jsonBytes, err := json.Marshal(tmp)
	if err != nil {
		hc.Logger.Error("Requests marshalling error", zap.Error(err))
		return
	}

	reader := bytes.NewReader(jsonBytes)
	sendData(hc.HitsAddr, hc.APIKey, "requests", reader, hc.Logger)
}
