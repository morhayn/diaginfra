package churl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUrl struct {
	Url    string
	Status int
}

func (m mockUrl) checkHttp(url string, res chan Url) {
	// defer wg_url.Done()
	if url == "http://test" {
		res <- Url{Url: url, Status: 200}
	} else {
		res <- Url{Url: url, Status: 404}
	}
}
func (m mockUrl) checkHttps(url string, res chan Url) {
	// defer wg_url.Done()
	if url == "https://test" {
		res <- Url{Url: url, Status: 200}
	} else {
		res <- Url{Url: url, Status: 404}
	}
}
func TestCheckhttp(t *testing.T) {
	httpGet = func(s string) (*http.Response, error) {
		if s == "test" {
			res := &http.Response{StatusCode: 404}
			return res, nil
		} else {
			return nil, http.ErrServerClosed
		}
	}
	t.Run("true", func(t *testing.T) {
		res := make(chan Url)
		var u Url
		go u.checkHttp("test", res)
		r := <-res
		assert.Equal(t, r, Url{Url: "test", Status: 404})
	})
	t.Run("error", func(t *testing.T) {
		res := make(chan Url)
		var u Url
		go u.checkHttp("error", res)
		r := <-res
		assert.Equal(t, r, Url{Url: "error", Status: 0})
	})
}
func TestCheck_url(t *testing.T) {
	t.Run("check http true", func(t *testing.T) {
		var u mockUrl
		res := CheckUrl([]string{"http://test"}, u)
		assert.Equal(t, res[0].Status, 200)
	})
	t.Run("check https true", func(t *testing.T) {
		var u mockUrl
		res := CheckUrl([]string{"https://test"}, u)
		assert.Equal(t, res[0].Status, 200)
	})
	t.Run("https failed", func(t *testing.T) {
		var u mockUrl
		res := CheckUrl([]string{"https://failed"}, u)
		assert.Equal(t, res[0].Status, 404)
	})
}
