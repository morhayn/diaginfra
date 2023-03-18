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
		if res != "sudo tail -n 300 /log/elastic.log" {
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
func TestElasticHandler(t *testing.T) {}
