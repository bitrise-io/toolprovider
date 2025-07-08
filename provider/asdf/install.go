package asdf

import (
	"fmt"

	"github.com/bitrise-io/toolprovider/provider/asdf/workarounds"
)


func (a *AsdfToolProvider) installToolVersion(
	toolName string,
	versionString string,
) error {
	// TODO: install plugin if not installed
	
	out, err := a.ExecEnv.RunAsdf("install", toolName, versionString)
	if err != nil {
		return fmt.Errorf("install %s %s: %w\n\nOutput:\n%s", toolName, versionString, err, out)
	}

	if toolName == "nodejs" {
		err = workarounds.SetupCorepack(a.ExecEnv, versionString)
		if err != nil {
			return fmt.Errorf("setup corepack for %s %s: %w", toolName, versionString, err)
		}
	}
	return nil
}
