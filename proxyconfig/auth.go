package proxyconfig

import "github.com/duratarskeyk/go-common-utils/authorizer"

func (c *Config) IPAuth(proxyIP, userIP string) authorizer.AuthResult {
	if packageID, ok := c.ipToAllowedIPs[proxyIP][userIP]; ok {
		return authorizer.AuthResult{
			OK:          true,
			SystemUser:  false,
			Backconnect: false,
			PackageID:   packageID,
			UserID:      c.packageIDsToUserIDs[packageID],
		}
	}
	return authorizer.BadAuth
}

func (c *Config) CredentialsAuth(proxyIP, username, password string) authorizer.AuthResult {
	credentials := username + ":" + password
	if c.allAccess[credentials] {
		return authorizer.AuthResult{
			OK:          true,
			SystemUser:  credentials != c.BackconnectUser,
			Backconnect: credentials == c.BackconnectUser,
			PackageID:   0,
			UserID:      0,
		}
	}
	if packageID, ok := c.ipToCredentials[proxyIP][credentials]; ok {
		return authorizer.AuthResult{
			OK:          true,
			SystemUser:  false,
			Backconnect: false,
			PackageID:   packageID,
			UserID:      c.packageIDsToUserIDs[packageID],
		}
	}
	return authorizer.BadAuth
}

func (c *Config) AllAccessAuth(username, password string) bool {
	return c.allAccess[username+":"+password]
}
