package mise

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/bitrise-io/toolprovider/provider"
)

var errNoMatchingVersion = errors.New("no matching version found")

func (m *MiseToolProvider) resolveToConcreteVersionAfterInstall(tool provider.ToolRequest) (string, error) {
	// Mise doesn't tell us what version it resolved to when installing the user-provided (and potentially fuzzy) version.
	// But we can use `mise latest` to find out the concrete version.
	switch tool.ResolutionStrategy {
	case provider.ResolutionStrategyLatestInstalled:
		return m.resolveToLatestInstalled(tool.ToolName, tool.UnparsedVersion)
	case provider.ResolutionStrategyLatestReleased, provider.ResolutionStrategyStrict:
		// Mise works with fuzzy versions by default, so it happily installs both node@20 and node@20.19.3.
		// Therefore, when the Bitrise config contains simply 20 (and not 20:latest), it actually behaves
		// as "latest released".
		return m.resolveToLatestReleased(tool.ToolName, tool.UnparsedVersion)
	default:
		return "", fmt.Errorf("unknown resolution strategy: %v", tool.ResolutionStrategy)
	}
}

func (m *MiseToolProvider) resolveToLatestReleased(toolName string, version string) (string, error) {
	// Even if version is empty string "sometool@" will not cause an error.
	cmd := exec.Command("mise", "latest", fmt.Sprintf("%s@%s", toolName, version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("mise latest %s@%s: %w", toolName, version, err)
	}

	v := strings.TrimSpace(string(output))
	if v == "" {
		return "", errNoMatchingVersion
	}

	return strings.TrimSpace(string(output)), nil
}

func (m *MiseToolProvider) resolveToLatestInstalled(toolName string, version string) (string, error) {
	// Even if version is empty string "sometool@" will not cause an error.
	cmd := exec.Command("mise", "latest", "--installed", fmt.Sprintf("%s@%s", toolName, version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("mise latest --installed %s@%s: %w", toolName, version, err)
	}

	v := strings.TrimSpace(string(output))
	if v == "" {
		return "", errNoMatchingVersion
	}

	return v, nil
}
