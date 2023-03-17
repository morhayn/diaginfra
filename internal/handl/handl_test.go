package handl

//import (
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestHandleTomcat(t *testing.T) {
//	t.Run("HandleTomcat", func(t *testing.T) {
//		t.Run("true", func(t *testing.T) {
//			in := `OK - Listed applications for virtual host [localhost]
//				/test:running:0:test
//				/manager:stopped:0:test`
//			r, err := handleTomcatStatus(in)
//			assert.NoError(t, err)
//			assert.True(t, r["test"] == "running")
//			assert.True(t, r["manager"] == "stopped")
//		})
//		t.Run("path to out", func(t *testing.T) {
//			in := `OK - Listed applications for virtual host [localhost]
//				/test:running:0:/home/local/test
//				/manager:stopped:0:/d01/test/test`
//			r, err := handleTomcatStatus(in)
//			assert.NoError(t, err)
//			assert.True(t, r["test"] == "running")
//			assert.True(t, r["manager"] == "stopped")
//		})
//		t.Run("not running", func(t *testing.T) {
//			in := `ERROR - Listed applications for virtual host [localhost]
//				/test:running:0:test
//				/manager:stopped:0:test`
//			r, err := handleTomcatStatus(in)
//			assert.Errorf(t, err, err.Error())
//			assert.Equal(t, ErrTomcatService, err)
//			assert.Nil(t, r)
//		})
//		t.Run("error parse", func(t *testing.T) {
//			in := `OK - Listed applications for virtual host [localhost]
//				/test
//				/manager`
//			r, err := handleTomcatStatus(in)
//			assert.Error(t, err)
//			assert.Equal(t, err, ErrTomcatParse)
//			assert.Nil(t, r)
//		})
//	})
//	t.Run("AddDataWars", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			war := make(map[string]string)
//			info := make(map[string]string)
//			war["test"] = "running"
//			war["manager"] = "running"
//			info["test"] = "19.01"
//			r, err := addDataWars(war, info)
//			assert.NoError(t, err)
//			assert.True(t, len(r) == 2)
//			for _, res := range r {
//				if res.Status == "test" {
//					assert.Equal(t, res, Result{Status: "test", Service: "tomcat", Result: "running", Alarm: false, Tooltip: "19.01"})
//				}
//			}
//		})
//	})
//	t.Run("HandleTomcatInfo", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			in := `02.12.2022 12:20:52
//				wap-logging
//				0.1.0-66a9c9a`
//			r, err := handleTomcatInfo(in)
//			assert.NoError(t, err)
//			assert.True(t, r["wap-logging"] == "02.12.2022 12:20:52, ver: 0.1.0-66a9c9a")
//		})
//		t.Run("small data", func(t *testing.T) {
//			in := `02.12.2022
//					wap-log`
//			r, err := handleTomcatInfo(in)
//			assert.Error(t, err)
//			assert.Equal(t, err, ErrTomcatData)
//			assert.Nil(t, r)
//		})
//		t.Run("not full data", func(t *testing.T) {
//			in := `02.12.2022
//					wap-log`
//			r, err := handleTomcatInfo(in)
//			assert.Error(t, err)
//			assert.Equal(t, err, ErrTomcatData)
//			assert.Nil(t, r)
//		})
//	})
//	t.Run("HandleDocker", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			in := `{"name":"test", "status":"Up 6 days"}`
//			r, err := handleDocker(in)
//			assert.NoError(t, err)
//			assert.Equal(t, r[0], Result{Status: "test", Service: "docker", Result: "running", Tooltip: "", Alarm: false})
//		})
//		t.Run("json error", func(t *testing.T) {
//			in := `{"name:"error", "st":"err"}`
//			r, err := handleDocker(in)
//			assert.Error(t, err)
//			assert.Nil(t, r)
//		})
//	})
//	t.Run("HandleCeph", func(t *testing.T) {
//		t.Run("simple", func(t *testing.T) {
//			in := "HEALTH_OK Ceph"
//			r, err := handleCeph(in)
//			assert.NoError(t, err)
//			assert.NotNil(t, r)
//			assert.True(t, r[0] == Result{Status: "Ceph: OK", Service: "ceph", Result: "running", Tooltip: "", Alarm: false})
//		})
//		t.Run("error", func(t *testing.T) {
//			in := "NOT_OK Ceph"
//			r, err := handleCeph(in)
//			assert.Error(t, err)
//			assert.Nil(t, r)
//			assert.Equal(t, err, ErrCephCheck)
//		})
//	})
//}
//
