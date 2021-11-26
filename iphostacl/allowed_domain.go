package iphostacl

func (acl *Acl) AllowedDomain(pkgID, userID int, host string, proxy bool) bool {
	ok, _ := acl.AllowedDomainLevel(pkgID, userID, host, proxy)

	return ok
}

func (acl *Acl) AllowedDomainLevel(pkgID, userID int, host string, proxy bool) (bool, uint) {
	if proxy {
		return acl.allowedDomainProxyLevel(pkgID, userID, host)
	}

	return acl.allowedDomainBackconnectLevel(pkgID, userID, host)
}

func (acl *Acl) allowedDomainProxyLevel(pkgID, userID int, host string) (bool, uint) {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.proxy[pkgID]
	if ok {
		if sAcl.allowed.domainPresent(host) {
			return true, LevelPackage
		}
		if sAcl.banned.domainPresent(host) {
			return false, LevelPackage
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.domainPresent(host) {
			return true, LevelUser
		}
		if sAcl.banned.domainPresent(host) {
			return false, LevelUser
		}
		sAcl = c.proxy
		if sAcl.allowed.domainPresent(host) {
			return true, LevelUserProxy
		}
		if sAcl.banned.domainPresent(host) {
			return false, LevelUserProxy
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.domainPresent(host) {
		return true, LevelServer
	}
	if sAcl.banned.domainPresent(host) {
		return false, LevelServer
	}
	sAcl = acl.servers.proxy
	if sAcl.allowed.domainPresent(host) {
		return true, LevelServerProxy
	}
	if sAcl.banned.domainPresent(host) {
		return false, LevelServerProxy
	}

	global := acl.globalBlacklist.all
	if global.domainPresent(host) {
		return false, LevelGlobal
	}
	global = acl.globalBlacklist.proxy
	if global.domainPresent(host) {
		return false, LevelGlobalProxy
	}
	return true, LevelNotPresent
}

func (acl *Acl) allowedDomainBackconnectLevel(pkgID, userID int, host string) (bool, uint) {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.backconnect[pkgID]
	if ok {
		if sAcl.allowed.domainPresent(host) {
			return true, LevelPackage
		}
		if sAcl.banned.domainPresent(host) {
			return false, LevelPackage
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.domainPresent(host) {
			return true, LevelUser
		}
		if sAcl.banned.domainPresent(host) {
			return false, LevelUser
		}
		sAcl = c.backconnect
		if sAcl.allowed.domainPresent(host) {
			return true, LevelUserBackconnect
		}
		if sAcl.banned.domainPresent(host) {
			return false, LevelUserBackconnect
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.domainPresent(host) {
		return true, LevelServer
	}
	if sAcl.banned.domainPresent(host) {
		return false, LevelServer
	}
	sAcl = acl.servers.backconnect
	if sAcl.allowed.domainPresent(host) {
		return true, LevelServerBackconnect
	}
	if sAcl.banned.domainPresent(host) {
		return false, LevelServerBackconnect
	}

	global := acl.globalBlacklist.all
	if global.domainPresent(host) {
		return false, LevelGlobal
	}
	global = acl.globalBlacklist.backconnect
	if global.domainPresent(host) {
		return false, LevelGlobalBackconnect
	}
	return true, LevelNotPresent
}
