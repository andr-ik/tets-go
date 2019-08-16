package user_data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type UserData struct {
	Geo Geo
}
type Geo struct {
	IP         string
	Connection string
}

func Start() {
	ips := loadIps()
	result := make(map[string][]string)

	for _, ip := range ips {
		response, _ := http.Get(fmt.Sprintf("http://1userdata01.rtty.in:2390/v3/user?ip=%s", ip))
		jsonRaw, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()

		var userData UserData
		_ = json.Unmarshal(jsonRaw, &userData)

		if _, ok := result[userData.Geo.Connection]; !ok {
			result[userData.Geo.Connection] = []string{}
		}

		result[userData.Geo.Connection] = append(result[userData.Geo.Connection], userData.Geo.IP)
	}

	for group, node := range result {
		fmt.Println(group, node)
	}
}

func loadIps() []string {
	bytes, _ := ioutil.ReadFile("propeller/user_data/ips.txt")
	ips := string(bytes)

	return strings.Split(ips, "\n")
}
