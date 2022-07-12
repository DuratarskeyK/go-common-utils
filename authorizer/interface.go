package authorizer

type AuthResult struct {
	OK bool

	SystemUser  bool
	Backconnect bool
	PackageID   int
	UserID      int
}

var BadAuthResult = AuthResult{}
var SystemUserResult = AuthResult{OK: true, SystemUser: true}
var BackconnectResult = AuthResult{OK: true, Backconnect: true}

type Authorizer interface {
	IPAuth(proxyIP, userIP string) AuthResult
	CredentialsAuth(proxyIP, username, password string) AuthResult
}
