package modules

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/morhayn/diaginfra/internal/global"
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
func (t Docker) Handler(in string) ([]global.Result, error) {
	var docker = Dock{}
	var res = []global.Result{}
	lines := strings.Split(in, "\n")
	for _, obj := range lines {
		if strings.TrimSpace(obj) != "" {
			err := json.Unmarshal([]byte(obj), &docker)
			if err != nil {
				return nil, err
			}
			alarm := false
			status := "running"
			if !(strings.HasPrefix(docker.Status, "Up")) {
				alarm = true
				status = "failed"
			}
			res = append(res, global.Result{
				Service: "Docker",
				Output:  docker.Name,
				Status:  status,
				Alarm:   alarm,
				Tooltip: "",
			})
		}
	}
	sort.Slice(res, func(p, q int) bool {
		return res[p].Status > res[q].Status
	})
	return res, nil
}
