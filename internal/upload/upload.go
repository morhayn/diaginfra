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
)

var (
	statusFile       = "conf/stend.status"
	nfsDisk          = ""
	uploadDir        = "upload"
	modUploadDir     = "sudo chmod 0777 %s"
	removeWebApps    = "sudo rm -rf /var/lib/tomcat8/webapps/*"
	removeLogs       = "sudo rm -rf /var/log/tomcat8/*"
	psqlManageStop   = "cd /d01/ && sudo /d01/script/manager_stop.sh"
	psqlManagerStart = "cd /d01/ && sudo /d01/script/manager_start.sh"
	stopTomcat       = "sudo systemctl stop tomcat8"
	startTomcat      = "sudo systemctl start tomcat8"
	stopCrone        = "sudo systemctl stop crone"
	startCrone       = "sudo systemctl start crone"
	stopNginx        = "sudo systemctl stop nginx"
	startNginx       = "sudo systemctl start nginx"
	copyWar          = "sudo cp %s/%s.war /var/lib/tomcat8/webapps/ && sudo chown tomcat8:tomcat8 /var/lib/tomcat8/webapps/%s.war"
)

func CopyWars(status global.Hosts, path string, conf sshcmd.Execer) error {
	err := saveStatus(status)
	if err != nil {
		fmt.Println("Error Save Stend status", err)
		return err
	}
	warRegEx, err := regexp.Compile("^.+\\.(war)$")
	if err != nil {
		return err
	}
	err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && warRegEx.MatchString(info.Name()) {
			// Copy file
		}
		return nil
	})
	return nil
}

func PreUploadDb() {
}

func ProdUploadDb() {
}
func UploadWars(conf sshcmd.Execer, port chport.Cheker) error {
	status, err := loadStatus()
	if err != nil {
		return err
	}
	for _, host := range status.Stend {
		if chport.CheckSshPort(host.Ip, conf.GetSshPort(), port) {
			wars := []string{}
			_ = exec(host.Ip, stopTomcat, conf)
			for _, st := range host.Status {
				if st.Service == "Tomcat" {
					if st.Output != "manager" && st.Output != "host-manager" {
						wars = append(wars, st.Output)
					}
				}
			}
			_ = exec(host.Ip, removeWebApps, conf)
			for _, war := range wars {
				_ = exec(host.Ip, fmt.Sprintf(copyWar, nfsDisk, war, war), conf)
			}
			_ = exec(host.Ip, startTomcat, conf)
		}
	}
	return nil
}

func clearWebapps(ip string, conf sshcmd.Execer) string {
	return exec(ip, removeWebApps, conf)
}
func clearLogs(ip string, conf sshcmd.Execer) string {
	return exec(ip, removeLogs, conf)
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
func copy(ip, scr, dest string, conf sshcmd.Execer) string {

}
