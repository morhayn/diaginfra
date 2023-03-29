package global

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
