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
		if res != "sudo tail -n 300 /log/jar/service.log" {
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
func TestJarHandler(t *testing.T) {}
