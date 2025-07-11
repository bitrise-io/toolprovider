package asdf

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/toolprovider/provider"
)

func (a AsdfToolProvider) ActivateEnv(result provider.ToolInstallResult) (provider.EnvironmentActivation, error) {
	envKey := fmt.Sprint("ASDF_", strings.ToUpper(result.ToolName), "_VERSION")
	return provider.EnvironmentActivation{
		ContributedEnvVars: map[string]string{
			envKey: result.ConcreteVersion,
		},
		ContributedPaths: []string{}, // TODO: shims dir?
	}, nil
}
