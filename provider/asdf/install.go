package asdf

import (
	"fmt"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/asdf/workarounds"
)

func (a *AsdfToolProvider) installToolVersion(
	toolName string,
	versionString string,
) error {
	if toolName == "" || versionString == "" {
		return fmt.Errorf("toolName and versionString must not be empty")
	}

	out, err := a.ExecEnv.RunAsdf("install", toolName, versionString)
	if err != nil {
		return provider.ToolInstallError{
			ToolName:         toolName,
			RequestedVersion: versionString,
			Cause:            fmt.Sprintf("asdf install %s %s: %s", toolName, versionString, err),
			RawOutput:        out,
		}
	}

	if toolName == "nodejs" {
		err = workarounds.SetupCorepack(a.ExecEnv, versionString)
		if err != nil {
			return fmt.Errorf("setup corepack for %s %s: %w", toolName, versionString, err)
		}
	}
	return nil
}
