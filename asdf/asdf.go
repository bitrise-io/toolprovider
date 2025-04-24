package asdf

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bitrise-io/toolprovider"
	"github.com/hashicorp/go-version"
	"golang.org/x/exp/slices"
)

var ErrStrictVersionNotInstalled = errors.New("requested version is not installed and resolution strategy is strict")
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
	// TODO: figure out how to persist selection (local tool-version file might make git dirty)

	return nil
}

// TODO: version lists must be strings (possibly not semver)
func ResolveVersion(
	request toolprovider.ToolRequest,
	releasedVersions []string,
	installedVersions []string,
) (VersionResolution, error) {

	if slices.Contains(installedVersions, strings.TrimSpace(request.UnparsedVersion)) {
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
			requestedSemVer, err := version.NewVersion(request.UnparsedVersion)
			if err != nil {
				return VersionResolution{}, fmt.Errorf("parse %s %s: %w", request.ToolName, request.UnparsedVersion, ErrRequestedVersionNotSemVer)
			}

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
				if v.Equal(requestedSemVer) {
					return VersionResolution{
						VersionString: v.String(),
						IsSemVer:      true,
						SemVer:        v,
						IsInstalled:   true,
					}, nil
				}

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

			// If no installed version matches, we check the released versions (despite the strategy being "latest installed").
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

			return VersionResolution{}, fmt.Errorf("TODO")

		} else {
			// TODO: heuristic for natural order of non-semver strings
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
		return VersionResolution{}, fmt.Errorf("TODO")
	}

	return VersionResolution{}, fmt.Errorf("TODO")
}

// isToolSemVer checks if the requested tool is semver compatible and returns a boolean
// and the semver-violating version string if it is not semver compatible.
// Note: it's possible that the requested version looks semver compatible but one released version is not.
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
