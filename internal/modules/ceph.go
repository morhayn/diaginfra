package modules

import (
	"fmt"
	"strings"
)

type Ceph struct{}

func (t Ceph) RunString(arg ...string) (string, error) {
	cmd := "sudo ceph status | awk '/health/ {print $2}'"
	return fmt.Sprint(cmd), nil
}

func (t Ceph) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		log = fmt.Sprintf("sudo tail -n %d %s", count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Ceph %s", arg)
}
func (t Ceph) Handler(in string) ([]Result, error) {
	var res = []Result{}
	if strings.HasPrefix(in, "HEALTH_OK") {
		res = append(res, Result{
			Service: "Ceph",
			Status:  "Ceph: OK",
			Result:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		return nil, fmt.Errorf("Error Ceph status %s", in)
	}
	return res, nil
}
