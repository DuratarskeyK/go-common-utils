package proxyconfig

import "github.com/duratarskeyk/go-common-utils/authorizer"

func (c *Config) IPAuth(proxyIP, userIP string) authorizer.AuthResult {
	if packageID := c.ipToAllowedIPs[proxyIP][userIP]; packageID != 0 {
		return authorizer.AuthResult{
			OK:        true,
			PackageID: packageID,
			UserID:    c.packageIDsToUserIDs[packageID],
		}
	}
	return authorizer.BadAuthResult
}

func (c *Config) CredentialsAuth(proxyIP, username, password string) authorizer.AuthResult {
	credentials := username + ":" + password
	if credentials == c.checkerUser || c.allAccess[credentials] {
		return authorizer.SystemUserResult
	} else if credentials == c.backconnectUser {
		return authorizer.BackconnectResult
	}
	if packageID := c.ipToCredentials[proxyIP][credentials]; packageID != 0 {
		return authorizer.AuthResult{
			OK:        true,
			PackageID: packageID,
			UserID:    c.packageIDsToUserIDs[packageID],
		}
	}
	return authorizer.BadAuthResult
}

func (c *Config) AllAccessAuth(username, password string) bool {
	return c.allAccess[username+":"+password]
}
