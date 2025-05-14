package asdf

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bitrise-io/bitrise/v2/log"
	"github.com/bitrise-io/toolprovider"
	"github.com/hashicorp/go-version"
	"golang.org/x/exp/slices"
)

var ErrStrictVersionNotInstalled = errors.New("requested version is not installed and resolution strategy is strict")

// TODO: better error when update is implemented
var ErrNoMatchingVersion = errors.New("no matching version found")

var ErrRequestedVersionNotSemVer = errors.New("requested version is not semver compatible while available tool versions follow semver")

type ProviderOptions struct {
	AsdfVersion string
}

type VersionResolution struct {
	VersionString string
	IsSemVer      bool
	SemVer        *version.Version
	IsInstalled   bool
}

type AsdfToolProvider struct {
}

func (a *AsdfToolProvider) Bootstrap() error {
	// TODO:
	// Check if asdf is installed
	// Check if asdf version satisfies the supported version range

	return nil
}

func (a *AsdfToolProvider) InstallTool(tool toolprovider.ToolRequest) error {
	// TODO
	return nil
}

func ResolveVersion(
	request toolprovider.ToolRequest,
	releasedVersions []string,
	installedVersions []string,
) (VersionResolution, error) {
	// Short-circuit for exact version match among installed versions
	if slices.Contains(installedVersions, strings.TrimSpace(request.UnparsedVersion)) {
		if request.ResolutionStrategy != toolprovider.ResolutionStrategyStrict {
			log.Warn("Request matches an installed version, but resolution strategy is not set to strict. You might want to use a partial version string as the requested version.")
		}
		requestedSemVer, err := version.NewVersion(request.UnparsedVersion)
		return VersionResolution{
			VersionString: request.UnparsedVersion,
			IsSemVer:      err == nil,
			SemVer:        requestedSemVer,
			IsInstalled:   true,
		}, nil
	}

	// TODO
	isToolSemVer, _ := isToolSemVer(releasedVersions)

	if request.ResolutionStrategy == toolprovider.ResolutionStrategyStrict {
		if slices.Contains(releasedVersions, request.UnparsedVersion) {
			requestedSemVer, _ := version.NewVersion(request.UnparsedVersion)
			return VersionResolution{
				VersionString: request.UnparsedVersion,
				IsSemVer:      isToolSemVer,
				SemVer:        requestedSemVer,
				IsInstalled:   slices.Contains(installedVersions, request.UnparsedVersion),
			}, nil
		}
		return VersionResolution{}, ErrNoMatchingVersion
	}

	if request.ResolutionStrategy == toolprovider.ResolutionStrategyLatestInstalled {
		if isToolSemVer {
			// Installed versions are checked first because strategy is "latest installed"
			var sortedInstalledVersions version.Collection
			for _, v := range installedVersions {
				installedV, err := version.NewVersion(v)
				if err != nil {
					return VersionResolution{}, fmt.Errorf("parse %s %s: %w", request.ToolName, v, err)
				}
				sortedInstalledVersions = append(sortedInstalledVersions, installedV)
			}
			sort.Sort(sort.Reverse(sortedInstalledVersions))

			for _, v := range sortedInstalledVersions {
				if strings.HasPrefix(v.String(), request.UnparsedVersion) {
					// Since versions are semver-compatible and `version.Collection`
					// guarantees correct ordering (even for pre-releases),
					// we can stop searching if the version prefix-matches the requested version.
					return VersionResolution{
						VersionString: v.String(),
						IsSemVer:      true,
						SemVer:        v,
						IsInstalled:   true,
					}, nil
				}
			}

			// If there is no match among installed versions, we check the released versions (despite the strategy being "latest installed").
			var sortedReleasedVersions version.Collection
			for _, v := range releasedVersions {
				releasedV, err := version.NewVersion(v)
				if err != nil {
					return VersionResolution{}, fmt.Errorf("parse %s %s: %w", request.ToolName, v, err)
				}
				sortedReleasedVersions = append(sortedReleasedVersions, releasedV)
			}
			sort.Sort(sort.Reverse(sortedReleasedVersions))

			for _, v := range sortedReleasedVersions {
				if strings.HasPrefix(v.String(), request.UnparsedVersion) {
					// Since versions are semver-compatible and `version.Collection`
					// guarantees correct ordering (even for pre-releases),
					// we can stop searching if the version prefix-matches the requested version.
					return VersionResolution{
						VersionString: v.String(),
						IsSemVer:      true,
						SemVer:        v,
						IsInstalled:   false,
					}, nil
				}
			}

			return VersionResolution{}, ErrNoMatchingVersion

		} else {
			sortedInstalledVersions := slices.Clone(installedVersions)
			slices.Sort(sortedInstalledVersions)
			slices.Reverse(sortedInstalledVersions)
			for _, v := range sortedInstalledVersions {
				if strings.HasPrefix(v, request.UnparsedVersion) {
					return VersionResolution{
						VersionString: v,
						IsSemVer:      false,
						SemVer:        nil,
						IsInstalled:   true,
					}, nil
				}
			}

			sortedReleasedVersions := slices.Clone(releasedVersions)
			slices.Sort(sortedReleasedVersions)
			slices.Reverse(sortedReleasedVersions)
			for _, v := range sortedReleasedVersions {
				if strings.HasPrefix(v, request.UnparsedVersion) {
					return VersionResolution{
						VersionString: v,
						IsSemVer:      false,
						SemVer:        nil,
						IsInstalled:   false,
					}, nil
				}
			}

			return VersionResolution{}, ErrNoMatchingVersion
		}
	} else if request.ResolutionStrategy == toolprovider.ResolutionStrategyLatestReleased {
		if isToolSemVer {
			var sortedReleasedVersions version.Collection
			for _, v := range releasedVersions {
				releasedV, err := version.NewVersion(v)
				if err != nil {
					return VersionResolution{}, fmt.Errorf("parse %s %s: %w", request.ToolName, v, err)
				}
				sortedReleasedVersions = append(sortedReleasedVersions, releasedV)
			}
			sort.Sort(sort.Reverse(sortedReleasedVersions))
			for _, v := range sortedReleasedVersions {
				if strings.HasPrefix(v.String(), request.UnparsedVersion) {
					// Since versions are semver-compatible and `version.Collection`
					// guarantees correct ordering (even for pre-releases),
					// we can stop searching if the version prefix-matches the requested version.

					// Even though we search the released versions primarily,
					// it's still possible that the latest released version is also installed.
					isInstalled := slices.Contains(installedVersions, v.String())

					return VersionResolution{
						VersionString: v.String(),
						IsSemVer:      true,
						SemVer:        v,
						IsInstalled:   isInstalled,
					}, nil
				}
			}
			return VersionResolution{}, ErrNoMatchingVersion
		} else {
			sortedReleasedVersions := slices.Clone(releasedVersions)
			slices.Sort(sortedReleasedVersions)
			slices.Reverse(sortedReleasedVersions)
			for _, v := range sortedReleasedVersions {
				if strings.HasPrefix(v, request.UnparsedVersion) {
					// Even though we search the released versions primarily,
					// it's still possible that the latest released version is also installed.
					isInstalled := slices.Contains(installedVersions, v)

					return VersionResolution{
						VersionString: v,
						IsSemVer:      false,
						SemVer:        nil,
						IsInstalled:   isInstalled,
					}, nil
				}
			}
			return VersionResolution{}, ErrNoMatchingVersion
		}
	}

	// TODO
	return VersionResolution{}, fmt.Errorf("TODO")
}

// isToolSemVer guesses if a tool follows semantic versioning.
// It returns true if all released versions are semver compatible,
// or the semver-violating version string if it is not semver compatible.
// TODO: can we guarantee that releasedVersions is always up to date? Or that it's good enough?
func isToolSemVer(releasedVersions []string) (bool, string) {
	for _, v := range releasedVersions {
		// Note: go-version pads missing segments with 0s, so "1.0" becomes "1.0.0".
		// We allow this because the real world is messy.
		// For example, Golang major releases prior to 1.21.0 had only major.minor segments (1.20, 1.19, etc.).
		_, err := version.NewVersion(v)
		if err != nil {
			return false, v
		}
	}
	return true, ""
}
