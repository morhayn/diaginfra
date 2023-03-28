package sshcmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/morhayn/diaginfra/internal/modules"
	"golang.org/x/crypto/ssh"
)

var (
	mapCmd = map[string]string{
		"DiskFree": "df / | awk 'FNR > 1 {print $5}'",
		"LoadAvg":  "awk '{print $1}' /proc/loadavg",
		"Systemd":  "sudo systemctl is-active %s",
	}
	reg = `%!.?\([EXTRA|MISSING]`
)

type Execer interface {
	Execute(string, CmdExec)
}

type SshConfig struct {
	*ssh.ClientConfig
	sshPort string
	recurs  int
}

type Comands struct {
	Comm []CmdExec
}
type CmdExec struct {
	Name    string
	PrgName string
	Chan    chan Out
	Cmd     string
}
type Out struct {
	Name    string `json:"name"`
	Result  string `json:"result"`
	PrgName string `json:"prgname"`
}

// Init_ssh Configure ssh connection to servers
func (s *SshConfig) Init_ssh(username, port string) {
	s.sshPort = port
	key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable read key")
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable parse private key")
	}
	s.ClientConfig = &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// NewOut  Create struscture with result shell command
func NewOut(name, prgName, res string) Out {
	return Out{Name: name, PrgName: prgName, Result: res}
}

// Create structure with failed shell command
func newOutFail(name, prgName string) Out {
	return Out{Name: name, PrgName: prgName, Result: "failed"}
}

// Run Runing gourutine with executing shell command check service
// Grouping result executing command
func Run(ip string, list []string, conf Execer) ([]Out, []Out, error) {
	var wg sync.WaitGroup
	list_srv := []Out{}
	list_prg := []Out{}
	done := make(chan string)
	cmd := Comands{}
	srv, prg := cmd.buildCmd(list)
	go func() {
		for _, c := range cmd.Comm {
			wg.Add(1)
			go func(c CmdExec) {
				defer wg.Done()
				conf.Execute(ip, c)
			}(c)
		}
		wg.Wait()
		done <- "true"
	}()
	for {
		select {
		case service := <-srv:
			list_srv = append(list_srv, service)
		case program := <-prg:
			list_prg = append(list_prg, program)
		case <-done:
			sort.Slice(list_srv, func(p, q int) bool {
				return list_srv[p].Name < list_srv[q].Name
			})
			sort.Slice(list_prg, func(p, q int) bool {
				return list_prg[p].Name < list_prg[q].Name
			})
			return list_srv, list_prg, nil
		}
	}
}

// Build structure for check service and programm
// chan srv for service (systemd check)
// chan prg for program running predetermined command
func (c *Comands) buildCmd(list []string) (_, _ chan Out) {
	srv := make(chan Out)
	prg := make(chan Out)
	for _, l := range list {
		if l == "Non" {
			return srv, prg
		}
	}
	df := CmdExec{
		Name: "DiskFree",
		Chan: srv,
		Cmd:  mapCmd["DiskFree"],
	}
	c.Comm = append(c.Comm, df)
	lavg := CmdExec{
		Name: "LoadAvg",
		Chan: srv,
		Cmd:  mapCmd["LoadAvg"],
	}
	c.Comm = append(c.Comm, lavg)
	for _, i := range list {
		ssh := CmdExec{
			Name: i,
		}
		ssh.swCmd(srv, prg)
		c.Comm = append(c.Comm, ssh)
	}
	return srv, prg
}

// Swith command for test service
func (s *CmdExec) swCmd(srv, prg chan Out) {
	var err error
	s.Chan = prg
	s.Cmd = ""
	splName := strings.Split(s.Name, ":")
	//Space in arg must inject command in test command
	if testSpaceInArg(splName) {
		return
	}
	s.Name = splName[0]
	if len(splName) == 2 {
		s.PrgName = splName[1]
	}
	test, ok := modules.MapService[splName[0]]
	if !ok {
		s.Chan = srv
		s.Cmd = fmt.Sprintf(mapCmd["Systemd"], s.Name)
	} else {
		sp := splName[1:]
		s.Cmd, err = test.RunString(sp...)
		if err != nil {
			s.Cmd = ""
		}
	}
	//If Sprintf print error in string MISSING or EXTRA arg
	r := regexp.MustCompile(reg)
	mutchErr := r.FindStringIndex(s.Cmd)
	if mutchErr != nil {
		s.Cmd = ""
	}
}

// Execute Run shell command in ssh
func (conf SshConfig) Execute(ip string, cmd CmdExec) {
	var d net.Dialer
	connStr := fmt.Sprintf("%s:%s", ip, conf.sshPort)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", connStr)
	defer conn.Close()
	if err != nil {
		fmt.Println("Error connect1", ip, cmd.Cmd)
		cmd.Chan <- newOutFail(cmd.Name, cmd.PrgName)
		return
	}
	c, chans, req, err := ssh.NewClientConn(conn, connStr, conf.ClientConfig)
	if err != nil {
		fmt.Println("Error connect2", ip, err)
		if conf.recurs < 3 {
			time.Sleep(1 * time.Second)
			conf.recurs += 1
			conf.Execute(ip, cmd)
		} else {
			cmd.Chan <- newOutFail(cmd.Name, cmd.PrgName)
		}
		return
	}
	client := ssh.NewClient(c, chans, req)
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		fmt.Println("Error connect3", ip, cmd.Cmd)
		cmd.Chan <- newOutFail(cmd.Name, cmd.PrgName)
		return
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	if err := session.Run(cmd.Cmd); err != nil {
		fmt.Println("!!!ERROR", err, cmd.Cmd, ip)
		cmd.Chan <- newOutFail(cmd.Name, cmd.PrgName)
		return
	}
	cmd.Chan <- NewOut(cmd.Name, cmd.PrgName, stdoutBuf.String())
}
func testSpaceInArg(str []string) bool {
	r := regexp.MustCompile(`\s+`)
	for _, s := range str {
		m := r.FindStringIndex(s)
		if m != nil {
			return true
		}
	}
	return false
}
