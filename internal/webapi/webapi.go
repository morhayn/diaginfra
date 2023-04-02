package webapi

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

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
	// Username string
	Status global.Hosts
)

type Terminal struct {
	Ip string `json:"ip"`
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
		result := handl.ServerHandler(loadData, port, url, conf)
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
		// Out object create
		res := getlog.GetErr(Status, loadData, port, conf)
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
