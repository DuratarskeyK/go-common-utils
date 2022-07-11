package authorizer

type AuthResult struct {
	OK bool

	SystemUser  bool
	Backconnect bool
	PackageID   int
	UserID      int
}

var BadAuth = AuthResult{OK: false}

type Authorizer interface {
	IPAuth(proxyIP, userIP string) AuthResult
	CredentialsAuth(proxyIP, username, password string) AuthResult
}
