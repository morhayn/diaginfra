package modules

import (
	"fmt"
	"strings"
)

type Postgresql struct{}

func (t *Postgresql) RunString(arg ...string) (string, error) {
	cmd := "pg_lsclusters | awk 'FNR > 1 {print $4}'"
	return fmt.Sprintf(cmd, arg), nil
}

func (t *Postgresql) Handler(in string) ([]Result, error) {
	var res = []Result{}
	if strings.HasPrefix(in, "online") {
		res = append(res, Result{
			Service: "Postgresql",
			Status:  "POSSTGERSQL: OK",
			Result:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		return nil, fmt.Errorf("Error Postgres status %s", in)
	}
	return res, nil
}
