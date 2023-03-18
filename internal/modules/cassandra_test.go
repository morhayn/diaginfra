package modules

import (
	"strings"
	"testing"
)

func TestCassRunString(t *testing.T) {
	cass := Cassandra{}
	res, err := cass.RunString()
	if err != nil {
		t.Fatal(err)
	}
	if res != "nodetool status" {
		t.Fatal("not right answer ", res)
	}
}
func TestCassLogs(t *testing.T) {
	cass := Cassandra{}
	t.Run("simple test", func(t *testing.T) {
		res, err := cass.Logs(300, "/log/cass.out")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo tail -n 300 /log/cass.out" {
			t.Fatal("not right answer ", res)
		}
	})
	t.Run("error short in", func(t *testing.T) {
		res, err := cass.Logs(300)
		if err == nil {
			t.Fatal("no error ", res, err)
		}
	})
}
func TestCassHandler(t *testing.T) {
	cass := Cassandra{}
	t.Run("simple", func(t *testing.T) {
		in := `Datacentre: datacenter1
		       =======================
			   Status=Up/Down
			   |/ State=Normal/Leaving/Joining/Moving
			   -- Address    Load   Tokens   Owns     Host         ID
			   UN 127.0.0.1  50.0KB  256     10%      a98fgtr567   rack1`
		res, err := cass.Handler(in)
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(res[0].Output) != "UN 127.0.0.1  50.0KB  256     10%      a98fgtr567   rack1" {
			t.Fatal("not right answer ", res[0].Output)
		}
	})
	t.Run("in short", func(t *testing.T) {
		in := `Datacentre: datacenter1
		       =======================
			   Status=Up/Down`
		_, err := cass.Handler(in)
		if err == nil {
			t.Fatal(" no error for short in ")
		}
	})
}
