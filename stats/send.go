package stats

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/duratarskeyk/go-common-utils/internal/globals"
	"go.uber.org/zap"
)

const tries = 4

func sendData(addr, apiKey, which string, r *bytes.Reader, logger *zap.Logger) {
	sleepTime := 1
	fail := true
	for i := 0; i < tries; i++ {
		r.Seek(0, io.SeekStart)
		req, err := http.NewRequest("POST", addr, r)
		if err != nil {
			logger.Error("Error creating HTTP request", zap.Int("try", i+1), zap.String("which", which), zap.Error(err))
			continue
		}
		req.SetBasicAuth("api", apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := globals.HTTPClient.Do(req)
		if err != nil {
			logger.Error("Error reporting %s: %s\n", zap.Int("try", i+1), zap.String("which", which), zap.Error(err))
			logger.Info("Sleeping before trying again", zap.Int("try", i+1), zap.String("which", which), zap.Int("sleep", sleepTime))
			time.Sleep(time.Duration(sleepTime) * time.Second)
			sleepTime *= 2
			continue
		}
		resp.Body.Close()
		fail = false
		break
	}
	if fail {
		logger.Error("Failed to report", zap.String("which", which), zap.Int("tries", tries))
	}
}
