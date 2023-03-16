package getlog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/morhayn/diaginfra/internal/sshcmd"
)

type GetLog struct {
	Host    string `json:"host"`
	Service string `json:"service"`
	Module  string `json:"module"`
	Errors  int    `json:"errors"`
}

// Get count ERROR in logs java programs
func (g GetLog) GetErrors(logs map[string]string, count int, conf sshcmd.Execer) GetLog {
	path, ok := logs[g.Service]
	if !ok {
		return g
	}
	cmd := ""
	if g.Service == "docker" {
		cmd = fmt.Sprintf(`sudo docker logs --tail %d %s | awk '/ERROR/ { err++ } END { print err }'`, count, g.Module)
	} else {
		cmd = fmt.Sprintf(`sudo tail -n %d %s%s.log | awk '/ERROR/ { err++ } END { print err }'`, count, path, g.Module)
	}
	out := strings.TrimSpace(g.runCmd(cmd, conf))
	if out != "" {
		if e, err := strconv.Atoi(out); err == nil {
			g.Errors = e
		} else {
			fmt.Println(err)
		}
	}
	return g
}

// Get tail logs service to displey it user
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
	if g.Service == "" {
		return "", errors.New("Service empty")
	}
	if _, ok := logs[g.Service]; !ok && g.Service != "docker" {
		return "", errors.New("Not logs path in config")
	}
	switch g.Service {
	case "Jar":
		cmd = fmt.Sprintf("sudo tail -n %d %s%s.log", tail, logs["Jar"], g.Module)
	case "tomcat":
		cmd = fmt.Sprintf("sudo tail -n %d %s%s.log", tail, logs["tomcat"], g.Module)
	case "Tomcat":
		cmd = fmt.Sprintf("sudo tail -n %d %scatalina.out", tail, logs["tomcat"])
	case "docker":
		cmd = fmt.Sprintf("sudo docker logs --tail %d %s", tail, g.Module)
	default:
		cmd = fmt.Sprintf("sudo tail -n %d %s", tail, logs[g.Service])
	}
	return cmd, nil
}
