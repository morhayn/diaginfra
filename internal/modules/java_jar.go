package modules

import (
	"fmt"
	"strings"

	"github.com/morhayn/diaginfra/internal/global"
)

type Jar struct {
	name string
}

func (t Jar) RunString(arg ...string) (string, error) {
	cmd := "sudo systemctl is-active %s"
	return fmt.Sprintf(cmd, iface(arg)...), nil
}
func (t Jar) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 1 {
		log = fmt.Sprintf("sudo tail -n %d %s%s.log", count, arg[0], arg[1])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Jar %s", arg)
}

func (t Jar) Handler(in string) ([]global.Result, error) {
	res := []global.Result{}
	st := "failed"
	if strings.TrimSpace(in) == "active" {
		st = "running"
	}
	res = append(res, global.Result{
		Service: "Jar",
		Output:  t.name,
		Status:  st,
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
