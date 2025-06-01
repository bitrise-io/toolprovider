package asdf

import (
	"fmt"

	"github.com/bitrise-io/toolprovider"
)

type ProviderOptions struct {
	AsdfVersion string
}

type AsdfToolProvider struct {
}

func (a *AsdfToolProvider) Bootstrap() error {
	// TODO:
	// Check if asdf is installed
	// Check if asdf version satisfies the supported version range

	return nil
}

func (a *AsdfToolProvider) InstallTool(tool toolprovider.ToolRequest) (toolprovider.ToolInstallResult, error) {
	installedVersions, err := listInstalled(tool.ToolName)
	if err != nil {
		return toolprovider.ToolInstallResult{}, fmt.Errorf("list installed versions: %w", err)
	}

	releasedVersions, err := listReleased(tool.ToolName)
	if err != nil {
		return toolprovider.ToolInstallResult{}, fmt.Errorf("list released versions: %w", err)
	}

	resolution, err := ResolveVersion(tool, releasedVersions, installedVersions)
	if err != nil {
		return toolprovider.ToolInstallResult{}, fmt.Errorf("resolve version: %w", err)
	}

	if resolution.IsInstalled {
		return toolprovider.ToolInstallResult{
			IsAlreadyInstalled: true,
		}, nil
	} else {
		err = installToolVersion(tool.ToolName, resolution.VersionString)
		if err != nil {
			return toolprovider.ToolInstallResult{}, err
		}

		// TODO: reshim workarounds

		return toolprovider.ToolInstallResult{
			IsAlreadyInstalled: false,
		}, nil
	}
}
