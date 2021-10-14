package iphostacl

func (acl *Acl) AllowedIP(pkgID, userID int, ip uint32, proxy bool) bool {
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
		return false
	}

	if proxy {
		return acl.allowedIPProxy(pkgID, userID, ip)
	}
	return acl.allowedIPBackconnect(pkgID, userID, ip)
}

func (acl *Acl) allowedIPProxy(pkgID, userID int, resolvedIP uint32) bool {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.proxy[pkgID]
	if ok {
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false
		}
		sAcl = c.proxy
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false
	}
	sAcl = acl.servers.proxy
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false
	}

	global := acl.globalBlacklist.all
	if global.ipPresent(resolvedIP) {
		return false
	}
	global = acl.globalBlacklist.proxy
	return !global.ipPresent(resolvedIP)
}

func (acl *Acl) allowedIPBackconnect(pkgID, userID int, resolvedIP uint32) bool {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.backconnect[pkgID]
	if ok {
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false
		}
		sAcl = c.backconnect
		if sAcl.allowed.ipPresent(resolvedIP) {
			return true
		}
		if sAcl.banned.ipPresent(resolvedIP) {
			return false
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false
	}
	sAcl = acl.servers.backconnect
	if sAcl.allowed.ipPresent(resolvedIP) {
		return true
	}
	if sAcl.banned.ipPresent(resolvedIP) {
		return false
	}

	global := acl.globalBlacklist.all
	if global.ipPresent(resolvedIP) {
		return false
	}
	global = acl.globalBlacklist.backconnect
	return !global.ipPresent(resolvedIP)
}
