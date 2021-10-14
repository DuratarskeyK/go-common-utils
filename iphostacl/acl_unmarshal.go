package iphostacl

import (
	"encoding/json"
	"strconv"
)

type singleAclJSON struct {
	Allowed *ipHostSet `json:"whitelist"`
	Banned  *ipHostSet `json:"blacklist"`
}

type caseAclJSON struct {
	All         *singleAclJSON `json:"all"`
	Proxy       *singleAclJSON `json:"proxy"`
	Backconnect *singleAclJSON `json:"backconnect"`
}

type packagesAclJSON struct {
	Proxy       map[string]*singleAclJSON `json:"proxy"`
	Backconnect map[string]*singleAclJSON `json:"backconnect"`
}

type globalBanlistJSON struct {
	All         *ipHostSet `json:"all"`
	Proxy       *ipHostSet `json:"proxy"`
	Backconnect *ipHostSet `json:"backconnect"`
}

type aclJSON struct {
	Packages        *packagesAclJSON        `json:"packages"`
	Users           map[string]*caseAclJSON `json:"users"`
	Servers         *caseAclJSON            `json:"servers"`
	GlobalBlacklist *globalBanlistJSON      `json:"global"`
}

func (acl *Acl) UnmarshalJSON(data []byte) error {
	var v aclJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	acl.packages = &packagesAcl{
		proxy:       make(map[int]*singleAcl, len(v.Packages.Proxy)),
		backconnect: make(map[int]*singleAcl, len(v.Packages.Backconnect)),
	}
	for key, packageList := range v.Packages.Proxy {
		intKey, _ := strconv.Atoi(key)
		acl.packages.proxy[intKey] = &singleAcl{
			allowed: packageList.Allowed,
			banned:  packageList.Banned,
		}
	}
	for key, packageList := range v.Packages.Backconnect {
		intKey, _ := strconv.Atoi(key)
		acl.packages.backconnect[intKey] = &singleAcl{
			allowed: packageList.Allowed,
			banned:  packageList.Banned,
		}
	}

	acl.users = make(map[int]*caseAcl, len(v.Users))
	for key, user := range v.Users {
		intKey, _ := strconv.Atoi(key)
		acl.users[intKey] = &caseAcl{
			all: &singleAcl{
				banned:  user.All.Banned,
				allowed: user.All.Allowed,
			},
			proxy: &singleAcl{
				banned:  user.Proxy.Banned,
				allowed: user.Proxy.Allowed,
			},
			backconnect: &singleAcl{
				banned:  user.Backconnect.Banned,
				allowed: user.Backconnect.Allowed,
			},
		}
	}
	acl.servers = &caseAcl{
		all: &singleAcl{
			banned:  v.Servers.All.Banned,
			allowed: v.Servers.All.Allowed,
		},
		proxy: &singleAcl{
			banned:  v.Servers.Proxy.Banned,
			allowed: v.Servers.Proxy.Allowed,
		},
		backconnect: &singleAcl{
			banned:  v.Servers.Backconnect.Banned,
			allowed: v.Servers.Backconnect.Allowed,
		},
	}
	acl.globalBlacklist = &globalBanlist{
		all:         v.GlobalBlacklist.All,
		proxy:       v.GlobalBlacklist.Proxy,
		backconnect: v.GlobalBlacklist.Backconnect,
	}

	return nil
}
