package modules

import (
	"fmt"

	"github.com/morhayn/diaginfra/internal/global"
)

type Kafka struct{}

func (t Kafka) RunString(arg ...string) (string, error) {
	cmd := "export KAFAK_OPTS='-Djava.security.auth.login.config=/etc/kafka/kafka_jaas.conf'; /d01/kafka/bin/kafka-topics.sh --list --zookeeper localhost:2181"
	return fmt.Sprint(cmd), nil
}
func (t Kafka) Logs(count int, arg ...string) (string, error) {
	log := ""
	if len(arg) > 0 {
		fileTest := fmt.Sprintf("sudo test -f %s &&", arg[0])
		log = fmt.Sprintf("%s sudo tail -n %d %s", fileTest, count, arg[0])
		return log, nil
	}
	return "", fmt.Errorf("not path to log Kafka %s", arg)
}

func (t Kafka) Handler(in string) ([]global.Result, error) {
	res := []global.Result{}
	res = append(res, global.Result{
		Service: "Kafka",
		Output:  in,
		Status:  "running",
		Alarm:   false,
		Tooltip: "",
	})
	return res, nil
}
