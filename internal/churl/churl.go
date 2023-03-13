package churl

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	wg_url  sync.WaitGroup
	httpGet = http.Get
)

type Churler interface {
	checkHttp(string, chan Url)
	checkHttps(string, chan Url)
}
type Url struct {
	Url    string `json:"url"`
	Status int    `json:"status"`
}

// Check response to http request if status code 200 true.
func (u Url) checkHttp(url string, res chan Url) {
	u.Url = url
	r, err := httpGet(u.Url)
	if err != nil {
		u.Status = 0
		res <- u
		return
	}
	u.Status = r.StatusCode
	res <- u
}

// Check response to https request if status code 200 true.
func (u Url) checkHttps(url string, res chan Url) {
	u.Url = url
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	r, err := client.Get(u.Url)
	if err != nil {
		u.Status = 0
		res <- u
		return
	}
	defer r.Body.Close()
	u.Status = r.StatusCode
	res <- u
}

// Run gourutine for all check urls
func CheckUrl(urls []string, u Churler) []Url {
	res := make(chan Url)
	result := []Url{}
	go func() {
		for _, url := range urls {
			wg_url.Add(1)
			fmt.Println(url)
			go func(url string) {
				defer wg_url.Done()
				if strings.HasPrefix(strings.ToLower(url), `https://`) {
					u.checkHttps(url, res)
				} else if strings.HasPrefix(strings.ToLower(url), `http://`) {
					u.checkHttp(url, res)
				}
			}(url)
		}
		wg_url.Wait()
		close(res)
	}()
	for r := range res {
		fmt.Print(r)
		result = append(result, r)
	}
	return result
}
