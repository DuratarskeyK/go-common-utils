package hashwatcher

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/duratarskeyk/go-common-utils/internal/globals"
	"go.uber.org/zap"
)

type Watcher struct {
	HashURL       string
	ContentURL    string
	Username      string
	Password      string
	WatchInterval time.Duration
	NotifyChan    chan []byte

	Logger *zap.Logger

	currentHash string

	stopChan chan struct{}
	stopped  chan struct{}
}

func (w *Watcher) Start() error {
	w.stopChan = make(chan struct{})
	w.stopped = make(chan struct{})
	err := w.updateConfig()
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(w.WatchInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := w.updateConfig()
				if err != nil {
					w.Logger.Error("Error updating config", zap.Error(err))
				}
			case <-w.stopChan:
				w.Logger.Info("Stopping watcher")
				close(w.stopped)
				return
			}
		}
	}()

	return nil
}

func (w *Watcher) Stop() {
	close(w.stopChan)
	<-w.stopped
}

func (w *Watcher) updateConfig() error {
	newHash, err := w.getCurrentConfigHash()
	if err != nil {
		return err
	}
	if newHash == w.currentHash {
		return nil
	}
	content, err := w.getCurrentContent()
	if err != nil {
		return err
	}
	w.currentHash = newHash
	w.Logger.Info("New config received", zap.String("hash", newHash))

	select {
	case w.NotifyChan <- content:
	case <-w.stopChan:
		break
	}

	return nil
}

func (w *Watcher) getCurrentConfigHash() (string, error) {
	req, err := http.NewRequest("GET", w.HashURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(w.Username, w.Password)

	resp, err := globals.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (w *Watcher) getCurrentContent() ([]byte, error) {
	req, err := http.NewRequest("GET", w.ContentURL, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(w.Username, w.Password)

	resp, err := globals.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}
