package modules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/morhayn/diaginfra/internal/global"
)

type Tomcat struct{}

func (t Tomcat) RunString(arg ...string) (string, error) {
	cmd := "curl -u %s:%s http://127.0.0.1:%s/manager/text/list"
	return fmt.Sprintf(cmd, iface(arg)...), nil
}
func (t Tomcat) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 1 {
		if arg[1] == "Tomcat" {
			arg[1] = "catalina.out"
			fileTest := fmt.Sprintf("sudo test -f %s%s &&", arg[0], arg[1])
			log = fmt.Sprintf("%s sudo tail -n %d %s%s", fileTest, count, arg[0], arg[1])
		} else {
			fileTest := fmt.Sprintf("sudo test -f %s%s.log &&", arg[0], arg[1])
			log = fmt.Sprintf("%s sudo tail -n %d %s%s.log", fileTest, count, arg[0], arg[1])
		}
		return log, nil
	}
	return "", fmt.Errorf("arg length less 2")
}
func (t Tomcat) Handler(in string) ([]global.Result, error) {
	res := []global.Result{}
	if strings.HasPrefix(strings.ToLower(in), "ok") {
		lines := strings.Split(in, "\n")
		//First string in array "OK - Listed applications for virtual host localhost" its only status tomcat
		for _, line := range lines[1:] {
			arr_out := strings.Split(line, ":")
			if len(arr_out) > 1 {
				name_war := strings.TrimPrefix(strings.TrimSpace(arr_out[0]), "/")
				res = append(res, global.Result{
					Service: "Tomcat",
					Output:  name_war,
					Status:  arr_out[1],
					Alarm:   false,
					Tooltip: "",
				})
			}
		}
		if len(res) == 0 {
			return nil, fmt.Errorf("Error response count line 0")
		}
		sort.Slice(res, func(p, q int) bool {
			return res[p].Status > res[q].Status
		})
		return res, nil
	} else {
		return nil, fmt.Errorf("Error Tomcat service failed response")
	}
}
