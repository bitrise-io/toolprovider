package asdf

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bitrise-io/bitrise/v2/log"
	"github.com/bitrise-io/toolprovider/provider"
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

func (a *AsdfToolProvider) InstallTool(tool provider.ToolRequest) (provider.ToolInstallResult, error) {
	installedVersions, err := a.listInstalled(tool.ToolName)
	if err != nil {
		return provider.ToolInstallResult{}, fmt.Errorf("list installed versions: %w", err)
	}

	// Short-circuit for exact version match among installed versions.
	// Fetching released versions is a slow operation that we want to avoid.
	v := strings.TrimSpace(tool.UnparsedVersion)
	if tool.ResolutionStrategy == provider.ResolutionStrategyStrict && slices.Contains(installedVersions, v) {
		return provider.ToolInstallResult{
			ToolName:           tool.ToolName,
			IsAlreadyInstalled: true,
			ConcreteVersion:    v,
		}, nil
	}

	releasedVersions, err := a.listReleased(tool.ToolName)
	if err != nil {
		return provider.ToolInstallResult{}, fmt.Errorf("list released versions: %w", err)
	}

	resolution, err := ResolveVersion(tool, releasedVersions, installedVersions)
	if err != nil {
		if errors.Is(err, &ErrNoMatchingVersion{}) {
			log.Warn("No matching version found, updating asdf-%s plugin and retrying...", tool.ToolName)
			// Some asdf plugins hardcode the list of installable versions and need a new plugin release to support new versions.
			_, err = a.ExecEnv.runAsdf("plugin", "update", tool.ToolName)
			if err != nil {
				return provider.ToolInstallResult{}, fmt.Errorf("update plugin: %w", err)
			}
			releasedVersions, err = a.listReleased(tool.ToolName)
			if err != nil {
				return provider.ToolInstallResult{}, fmt.Errorf("list released versions after plugin update: %w", err)
			}
			resolution, err = ResolveVersion(tool, releasedVersions, installedVersions)
			if err != nil {
				return provider.ToolInstallResult{}, fmt.Errorf("resolve version: %w", err)
			}
		}

		return provider.ToolInstallResult{}, fmt.Errorf("resolve version: %w", err)
	}

	if resolution.IsInstalled {
		return provider.ToolInstallResult{
			ToolName:           tool.ToolName,
			IsAlreadyInstalled: true,
			ConcreteVersion:    resolution.VersionString,
		}, nil
	} else {
		err = a.installToolVersion(tool.ToolName, resolution.VersionString)
		if err != nil {
			return provider.ToolInstallResult{}, err
		}

		return provider.ToolInstallResult{
			ToolName:           tool.ToolName,
			IsAlreadyInstalled: false,
			ConcreteVersion:    resolution.VersionString,
		}, nil
	}
}
