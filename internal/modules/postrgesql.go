package modules

import (
	"fmt"
	"strings"
)

type Postgresql struct{}

func (t Postgresql) RunString(arg ...string) (string, error) {
	cmd := "pg_lsclusters | awk 'FNR > 1 {print $4}'"
	return fmt.Sprint(cmd), nil
}
func (t Postgresql) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		log = fmt.Sprintf("sudo tail -n %d %s", count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Postgresql %s", arg)
}

func (t Postgresql) Handler(in string) ([]Result, error) {
	var res = []Result{}
	if strings.HasPrefix(in, "online") {
		res = append(res, Result{
			Service: "Postgresql",
			Output:  "POSSTGERSQL: OK",
			Status:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		return nil, fmt.Errorf("Error Postgres status %s", in)
	}
	return res, nil
}
