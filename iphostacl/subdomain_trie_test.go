package iphostacl

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestSubdomainTrie(t *testing.T) {
	jsonData, err := os.ReadFile("./testdata/subdomain_trie.json")
	if err != nil {
		t.Fatalf("ReadFile error: %s", err)
	}
	var trie subdomainTrie
	json.Unmarshal(jsonData, &trie)
	subdomains := []string{}
	jsonData, err = os.ReadFile("./testdata/subdomain_trie_cases.json")
	if err != nil {
		t.Fatalf("ReadFile error: %s", err)
	}
	json.Unmarshal(jsonData, &subdomains)
	for _, subdomain := range subdomains {
		testStr := "test" + subdomain
		res := trie.include(testStr)
		if !res {
			t.Errorf("Expected true for %s, but got false", testStr)
		}
		testStr = subdomain[1:]
		res = trie.include(testStr)
		if !res {
			t.Errorf("Expected true for %s, but got false", testStr)
		}
		testStr = "blah" + subdomain[1:]
		res = trie.include(testStr)
		if res {
			t.Errorf("Expected false for %s, but got true", testStr)
		}
	}
	dots := []string{}
	for _, subdomain := range subdomains {
		if strings.Count(subdomain, ".") > 2 {
			dots = append(dots, subdomain)
		}
	}
	for _, subdomain := range dots {
		testStr := strings.Join(strings.Split(subdomain, ".")[2:], ".")
		res := trie.include(testStr)
		if res {
			t.Errorf("Expected false for %s, but got true", testStr)
		}
	}

	var emptyTrie subdomainTrie
	if emptyTrie.include("example.org") {
		t.Error("Empty subtree must not include any domains")
	}
}

var trieResult bool

func BenchmarkTrie(b *testing.B) {
	jsonData, err := os.ReadFile("./testdata/subdomain_trie.json")
	if err != nil {
		b.Fatalf("ReadFile error: %s", err)
	}
	var trie subdomainTrie
	json.Unmarshal(jsonData, &trie)
	subdomains := []string{}
	jsonData, err = os.ReadFile("./testdata/subdomain_trie_cases.json")
	if err != nil {
		b.Fatalf("ReadFile error: %s", err)
	}
	json.Unmarshal(jsonData, &subdomains)
	benchmarkDomains := []string{}
	for _, subdomain := range subdomains {
		benchmarkDomains = append(benchmarkDomains, "test"+subdomain)
		benchmarkDomains = append(benchmarkDomains, subdomain[1:])
		benchmarkDomains = append(benchmarkDomains, "blah"+subdomain[1:])
		if strings.Count(subdomain, ".") > 2 {
			benchmarkDomains = append(benchmarkDomains, strings.Join(strings.Split(subdomain, ".")[2:], "."))
		}
	}

	length := len(benchmarkDomains)
	var res bool
	for i := 0; i < b.N; i++ {
		res = trie.include(benchmarkDomains[i%length])
		trieResult = res
	}
}
