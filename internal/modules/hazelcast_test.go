package modules

import (
	"testing"
)

func TestHazelcastRunString(t *testing.T) {
	hazel := Hazelcast{}
	res, err := hazel.RunString("user", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if res != `curl --data "user&pass" --silent "http://127.0.0.1:5701/hazelcast/rest/management/cluster/state"` {
		t.Fatal("result not right ", res)
	}
}
func TestHazelcastLogs(t *testing.T) {
	hazel := Hazelcast{}
	t.Run("simple ", func(t *testing.T) {
		res, err := hazel.Logs(300, "/log/hazel.log")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo test -f /log/hazel.log && sudo tail -n 300 /log/hazel.log" {
			t.Fatal("result not right ", res)
		}
	})
	t.Run("short in", func(t *testing.T) {
		_, err := hazel.Logs(300)
		if err == nil {
			t.Fatal("no error to short in")
		}
	})
}
func TestHazelcastHandler(t *testing.T) {}
