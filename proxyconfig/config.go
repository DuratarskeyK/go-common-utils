package proxyconfig

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/duratarskeyk/go-common-utils/iphostacl"
)

type Config struct {
	ServerName string

	BackconnectUser string
	CheckerUser     string

	PackageIDsToUserIDs         map[int]int
	UserIDToEmail               map[int]string
	UserPackageConnectionsLimit map[int]uint
	UserPackageAllowUDP         map[int]bool

	IPHostACL *iphostacl.Acl

	userPackageAllowedTCPPorts        map[int][][2]uint16
	backconnectPackageAllowedTCPPorts map[int][][2]uint16
	userAllowedTCPPorts               map[int][][2]uint16

	ipToCredentials map[string]map[string]int
	ipToAllowedIPs  map[string]map[string]int
	allAccess       map[string]bool
}

type configJSON struct {
	ServerName string `json:"server_name"`

	IPToCredentials map[string]map[string]int `json:"ips_to_credentials"`
	IPToAllowedIPs  map[string]map[string]int `json:"ips_to_authorized_ips"`
	BackconnectUser string                    `json:"backconnect_user"`
	CheckerUser     string                    `json:"checker_user"`
	AllAccess       map[string]bool           `json:"all_access"`

	PackageIDsToUserIDs         map[string]int    `json:"package_ids_to_user_ids"`
	UserIDToEmail               map[string]string `json:"user_id_to_email"`
	UserPackageConnectionsLimit map[string]uint   `json:"user_package_connection_limits"`
	UserPackageAllowUDP         map[string]bool   `json:"user_package_allow_udp"`

	UserPackageAllowedTCPPorts        map[string]string `json:"user_package_allowed_tcp_ports"`
	BackconnectPackageAllowedTCPPorts map[string]string `json:"backconnect_package_allowed_tcp_ports"`
	UserAllowedTCPPorts               map[string]string `json:"user_allowed_tcp_ports"`

	IPHostACL *iphostacl.Acl `json:"graylist"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var cj configJSON
	if err := json.Unmarshal(data, &cj); err != nil {
		return err
	}

	c.ServerName = cj.ServerName

	c.BackconnectUser = cj.BackconnectUser
	c.CheckerUser = cj.CheckerUser

	c.ipToCredentials = cj.IPToCredentials
	c.ipToAllowedIPs = cj.IPToAllowedIPs
	c.allAccess = cj.AllAccess

	c.IPHostACL = cj.IPHostACL
	c.PackageIDsToUserIDs = make(map[int]int)
	for k, v := range cj.PackageIDsToUserIDs {
		intK, _ := strconv.Atoi(k)
		c.PackageIDsToUserIDs[intK] = v
	}
	c.UserIDToEmail = make(map[int]string)
	for k, v := range cj.UserIDToEmail {
		intK, _ := strconv.Atoi(k)
		c.UserIDToEmail[intK] = v
	}
	c.UserPackageConnectionsLimit = make(map[int]uint)
	for k, v := range cj.UserPackageConnectionsLimit {
		intK, _ := strconv.Atoi(k)
		c.UserPackageConnectionsLimit[intK] = v
	}
	c.UserPackageAllowUDP = make(map[int]bool)
	for k := range cj.UserPackageAllowUDP {
		intK, _ := strconv.Atoi(k)
		c.UserPackageAllowUDP[intK] = true
	}

	c.userPackageAllowedTCPPorts = getPortRanges(cj.UserPackageAllowedTCPPorts)
	c.backconnectPackageAllowedTCPPorts = getPortRanges(cj.BackconnectPackageAllowedTCPPorts)
	c.userAllowedTCPPorts = getPortRanges(cj.UserAllowedTCPPorts)

	return nil
}

func getPortRanges(ports map[string]string) map[int][][2]uint16 {
	res := make(map[int][][2]uint16)
	for id, allowedPorts := range ports {
		idInt, _ := strconv.Atoi(id)
		res[idInt] = make([][2]uint16, 0)
		ranges := strings.Split(allowedPorts, ",")
		for _, r := range ranges {
			if strings.IndexByte(r, '-') != -1 {
				tmp := strings.Split(r, "-")
				start, _ := strconv.Atoi(tmp[0])
				end, _ := strconv.Atoi(tmp[1])
				res[idInt] = append(res[idInt], [2]uint16{uint16(start), uint16(end)})
			} else {
				start, _ := strconv.Atoi(r)
				res[idInt] = append(res[idInt], [2]uint16{uint16(start), uint16(start)})
			}
		}
	}

	return res
}
