package asdf

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
)

func installedAsdfVersion() (*version.Version, error) {
	// We spawn a login shell because "classic" asdf is implemented in Bash and sourced in the shell config.
	cmd := exec.Command("bash", "-lc", "asdf", "version")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("exec asdf version: %w, %s", err, output)
	}
	
	versionStr := strings.TrimSpace(string(output))
	return version.NewVersion(versionStr)
}

// TODO: check if tool-plugin is installed
func listInstalled(toolName string) ([]string, error) {
	// We spawn a login shell because "classic" asdf is implemented in Bash and sourced in the shell config.
	cmd := exec.Command("bash", "-lc", fmt.Sprintf("asdf list %s", toolName))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("exec asdf list %s: %w, %s", toolName, err, output)
	}
	
	installedVersions := parseAsdfListOutput(output)
	return installedVersions, nil
}

// TODO: check if tool-plugin is installed
func listReleased(toolName string) ([]string, error) {
	// We spawn a login shell because "classic" asdf is implemented in Bash and sourced in the shell config.
	cmd := exec.Command("bash", "-lc", fmt.Sprintf("asdf list-all %s", toolName))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("exec asdf list-all %s: %w, %s", toolName, err, output)
	}

	releasedVersions := parseAsdfListOutput(output)
	return releasedVersions, nil
}

func parseAsdfListOutput(output []byte) []string {
	// There is no machine-readable output, we are parsing this:
	//   1.21.0
	//   1.21.11
	//   1.21
	//   1.22.0
	//  *1.22
	//   1.23.5
	//   1.23.7
	//   1.23
	//   1.24.0
	//   1

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var versions []string = []string{}
	for i := range lines {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		versions = append(versions, strings.TrimSpace(strings.Replace(lines[i], "*", "", 1)))
	}
	return versions
}
