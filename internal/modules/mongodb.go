package modules

import (
	"fmt"

	"github.com/morhayn/diaginfra/internal/global"
)

type Mongodb struct{}

func (t Mongodb) RunString(arg ...string) (string, error) {
	cmd := `mongo -u %s -p "%s"  --eval 'db.serverStatus()'`
	return fmt.Sprintf(cmd, iface(arg)...), nil
}
func (t Mongodb) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		fileTest := fmt.Sprintf("sudo test -f %s &&", arg[0])
		log = fmt.Sprintf("%s sudo tail -n %d %s", fileTest, count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Mongodb %s", arg)
}

func (t Mongodb) Handler(in string) ([]global.Result, error) {
	res := []global.Result{}
	res = append(res, global.Result{
		Service: "Kafka",
		Output:  in,
		Status:  "running",
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
