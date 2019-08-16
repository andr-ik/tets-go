package old

import (
	"encoding/json"
	"fmt"
	"strings"
)
import "github.com/hashicorp/consul/api"

type ServiceConfig struct {
	Name string
	Host string
	Port int
	Tags []string
}

func NewServiceConfig(name, host string, port int, tags []string) *ServiceConfig {
	return &ServiceConfig{
		Host: host,
		Name: name,
		Port: port,
		Tags: tags,
	}
}

func (service *ServiceConfig) Id() string {
	return fmt.Sprintf("%s_%s_%d_%s", service.Name, service.Host, service.Port, strings.Join(service.Tags, "_"))
}

func Finish(consul *api.Client, id string) {
	err := consul.Agent().ServiceDeregister(id)
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	serviceConfig := NewServiceConfig("test", "127.0.0.1", 1234, []string{"a", "b", "c"})

	service := api.AgentServiceRegistration{
		ID:      serviceConfig.Id(),
		Name:    serviceConfig.Name,
		Address: serviceConfig.Host,
		Port:    serviceConfig.Port,
		Tags:    serviceConfig.Tags,
	}

	err = consul.Agent().ServiceRegister(&service)
	if err != nil {
		panic(err)
	}
	defer Finish(consul, serviceConfig.Id())

	kv := consul.KV()
	dataRaw := []int{1, 2, 4}
	data, _ := json.Marshal(dataRaw)
	key := &api.KVPair{Key: "a/", Value: data}
	kv.Put(key, nil)

	var enter int
	fmt.Println("Hello world!")

	for {
		keys, _, _ := kv.Keys("a", "", nil)
		for _, k := range keys {
			v, _, _ := kv.Get(k, nil)
			if v != nil {
				fmt.Print(k)
				fmt.Print(" ")
				fmt.Println(string(v.Value))
			}
		}
		_, err = fmt.Scan(&enter)
		if enter == 1 {
			break
		}

		fmt.Println(enter)
	}

	_, err = fmt.Scan(&enter)
	if err != nil {
		panic(err)
	}
}
