package asdf

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
)

func (a *AsdfToolProvider) asdfVersion() (*version.Version, error) {
	output, err := a.ExecEnv.runAsdf("--version")
	if err != nil {
		return nil, err
	}

	versionStr := strings.TrimSpace(string(output))
	ver, err := version.NewVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("parse asdf version: %w", err)
	}
	return ver, nil
}

// TODO: check if tool-plugin is installed
func (a *AsdfToolProvider) listInstalled(toolName string) ([]string, error) {
	output, err := a.ExecEnv.runAsdf("list", toolName)
	if err != nil {
		// asdf 0.16.0+ returns exit code 1 if no versions are installed
		if strings.Contains(err.Error(), "No compatible versions installed") {
			return []string{}, nil
		}
		return nil, err
	}

	installedVersions := parseAsdfListOutput(output)
	filteredVersions, err := filterAliasVersions(toolName, installedVersions)
	if err != nil {
		return nil, fmt.Errorf("filter alias versions: %w", err)
	}
	return filteredVersions, nil
}

// TODO: check if tool-plugin is installed
func (a *AsdfToolProvider) listReleased(toolName string) ([]string, error) {
	asdfVer, err := a.asdfVersion()
	if err != nil {
		return nil, err
	}
	var subcommand string
	if asdfVer.GreaterThanOrEqual(version.Must(version.NewVersion("0.16.0"))) {
		subcommand = "list all"
	} else {
		subcommand = "list-all"
	}

	output, err := a.ExecEnv.runAsdf(subcommand, toolName)
	if err != nil {
		return nil, err
	}

	releasedVersions := parseAsdfListOutput(output)
	return releasedVersions, nil
}

func parseAsdfListOutput(output string) []string {
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

func filterAliasVersions(tool string, versions []string) ([]string, error) {
	// Filter out versions that are symlinks created by the asdf-alias plugin.
	var filtered []string
	for _, v := range versions {
		out, err := exec.Command("asdf", "where", tool, v).Output()
		if err != nil {
			return nil, fmt.Errorf("asdf where %s %s: %w", tool, v, err)
		}

		fileInfo, err := os.Lstat(strings.TrimSpace(string(out)))
		if err != nil {
			return nil, fmt.Errorf("lstat %s: %w", strings.TrimSpace(string(out)), err)
		}

		if fileInfo.Mode()&os.ModeSymlink == 0 {
			filtered = append(filtered, v)
		}
	}
	return filtered, nil
}
