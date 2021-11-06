package banlist

import (
	"testing"
	"time"
)

func TestBannedIPs(t *testing.T) {
	bannedIPs := &BannedIPs{
		DurationSec: 5,
		Count:       3,
		IPBanChan:   make(chan string),
		IPUnbanChan: make(chan string),
	}
	go bannedIPs.Start()
	for i := 0; i < 5; i++ {
		bannedIPs.IPBanChan <- "1.1.1.1"
		time.Sleep(time.Second)
	}
	if !bannedIPs.IsBanned("1.1.1.1") {
		t.Fatal("Expected 1.1.1.1 to be banned")
	}
	time.Sleep(6 * time.Second)
	if bannedIPs.IsBanned("1.1.1.1") {
		t.Fatal("Expected 1.1.1.1 to not be banned")
	}
	if bannedIPs.IsBanned("1.1.1.1") {
		t.Fatal("Expected 1.1.1.1 to not be banned")
	}
	for i := 0; i < 2; i++ {
		bannedIPs.IPBanChan <- "2.2.2.2"
	}
	if bannedIPs.IsBanned("2.2.2.2") {
		t.Fatal("Expected 2.2.2.2 to not be banned")
	}
	bannedIPs.IPBanChan <- "2.2.2.2"
	if bannedIPs.IsBanned("2.2.2.2") {
		t.Fatal("Expected 2.2.2.2 to be banned")
	}
	bannedIPs.IPUnbanChan <- "2.2.2.2"
	if bannedIPs.IsBanned("2.2.2.2") {
		t.Fatal("Expected 2.2.2.2 to not be banned")
	}

	bannedIPs.StopAndWait()
}

func TestCleanup(t *testing.T) {
	BannedIPs := &BannedIPs{
		DurationSec:    60,
		Count:          3,
		infractions:    map[string]int{"1.1.1.1": 2, "2.2.2.2": 1},
		lastInfraction: map[string]int64{"1.1.1.1": time.Now().Unix() - 3601, "2.2.2.2": time.Now().Unix() - 1800},
		bannedUntil:    map[string]int64{"3.3.3.3": time.Now().Unix() - 5, "4.4.4.4": time.Now().Unix() + 30},
	}
	BannedIPs.cleanup()
	if _, ok := BannedIPs.infractions["1.1.1.1"]; ok {
		t.Fatalf("1.1.1.1 should be removed from infractions")
	}
	if _, ok := BannedIPs.infractions["2.2.2.2"]; !ok {
		t.Fatalf("2.2.2.2 should not be removed from infractions")
	}
	if _, ok := BannedIPs.lastInfraction["1.1.1.1"]; ok {
		t.Fatalf("1.1.1.1 should be removed from lastInfractions")
	}
	if _, ok := BannedIPs.lastInfraction["2.2.2.2"]; !ok {
		t.Fatalf("2.2.2.2 should not be removed from lastInfractions")
	}
	if _, ok := BannedIPs.bannedUntil["3.3.3.3"]; ok {
		t.Fatalf("3.3.3.3 should be removed from bannedUntil")
	}
	if _, ok := BannedIPs.bannedUntil["4.4.4.4"]; !ok {
		t.Fatalf("4.4.4.4 should not be removed from bannedUntil")
	}
}
