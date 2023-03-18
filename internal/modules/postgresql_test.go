package modules

import (
	"testing"
)

func TestPostgresRunString(t *testing.T) {
	postgres := Postgresql{}
	res, err := postgres.RunString()
	if err != nil {
		t.Fatal(err)
	}
	if res != "sudo pg_lsclusters | awk 'FNR > 1 {print $4}'" {
		t.Fatal("answer not right ", res)
	}
}
func TestPostgresLogs(t *testing.T) {
	postgre := Postgresql{}
	t.Run("simple", func(t *testing.T) {
		res, err := postgre.Logs(300, "/log/postgres.log")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo tail -n 300 /log/postgres.log" {
			t.Fatal("answer not right ", res)
		}
	})
	t.Run("error short arg in", func(t *testing.T) {
		_, err := postgre.Logs(300)
		if err == nil {
			t.Fatal("no error to short in")
		}
	})
}
func TestPostgresHandler(t *testing.T) {}
