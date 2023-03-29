package webapi

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"sync"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/churl"
	"github.com/morhayn/diaginfra/internal/getlog"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/handl"
	"github.com/morhayn/diaginfra/internal/sshcmd"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	RunOps string
	wg     sync.WaitGroup
	// Username string
	Status global.Hosts
)

type Terminal struct {
	Ip string `json:"ip"`
}

// type Hosts struct {
// ListUrls []churl.Url `json:"list_url"`
// Stend    []Host      `josn:"stand"`
// }
// type Host struct {
// Name     string           `json:"name"`
// Ip       string           `json:"ip"`
// ListPort []chport.Port    `json:"list_port"`
// ListSsh  []sshcmd.Out     `json:"list_ssh"`
// Status   []modules.Result `json:"status"`
// }

// Create new structure for loading config file
func newHost(name, ip string) global.Host {
	return global.Host{
		Name:    name,
		Ip:      ip,
		ListSsh: []global.Out{},
		Status:  []global.Result{},
	}
}

// Check ssh port if ssh port failed not nid run ssh command to server
func checkSshPort(ip, sshPort string, port chport.Cheker) bool {
	if check := chport.CheckPort(ip, []string{sshPort}, port); check[0].Status == "failed" {
		return false
	}
	return true
}

// Run test command to one server
func checkHost(host global.Init, ch chan global.Host, port chport.Cheker, conf sshcmd.Execer) {
	h := newHost(host.Name, host.Ip)
	h.ListPort = chport.CheckPort(host.Ip, host.ListPorts, port)
	if checkSshPort(host.Ip, conf.GetSshPort(), port) {
		srv, prg, _ := sshcmd.Run(host.Ip, host.ListService, conf)
		h.ListSsh = srv
		h.Status = handl.HandleResult(prg)
		sort.Slice(h.Status, func(p, q int) bool {
			return h.Status[p].Status < h.Status[q].Status
		})

	}
	ch <- h
}

// Run gorutine for all servers in config file and grouping result
func serverHandler(loadData global.YumInit, port chport.Cheker, url churl.Churler, conf sshcmd.Execer) global.Hosts {
	result := global.Hosts{}
	ch := make(chan global.Host)
	go func() {
		for _, host := range loadData.Hosts {
			wg.Add(1)
			go func(host global.Init) {
				defer wg.Done()
				checkHost(host, ch, port, conf)
			}(host)
		}
		wg.Wait()
		close(ch)
	}()
	for c := range ch {
		result.Stend = append(result.Stend, c)
	}
	sort.Slice(result.Stend, func(p, q int) bool {
		return result.Stend[p].Name < result.Stend[q].Name
	})
	result.ListUrls = churl.CheckUrl(loadData.ListUrls, url)
	return result
}

// OpenTerminal Open terminal for administrating servers
func OpenTerminal(t Terminal, username string) {
	if RunOps != "server" {
		desktop := os.Getenv("DESKTOP_SESSION")
		if desktop == "gnome" {
			cmd := exec.Command("gnome-terminal", "--", "ssh", username+"@"+t.Ip, "-tt", "sudo -i")
			cmd.Run()
		} else if desktop == "fly" {
			// sudo not work... Need testing
			cmd := exec.Command("fly-term", "-e", "ssh", username+"@"+t.Ip)
			cmd.Run()
		}
	}
}

// RunGin Main package run this function
// Run web server.
func RunGin(port chport.Cheker, url churl.Churler, conf sshcmd.Execer, loadData global.YumInit) {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./build", true)))
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
	api.POST("/get", func(c *gin.Context) {
		result := serverHandler(loadData, port, url, conf)
		Status = result
		c.Header("Context-Type", "application/json")
		c.JSON(http.StatusOK, result)
	})
	api.POST("/terminal", func(c *gin.Context) {
		var terminal Terminal
		err := c.BindJSON(&terminal)
		if err != nil {
			c.JSON((http.StatusBadRequest), "")
			return
		}
		OpenTerminal(terminal, loadData.UserName)
		c.Header("Context-Type", "application/json")
		c.JSON(http.StatusOK, "")
	})
	api.POST("/errorlogs", func(c *gin.Context) {
		var wg_l sync.WaitGroup
		var ch = make(chan getlog.GetLog)
		res := []getlog.GetLog{}
		go func() {
			for _, host := range Status.Stend {
				if checkSshPort(host.Ip, conf.GetSshPort(), port) {
					for _, st := range host.Status {
						wg_l.Add(1)
						go func(host global.Host, st global.Result) {
							defer wg_l.Done()
							get := getlog.GetLog{
								Host:    host.Ip,
								Service: st.Service,
								Module:  st.Output,
							}
							out := get.GetErrors(loadData.Logs, loadData.CountLog, conf)
							ch <- out
						}(host, st)
					}
				}
			}
			wg_l.Wait()
			close(ch)
		}()
		for o := range ch {
			res = append(res, o)
		}
		// Out object create
		c.Header("Context-Type", "application/json")
		c.JSON(http.StatusOK, res)
	})
	api.POST("/warlog", func(c *gin.Context) {
		var getlog getlog.GetLog
		if err := c.BindJSON(&getlog); err != nil {
			fmt.Println(err)
		}
		logs := getlog.GetLogs(loadData.Logs, loadData.CountLog, conf)
		c.Header("Context-Type", "application/json")
		c.JSON(http.StatusOK, logs)
	})
	cmd := exec.Command("firefox", "http://localhost:3000/")
	go cmd.Run()
	router.Run(":3000")
}
