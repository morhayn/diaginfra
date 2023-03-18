package modules

import (
	"testing"
)

func TestCephRunString(t *testing.T) {
	ceph := Ceph{}
	res, err := ceph.RunString()
	if err != nil {
		t.Fatal(err)
	}
	if res != "sudo ceph status | awk '/health/ {print $2}'" {
		t.Fatal("answer not right ", res)
	}
}
func TestCephLogs(t *testing.T) {
	ceph := Ceph{}
	t.Run("simple", func(t *testing.T) {
		res, err := ceph.Logs(300, "/log/ceph.log")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo tail -n 300 /log/ceph.log" {
			t.Fatal(" answer not right ", res)
		}
	})
	t.Run("err short in", func(t *testing.T) {
		_, err := ceph.Logs(300)
		if err == nil {
			t.Fatal("not error to short in")
		}
	})
}
func TestCephHandler(t *testing.T) {
	ceph := Ceph{}
	t.Run("simple", func(t *testing.T) {
		in := `HEALTH_OK`
		res, err := ceph.Handler(in)
		if err != nil {
			t.Fatal(err)
		}
		if res[0].Status != "running" {
			t.Fatal("wrog status ", res[0].Status)
		}
	})
	t.Run("ceph warn", func(t *testing.T) {
		in := `HEALTH_WARN`
		_, err := ceph.Handler(in)
		if err == nil {
			t.Fatal(" not error on failed status ")
		}
	})
}
