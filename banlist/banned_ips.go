package banlist

import (
	"sort"
	"sync"
	"time"
)

type BannedIPs struct {
	sync.Mutex

	DurationSec int64
	Count       int

	IPChan  chan string
	Stopped chan struct{}

	counter map[string][]int64
}

func (b *BannedIPs) Start() {
	b.Stopped = make(chan struct{})
	defer close(b.Stopped)

	b.counter = make(map[string][]int64)
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.Lock()
			for key := range b.counter {
				b.cleanup(key)
			}
			b.Unlock()
		case ip, ok := <-b.IPChan:
			if !ok {
				return
			}
			b.Lock()
			b.counter[ip] = append(b.counter[ip], time.Now().Unix())
			b.Unlock()
		}
	}
}

func (b *BannedIPs) IsBanned(ip string) bool {
	b.Lock()
	defer b.Unlock()
	if len(b.counter[ip]) < b.Count {
		return false
	}
	b.cleanup(ip)
	return len(b.counter[ip]) >= b.Count
}

func (b *BannedIPs) cleanup(key string) {
	points := b.counter[key]
	curTime := time.Now().Unix()
	pointsLen := len(points)
	n := sort.Search(pointsLen, func(i int) bool {
		return curTime-points[i] < b.DurationSec
	})
	if n > 0 && n < pointsLen {
		b.counter[key] = points[n:]
	} else if n == pointsLen {
		delete(b.counter, key)
	}
}
