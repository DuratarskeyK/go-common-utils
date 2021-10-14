package iphostacl

import (
	"encoding/json"
	"strconv"
)

type ipHostSet struct {
	empty bool

	ips     map[uint32]bool
	subnets *subnetList

	domains    map[string]bool
	subdomains subdomainTrie
}

type ipHostSetJSON struct {
	Empty bool `json:"empty"`

	IPs     map[string]bool `json:"ips"`
	Subnets *subnetList     `json:"subnets"`

	Domains    map[string]bool `json:"domains"`
	Subdomains subdomainTrie   `json:"subdomains"`
}

func (l *ipHostSet) UnmarshalJSON(data []byte) error {
	var v ipHostSetJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	l.empty = v.Empty
	if v.Empty {
		return nil
	}
	l.domains = v.Domains
	l.subdomains = v.Subdomains
	l.subnets = v.Subnets
	l.ips = make(map[uint32]bool)
	for k := range v.IPs {
		ip, _ := strconv.Atoi(k)
		l.ips[uint32(ip)] = true
	}
	return nil
}

func (l *ipHostSet) ipPresent(ip uint32) bool {
	if l.empty {
		return false
	}

	_, ok := l.ips[ip]
	if ok {
		return true
	}
	return l.subnets.include(ip)
}

func (l *ipHostSet) domainPresent(domain string) bool {
	if l.empty {
		return false
	}

	if _, ok := l.domains[domain]; ok {
		return true
	}

	return l.subdomains.include(domain)
}
