package proxyconfig

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/duratarskeyk/go-common-utils/internal/globals"
	"go.uber.org/zap"
)

type Watcher struct {
	APIAddr       string
	APIKey        string
	ServerID      int
	WatchInterval time.Duration
	Logger        *zap.Logger

	StopChan chan struct{}
	Stopped  chan struct{}

	configHashURL string
	configURL     string

	credentialsHash string
	configSync      sync.RWMutex
	currentConfig   *Config

	started     bool
	notifyChans []chan *Config
}

func (w *Watcher) Start() {
	w.started = true
	defer close(w.Stopped)
	apiAddr := strings.TrimRight(w.APIAddr, "/")
	w.configHashURL = fmt.Sprintf("%s/server/%d/auth_config_hash", apiAddr, w.ServerID)
	w.configURL = fmt.Sprintf("%s/server/%d/auth_config", apiAddr, w.ServerID)
	err := w.updateConfig()
	if err != nil {
		w.Logger.Fatal("Error getting auth config for the first time", zap.Error(err))
	}
	ticker := time.NewTicker(w.WatchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := w.updateConfig()
			if err != nil {
				w.Logger.Error("Error updating config", zap.Error(err))
			}
		case <-w.StopChan:
			return
		}
	}
}

func (w *Watcher) GetConfig() *Config {
	w.configSync.RLock()
	res := w.currentConfig
	w.configSync.RUnlock()

	return res
}

func (w *Watcher) AddNotifyChan(ch chan *Config) {
	if w.started {
		return
	}
	w.notifyChans = append(w.notifyChans, ch)
}

func (w *Watcher) updateConfig() error {
	newHash, err := w.getCurrentConfigHash()
	if err != nil {
		return err
	}
	if newHash == w.credentialsHash {
		return nil
	}
	configDataRaw, err := w.getCurrentConfigData()
	if err != nil {
		return err
	}
	var newConfig Config
	err = json.Unmarshal(configDataRaw, &newConfig)
	if err != nil {
		return err
	}
	w.configSync.Lock()
	w.currentConfig = &newConfig
	w.credentialsHash = newHash
	w.configSync.Unlock()
	w.Logger.Info("New config received", zap.String("hash", newHash))

	for _, ch := range w.notifyChans {
		select {
		case ch <- w.currentConfig:
		case <-w.StopChan:
			return nil
		}
	}

	return nil
}

func (w *Watcher) getCurrentConfigData() ([]byte, error) {
	req, err := http.NewRequest("GET", w.configURL, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("api", w.APIKey)

	resp, err := globals.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code %d", resp.StatusCode)
	}

	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return resB.Bytes(), nil
}

func (w *Watcher) getCurrentConfigHash() (string, error) {
	req, err := http.NewRequest("GET", w.configHashURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("api", w.APIKey)

	resp, err := globals.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
