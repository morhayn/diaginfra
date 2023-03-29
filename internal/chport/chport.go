package chport

import (
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/morhayn/diaginfra/internal/global"
)

// Interface for testing
type Cheker interface {
	Check(string, string, chan global.Port)
}
type Port struct {
	// Port   string `json:"port"`
	// Status string `json:"status"`
}

// After gourutine work ports nids sort.
func sortPorts(list_port []global.Port) {
	sort.Slice(list_port, func(i, j int) bool {
		port_i, _ := strconv.Atoi(list_port[i].Port)
		port_j, _ := strconv.Atoi(list_port[j].Port)
		return port_i < port_j
	})
}

// Check one port and return result in chanel
// ip - address server '10.0.0.1'
// port - '9200'
// res - channel to CheckPort
func (p Port) Check(ip, port string, res chan global.Port) {
	result := global.Port{
		Port:   port,
		Status: "failed",
	}
	// p.Port = port
	// p.Status = "failed"
	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err == nil && conn != nil {
		result.Status = "success"
		_ = conn.Close()
	}
	res <- result
}

// Check ssh port if ssh port failed not nid run ssh command to server
func CheckSshPort(ip, sshPort string, port Cheker) bool {
	if check := CheckPort(ip, []string{sshPort}, port); check[0].Status == "failed" {
		return false
	}
	return true
}

// CheckPort  - run goroutine to check all ports on server
// ip - address server, ports - array number check ports
// p - interface
func CheckPort(ip string, ports []string, p Cheker) []global.Port {
	var wg_p sync.WaitGroup
	res := make(chan global.Port)
	result := []global.Port{}
	go func() {
		for _, port := range ports {
			wg_p.Add(1)
			go func(port string) {
				defer wg_p.Done()
				p.Check(ip, port, res)
			}(port)
		}
		wg_p.Wait()
		close(res)
	}()
	for r := range res {
		result = append(result, r)
	}
	sortPorts(result)
	return result
}
