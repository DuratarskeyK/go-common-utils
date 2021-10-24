package proxyconfig

import (
	"encoding/json"
	"strconv"

	"github.com/duratarskeyk/go-common-utils/iphostacl"
)

type Config struct {
	PackageIDsToUserIDs map[int]int

	IPHostACL *iphostacl.Acl

	ipToCredentials map[string]map[string]int
	ipToAllowedIPs  map[string]map[string]int
	backconnectUser string
	checkerUser     string
	allAccess       map[string]bool
}

type configJSON struct {
	IPToCredentials map[string]map[string]int `json:"ips_to_credentials"`
	IPToAllowedIPs  map[string]map[string]int `json:"ips_to_authorized_ips"`
	BackconnectUser string                    `json:"backconnect_user"`
	CheckerUser     string                    `json:"checker_user"`
	AllAccess       map[string]bool           `json:"all_access"`

	PackageIDsToUserIDs map[string]int `json:"package_ids_to_user_ids"`

	IPHostACL *iphostacl.Acl `json:"graylist"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var cj configJSON
	if err := json.Unmarshal(data, &cj); err != nil {
		return err
	}

	c.ipToCredentials = cj.IPToCredentials
	c.ipToAllowedIPs = cj.IPToAllowedIPs
	c.backconnectUser = cj.BackconnectUser
	c.checkerUser = cj.CheckerUser
	c.allAccess = cj.AllAccess

	c.IPHostACL = cj.IPHostACL
	c.PackageIDsToUserIDs = make(map[int]int)
	for k, v := range cj.PackageIDsToUserIDs {
		intK, _ := strconv.Atoi(k)
		c.PackageIDsToUserIDs[intK] = v
	}

	return nil
}
