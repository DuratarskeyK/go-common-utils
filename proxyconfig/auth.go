package proxyconfig

func (c *Config) IPAuth(proxyIP, userIP string) (bool, int) {
	if p, ok := c.ipToAllowedIPs[proxyIP]; ok {
		if packageID, ok := p[userIP]; ok {
			return true, packageID
		}
	}
	return false, 0
}

func (c *Config) CredentialsAuth(proxyIP, username, password string) (bool, int) {
	credentials := username + ":" + password
	if _, ok := c.allAccess[credentials]; ok {
		return true, 0
	}
	if val, ok := c.ipToCredentials[proxyIP]; ok {
		if id, ok := val[credentials]; ok {
			return true, id
		}
	}
	return false, 0
}

func (c *Config) AllAccessAuth(username, password string) bool {
	credentials := username + ":" + password
	_, ok := c.allAccess[credentials]
	return ok
}
