package proxyconfig

func (c *Config) AllowedPort(port uint16, packageID int, userID int, backconnect bool) bool {
	ranges, ok := c.userAllowedTCPPorts[userID]
	if ok {
		for _, r := range ranges {
			if r[0] <= port && port <= r[1] {
				return true
			}
		}
	}
	if backconnect {
		ranges, ok := c.backconnectPackageAllowedTCPPorts[packageID]
		if ok {
			for _, r := range ranges {
				if r[0] <= port && port <= r[1] {
					return true
				}
			}
		}
	} else {
		ranges, ok := c.userPackageAllowedTCPPorts[packageID]
		if ok {
			for _, r := range ranges {
				if r[0] <= port && port <= r[1] {
					return true
				}
			}
		}
	}

	return false
}
