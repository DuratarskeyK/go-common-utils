package iphostacl

type subnetList struct {
	Min []uint32 `json:"min"`
	Max []uint32 `json:"max"`
}

func (slist *subnetList) include(ip uint32) bool {
	length := len(slist.Min)
	if length == 0 {
		return false
	}
	l := 0
	r := length
	for l < r {
		mid := (l + r) / 2
		if slist.Min[mid] > ip {
			l = mid + 1
		} else {
			r = mid
		}
	}
	if l < length && slist.Min[l] <= ip && ip <= slist.Max[l] {
		return true
	}
	return false
}
