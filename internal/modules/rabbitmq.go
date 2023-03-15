package modules

import "fmt"

type Rabbitmq struct{}

func (t *Rabbitmq) RunString(arg ...string) (string, error) {
	cmd := "rabbitmqctl status"
	return fmt.Sprintf(cmd, arg), nil
}

func (t *Rabbitmq) Handler(in string) ([]Result, error) {
	return []Result{}, nil
}
