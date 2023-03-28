package handl

import (
	"github.com/morhayn/diaginfra/internal/modules"
	"github.com/morhayn/diaginfra/internal/sshcmd"
)

// HandkeResult -
func HandleResult(list []sshcmd.Out) []modules.Result {
	var result modules.Results
	for _, res := range list {
		if mod, ok := modules.MapService[res.Name]; ok {
			result.AddResults(res.Result, res.Name, res.PrgName, mod.Handler)
		}
	}
	return result.Res
}
