package global

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/morhayn/diaginfra/internal/churl"
	"gopkg.in/yaml.v2"
)

type YumInit struct {
	// Stend     string            `yaml:"stend"`
	UserName string            `yaml:"user"`
	SshPort  string            `yaml:"ssh_port"`
	CountLog int               `yaml:"countlog"`
	ListUrls []string          `yaml:"list_urls"`
	Logs     map[string]string `yaml:"logs"`
	Hosts    []Init            `yaml:"hosts"`
}
type Init struct {
	Name        string   `yaml:"name"`
	Ip          string   `yaml:"ip"`
	ListPorts   []string `yaml:"list_ports"`
	ListService []string `yaml:"list_service"`
	Wars        []string `yaml:"wars"`
	// Jars         []string `yaml:"jars"`
}
type Hosts struct {
	ListUrls []churl.Url `json:"list_url"`
	Stend    []Host      `josn:"stand"`
}
type Host struct {
	Name     string   `json:"name"`
	Ip       string   `json:"ip"`
	ListPort []Port   `json:"list_port"`
	ListSsh  []Out    `json:"list_ssh"`
	Status   []Result `json:"status"`
}
type Out struct {
	Name    string `json:"name"`
	Result  string `json:"result"`
	PrgName string `json:"prgname"`
}
type Port struct {
	Port   string `json:"port"`
	Status string `json:"status"`
}
type Result struct {
	Service string `json:"service"`
	Output  string `json:"status"`
	Status  string `json:"result"`
	Alarm   bool   `json:"alarm"`
	Tooltip string `json:"tooltip"`
}

// ReadConfig Read Config file and unmarshall data in structure
func (y YumInit) ReadConfig(file string) YumInit {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error open file")
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, &y)
	if err != nil {
		fmt.Println("Error unmarshal")
		log.Fatal(err)
	}
	return y
}
