package banlist

import (
	"sync"
	"time"
)

type BannedIPs struct {
	guard sync.RWMutex

	DurationSec int64
	Count       int

	IPBanChan   chan string
	IPUnbanChan chan string

	stop    chan struct{}
	stopped chan struct{}

	infractions    map[string]int
	lastInfraction map[string]int64
	bannedUntil    map[string]int64
}

const cleanUpDuration int64 = 3600

func (b *BannedIPs) Start() {
	b.stopped = make(chan struct{})
	b.stop = make(chan struct{})
	go b.realStart()
}

func (b *BannedIPs) StopAndWait() {
	close(b.stop)
	<-b.stopped
}

func (b *BannedIPs) realStart() {
	defer close(b.stopped)
	b.infractions = make(map[string]int)
	b.lastInfraction = make(map[string]int64)
	b.bannedUntil = make(map[string]int64)

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.cleanup()
		case ip := <-b.IPBanChan:
			b.guard.Lock()
			if !b.isBannedInternal(ip) {
				b.infractions[ip]++
				if b.infractions[ip] == b.Count {
					delete(b.infractions, ip)
					delete(b.lastInfraction, ip)
					b.bannedUntil[ip] = time.Now().Unix() + b.DurationSec
				} else {
					b.lastInfraction[ip] = time.Now().Unix()
				}
			} else {
				b.bannedUntil[ip] = time.Now().Unix() + b.DurationSec
			}
			b.guard.Unlock()
		case ip := <-b.IPUnbanChan:
			b.guard.Lock()
			delete(b.bannedUntil, ip)
			delete(b.infractions, ip)
			delete(b.lastInfraction, ip)
			b.guard.Unlock()
		case <-b.stop:
			return
		}
	}
}

func (b *BannedIPs) IsBanned(ip string) bool {
	b.guard.Lock()
	defer b.guard.Unlock()
	return b.isBannedInternal(ip)
}

func (b *BannedIPs) isBannedInternal(ip string) bool {
	until, ok := b.bannedUntil[ip]
	if !ok {
		return false
	}
	banned := time.Now().Unix() <= until
	if !banned {
		delete(b.bannedUntil, ip)
	}
	return banned
}

func (b *BannedIPs) cleanup() {
	b.guard.Lock()
	curTime := time.Now().Unix()
	for ip, until := range b.bannedUntil {
		if curTime > until {
			delete(b.bannedUntil, ip)
		}
	}
	for ip, last := range b.lastInfraction {
		if curTime-last > cleanUpDuration {
			delete(b.lastInfraction, ip)
			delete(b.infractions, ip)
		}
	}
	b.guard.Unlock()
}
