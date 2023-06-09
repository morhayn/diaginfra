package modules

import (
	"encoding/json"
	"fmt"

	"github.com/morhayn/diaginfra/internal/global"
)

type Elastic struct{}

func (t Elastic) RunString(arg ...string) (string, error) {
	cmd := "curl -X GET http://127.0.0.1:9200/_cluster/health"
	return fmt.Sprint(cmd), nil
}
func (t Elastic) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		fileTest := fmt.Sprintf("sudo test -f %s &&", arg[0])
		log = fmt.Sprintf("%s sudo tail -n %d %s", fileTest, count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Elasticsearch %s", arg)
}

func (t Elastic) Handler(in string) ([]global.Result, error) {
	var res = []global.Result{}
	var elastic = El{}
	err := json.Unmarshal([]byte(in), &elastic)
	if err != nil {
		return nil, err
	}
	result := fmt.Sprintf("ELASTIC: %s NODES: %v  STATUS: %s  Waiting in QUEUE: %v",
		elastic.Cluster_name, elastic.Number_of_nodes, elastic.Status, elastic.Task_max_waiting_in_queue_millis)
	res = append(res, global.Result{
		Service: "Elastic",
		Output:  result,
		Status:  "running",
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
