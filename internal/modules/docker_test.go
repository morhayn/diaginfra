package modules

import (
	"testing"
)

func TestDockerRunString(t *testing.T) {
	docker := Docker{}
	res, err := docker.RunString()
	if err != nil {
		t.Fatal(err)
	}
	if res != `sudo docker ps --format '{"name":"{{.Names}}", "status":"{{.Status}}"}'` {
		t.Fatal("result not right ", res)
	}
}
func TestDockerLogs(t *testing.T) {
	docker := Docker{}
	t.Run("simple ", func(t *testing.T) {
		res, err := docker.Logs(300, "test")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo docker logs --tail 300 test" {
			t.Fatal("result not right ", res)
		}
	})
	t.Run("short in", func(t *testing.T) {
		_, err := docker.Logs(300)
		if err == nil {
			t.Fatal("no error to short in")
		}
	})
}
func TestDockerHandler(t *testing.T) {
	docker := Docker{}
	in := `{"name":"cassandra","status":"Up 6 seconds ago"}
	{"name":"elasticsearch","status":"Restarted 3 seconds ago"}`
	res, err := docker.Handler(in)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		if r.Output == "cassandra" && r.Status != "running" {
			t.Fatal("RUNNING result not right ", r)
		}
		if r.Output == "elasticsearch" && r.Status != "failed" {
			t.Fatal("FAILED result not right ", r)
		}
	}
}
