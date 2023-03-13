package handl

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/morhayn/diaginfra/internal/sshcmd"
)

var (
	ErrTomcatService = errors.New("Tomcat not started")
	ErrTomcatParse   = errors.New("Tomcat out not parse")
	ErrTomcatData    = errors.New("Error data from tomcat service")
	ErrCephCheck     = errors.New("Ceph out not HELTH OK")
	ErrPostgresCheck = errors.New("Postgres not running")
	handlers         = map[string]Handlers{
		"Elastic":    handleElastic,
		"Docker":     handleDocker,
		"Postgresql": handlePostgresql,
		"Cassandra":  handleCassandra,
		"Hazelcast":  handleHazelcast,
		"Ceph":       handleCeph,
	}
)

type (
	// Handle   = func(string) (*Result, error)
	Handlers = func(string) ([]Result, error)
)

type Docker struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type Elastic struct {
	Cluster_name                     string  `json:"cluster_name"`
	Status                           string  `json:"status"`
	Timed_out                        bool    `json:"time_out"`
	Number_of_nodes                  int     `json:"number_of_nodes"`
	Number_of_data_nodes             int     `json:"number_of_data_nodes"`
	Active_primary_shards            int     `json:"active_primary_shards"`
	Active_shards                    int     `json:"active_shards"`
	Relocating_shards                int     `json:"relocating_shards"`
	Initializing_shards              int     `json:"initializing_shards"`
	Unassigned_shards                int     `json:"unassigned_shards"`
	Delayed_unassigned_shards        int     `json:"delayed_unassigned_shards"`
	Number_of_pending_tasks          int     `json:"number_of_pending_tasks"`
	Number_of_in_flight_fetch        int     `json:"number_of_in_flight_fetch"`
	Task_max_waiting_in_queue_millis int     `json:"task_max_waiting_in_queue_millis"`
	Active_shards_percent_as_number  float64 `json:"active_shards_percent_as_number"`
}
type Hazelcast struct {
	Status string `json:"status"`
	State  string `json:"state"`
}
type Results struct {
	Res []Result
}
type Result struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Result  string `json:"result"`
	Alarm   bool   `json:"alarm"`
	Tooltip string `json:"tooltip"`
}

// func (r *Results) addResult(o sshcmd.Out, fn Handle) {
// out, err := fn(o.Result)
// if err != nil {
// r.Res = append(r.Res, *resultFail(o.Name))
// } else {
// r.Res = append(r.Res, *out)
// }
// }
func (r *Results) addResults(o sshcmd.Out, fn Handlers) {
	out, err := fn(o.Result)
	if err != nil {
		r.Res = append(r.Res, *resultFail(o.Name))
	} else {
		r.Res = append(r.Res, out...)
	}
}

func newResult(service, status, result, tooltip string, alarm bool) *Result {
	return &Result{
		Service: service,
		Status:  status,
		Result:  result,
		Alarm:   alarm,
		Tooltip: tooltip,
	}
}
func resultFail(name string) *Result {
	return newResult(name, name, "failed", "", true)
}

