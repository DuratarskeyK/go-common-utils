package serverid

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"files.wooservers.com/sergey/go-common-utils/internal/globals"
)

func Get(apiAddr, apiKey string) int {
	addr := fmt.Sprintf("%s/server/id", strings.TrimRight(apiAddr, "/"))
	return tryGet(addr, apiKey)
}

func tryGet(idAddr, apiKey string) int {
	req, err := http.NewRequest("GET", idAddr, nil)
	if err != nil {
		return 0
	}
	req.SetBasicAuth("api", apiKey)

	resp, err := globals.HTTPClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	serverID, err := strconv.Atoi(string(bodyBytes))
	if err != nil {
		return 0
	}

	return serverID
}
