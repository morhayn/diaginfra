package getlog

import (
	"testing"

	"github.com/morhayn/diaginfra/internal/sshcmd"
	"github.com/stretchr/testify/assert"
)

type mockExec struct {
}

func (m mockExec) Execute(ip string, c sshcmd.CmdExec) {
	c.Chan <- sshcmd.NewOut("test", "test", c.Cmd)
}
func TestGetLogs(t *testing.T) {
	mock := mockExec{}
	g := GetLog{
		Host:   "127.0.0.1",
		Module: "wap-test",
	}
	t.Run("test fpko", func(t *testing.T) {
		g.Service = "tomcat"
		logs := map[string]string{"tomcat": "/d01/tomcat/tomcat1/logs/"}
		res := g.GetLogs(logs, mock)
		assert.Equal(t, res, "sudo tail -n 300 /d01/tomcat/tomcat1/logs/wap-test.log")
	})
	t.Run("test scuo", func(t *testing.T) {
		g.Service = "tomcat"
		logs := map[string]string{"tomcat": "/var/log/tomcat8/"}
		res := g.GetLogs(logs, mock)
		assert.Equal(t, res, "sudo tail -n 300 /var/log/tomcat8/wap-test.log")
	})
	t.Run("test dockeer log", func(t *testing.T) {
		g.Service = "docker"
		g.Module = "test"
		logs := map[string]string{}
		res := g.GetLogs(logs, mock)
		assert.Equal(t, res, "sudo docker logs --tail 300 test")
	})
	t.Run("test cassandra", func(t *testing.T) {
		g.Service = "cassandra"
		logs := map[string]string{"cassandra": "/var/log/cassandra/system.log"}
		res := g.GetLogs(logs, mock)
		assert.Equal(t, res, "sudo tail -n 300 /var/log/cassandra/system.log")
	})
	t.Run("test hazelcast", func(t *testing.T) {
		g.Service = "hazelcast"
		logs := map[string]string{"hazelcast": "/var/log/hazelcast/hazelcast.log"}
		res := g.GetLogs(logs, mock)
		assert.Equal(t, res, "sudo tail -n 300 /var/log/hazelcast/hazelcast.log")
	})
}
