package modules

import "fmt"

type Rabbitmq struct{}

func (t Rabbitmq) RunString(arg ...string) (string, error) {
	cmd := "sudo rabbitmqctl status"
	return fmt.Sprint(cmd), nil
}
func (t Rabbitmq) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		log = fmt.Sprintf("sudo tail -n %d %s", count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log RabbitMq %s", arg)
}

func (t Rabbitmq) Handler(in string) ([]Result, error) {
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
