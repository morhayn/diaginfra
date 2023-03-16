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

func (t Cassandra) Handler(in string) ([]Result, error) {
	var res = []Result{}
	spl_res := strings.Split(in, "\n")
	if len(spl_res) > 4 {
		res = append(res, Result{
			Service: "Cassandra",
			Status:  spl_res[5],
			Result:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		return []Result{}, fmt.Errorf("Cassandre parse failed")
	}
	return res, nil
}
