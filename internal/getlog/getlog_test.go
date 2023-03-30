package getlog

import (
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
func TestGetLogs(t *testing.T) {
	mock := mockExec{}
	g := GetLog{
		Host:   "127.0.0.1",
		Module: "wap-test",
	}
	count := 300
	t.Run("test fpko", func(t *testing.T) {
		g.Service = "Tomcat"
		logs := map[string]string{"Tomcat": "/d01/tomcat/tomcat1/logs/"}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, "sudo tail -n 300 /d01/tomcat/tomcat1/logs/wap-test.log")
	})
	t.Run("test scuo", func(t *testing.T) {
		g.Service = "Tomcat"
		logs := map[string]string{"Tomcat": "/var/log/tomcat8/"}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, "sudo tail -n 300 /var/log/tomcat8/wap-test.log")
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
		logs := map[string]string{"Cassandra": "/var/log/cassandra/system.log"}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, "sudo tail -n 300 /var/log/cassandra/system.log")
	})
	t.Run("test hazelcast", func(t *testing.T) {
		g.Service = "Hazelcast"
		logs := map[string]string{"Hazelcast": "/var/log/hazelcast/hazelcast.log"}
		res := g.GetLogs(logs, count, mock)
		assert.Equal(t, res, "sudo tail -n 300 /var/log/hazelcast/hazelcast.log")
	})
}
func TestGetErr(t *testing.T) {
	port := chport.Port{}
	ssh := sshcmd.SshConfig{}
	ssh.Init_ssh("user", "62222")

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
