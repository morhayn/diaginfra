package webapi

import (
	"testing"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/churl"
	"github.com/morhayn/diaginfra/internal/modules"
	"github.com/morhayn/diaginfra/internal/sshcmd"
	"github.com/stretchr/testify/assert"
)

type mockUrl struct{}

func (m mockUrl) Http(url string, res chan churl.Url) {
	res <- churl.Url{Url: url, Status: 404}
}

type mockPort struct{}

func (m mockPort) Check(ip, port string, res chan chport.Port) {
	if port == "22" {
		res <- chport.Port{Port: port, Status: "success"}
	} else {
		res <- chport.Port{Port: port, Status: "failed"}
	}
}

type mockExec struct{}

func (m mockExec) Execute(ip string, c sshcmd.CmdExec) {
	if c.Name == "sshd" {
		c.Chan <- sshcmd.NewOut(c.Name, c.Name, "active")

	} else {
		c.Chan <- sshcmd.NewOut(c.Name, c.Name, "no-active")
	}
}

func TestCheckHost(t *testing.T) {
	// mocksshcmdRun = func(ip, stend string, list []string, conf sshcmd.Execer) ([]sshcmd.Out, []sshcmd.Out, error) {
	// srv := []sshcmd.Out{
	// {
	// Name:   "sshd",
	// Result: "active",
	// },
	// }
	// prg := []sshcmd.Out{}
	// return srv, prg, nil
	// }
	// mockCheckSshPort = func(list []chport.Port) bool {
	// return true
	// }
	// mockChportCheck = func(ip string, ports []string) []chport.Port {
	// res := []chport.Port{}
	// for _, p := range ports {
	// res = append(res, chport.Port{Port: p, Status: "true"})
	// }
	// return res
	// }
	var port mockPort
	// var url mockUrl
	var conf mockExec
	t.Run("simple", func(t *testing.T) {
		// mock := mockExec{}
		h := Init{
			Name:        "test",
			Ip:          "127.0.0.1",
			ListPorts:   []string{"3000", "22"},
			ListService: []string{"tomcat", "sshd"},
		}
		ch := make(chan Host)
		go checkHost(h, ch, port, conf)
		res := <-ch
		assert.Equal(t, res, Host{
			Name: "test",
			Ip:   "127.0.0.1",
			ListPort: []chport.Port{
				{Port: "22", Status: "success"},
				{Port: "3000", Status: "failed"},
			},
			ListSsh: []sshcmd.Out{
				{Name: "DiskFree", PrgName: "DiskFree", Result: "no-active"},
				{Name: "LoadAvg", PrgName: "LoadAvg", Result: "no-active"},
				{Name: "sshd", PrgName: "sshd", Result: "active"},
				{Name: "tomcat", PrgName: "tomcat", Result: "no-active"},
			},
			Status: []modules.Result(nil),
		})
	})
}
