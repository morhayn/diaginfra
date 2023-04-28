package handl

import (
	"testing"
	"time"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/churl"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/sshcmd"
	"github.com/stretchr/testify/assert"
)

type mockUrl struct{}

func (m mockUrl) Http(url string, res chan churl.Url) {
	res <- churl.Url{Url: url, Status: 404}
}

type mockPort struct{}

func (m mockPort) Check(ip, port string, res chan global.Port) {
	if port == "22" {
		res <- global.Port{Port: port, Status: "success"}
	} else {
		res <- global.Port{Port: port, Status: "failed"}
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
func (m mockExec) GetSshPort() string {
	return "22"
}
func (m mockExec) Scp(ip, src, dest string) error {
	return nil
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
		h := global.Init{
			Name:        "test",
			Ip:          "127.0.0.1",
			ListPorts:   []string{"3000", "22"},
			ListService: []string{"tomcat", "sshd"},
		}
		ch := make(chan global.Host)
		go checkHost(h, ch, port, conf)
		res := <-ch
		assert.Equal(t, res, global.Host{
			Name: "test",
			Ip:   "127.0.0.1",
			ListPort: []global.Port{
				{Port: "22", Status: "success"},
				{Port: "3000", Status: "failed"},
			},
			ListSsh: []global.Out{
				{Name: "DiskFree", PrgName: "DiskFree", Result: "no-active"},
				{Name: "LoadAvg", PrgName: "LoadAvg", Result: "no-active"},
				{Name: "sshd", PrgName: "sshd", Result: "active"},
				{Name: "tomcat", PrgName: "tomcat", Result: "no-active"},
			},
			Status: []global.Result(nil),
		})
	})
}

func TestCheckHostTimeOut(t *testing.T) {
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
		Hosts: []global.Init{
			{
				Name:        "test-server",
				Ip:          "127.0.0.1",
				ListPorts:   []string{"22"},
				ListService: []string{"test", "service"},
				Wars:        []string{},
			},
		},
	}
	t.Run("Ssh Port Check", func(t *testing.T) {
		ch := make(chan global.Host)
		go func() {
			checkHost(loadData.Hosts[0], ch, port, ssh)
		}()
		select {
		case <-ch:
		case <-time.After(3 * time.Second):
			t.Fatal("TIMEOUT 3 second")
		}
	})
}

//
//func TestHandleTomcat(t *testing.T) {
//	t.Run("HandleTomcat", func(t *testing.T) {
//		t.Run("true", func(t *testing.T) {
//			in := `OK - Listed applications for virtual host [localhost]
//				/test:running:0:test
//				/manager:stopped:0:test`
//			r, err := handleTomcatStatus(in)
//			assert.NoError(t, err)
//			assert.True(t, r["test"] == "running")
//			assert.True(t, r["manager"] == "stopped")
//		})
//		t.Run("path to out", func(t *testing.T) {
//			in := `OK - Listed applications for virtual host [localhost]
//				/test:running:0:/home/local/test
//				/manager:stopped:0:/d01/test/test`
//			r, err := handleTomcatStatus(in)
//			assert.NoError(t, err)
//			assert.True(t, r["test"] == "running")
//			assert.True(t, r["manager"] == "stopped")
//		})
//		t.Run("not running", func(t *testing.T) {
//			in := `ERROR - Listed applications for virtual host [localhost]
//				/test:running:0:test
//				/manager:stopped:0:test`
//			r, err := handleTomcatStatus(in)
//			assert.Errorf(t, err, err.Error())
//			assert.Equal(t, ErrTomcatService, err)
//			assert.Nil(t, r)
//		})
//		t.Run("error parse", func(t *testing.T) {
//			in := `OK - Listed applications for virtual host [localhost]
//				/test
//				/manager`
//			r, err := handleTomcatStatus(in)
//			assert.Error(t, err)
//			assert.Equal(t, err, ErrTomcatParse)
//			assert.Nil(t, r)
//		})
//	})
//	t.Run("AddDataWars", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			war := make(map[string]string)
//			info := make(map[string]string)
//			war["test"] = "running"
//			war["manager"] = "running"
//			info["test"] = "19.01"
//			r, err := addDataWars(war, info)
//			assert.NoError(t, err)
//			assert.True(t, len(r) == 2)
//			for _, res := range r {
//				if res.Status == "test" {
//					assert.Equal(t, res, Result{Status: "test", Service: "tomcat", Result: "running", Alarm: false, Tooltip: "19.01"})
//				}
//			}
//		})
//	})
//	t.Run("HandleTomcatInfo", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			in := `02.12.2022 12:20:52
//				wap-logging
//				0.1.0-66a9c9a`
//			r, err := handleTomcatInfo(in)
//			assert.NoError(t, err)
//			assert.True(t, r["wap-logging"] == "02.12.2022 12:20:52, ver: 0.1.0-66a9c9a")
//		})
//		t.Run("small data", func(t *testing.T) {
//			in := `02.12.2022
//					wap-log`
//			r, err := handleTomcatInfo(in)
//			assert.Error(t, err)
//			assert.Equal(t, err, ErrTomcatData)
//			assert.Nil(t, r)
//		})
//		t.Run("not full data", func(t *testing.T) {
//			in := `02.12.2022
//					wap-log`
//			r, err := handleTomcatInfo(in)
//			assert.Error(t, err)
//			assert.Equal(t, err, ErrTomcatData)
//			assert.Nil(t, r)
//		})
//	})
//	t.Run("HandleDocker", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			in := `{"name":"test", "status":"Up 6 days"}`
//			r, err := handleDocker(in)
//			assert.NoError(t, err)
//			assert.Equal(t, r[0], Result{Status: "test", Service: "docker", Result: "running", Tooltip: "", Alarm: false})
//		})
//		t.Run("json error", func(t *testing.T) {
//			in := `{"name:"error", "st":"err"}`
//			r, err := handleDocker(in)
//			assert.Error(t, err)
//			assert.Nil(t, r)
//		})
//	})
//	t.Run("HandleCeph", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			in := "HEALTH_OK Ceph"
//			r, err := handleCeph(in)
//			assert.NoError(t, err)
//			assert.NotNil(t, r)
//			assert.True(t, r[0] == Result{Status: "Ceph: OK", Service: "ceph", Result: "running", Tooltip: "", Alarm: false})
//		})
//		t.Run("error", func(t *testing.T) {
//			in := "NOT_OK Ceph"
//			r, err := handleCeph(in)
//			assert.Error(t, err)
//			assert.Nil(t, r)
//			assert.Equal(t, err, ErrCephCheck)
//		})
//	})
//}
//
