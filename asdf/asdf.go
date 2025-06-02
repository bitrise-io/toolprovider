package asdf

import (
	"fmt"

	"github.com/bitrise-io/toolprovider"
)

type ProviderOptions struct {
	AsdfVersion string
}

type AsdfToolProvider struct {
	ExecEnv ExecEnv
}

func (a *AsdfToolProvider) Bootstrap() error {
	// TODO:
	// Check if asdf is installed
	// Check if asdf version satisfies the supported version range

	return nil
}

func (a *AsdfToolProvider) InstallTool(tool toolprovider.ToolRequest) (toolprovider.ToolInstallResult, error) {
	installedVersions, err := a.listInstalled(tool.ToolName)
	if err != nil {
		return toolprovider.ToolInstallResult{}, fmt.Errorf("list installed versions: %w", err)
	}

	releasedVersions, err := a.listReleased(tool.ToolName)
	if err != nil {
		return toolprovider.ToolInstallResult{}, fmt.Errorf("list released versions: %w", err)
	}

	resolution, err := ResolveVersion(tool, releasedVersions, installedVersions)
	if err != nil {
		return toolprovider.ToolInstallResult{}, fmt.Errorf("resolve version: %w", err)
	}

	if resolution.IsInstalled {
		return toolprovider.ToolInstallResult{
			ToolName:           tool.ToolName,
			IsAlreadyInstalled: true,
			ConcreteVersion:    resolution.VersionString,
		}, nil
	} else {
		err = installToolVersion(tool.ToolName, resolution.VersionString)
		if err != nil {
			return toolprovider.ToolInstallResult{}, err
		}

		// TODO: reshim workarounds

		return toolprovider.ToolInstallResult{
			ToolName:           tool.ToolName,
			IsAlreadyInstalled: false,
			ConcreteVersion:    resolution.VersionString,
		}, nil
	}
}
