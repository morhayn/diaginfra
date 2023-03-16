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
	return fmt.Sprintf(cmd, arg), nil
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
