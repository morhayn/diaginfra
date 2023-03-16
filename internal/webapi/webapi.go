package webapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"sync"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/churl"
	"github.com/morhayn/diaginfra/internal/getlog"
	"github.com/morhayn/diaginfra/internal/handl"
	"github.com/morhayn/diaginfra/internal/sshcmd"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var (
	RunOps string
	wg     sync.WaitGroup
	// Username string
	Status Hosts
)

type Terminal struct {
	Ip string `json:"ip"`
}

type YumInit struct {
	// Stend     string            `yaml:"stend"`
	UserName string            `yaml:"user"`
	SshPort  string            `yaml:"ssh_port"`
	CountLog int               `yaml:"countlog"`
	ListUrls []string          `yaml:"list_urls"`
	Logs     map[string]string `yaml:"logs"`
	Hosts    []Init            `yaml:"hosts"`
}
type Init struct {
	Name        string   `yaml:"name"`
	Ip          string   `yaml:"ip"`
	ListPorts   []string `yaml:"list_ports"`
	ListService []string `yaml:"list_service"`
	Wars        []string `yaml:"wars"`
	// Jars         []string `yaml:"jars"`
}
type Hosts struct {
	ListUrls []churl.Url `json:"list_url"`
	Stend    []Host      `josn:"stand"`
}
type Host struct {
	Name     string         `json:"name"`
	Ip       string         `json:"ip"`
	ListPort []chport.Port  `json:"list_port"`
	ListSsh  []sshcmd.Out   `json:"list_ssh"`
	Status   []handl.Result `json:"status"`
}

// Create new structure for loading config file
func newHost(name, ip string) Host {
	return Host{
		Name:    name,
		Ip:      ip,
		ListSsh: []sshcmd.Out{},
		Status:  []handl.Result{},
	}
}

// Read Config file and unmarshall data in structure
func (y YumInit) ReadConfig(file string) YumInit {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error open file")
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, &y)
	if err != nil {
		fmt.Println("Error unmarshal")
		log.Fatal(err)
	}
	return y
}

// Check ssh port if ssh port failed not nid run ssh command to server
func checkSshPort(ports []chport.Port) bool {
	for _, p := range ports {
		if p.Port == "22" && p.Status == "failed" {
			return false
		}
	}
	return true
}

// Run test command to one server
func checkHost(host Init, ch chan Host, port chport.Cheker, conf sshcmd.Execer) {
	h := newHost(host.Name, host.Ip)
	h.ListPort = chport.CheckPort(host.Ip, host.ListPorts, port)
	if checkSshPort(h.ListPort) {
		srv, prg, _ := sshcmd.Run(host.Ip, host.ListService, conf)
		h.ListSsh = srv
		h.Status = handl.HandleResult(prg)
	}
	ch <- h
}

// Run gorutine for all servers in config file and grouping result
func serverHandler(loadData YumInit, port chport.Cheker, url churl.Churler, conf sshcmd.Execer) Hosts {
	result := Hosts{}
	ch := make(chan Host)
	go func() {
		for _, host := range loadData.Hosts {
			wg.Add(1)
			go func(host Init) {
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

// Open terminal for administrating servers
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

// Main package run this function
// Run web server.
func RunGin(port chport.Cheker, url churl.Churler, conf sshcmd.Execer, loadData YumInit) {
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
				for _, st := range host.Status {
					wg_l.Add(1)
					go func(ip string, st handl.Result) {
						defer wg_l.Done()
						get := getlog.GetLog{
							Host:    ip,
							Service: st.Service,
							Module:  st.Status,
						}
						out := get.GetErrors(loadData.Logs, loadData.CountLog, conf)
						ch <- out
					}(host.Ip, st)
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
