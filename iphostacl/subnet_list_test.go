package iphostacl

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSubnetList(t *testing.T) {
	var emptyList subnetList
	if emptyList.include(192<<24 | 1<<16 | 5) {
		t.Errorf("Empty list should not contain any ips")
	}

	jsonData, err := os.ReadFile("./testdata/subnet_list.json")
	if err != nil {
		t.Fatalf("ReadFile error: %s", err)
	}
	var list subnetList
	json.Unmarshal(jsonData, &list)
	ending := uint32((6 << 8) | 5)
	for a := 1; a <= 255; a++ {
		if a&7 != 1 {
			continue
		}
		for b := 1; b <= 255; b++ {
			if b&7 != 1 {
				continue
			}
			ip := (uint32(a) << 24) | (uint32(b) << 16) | ending
			res := list.include(ip)
			if !res {
				t.Errorf("Expected true for IP %d.%d.5.6, but got false", a, b)
			}
			ip = (uint32(a+1) << 24) | (uint32(b-1) << 16) | ending
			res = list.include(ip)
			if res {
				t.Errorf("Expected false for IP %d.%d.5.6, but got true", a+1, b-1)
			}
		}
	}
}

var result bool

func BenchmarkSubnetList(b *testing.B) {
	jsonData, err := os.ReadFile("./testdata/subnet_list.json")
	if err != nil {
		b.Fatalf("ReadFile error: %s", err)
	}
	var list subnetList
	json.Unmarshal(jsonData, &list)
	ending := uint32((6 << 8) | 5)
	var first, second uint32
	first = 1
	second = 1
	for i := 0; i < b.N; i++ {
		ip := (first << 24) | (second << 16) | ending
		result = list.include(ip)
		second++
		if second == 254 {
			second = 1
			first++
			if first == 254 {
				first = 1
			}
		}
	}
}
