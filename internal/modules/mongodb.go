package modules

import "fmt"

type Mongodb struct{}

func (t Mongodb) RunString(arg ...string) (string, error) {
	cmd := `mongo -u %s -p "%s"  --eval 'db.stats()'`
	return fmt.Sprintf(cmd, iface(arg)...), nil
}

func (t Mongodb) Handler(in string) ([]Result, error) {
	res := []Result{}
	res = append(res, Result{
		Service: "Kafka",
		Status:  in,
		Result:  "running",
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
