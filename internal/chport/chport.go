package chport

import (
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

var wg_p sync.WaitGroup

// Interface for testing
type Cheker interface {
	Check(string, string, chan Port)
}
type Port struct {
	Port   string `json:"port"`
	Status string `json:"status"`
}

// After gourutine work ports nids sort.
func sortPorts(list_port []Port) {
	sort.Slice(list_port, func(i, j int) bool {
		port_i, _ := strconv.Atoi(list_port[i].Port)
		port_j, _ := strconv.Atoi(list_port[j].Port)
		return port_i < port_j
	})
}

// Check one port and return result in chanel
func (p Port) Check(ip, port string, res chan Port) {
	p.Port = port
	p.Status = "failed"
	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err == nil && conn != nil {
		p.Status = "success"
		_ = conn.Close()
	}
	res <- p
}

// Run goroutine to check all ports on server
func CheckPort(ip string, ports []string, p Cheker) []Port {
	res := make(chan Port)
	result := []Port{}
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
