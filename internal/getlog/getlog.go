package getlog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/morhayn/diaginfra/internal/modules"
	"github.com/morhayn/diaginfra/internal/sshcmd"
)

type GetLog struct {
	Host    string `json:"host"`
	Service string `json:"service"`
	Module  string `json:"module"`
	Errors  int    `json:"errors"`
}

// GetErrors - get count ERROR in logs java programs
func (g GetLog) GetErrors(logs map[string]string, count int, conf sshcmd.Execer) GetLog {
	path, ok := logs[g.Service]
	if !ok {
		return g
	}
	cmd := ""
	x := g.Module == "host-manager"
	if (x) || (g.Module == "manager") {
		g.Errors = 0
		return g
	}
	awk := `awk 'BEGIN { err = 0 } /ERROR/ { err++ } END { print err }'`
	if g.Service == "Docker" {
		cmd = fmt.Sprintf(`sudo docker logs --tail %d %s | %s`, count, g.Module, awk)
	} else if g.Service == "Tomcat" || g.Service == "Jar" {
		fileTest := fmt.Sprintf("sudo test -f %s%s.log && ", path, g.Module)
		cmd = fmt.Sprintf(`%s sudo tail -n %d %s%s.log | %s`, fileTest, count, path, g.Module, awk)
	} else {
		fileTest := fmt.Sprintf("sudo test -f %s && ", path)
		cmd = fmt.Sprintf(`%s sudo tail -n %d %s | %s`, fileTest, count, path, awk)
	}
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
		Chan: make(chan sshcmd.Out),
	}
	c.Cmd = cmd
	go conf.Execute(g.Host, c)
	out := <-c.Chan
	return out.Result
}

// Parse log to map[LEVEL]COUNT
func parse(log string) map[string]int {
	regPatern := `.*(?P<level>(DEBUG|ERROR|WARN|INFO)+) .*`
	re := regexp.MustCompile(regPatern)
	parsedMap := make(map[string]int)
	lines := strings.Split(log, "\n")
	for _, l := range lines {
		match := re.FindStringSubmatch(l)
		if match != nil {
			parsedMap[match[1]] += 1
		}
	}
	return parsedMap
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
