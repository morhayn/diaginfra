package handl

import (
	"github.com/morhayn/diaginfra/internal/modules"
	"github.com/morhayn/diaginfra/internal/sshcmd"
)

func HandleResult(list []sshcmd.Out) []modules.Result {
	var result modules.Results
	for _, res := range list {
		if mod, ok := modules.MapCmd[res.Name]; ok {
			result.AddResults(res.Result, res.Name, mod.Handler)
		}
	}
	return result.Res
}
