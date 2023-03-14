package sshcmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockExec struct {
}

func (m mockExec) Execute(ip string, c CmdExec) {
	c.Chan <- NewOut("test", "test", c.Cmd)
}

func TestSwSmd(t *testing.T) {
	srv := make(chan Out)
	prg := make(chan Out)
	t.Run("service", func(t *testing.T) {
		ssh := CmdExec{
			Name: "tomcat1",
		}
		ssh.swCmd(srv, prg)
		assert.Equal(t, ssh.Cmd, "sudo systemctl is-active tomcat1")
		assert.Equal(t, ssh.Chan, srv)
	})
	t.Run("test Tomcat", func(t *testing.T) {
		ssh := CmdExec{
			Name: "Tomcat:test:1234567:8081",
		}
		ssh.swCmd(srv, prg)
		assert.Equal(t, ssh.Name, "Tomcat")
		assert.Equal(t, ssh.Cmd, "curl -u test:1234567 http://127.0.0.1:8081/manager/text/list")
		assert.Equal(t, ssh.Chan, prg)
	})
	t.Run("test Elastic", func(t *testing.T) {
		ssh := CmdExec{
			Name: "Elastic",
		}
		ssh.swCmd(srv, prg)
		assert.Equal(t, ssh.Cmd, "curl -X GET http://127.0.0.1:9200/_cluster/health")
		assert.Equal(t, ssh.Chan, prg)
	})
	t.Run("test space in arg swCmd", func(t *testing.T) {
		ssh := CmdExec{
			Name: "Systemd:Test Space",
		}
		ssh.swCmd(srv, prg)
		assert.Equal(t, ssh.Cmd, "")
	})
	t.Run("test error many arg to swCmd", func(t *testing.T) {
		ssh := CmdExec{
			Name: "Systemd:Test:ARG",
		}
		ssh.swCmd(srv, prg)
		assert.Equal(t, ssh.Cmd, "")
	})
	t.Run("test error few arg swCmd", func(t *testing.T) {
		ssh := CmdExec{
			Name: "Tomcat:admin:pass",
		}
		ssh.swCmd(srv, prg)
		assert.Equal(t, ssh.Cmd, "")
	})
}
func TestBuildCmd(t *testing.T) {
	t.Run("test Build one prg commands", func(t *testing.T) {
		c := Comands{}
		list := []string{
			"Cassandra",
		}
		_, p := c.buildCmd(list)
		assert.True(t, len(c.Comm) == 3)
		assert.Equal(t, c.Comm[0].Name, "DiskFree")
		assert.Equal(t, c.Comm[2], CmdExec{
			Name: "Cassandra",
			Chan: p,
			Cmd:  "nodetool status",
		})

	})
	t.Run("test Build one srv command", func(t *testing.T) {
		c := Comands{}
		list := []string{
			"sshd",
		}
		s, _ := c.buildCmd(list)
		assert.True(t, len(c.Comm) == 3)
		assert.Equal(t, c.Comm[0].Name, "DiskFree")
		assert.Equal(t, c.Comm[1].Name, "LoadAvg")
		assert.Equal(t, c.Comm[2], CmdExec{
			Name: "sshd",
			Chan: s,
			Cmd:  "sudo systemctl is-active sshd",
		})
	})
}
