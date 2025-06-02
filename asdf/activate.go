package asdf

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/toolprovider"
)

func (a *AsdfToolProvider) ActivateEnv(result toolprovider.ToolInstallResult) (toolprovider.EnvironmentActivation, error) {
	envKey := fmt.Sprint("ASDF_", strings.ToUpper(result.ToolName), "_VERSION")
	return toolprovider.EnvironmentActivation{
		ContributedEnvVars: map[string]string{
			envKey: result.ConcreteVersion,
		},
		ContributedPaths: []string{},
	}, nil
}
