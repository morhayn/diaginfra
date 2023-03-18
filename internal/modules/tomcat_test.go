package modules

import (
	"fmt"
	"testing"
)

func TestTomcatHandler(t *testing.T) {
	tomcat := Tomcat{}
	t.Run("right data", func(t *testing.T) {
		in := `OK - Listed application for virtual host localhost
		/manager:running:0:examples
		/test:running:0:examples
		/error:stopped:0:examples`
		res, err := tomcat.Handler(in)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range res {
			if r.Output == "manager" && r.Status != "running" {
				t.Fatal("error Tomcat check - manager not running result")
			}
			if r.Output == "test" && r.Status != "running" {
				t.Fatal("error Tomcat check - test not running result")
			}
			if r.Output == "error" && r.Status != "stopped" {
				t.Fatal("error Tomcat check - test not sttoped result")
			}
		}
	})
	t.Run("test error empty virtual host in", func(t *testing.T) {
		in := "OK - onli one string in output"
		_, err := tomcat.Handler(in)
		if err == nil {
			t.Fatal("not get error to not virtual hosts")
		}
	})
	t.Run("test  not OK status", func(t *testing.T) {
		in := `ERR - Listed application for virtual host localhost
		/manager:running:0:examples
		/test:running:0:examples
		/error:stopped:0:examples`
		_, err := tomcat.Handler(in)
		if err == nil {
			t.Fatal("Not get error to status not OK")
		}

	})
}
func TestTomcatLogs(t *testing.T) {
	tomcat := Tomcat{}
	t.Run("test right input", func(t *testing.T) {
		log, err := tomcat.Logs(300, "/test/", "manager")
		if err != nil {
			t.Fatal(err)
		}
		if log != "sudo tail -n 300 /test/manager.log" {
			t.Fatal("Command not right ", log)
		}
	})
	t.Run("test input short", func(t *testing.T) {
		log, err := tomcat.Logs(300, "/test/")
		if err == nil {
			fmt.Println("not right output", err, log)
		}
	})
}
func TestTomcatRunString(t *testing.T) {
	tomcat := Tomcat{}
	res, err := tomcat.RunString("user", "pass", "8080")
	if err != nil {
		t.Fatal(err)
	}
	if res != "curl -u user:pass http://127.0.0.1:8080/manager/text/list" {
		t.Fatal("output not right ", res)
	}
}
