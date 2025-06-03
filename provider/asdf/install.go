package asdf

import "fmt"


func (a *AsdfToolProvider) installToolVersion(
	toolName string,
	versionString string,
) error {
	// TODO: install plugin if not installed
	
	out, err := a.ExecEnv.runAsdf("install", toolName, versionString)
	if err != nil {
		return fmt.Errorf("install %s %s: %w\n\nOutput:\n%s", toolName, versionString, err, out)
	}


	// TODO: reshim workarounds after install
	return nil
}
