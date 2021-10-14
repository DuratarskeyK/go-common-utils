package iphostacl

func (acl *Acl) AllowedDomain(pkgID, userID int, host string, proxy bool) bool {
	if proxy {
		return acl.allowedDomainProxy(pkgID, userID, host)
	}

	return acl.allowedDomainBackconnect(pkgID, userID, host)
}

func (acl *Acl) allowedDomainProxy(pkgID, userID int, host string) bool {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.proxy[pkgID]
	if ok {
		if sAcl.allowed.domainPresent(host) {
			return true
		}
		if sAcl.banned.domainPresent(host) {
			return false
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.domainPresent(host) {
			return true
		}
		if sAcl.banned.domainPresent(host) {
			return false
		}
		sAcl = c.proxy
		if sAcl.allowed.domainPresent(host) {
			return true
		}
		if sAcl.banned.domainPresent(host) {
			return false
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.domainPresent(host) {
		return true
	}
	if sAcl.banned.domainPresent(host) {
		return false
	}
	sAcl = acl.servers.proxy
	if sAcl.allowed.domainPresent(host) {
		return true
	}
	if sAcl.banned.domainPresent(host) {
		return false
	}

	global := acl.globalBlacklist.all
	if global.domainPresent(host) {
		return false
	}
	global = acl.globalBlacklist.proxy
	return !global.domainPresent(host)
}

func (acl *Acl) allowedDomainBackconnect(pkgID, userID int, host string) bool {
	var sAcl *singleAcl
	var ok bool
	sAcl, ok = acl.packages.backconnect[pkgID]
	if ok {
		if sAcl.allowed.domainPresent(host) {
			return true
		}
		if sAcl.banned.domainPresent(host) {
			return false
		}
	}

	c, ok := acl.users[userID]
	if ok {
		sAcl = c.all
		if sAcl.allowed.domainPresent(host) {
			return true
		}
		if sAcl.banned.domainPresent(host) {
			return false
		}
		sAcl = c.backconnect
		if sAcl.allowed.domainPresent(host) {
			return true
		}
		if sAcl.banned.domainPresent(host) {
			return false
		}
	}

	sAcl = acl.servers.all
	if sAcl.allowed.domainPresent(host) {
		return true
	}
	if sAcl.banned.domainPresent(host) {
		return false
	}
	sAcl = acl.servers.backconnect
	if sAcl.allowed.domainPresent(host) {
		return true
	}
	if sAcl.banned.domainPresent(host) {
		return false
	}

	global := acl.globalBlacklist.all
	if global.domainPresent(host) {
		return false
	}
	global = acl.globalBlacklist.backconnect
	return !global.domainPresent(host)
}
