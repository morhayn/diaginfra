package modules

import (
	"encoding/json"
	"fmt"
)

type Elastic struct{}

func (t *Elastic) RunString(arg ...string) (string, error) {
	cmd := "curl -X GET http://127.0.0.1:9200/_cluster/health"
	return fmt.Sprintf(cmd, arg), nil
}

func (t *Elastic) Handler(in string) ([]Result, error) {
	var res = []Result{}
	var elastic = El{}
	err := json.Unmarshal([]byte(in), &elastic)
	if err != nil {
		return nil, err
	}
	result := fmt.Sprintf("ELASTIC: %s NODES: %v  STATUS: %s  Waiting in QUEUE: %v",
		elastic.Cluster_name, elastic.Number_of_nodes, elastic.Status, elastic.Task_max_waiting_in_queue_millis)
	res = append(res, Result{
		Service: "Elastic",
		Status:  result,
		Result:  "running",
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
