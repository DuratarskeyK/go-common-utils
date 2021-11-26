package iphostacl

func (acl *Acl) AllowedIP(pkgID, userID int, ip uint32, proxy bool) bool {
	ok, _ := acl.AllowedIPLevel(pkgID, userID, ip, proxy)

	return ok
}

func (acl *Acl) AllowedIPLevel(pkgID, userID int, ip uint32, proxy bool) (bool, uint) {
	// block anyaddr, private network addresses and zeroconf addresses:
	// 0.0.0.0
	// 10.0.0.0/8
	// 172.16.0.0/12
	// 192.168.0.0/16
	// 127.0.0.0/8
	// 169.254.0.0/16
	if ip == 0 || (167772160 <= ip && ip <= 184549375) || (2886729728 <= ip && ip <= 2887778303) ||
		(3232235520 <= ip && ip <= 3232301055) || (2130706432 <= ip && ip <= 2147483647) ||
		(2851995648 <= ip && ip <= 2852061183) {
		return false, LevelGlobal
	}

	if proxy {
		return acl.allowedIPProxyLevel(pkgID, userID, ip)
	}
	return acl.allowedIPBackconnectLevel(pkgID, userID, ip)
}

func (acl *Acl) allowedIPProxyLevel(pkgID, userID int, resolvedIP uint32) (bool, uint) {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.proxy[pkgID]
	if ok {
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true, LevelPackage
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false, LevelPackage
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true, LevelUser
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false, LevelUser
		}
		sAcl = c.proxy
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true, LevelUserProxy
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false, LevelUserProxy
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true, LevelServer
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false, LevelServer
	}
	sAcl = acl.servers.proxy
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true, LevelServerProxy
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false, LevelServerProxy
	}

	global := acl.globalBlacklist.all
	if global.ipPresent(resolvedIP) {
		return false, LevelGlobal
	}
	global = acl.globalBlacklist.proxy
	if global.ipPresent(resolvedIP) {
		return false, LevelGlobalProxy
	}
	return true, LevelNotPresent
}

func (acl *Acl) allowedIPBackconnectLevel(pkgID, userID int, resolvedIP uint32) (bool, uint) {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.backconnect[pkgID]
	if ok {
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true, LevelPackage
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false, LevelPackage
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true, LevelUser
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false, LevelUser
		}
		sAcl = c.backconnect
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true, LevelUserBackconnect
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false, LevelUserBackconnect
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true, LevelServer
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false, LevelServer
	}
	sAcl = acl.servers.backconnect
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true, LevelServerBackconnect
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false, LevelServerBackconnect
	}

	global := acl.globalBlacklist.all
	if global.ipPresent(resolvedIP) {
		return false, LevelGlobal
	}
	global = acl.globalBlacklist.backconnect
	if global.ipPresent(resolvedIP) {
		return false, LevelGlobalBackconnect
	}
	return true, LevelNotPresent
}
