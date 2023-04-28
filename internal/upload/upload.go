package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/sshcmd"
	"gopkg.in/yaml.v2"
)

var (
	conf             = "conf/upload.yml"
	statusFile       = "conf/stend.status"
	uploadDir        = "upload"
	modRelDir        = "sudo chmod 0777 %s"
	removeWebApps    = "sudo rm -rf /var/lib/tomcat8/webapps/*"
	removeLogs       = "sudo rm -rf /var/log/tomcat8/*"
	psqlManageStop   = "cd /d01/ && sudo /d01/script/manager_stop.sh"
	psqlManagerStart = "cd /d01/ && sudo /d01/script/manager_start.sh"
	stopTomcat       = "sudo systemctl stop tomcat8"
	startTomcat      = "sudo systemctl start tomcat8"
	stopCron         = "sudo systemctl stop cron"
	startCron        = "sudo systemctl start cron"
	stopNginx        = "sudo systemctl stop nginx"
	startNginx       = "sudo systemctl start nginx"
	copyWar          = "sudo cp %s%s.war /var/lib/tomcat8/webapps/ && sudo chown tomcat8:tomcat8 /var/lib/tomcat8/webapps/%s.war"
)

type Upload struct {
	nfsServer   string `yaml:"nfs_server"`
	preRelease  string `yaml:"pre_release"`
	prodRelease string `yaml:"prod_release"`
	nfsPre      string `yaml:"nfs_pre"`
	nfsProd     string `yaml:"nfs_prod"`
	preUpDb     []DbUp `yaml:"pre_up_db"`
	prodUpDb    []DbUp `yaml:"prod_up_db"`
	offOn       OffOn  `yaml:"off_on"`
}
type DbUp struct {
	Server string   `yaml:"server"`
	Wars   []string `yaml:"wars"`
}
type OffOn struct {
	Nginx    []string `yaml:"nginx"`
	Cron     []string `yaml:"cron"`
	DbManage []string `yaml:"db_manage"`
}

func CopyWars(status global.Hosts, stend string, conf sshcmd.Execer) error {
	destDir := ""
	upl, err := readConfig()
	if err != nil {
		return err
	}
	if stend == "prod" {
		destDir = upl.prodRelease
		for _, ip := range upl.offOn.Nginx {
			_ = exec(ip, stopNginx, conf)
		}
		for _, ip := range upl.offOn.Cron {
			_ = exec(ip, stopCron, conf)
		}
		for _, ip := range upl.offOn.DbManage {
			_ = exec(ip, psqlManageStop, conf)
		}
	}
	if stend == "pre" {
		destDir = upl.preRelease
	}
	err = saveStatus(status)
	if err != nil {
		fmt.Println("Error Save Stend status", err)
		return err
	}
	warRegEx, err := regexp.Compile("^.+\\.(war)$")
	if err != nil {
		return err
	}
	_ = exec(upl.nfsServer, fmt.Sprintf(modRelDir, destDir), conf)
	err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && warRegEx.MatchString(info.Name()) {
			err = conf.Scp(upl.nfsServer, path, destDir)
		}
		return nil
	})
	return nil
}

func UploadDb(stend string, conf sshcmd.Execer, port chport.Cheker) error {
	upDb := []DbUp{}
	nfsDir := ""
	upl, err := readConfig()
	if err != nil {
		return err
	}
	if stend == "prod" {
		upDb = upl.prodUpDb
		nfsDir = upl.nfsProd
	}
	if stend == "pre" {
		upDb = upl.preUpDb
		nfsDir = upl.nfsPre
	}
	status, err := loadStatus()
	if err != nil {
		return err
	}
	for _, host := range status.Stend {
		if chport.CheckSshPort(host.Ip, conf.GetSshPort(), port) {
			for _, st := range host.Status {
				if st.Service == "Tomcat" {
					_ = exec(host.Ip, stopTomcat, conf)
					break
				}
			}
		}
	}
	for _, srv := range upDb {
		_ = exec(srv.Server, stopTomcat, conf)
		_ = exec(srv.Server, removeWebApps, conf)
		_ = exec(srv.Server, removeLogs, conf)
		for _, war := range srv.Wars {
			_ = exec(srv.Server, fmt.Sprintf(copyWar, nfsDir, war, war), conf)
		}
		_ = exec(srv.Server, startTomcat, conf)
	}
	return nil
}
func UploadWars(stend string, conf sshcmd.Execer, port chport.Cheker) error {
	nfsDir := ""
	status, err := loadStatus()
	if err != nil {
		return err
	}
	upl, err := readConfig()
	if err != nil {
		return err
	}
	if stend == "prod" {
		nfsDir = upl.nfsProd
	}
	if stend == "pre" {
		nfsDir = upl.nfsPre
	}
	for _, host := range status.Stend {
		if chport.CheckSshPort(host.Ip, conf.GetSshPort(), port) {
			wars := []string{}
			for _, st := range host.Status {
				if st.Service == "Tomcat" {
					if st.Output != "manager" && st.Output != "host-manager" {
						wars = append(wars, st.Output)
					}
				}
			}
			if len(wars) > 0 {
				_ = exec(host.Ip, stopTomcat, conf)
				_ = exec(host.Ip, removeWebApps, conf)
				for _, war := range wars {
					_ = exec(host.Ip, fmt.Sprintf(copyWar, nfsDir, war, war), conf)
				}
				_ = exec(host.Ip, startTomcat, conf)
			}
		}
	}
	if stend == "prod" {
		for _, ip := range upl.offOn.Nginx {
			_ = exec(ip, startNginx, conf)
		}
		for _, ip := range upl.offOn.Cron {
			_ = exec(ip, startCron, conf)
		}
		for _, ip := range upl.offOn.DbManage {
			_ = exec(ip, psqlManagerStart, conf)
		}
	}
	return nil
}

func saveStatus(status global.Hosts) error {
	if checkFileExist() {
		return errors.New("Status File Exists")
	}
	b, err := json.Marshal(status)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(statusFile, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
func loadStatus() (global.Hosts, error) {
	gl := global.Hosts{}
	if !checkFileExist() {
		return gl, errors.New("Status File Not Exest")
	}
	jsonFile, err := os.Open(statusFile)
	if err != nil {
		return gl, err
	}
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return gl, err
	}
	err = json.Unmarshal(b, &gl)
	if err != nil {
		return gl, err
	}
	return gl, nil
}
func isOldOneDay(t time.Time) bool {
	return time.Now().Sub(t) > 24*time.Hour
}
func checkFileExist() bool {
	st, err := os.Stat(statusFile)
	if err != nil {
		return true
	}
	if st.Mode().IsRegular() && isOldOneDay(st.ModTime()) {
		err := os.Remove(statusFile)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
func exec(ip, cmd string, conf sshcmd.Execer) string {
	c := sshcmd.CmdExec{
		Name: "Upload",
		Chan: make(chan global.Out),
	}
	c.Cmd = cmd
	go conf.Execute(ip, c)
	out := <-c.Chan
	return out.Result
}

func readConfig() (Upload, error) {
	u := Upload{}
	f, err := ioutil.ReadFile(conf)
	if err != nil {
		fmt.Println("Error open file")
		return u, err
	}
	err = yaml.Unmarshal(f, &u)
	if err != nil {
		fmt.Println("Error unmarshal")
		return u, err
	}
	return u, nil
}
