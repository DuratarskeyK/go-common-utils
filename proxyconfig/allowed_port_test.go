package proxyconfig

import (
	"encoding/json"
	"testing"
)

type portTestData struct {
	port        uint16
	packageID   int
	userID      int
	backconnect bool
	allowed     bool
}

func TestAllowedPort(t *testing.T) {
	jsonData := []byte(`{
		"user_package_allowed_tcp_ports": {"1": "22,30-40", "2": "23,41-51"},
		"backconnect_package_allowed_tcp_ports": {"1": "51,55-100", "2": "121,130-140"},
		"user_allowed_tcp_ports": {"3":"200,300", "4": "400,500-600"}
		}`)
	var c Config
	json.Unmarshal(jsonData, &c)
	testData := []portTestData{
		{22, 1, 3, false, true},
		{35, 1, 3, false, true},
		{41, 1, 3, false, false},
		{200, 1, 3, false, true},
		{23, 2, 3, false, true},
		{46, 2, 3, false, true},
		{78, 2, 3, false, false},
		{300, 2, 3, false, true},
		{400, 2, 4, false, true},
		{550, 2, 4, false, true},
		{51, 1, 3, true, true},
		{68, 1, 3, true, true},
		{150, 1, 3, true, false},
		{200, 1, 3, true, true},
		{300, 1, 4, true, false},
		{400, 1, 4, true, true},
		{550, 1, 4, true, true},
		{121, 2, 3, true, true},
		{122, 2, 4, true, false},
		{550, 2, 4, true, true},
	}

	for nr, test := range testData {
		if c.AllowedPort(test.port, test.packageID, test.userID, test.backconnect) != test.allowed {
			t.Errorf(
				"Test #%d: Expected AllowedPort to return %v, got %v, data: port=%d packageID=%d userID=%d backconnect=%v",
				nr+1,
				test.allowed,
				!test.allowed,
				test.port,
				test.packageID,
				test.userID,
				test.backconnect,
			)
		}
	}
}
