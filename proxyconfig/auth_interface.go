package proxyconfig

type Authorizer interface {
	IPAuth(proxyIP, userIP string) *AuthResult
	CredentialsAuth(proxyIP, username, password string) *AuthResult
}
