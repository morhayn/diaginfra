package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/churl"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/sshcmd"
	"github.com/morhayn/diaginfra/internal/webapi"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		webapi.RunOps = args[1]
	} else {
		webapi.RunOps = "local"
	}
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	prDir := filepath.Dir(path)
	os.Chdir(prDir)
	var (
		conf sshcmd.SshConfig
		port chport.Port
		url  churl.Url
		sp   global.YumInit
	)
	loadData, err := sp.ReadConfig("conf/config.yaml")
	if err != nil {
		fmt.Printf("error reading config file: %s", err)
		os.Exit(1)
	}
	if loadData.UserName == "" {
		loadData.UserName = os.Getenv("USER")
	}
	if loadData.SshPort == "" {
		loadData.SshPort = "22"
	}
	conf.Init_ssh(loadData.UserName, loadData.SshPort)
	webapi.RunGin(port, url, conf, loadData)
}
