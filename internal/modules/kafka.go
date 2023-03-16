package modules

import "fmt"

type Kafka struct{}

func (t Kafka) RunString(arg ...string) (string, error) {
	cmd := "export KAFAK_OPTS='-Djava.security.auth.login.config=/etc/kafka/kafka_jaas.conf'; /d01/kafka/bin/kafka-topics.sh --list --zookeeper localhost:2181"
	return fmt.Sprint(cmd), nil
}

func (t Kafka) Handler(in string) ([]Result, error) {
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
