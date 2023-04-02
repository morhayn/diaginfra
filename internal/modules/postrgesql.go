package modules

import (
	"fmt"
	"strings"

	"github.com/morhayn/diaginfra/internal/global"
)

type Postgresql struct{}

func (t Postgresql) RunString(arg ...string) (string, error) {
	cmd := "sudo pg_lsclusters | awk 'FNR > 1 {print $4}'"
	return fmt.Sprint(cmd), nil
}
func (t Postgresql) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		fileTest := fmt.Sprintf("sudo test -f %s &&", arg[0])
		log = fmt.Sprintf("%s sudo tail -n %d %s", fileTest, count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Postgresql %s", arg)
}

func (t Postgresql) Handler(in string) ([]global.Result, error) {
	var res = []global.Result{}
	if strings.HasPrefix(in, "online") {
		res = append(res, global.Result{
			Service: "Postgresql",
			Output:  "POSTGRESQL: OK",
			Status:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		return nil, fmt.Errorf("Error Postgres status %s", in)
	}
	return res, nil
}
