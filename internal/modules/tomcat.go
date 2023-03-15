package modules

import (
	"fmt"
	"strings"
)

type Tomcat struct{}

func (t *Tomcat) RunString(arg ...string) (string, error) {
	cmd := "curl -u %s:%s http://127.0.0.1:%s/manager/text/list"
	return fmt.Sprintf(cmd, arg), nil
}
func (t *Tomcat) Handler(in string) ([]Result, error) {
	res := []Result{}
	if strings.HasPrefix(strings.ToLower(in), "ok") {
		lines := strings.Split(in, "\n")
		for _, line := range lines[1:] {
			arr_out := strings.Split(line, ":")
			if len(arr_out) > 1 {
				name_war := strings.TrimPrefix(strings.TrimSpace(arr_out[0]), "/")
				res = append(res, Result{
					Service: "Tomcat",
					Status:  name_war,
					Result:  arr_out[1],
					Alarm:   false,
					Tooltip: "",
				})
			}
		}
		if len(res) == 0 {
			return nil, fmt.Errorf("Error response count line 0")
		}
		return res, nil
	} else {
		return nil, fmt.Errorf("Error Tomcat service failed response")
	}
}