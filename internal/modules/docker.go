package modules

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Docker struct{}

func (t Docker) RunString(arg ...string) (string, error) {
	cmd := `sudo docker ps --format '{"name":"{{.Names}}", "status":"{{.Status}}"}'`
	return fmt.Sprint(cmd), nil
}
func (t Docker) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		log = fmt.Sprintf("sudo docker logs --tail %d %s", count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Docker %s", arg)
}
func (t Docker) Handler(in string) ([]Result, error) {
	var docker = Dock{}
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
			res = append(res, Result{
				Service: "Docker",
				Status:  docker.Name,
				Result:  "running",
				Alarm:   alarm,
				Tooltip: "",
			})
		}
	}
	sort.Slice(res, func(p, q int) bool {
		return res[p].Result > res[q].Result
	})
	return res, nil
}
