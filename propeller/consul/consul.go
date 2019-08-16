package consul

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Node struct {
	ID   string
	Node string
}

func Start() {
	response, _ := http.Get("http://consulweb.rtty.in/v1/catalog/nodes")
	jsonRaw, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var nodes []Node
	_ = json.Unmarshal(jsonRaw, &nodes)

	parseNode(nodes[60])
	result := make(map[string][]string)
	for _, node := range nodes {
		name := node.Node
		name, index := parseNode(node)

		if _, ok := result[name]; !ok {
			result[name] = []string{}
		}

		result[name] = append(result[name], index)
	}

	for group, node := range result {
		fmt.Println(group, node)
	}
}

func parseNode(node Node) (string, string) {
	r := regexp.MustCompile(`(?P<name>.*?)(?P<index>\d+)\.rtty\.in`)
	values := r.FindStringSubmatch(node.Node)
	if values == nil {
		return "all", node.Node
	}

	return values[1], values[2]
}
