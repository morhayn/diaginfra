package upload

import (
	"os"
	"testing"

	"github.com/morhayn/diaginfra/internal/churl"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/sshcmd"
)

var cmdList []string

type MockExec struct {
}

func (m MockExec) Execute(ip string, cmd sshcmd.CmdExec) {
	cmdList = append(cmdList, "exec "+cmd.Cmd)
	cmd.Chan <- sshcmd.NewOut("test", "test", cmd.Cmd)
}
func (m MockExec) GetSshPort() string {
	return "22"
}
func (m MockExec) Scp(ip, src, dest string) error {
	cmdList = append(cmdList, "scp "+ip+" "+src+" "+dest)
	return nil
}

type MockPort struct{}

func (p MockPort) Check(p1 string, p2 string, ch chan global.Port) {
	ch <- global.Port{
		Port:   "22",
		Status: "success",
	}
}

var MockStatus = global.Hosts{
	ListUrls: []churl.Url{},
	Stend: []global.Host{
		{
			Name: "test-nginx",
			Ip:   "10.2.3.50",
			ListPort: []global.Port{
				{
					Port:   "22",
					Status: "success",
				},
				{
					Port:   "80",
					Status: "success",
				},
			},
			ListSsh: []global.Out{
				{
					Name:    "nginx",
					Result:  "active",
					PrgName: "nginx",
				},
			},
			Status: []global.Result{},
		},
		{
			Name: "test-tomcat",
			Ip:   "10.2.3.100",
			ListPort: []global.Port{
				{
					Port:   "22",
					Status: "success",
				},
				{
					Port:   "8080",
					Status: "success",
				},
			},
			ListSsh: []global.Out{
				{
					Name:    "tomcat",
					Result:  "active",
					PrgName: "tomcat",
				},
			},
			Status: []global.Result{
				{
					Service: "Tomcat",
					Output:  "mod-fs",
					Status:  "running",
					Alarm:   false,
					Tooltip: "10.2.0",
				},
				{
					Service: "Tomcat",
					Output:  "mod-gateway",
					Status:  "running",
					Alarm:   false,
					Tooltip: "10.2.0",
				},
			},
		},
	},
}

func TestSaveLoadState(t *testing.T) {
	statusFile = "/tmp/status.test"
	t.Run("Save status", func(t *testing.T) {
		err := saveStatus(MockStatus)
		if err != nil {
			t.Fatal(err)
		}
		st, err := os.Stat(statusFile)
		if err != nil {
			t.Fatal(err)
		}
		if !st.Mode().IsRegular() && st.Size() == 0 {
			t.Fatal("File empty or not file")
		}
	})
	t.Run("Load status", func(t *testing.T) {
		st, err := loadStatus()
		if err != nil {
			t.Fatal(err)
		}
		if len(st.Stend) != 2 {
			t.Fatal("Error count load servers")
		}
		for _, stend := range st.Stend {
			for _, s := range stend.Status {
				if s.Service == "Tomcat" && s.Status != "running" {
					t.Fatal("Wrong load status")
				}
			}
		}
	})
	t.Run("Read upload conf", func(t *testing.T) {
		upl, err := readConfig("/tmp/upload.yaml")
		if err != nil {
			t.Fatal(err)
		}
		t.Log("!!! NFS", upl.NfsServer)
		t.Log("!!! ", upl)
	})
	err := os.Remove(statusFile)
	if err != nil {
		t.Fatal(err)
	}
}
func TestCopyWar(t *testing.T) {
	statusFile = "/tmp/status.test"
	uploadFile = "/tmp/upload.yaml"
	uploadDir = "/d01/wcs-app/"
	err := CopyWars(MockStatus, "pre", MockExec{})
	if err != nil {
		t.Fatal("FFFFF", err)
	}
	// t.Fatal("!!!!!!", len(cmdList))
	for _, cmd := range cmdList {
		t.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", cmd)
	}
}
