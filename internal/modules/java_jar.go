package modules

import (
	"fmt"
	"strings"
)

type Jar struct{}

func (t Jar) RunString(arg ...string) (string, error) {
	cmd := "sudo systemctl is-active %s"
	return fmt.Sprintf(cmd, iface(arg)), nil
}
func (t Jar) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 1 {
		log = fmt.Sprintf("sudo tail -n %d %s%s.log", count, arg[0], arg[1])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Kafka %s", arg)
}

func (t Jar) Handler(in string) ([]Result, error) {
	res := []Result{}
	st := "failed"
	if strings.TrimSpace(in) == "active" {
		st = "running"
	}
	res = append(res, Result{
		Service: "Jar",
		Output:  in,
		Status:  st,
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
