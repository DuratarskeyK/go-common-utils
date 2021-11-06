package banlist

import (
	"sort"
)

type Ban struct {
	IP    string
	Until int64
}

type Infraction struct {
	IP               string
	Count            int
	LastInfractionAt int64
}

func (b *BannedIPs) Statistics() ([]Ban, []Infraction) {
	b.cleanup()
	b.guard.RLock()
	banned := make([]Ban, 0, len(b.bannedUntil))
	for ip, until := range b.bannedUntil {
		banned = append(banned, Ban{IP: ip, Until: until})
	}
	infractions := make([]Infraction, 0, len(b.infractions))
	for ip, count := range b.infractions {
		infractions = append(infractions, Infraction{IP: ip, Count: count, LastInfractionAt: b.lastInfraction[ip]})
	}
	b.guard.RUnlock()
	sort.Slice(banned, func(i, j int) bool {
		return banned[i].Until > banned[j].Until
	})
	sort.Slice(infractions, func(i, j int) bool {
		a := infractions[i]
		b := infractions[j]
		return a.Count > b.Count || (a.Count == b.Count && a.LastInfractionAt > b.LastInfractionAt)
	})

	return banned, infractions
}
