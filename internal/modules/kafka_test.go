package modules

import (
	"testing"
)

func TestKafkaRunString(t *testing.T) {
	kafka := Kafka{}
	res, err := kafka.RunString()
	if err != nil {
		t.Fatal(err)
	}
	if res != "export KAFAK_OPTS='-Djava.security.auth.login.config=/etc/kafka/kafka_jaas.conf'; /d01/kafka/bin/kafka-topics.sh --list --zookeeper localhost:2181" {
		t.Fatal("result not right ", res)
	}
}
func TestKafkaLogs(t *testing.T) {
	kafka := Kafka{}
	t.Run("simple ", func(t *testing.T) {
		res, err := kafka.Logs(300, "/log/kafka.log")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo test -f /log/kafka.log && sudo tail -n 300 /log/kafka.log" {
			t.Fatal("result not right ", res)
		}
	})
	t.Run("short in", func(t *testing.T) {
		_, err := kafka.Logs(300)
		if err == nil {
			t.Fatal("no error to short in")
		}
	})
}
func TestKafkaHandler(t *testing.T) {}