func HandleResult(list []sshcmd.Out) []Result {
	var result Results
	status := make(map[string]string)
	info := make(map[string]string)
	var err error
	for _, res := range list {
		switch {
		case res.Name == "WarTomcat":
			info, err = handleTomcatInfo(res.Result)
			if err != nil {
				result.Res = append(result.Res, *resultFail(res.Name))
			}
		case res.Name == "Tomcat":
			status, err = handleTomcatStatus(res.Result)
			if err != nil {
				result.Res = append(result.Res, *resultFail(res.Name))
			}
		default:
			if h, ok := handlers[res.Name]; ok {
				result.addResults(res, h)
			} else {
				st := "failed"
				if strings.TrimSpace(res.Result) == "active" {
					st = "running"
				}
				r := Result{
					Service: res.Name,
					Status:  res.PrgName,
					Result:  st,
					Alarm:   false,
					Tooltip: res.Result,
				}
				result.Res = append(result.Res, r)
			}
		}
	}
	if len(status) != 0 {
		res, err := addDataWars(status, info)
		if err != nil {
			result.Res = append(result.Res, *resultFail("Tomcat"))
		} else {
			result.Res = append(result.Res, res...)
		}
	}
	return result.Res
}
func handleCeph(in string) ([]Result, error) {
	var res = []Result{}
	if strings.HasPrefix(in, "HEALTH_OK") {
		res = append(res, *newResult("ceph", "Ceph: OK", "running", "", false))
	} else {
		return nil, ErrCephCheck
	}
	return res, nil
}
func handlePostgresql(in string) ([]Result, error) {
	var res = []Result{}
	if strings.HasPrefix(in, "online") {
		res = append(res, *newResult("postgresql", "POSTGRES: OK", "running", "", false))
	} else {
		return nil, ErrPostgresCheck
	}
	return res, nil
}
func handleHazelcast(in string) ([]Result, error) {
	var res = []Result{}
	var hazel = Hazelcast{}
	err := json.Unmarshal([]byte(in), &hazel)
	if err != nil {
		return nil, err
	}
	if hazel.State == "active" && hazel.Status == "success" {
		res = append(res, *newResult("hazelcast", "HAZEL: "+hazel.Status, "running", "", false))
	} else {
		res = append(res, *newResult("hazelcast", "HAZEL: "+hazel.Status, "running", "", false))
	}
	return res, nil
}
func handleElastic(in string) ([]Result, error) {
	var res = []Result{}
	var elastic = Elastic{}
	err := json.Unmarshal([]byte(in), &elastic)
	if err != nil {
		return nil, err
	}
	result := fmt.Sprintf("ELASTIC: %s NODES: %v  STATUS: %s  Waiting in QUEUE: %v",
		elastic.Cluster_name, elastic.Number_of_nodes, elastic.Status, elastic.Task_max_waiting_in_queue_millis)
	res = append(res, *newResult("elasticsearch", result, "running", "", false))
	return res, nil
}
func handleCassandra(in string) ([]Result, error) {
	var res = []Result{}
	spl_res := strings.Split(in, "\n")
	if len(spl_res) > 4 {
		res = append(res, *newResult("cassandra", spl_res[5], "running", "", false))
	} else {
		res = append(res, *newResult("cassandra", "Cassandra", "failed", "", true))
	}
	return res, nil
}
func handleDocker(in string) ([]Result, error) {
	var docker = Docker{}
	var res = []Result{}
	spl_res := strings.Split(in, "\n")
	for _, obj := range spl_res {
		if strings.TrimSpace(obj) != "" {
			err := json.Unmarshal([]byte(obj), &docker)
			if err != nil {
				return nil, err
			}
			alarm := false
			if !(strings.HasPrefix(docker.Status, "Up")) {
				alarm = true
			}
			res = append(res, *newResult("docker", docker.Name, "running", "", alarm))
		}
	}
	return res, nil
}

func handleTomcatInfo(in string) (map[string]string, error) {
	fmt.Println(in)
	res := make(map[string]string)
	lines := strings.Split(in, "\n")
	if len(lines) == 0 || len(lines) < 3 {
		// fmt.Println(len(lines), len(lines)%3)
		return nil, ErrTomcatData
	}
	for i := 0; i < (len(lines) / 3); i++ {
		name := strings.TrimSpace(lines[i*3+1])
		data := strings.TrimSpace(lines[i*3])
		vers := strings.TrimSpace(lines[i*3+2])
		res[name] = fmt.Sprintf("%s, ver: %s", data, vers)
	}
	return res, nil
}
func handleTomcatStatus(in string) (map[string]string, error) {
	fmt.Println(in)
	res := make(map[string]string)
	if strings.HasPrefix(strings.ToLower(in), "ok") {
		lines := strings.Split(in, "\n")
		for _, line := range lines[1:] {
			arr_out := strings.Split(line, ":")
			if len(arr_out) > 1 {
				name_war := strings.TrimPrefix(strings.TrimSpace(arr_out[0]), "/")
				res[name_war] = arr_out[1]
			}
		}
		if len(res) == 0 {
			return nil, ErrTomcatParse
		}
		return res, nil
	} else {
		return nil, ErrTomcatService
	}
}
func addDataWars(stat, info map[string]string) ([]Result, error) {
	res := []Result{}
	for k, v := range stat {
		alarm := false
		if v != "running" {
			alarm = true
		}
		res = append(res, *newResult("tomcat", k, v, info[k], alarm))
	}
	return res, nil
}
