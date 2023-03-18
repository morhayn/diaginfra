package modules

import (
	"fmt"
	"strings"
)

type Cassandra struct{}

func (t Cassandra) RunString(arg ...string) (string, error) {
	cmd := "nodetool status"
	return fmt.Sprint(cmd), nil
}

func (t Cassandra) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		log = fmt.Sprintf("sudo tail -n %d %s", count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Cassandra %s", arg)
}
func (t Cassandra) Handler(in string) ([]Result, error) {
	var res = []Result{}
	spl_res := strings.Split(in, "\n")
	if len(spl_res) > 4 {
		res = append(res, Result{
			Service: "Cassandra",
			Output:  spl_res[5],
			Status:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		return []Result{}, fmt.Errorf("Cassandre parse failed")
	}
	return res, nil
}
