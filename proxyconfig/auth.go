package proxyconfig

type AuthResult struct {
	SystemUser  bool
	Backconnect bool
	PackageID   int
	UserID      int
}

func (c *Config) IPAuth(proxyIP, userIP string) *AuthResult {
	if packageID, ok := c.ipToAllowedIPs[proxyIP][userIP]; ok {
		return &AuthResult{
			SystemUser:  false,
			Backconnect: false,
			PackageID:   packageID,
			UserID:      c.packageIDsToUserIDs[packageID],
		}
	}
	return nil
}

func (c *Config) CredentialsAuth(proxyIP, username, password string) *AuthResult {
	credentials := username + ":" + password
	if _, ok := c.allAccess[credentials]; ok {
		return &AuthResult{
			SystemUser:  credentials != c.BackconnectUser,
			Backconnect: credentials == c.BackconnectUser,
			PackageID:   0,
			UserID:      0,
		}
	}
	if packageID, ok := c.ipToCredentials[proxyIP][credentials]; ok {
		return &AuthResult{
			SystemUser:  false,
			Backconnect: false,
			PackageID:   packageID,
			UserID:      c.packageIDsToUserIDs[packageID],
		}
	}
	return nil
}

func (c *Config) AllAccessAuth(username, password string) bool {
	return c.allAccess[username+":"+password]
}
