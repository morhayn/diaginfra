package modules

import "fmt"

type Rabbitmq struct{}

func (t Rabbitmq) RunString(arg ...string) (string, error) {
	cmd := "rabbitmqctl status"
	return fmt.Sprint(cmd), nil
}

func (t Rabbitmq) Handler(in string) ([]Result, error) {
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
