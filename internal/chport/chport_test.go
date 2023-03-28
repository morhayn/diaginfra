package chport

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPort struct {
	Port   string
	Status string
}

func (p mockPort) Check(ip, port string, res chan Port) {
	// defer wg_p.Done()
	if port == "22" {
		res <- Port{Port: port, Status: "success"}
	} else {
		res <- Port{Port: port, Status: "failed"}
	}
}
func TestCheckPort(t *testing.T) {
	t.Run("check ports", func(t *testing.T) {
		var ports = []string{
			"2222",
			"22",
			"3000",
		}
		var p mockPort
		res := CheckPort("127.0.0.1", ports, p)
		assert.Equal(t, len(res), 3)
		assert.Equal(t, res[0], Port{Port: "22", Status: "success"})
		assert.Equal(t, res[1].Status, "failed")
		assert.Equal(t, res[2].Status, "failed")
	})
}
func TestCheck(t *testing.T) {
	t.Run("Sort", func(t *testing.T) {
		ports := []Port{
			{
				Port:   "222",
				Status: "success",
			},
			{
				Port:   "5013",
				Status: "success",
			},
			{
				Port:   "2",
				Status: "failed",
			},
		}
		sortPorts(ports)
		assert.Equal(t, ports[0].Port, "2")
		assert.Equal(t, ports[1].Port, "222")
		assert.Equal(t, ports[2].Port, "5013")
	})
	t.Run("Check port", func(t *testing.T) {
		ch := make(chan Port)
		var p = Port{Port: "5000", Status: "failed"}
		go p.Check("127.0.0.1", p.Port, ch)
		r := <-ch
		assert.Equal(t, r, Port{Port: "5000", Status: "failed"})
	})
	t.Run("Check successfull ports check", func(t *testing.T) {
		ch := make(chan Port)
		var p = Port{Port: "5000", Status: "failed"}
		go p.Check("127.0.0.1", p.Port, ch)
		l, err := net.Listen("tcp", "127.0.0.1:5000")
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		r := <-ch
		assert.Equal(t, r.Status, "success")
	})
}
