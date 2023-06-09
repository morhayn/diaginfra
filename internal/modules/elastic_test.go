package modules

import (
	"testing"
)

func TestElasticRunString(t *testing.T) {
	elastic := Elastic{}
	res, err := elastic.RunString()
	if err != nil {
		t.Fatal(err)
	}
	if res != "curl -X GET http://127.0.0.1:9200/_cluster/health" {
		t.Fatal("not wrong answer ", res)
	}
}
func TestElasticLogs(t *testing.T) {
	elastic := Elastic{}
	t.Run("simple", func(t *testing.T) {
		res, err := elastic.Logs(300, "/log/elastic.log")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo test -f /log/elastic.log && sudo tail -n 300 /log/elastic.log" {
			t.Fatal("not right answer ", res)
		}
	})
	t.Run("error short in", func(t *testing.T) {
		_, err := elastic.Logs(300)
		if err == nil {
			t.Fatal("no error to short in")
		}
	})
}
func TestElasticHandler(t *testing.T) {
	elastic := Elastic{}
	in := `{"cluster_name":"TestCluster","status":"ok","number_of_nodes": 1,"task_max_waiting_in_queue_millis": 10}`
	res, err := elastic.Handler(in)
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Output != "ELASTIC: TestCluster NODES: 1  STATUS: ok  Waiting in QUEUE: 10" {
		t.Fatal("result not right ", res[0])
	}
}
