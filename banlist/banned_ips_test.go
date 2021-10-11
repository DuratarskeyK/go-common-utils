package banlist_test

import (
	"testing"
	"time"

	"github.com/duratarskeyk/go-common-utils/banlist"
)

func TestBannedIPs(t *testing.T) {
	banned_ips := &banlist.BannedIPs{
		DurationSec: 5,
		Count:       3,
		IPChan:      make(chan string),
		Stopped:     make(chan struct{}),
	}
	go banned_ips.Start()
	for i := 0; i < 5; i++ {
		banned_ips.IPChan <- "1.1.1.1"
		time.Sleep(time.Second)
	}
	if !banned_ips.IsBanned("1.1.1.1") {
		t.Fatal("Expected 1.1.1.1 to be banned")
	}
	time.Sleep(2 * time.Second)
	if banned_ips.IsBanned("1.1.1.1") {
		t.Fatal("Expected 1.1.1.1 to not be banned")
	}
	time.Sleep(6 * time.Second)
	if banned_ips.IsBanned("1.1.1.1") {
		t.Fatal("Expected 1.1.1.1 to not be banned")
	}
	for i := 0; i < 2; i++ {
		banned_ips.IPChan <- "2.2.2.2"
	}
	if banned_ips.IsBanned("2.2.2.2") {
		t.Fatal("Expected 2.2.2.2 to not be banned")
	}
	banned_ips.IPChan <- "2.2.2.2"

	close(banned_ips.IPChan)
	<-banned_ips.Stopped
}
