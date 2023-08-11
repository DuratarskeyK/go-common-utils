package iphostacl

type singleAcl struct {
	banned  *ipHostSet
	allowed *ipHostSet
}

type caseAcl struct {
	all         *singleAcl
	proxy       *singleAcl
	backconnect *singleAcl
}

type packagesAcl struct {
	proxy       map[int]*singleAcl
	backconnect map[int]*singleAcl
}

type globalBanlist struct {
	all         *ipHostSet
	proxy       *ipHostSet
	backconnect *ipHostSet
}

type Acl struct {
	ServerID int

	packages        *packagesAcl
	users           map[int]*caseAcl
	servers         *caseAcl
	globalBlacklist *globalBanlist
}
