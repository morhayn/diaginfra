package modules

import (
	"encoding/json"
	"fmt"
)

type Hazelcast struct{}

func (t Hazelcast) RunString(arg ...string) (string, error) {
	cmd := `curl --data "%s&%s" --silent "http://127.0.0.1:5701/hazelcast/rest/management/cluster/state"`
	return fmt.Sprintf(cmd, iface(arg)...), nil
}

func (t Hazelcast) Handler(in string) ([]Result, error) {
	var res = []Result{}
	var hazel = Hazel{}
	err := json.Unmarshal([]byte(in), &hazel)
	if err != nil {
		return nil, err
	}
	if hazel.State == "active" && hazel.Status == "success" {
		res = append(res, Result{
			Service: "Hazelcast",
			Status:  "HAZEL:" + hazel.Status,
			Result:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	} else {
		res = append(res, Result{
			Service: "Hazelcast",
			Status:  "HAZEL:" + hazel.Status,
			Result:  "running",
			Alarm:   false,
			Tooltip: "",
		})
	}
	return res, nil
}
