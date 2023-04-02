package modules

import (
	"encoding/json"
	"fmt"

	"github.com/morhayn/diaginfra/internal/global"
)

type Hazelcast struct{}

func (t Hazelcast) RunString(arg ...string) (string, error) {
	cmd := `curl --data "%s&%s" --silent "http://127.0.0.1:5701/hazelcast/rest/management/cluster/state"`
	return fmt.Sprintf(cmd, iface(arg)...), nil
}
func (t Hazelcast) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		fileTest := fmt.Sprintf("sudo test -f %s &&", arg[0])
		log = fmt.Sprintf("%s sudo tail -n %d %s", fileTest, count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Hazelcast %s", arg)
}

func (t Hazelcast) Handler(in string) ([]global.Result, error) {
	var res = []global.Result{}
	var hazel = Hazel{}
	err := json.Unmarshal([]byte(in), &hazel)
	if err != nil {
		return nil, err
	}
	if hazel.State == "active" && hazel.Status == "success" {
		res = append(res, global.Result{
			Service: "Hazelcast",
			Output:  "HAZEL: " + hazel.Status,
			Status:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		res = append(res, global.Result{
			Service: "Hazelcast",
			Output:  "HAZEL:" + hazel.Status,
			Status:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	}
	return res, nil
}
