package modules

import "fmt"

type Mongodb struct{}

func (t Mongodb) RunString(arg ...string) (string, error) {
	cmd := `mongo -u %s -p "%s"  --eval 'db.stats()'`
	return fmt.Sprintf(cmd, iface(arg)...), nil
}
func (t Mongodb) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		log = fmt.Sprintf("sudo tail -n %d %s", count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Mongodb %s", arg)
}

func (t Mongodb) Handler(in string) ([]Result, error) {
	res := []Result{}
	res = append(res, Result{
		Service: "Kafka",
		Output:  in,
		Status:  "running",
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
