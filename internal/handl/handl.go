package handl

import (
	"sort"
	"sync"

	"github.com/morhayn/diaginfra/internal/chport"
	"github.com/morhayn/diaginfra/internal/churl"
	"github.com/morhayn/diaginfra/internal/global"
	"github.com/morhayn/diaginfra/internal/modules"
	"github.com/morhayn/diaginfra/internal/sshcmd"
)

// Create new structure for loading config file
func newHost(name, ip string) global.Host {
	return global.Host{
		Name:    name,
		Ip:      ip,
		ListSsh: []global.Out{},
		Status:  []global.Result{},
	}
}

// HandkeResult -
func HandleResult(list []global.Out) []global.Result {
	var result modules.Results
	for _, res := range list {
		if mod, ok := modules.MapService[res.Name]; ok {
			result.AddResults(res.Result, res.Name, res.PrgName, mod.Handler)
		}
	}
	return result.Res
}

// Run test command to one server
func checkHost(host global.Init, ch chan global.Host, port chport.Cheker, conf sshcmd.Execer) {
	h := newHost(host.Name, host.Ip)
	h.ListPort = chport.CheckPort(host.Ip, host.ListPorts, port)
	if chport.CheckSshPort(host.Ip, conf.GetSshPort(), port) {
		srv, prg, _ := sshcmd.Run(host.Ip, host.ListService, conf)
		h.ListSsh = srv
		h.Status = HandleResult(prg)
		sort.Slice(h.Status, func(p, q int) bool {
			return h.Status[p].Status < h.Status[q].Status
		})

	}
	ch <- h
}

// Run gorutine for all servers in config file and grouping result
func ServerHandler(loadData global.YumInit, port chport.Cheker, url churl.Churler, conf sshcmd.Execer) global.Hosts {
	var wg sync.WaitGroup
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
