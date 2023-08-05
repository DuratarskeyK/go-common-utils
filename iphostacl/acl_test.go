package iphostacl

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testItem struct {
	PkgID  int    `json:"pkg_id"`
	UserID int    `json:"user_id"`
	Val    string `json:"val"`
	IP     bool   `json:"ip"`
	Proxy  bool   `json:"proxy"`
	Result bool   `json:"res"`
}

func (t testItem) String() string {
	res := fmt.Sprintf("pkgID=%d userID=%d ", t.PkgID, t.UserID)
	if t.IP {
		res += fmt.Sprintf("ip=%s ", t.Val)
	} else {
		res += fmt.Sprintf("ip=%s ", t.Val)
	}
	if t.Proxy {
		res += "proxy"
	} else {
		res += "backconnect"
	}

	return res
}

func TestAllowedIP(t *testing.T) {
	jsonStr, err := os.ReadFile("./testdata/acl.json")
	if err != nil {
		t.Fatalf("ReadFile error: %s", err)
	}
	var acl Acl
	json.Unmarshal(jsonStr, &acl)
	var testData []testItem
	jsonStr, err = os.ReadFile("./testdata/acl_test_cases.json")
	if err != nil {
		t.Fatalf("ReadFile error: %s", err)
	}
	json.Unmarshal(jsonStr, &testData)

	for i, item := range testData {
		var r bool
		if item.IP {
			parts := strings.Split(item.Val, ".")
			var ip uint32
			for i := 0; i < 4; i++ {
				part, _ := strconv.Atoi(parts[i])
				ip = (ip << 8) | uint32(part)
			}
			r = acl.AllowedIP(item.PkgID, item.UserID, ip, item.Proxy)
		} else {
			r = acl.AllowedDomain(item.PkgID, item.UserID, item.Val, item.Proxy)
		}
		if r != item.Result {
			t.Errorf("Test#%d \"%s\": Expected %v but got %v", i+1, item, item.Result, r)
		}
	}
}
