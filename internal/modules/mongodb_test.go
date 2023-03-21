package modules

import (
	"testing"
)

func TestMongoRunString(t *testing.T) {
	mongo := Mongodb{}
	res, err := mongo.RunString("user", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if res != `mongo -u user -p "pass"  --eval 'db.serverStatus()'` {
		t.Fatal("result not right ", res)
	}
}
func TestMongoLogs(t *testing.T) {
	mongo := Mongodb{}
	t.Run("simple ", func(t *testing.T) {
		res, err := mongo.Logs(300, "/log/mongo.log")
		if err != nil {
			t.Fatal(err)
		}
		if res != "sudo tail -n 300 /log/mongo.log" {
			t.Fatal("result not right ", res)
		}
	})
	t.Run("short in", func(t *testing.T) {
		_, err := mongo.Logs(300)
		if err == nil {
			t.Fatal("no error to short in")
		}
	})
}
func TestMongoHandler(t *testing.T) {}
