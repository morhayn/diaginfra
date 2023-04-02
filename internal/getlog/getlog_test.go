package getlog

import (
	"fmt"
	"testing"
	"time"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/sshcmd"
	"github.com/stretchr/testify/assert"
)

type mockExec struct{}

func (m mockExec) Execute(ip string, c sshcmd.CmdExec) {
	c.Chan <- sshcmd.NewOut("test", "test", c.Cmd)
}
func (m mockExec) GetSshPort() string {
	return "22"
}

type mockErr struct{}

func (e mockErr) Execute(ip string, c sshcmd.CmdExec) {
	c.Chan <- sshcmd.NewOut("test", "test", "10")
}
func (e mockErr) GetSshPort() string {
	return "22"
}

func TestGetLogs(t *testing.T) {
	mock := mockExec{}
	g := GetLog{
		Host:   "127.0.0.1",
		Module: "wap-test",
	}
	count := 300
	t.Run("test fpko", func(t *testing.T) {
		g.Service = "Tomcat"
		path := "/d01/tomcat/logs/"
		logs := map[string]string{"Tomcat": path}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, fmt.Sprintf("sudo test -f %swap-test.log && sudo tail -n 300 %swap-test.log", path, path))
	})
	t.Run("test scuo", func(t *testing.T) {
		g.Service = "Tomcat"
		path := "/var/log/tomcat/"
		logs := map[string]string{"Tomcat": path}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, fmt.Sprintf("sudo test -f %swap-test.log && sudo tail -n 300 %swap-test.log", path, path))
	})
	t.Run("test dockeer log", func(t *testing.T) {
		g.Service = "Docker"
		g.Module = "test"
		logs := map[string]string{}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, "sudo docker logs --tail 300 test")
	})
	t.Run("test cassandra", func(t *testing.T) {
		g.Service = "Cassandra"
		path := "/var/log/cassandra/system.log"
		logs := map[string]string{"Cassandra": path}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, fmt.Sprintf("sudo test -f %s && sudo tail -n 300 %s", path, path))
	})
	t.Run("test hazelcast", func(t *testing.T) {
		g.Service = "Hazelcast"
		path := "/var/log/hazelcast/hazelcast.log"
		logs := map[string]string{"Hazelcast": path}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, fmt.Sprintf("sudo test -f %s && sudo tail -n 300 %s", path, path))
	})
}
func TestGetErr(t *testing.T) {
	port := chport.Port{}
	ssh := sshcmd.SshConfig{}
	loadData := global.YumInit{
		UserName: "user",
		SshPort:  "62222",
		CountLog: 400,
		ListUrls: []string{},
		Logs: map[string]string{
			"test": "/var/log/test",
		},
		Hosts: []global.Init{},
	}
	Status := global.Hosts{
		Stend: []global.Host{
			{
				Name:     "test1",
				Ip:       "127.0.0.1",
				ListPort: []global.Port{},
				ListSsh:  []global.Out{},
				Status: []global.Result{
					{
						Service: "test",
						Output:  "Test Service",
						Status:  "running",
						Alarm:   false,
						Tooltip: "",
					},
				},
			},
		},
	}
	t.Run("Ssh Port Check", func(t *testing.T) {
		ch := make(chan []GetLog)
		go func() {
			ch <- GetErr(Status, loadData, port, ssh)
		}()
		select {
		case <-ch:
		case <-time.After(3 * time.Second):
			t.Fatal("TIMEOUT 3 second")
		}
	})
}
func TestErrBuildCmd(t *testing.T) {
	getLog := GetLog{
		Host:    "127.0.0.1",
		Service: "Tomcat",
		Module:  "test",
		Errors:  0,
	}
	t.Run("test Tomcat module", func(t *testing.T) {
		getLog.Service = "Tomcat"
		logs := map[string]string{
			"Tomcat": "/var/log/tomcat/",
		}
		res := getLog.errBuildCmd(logs, 300, mockErr{})
		if res.Errors != 10 {
			t.Fatal("Errors from function not 10 ", res)
		}
	})
	t.Run("test Service Cassandra", func(t *testing.T) {
		getLog.Service = "Cassandra"
		logs := map[string]string{
			"Cassandra": "/var/log/cassandra/cassandra.log",
		}
		res := getLog.errBuildCmd(logs, 300, mockErr{})
		if res.Errors != 10 {
			t.Fatal("Errors from function not 10 ", res)
		}
	})
	t.Run("server return string", func(t *testing.T) {
		getLog.Service = "Tomcat"
		logs := map[string]string{
			"Tomcat": "/var/log/tomcat/",
		}
		res := getLog.errBuildCmd(logs, 300, mockExec{})
		if res.Errors != 999 {
			t.Fatal("Errors not 999 ", res)
		}
	})
	t.Run("docker logs", func(t *testing.T) {
		logs := map[string]string{}
		getLog.Service = "Docker"
		res := getLog.errBuildCmd(logs, 300, mockErr{})
		if res.Errors != 10 {
			t.Fatal("Errors from function not 10 ", res)
		}
	})
}
