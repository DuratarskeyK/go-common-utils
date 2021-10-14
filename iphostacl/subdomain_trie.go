package iphostacl

type subdomainTrie []map[string]int

func (t subdomainTrie) include(domain string) bool {
	if len(t) == 0 {
		return false
	}

	root := t[0]
	start := len(domain)
	var cur int
	for {
		for cur = start - 1; cur >= 0 && domain[cur] != '.'; cur-- {
		}
		idx, ok := root[domain[cur+1:start]]
		if !ok {
			return false
		}
		if idx == 0 {
			return true
		}
		if cur == -1 {
			break
		}
		root = t[idx]
		start = cur
	}

	return false
}
