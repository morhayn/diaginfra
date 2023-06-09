package getlog

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/modules"
	"github.com/morhayn/diaginfra/internal/sshcmd"
)

type GetLog struct {
	Host    string `json:"host"`
	Service string `json:"service"`
	Module  string `json:"module"`
	Errors  int    `json:"errors"`
}

func GetErr(Status global.Hosts, loadData global.YumInit, port chport.Cheker, conf sshcmd.Execer) []GetLog {
	var wg_l sync.WaitGroup
	var ch = make(chan GetLog)
	res := []GetLog{}
	go func() {
		for _, host := range Status.Stend {
			if chport.CheckSshPort(host.Ip, conf.GetSshPort(), port) {
				for _, st := range host.Status {
					wg_l.Add(1)
					go func(host global.Host, st global.Result) {
						defer wg_l.Done()
						get := GetLog{
							Host:    host.Ip,
							Service: st.Service,
							Module:  st.Output,
						}
						out := get.errBuildCmd(loadData.Logs, loadData.CountLog, conf)
						ch <- out
					}(host, st)
				}
			}
		}
		wg_l.Wait()
		close(ch)
	}()
	for o := range ch {
		res = append(res, o)
	}
	return res
}

// GetErrors - get count ERROR in logs java programs
func (g GetLog) errBuildCmd(logs map[string]string, count int, conf sshcmd.Execer) GetLog {
	cmd := ""
	x := g.Module == "host-manager"
	if (x) || (g.Module == "manager") {
		g.Errors = 0
		return g
	}
	log, err := g.cmdReadLog(logs, count)
	if err != nil {
		g.Errors = 0
		return g
	}
	awk := `awk 'BEGIN { err = 0 } /ERROR/ { err++ } END { print err }'`
	cmd = fmt.Sprintf(`%s | %s`, log, awk)
	out := strings.TrimSpace(g.runCmd(cmd, conf))
	if out != "" {
		if e, err := strconv.Atoi(out); err == nil {
			g.Errors = e
		} else {
			g.Errors = 999
		}
	}
	return g
}

// GetLogs get tail logs service to displey it user
func (g GetLog) GetLogs(logs map[string]string, count int, conf sshcmd.Execer) string {
	cmd, err := g.cmdReadLog(logs, count)
	if err != nil {
		return "no logs"
	}
	out := g.runCmd(cmd, conf)
	return out
}

// Run ssh command on server
func (g GetLog) runCmd(cmd string, conf sshcmd.Execer) string {
	c := sshcmd.CmdExec{
		Name: "logs",
		Chan: make(chan global.Out),
	}
	c.Cmd = cmd
	go conf.Execute(g.Host, c)
	out := <-c.Chan
	return out.Result
}

// Get command for service log
func (g GetLog) cmdReadLog(logs map[string]string, tail int) (string, error) {
	var cmd string
	var err error
	if g.Service == "" {
		return "", errors.New("Service empty")
	}
	if _, ok := logs[g.Service]; !ok && g.Service != "Docker" {
		return "", errors.New("Not logs path in config")
	}
	if mod, ok := modules.MapService[g.Service]; ok {
		if path, ok := logs[g.Service]; ok {
			cmd, err = mod.Logs(tail, path, g.Module)
			if err != nil {
				return "", err
			}
		} else {
			cmd, err = mod.Logs(tail, g.Module)
			if err != nil {
				return "", err
			}
		}
	}
	return cmd, nil
}
