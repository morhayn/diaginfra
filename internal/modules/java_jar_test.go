package modules

import (
	"testing"
)

func TestJarRunString(t *testing.T) {
	jar := Jar{}
	res, err := jar.RunString("jar-service")
	if err != nil {
		t.Fatal(err)
	}
	if res != "sudo systemctl is-active jar-service" {
		t.Fatal("result not right ", res)
	}
}
func TestJarLogs(t *testing.T) {
	jar := Jar{}
	t.Run("simple ", func(t *testing.T) {
		res, err := jar.Logs(300, "/log/jar/", "service")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo test -f /log/jar/service.log && sudo tail -n 300 /log/jar/service.log" {
			t.Fatal("result not right ", res)
		}
	})
	t.Run("short in", func(t *testing.T) {
		_, err := jar.Logs(300)
		if err == nil {
			t.Fatal("no error to short in")
		}
	})
}
func TestJarHandler(t *testing.T) {
	jar := Jar{}
	t.Run("running", func(t *testing.T) {
		in := "active"
		res, err := jar.Handler(in)
		if err != nil {
			t.Fatal(err)
		}
		if res[0].Status != "running" {
			t.Fatal("result not right ", res[0])
		}
	})
	t.Run("not running", func(t *testing.T) {
		in := "no active"
		res, err := jar.Handler(in)
		if err != nil {
			t.Fatal(err)
		}
		if res[0].Status != "failed" {
			t.Fatal("result not right ", res[0])
		}
	})
}
