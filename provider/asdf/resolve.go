package asdf

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/hashicorp/go-version"
)

type ErrNoMatchingVersion struct {
	RequestedVersion  string
	AvailableVersions []string
}

func (e ErrNoMatchingVersion) Error() string {
	if len(e.AvailableVersions) == 0 {
		return "no match for requested version " + e.RequestedVersion
	}

	versionList := ""
	for _, v := range e.AvailableVersions {
		if strings.HasPrefix(v, e.RequestedVersion) {
			versionList += fmt.Sprintf("- %s\n", v)
		}
	}
	if versionList == "" {
		return fmt.Sprintf("no match for requested version %s", e.RequestedVersion)
	} else {
		return fmt.Sprintf("no match for requested version %s. Similar versions: \n%s", e.RequestedVersion, versionList)
	}
}

type VersionResolution struct {
	VersionString string
	IsSemVer      bool
	SemVer        *version.Version
	IsInstalled   bool
}

func ResolveVersion(
	request provider.ToolRequest,
	releasedVersions []string,
	installedVersions []string,
) (VersionResolution, error) {
	// Short-circuit for exact version match among installed versions
	if slices.Contains(installedVersions, strings.TrimSpace(request.UnparsedVersion)) {
		requestedSemVer, err := version.NewVersion(request.UnparsedVersion)
		return VersionResolution{
			VersionString: request.UnparsedVersion,
			IsSemVer:      err == nil,
			SemVer:        requestedSemVer,
			IsInstalled:   true,
		}, nil
	}

	if request.ResolutionStrategy == provider.ResolutionStrategyStrict {
		if request.UnparsedVersion == "" || request.UnparsedVersion == "latest" {
			// Find absolute latest released version
			sortedReleasedVersions := logicallySortedVersions(releasedVersions)
			latestReleased := sortedReleasedVersions[0]
			if latestReleased == "" {
				return VersionResolution{}, &ErrNoMatchingVersion{
					AvailableVersions: releasedVersions,
					RequestedVersion:  "latest",
				}
			}
			isInstalled := slices.Contains(installedVersions, latestReleased)
			semverV, err := version.NewVersion(latestReleased)
			return VersionResolution{
				VersionString: latestReleased,
				IsSemVer:      err == nil,
				SemVer:        semverV,
				IsInstalled:   isInstalled,
			}, nil
		}

		if request.UnparsedVersion == "installed" {
			// Find absolute latest installed version
			sortedInstalledVersions := logicallySortedVersions(installedVersions)
			latestInstalled := sortedInstalledVersions[0]
			if latestInstalled == "" {
				return VersionResolution{}, &ErrNoMatchingVersion{
					AvailableVersions: installedVersions,
					RequestedVersion:  "installed",
				}
			}
			semverV, err := version.NewVersion(latestInstalled)
			return VersionResolution{
				VersionString: latestInstalled,
				IsSemVer:      err == nil,
				SemVer:        semverV,
				IsInstalled:   true,
			}, nil
		}

		if slices.Contains(releasedVersions, request.UnparsedVersion) {
			requestedSemVer, err := version.NewVersion(request.UnparsedVersion)
			return VersionResolution{
				VersionString: request.UnparsedVersion,
				IsSemVer:      err == nil,
				SemVer:        requestedSemVer,
				IsInstalled:   slices.Contains(installedVersions, request.UnparsedVersion),
			}, nil
		}
		return VersionResolution{}, &ErrNoMatchingVersion{AvailableVersions: releasedVersions, RequestedVersion: request.UnparsedVersion}
	}

	switch request.ResolutionStrategy {
	case provider.ResolutionStrategyLatestInstalled:
		// Installed versions are checked first because strategy is "latest installed"
		sortedInstalledVersions := logicallySortedVersions(installedVersions)

		if request.UnparsedVersion == "" {
			latestInstalled := sortedInstalledVersions[0]
			if latestInstalled == "" {
				return VersionResolution{}, &ErrNoMatchingVersion{
					AvailableVersions: installedVersions,
					RequestedVersion:  "installed",
				}
			}
			semverV, err := version.NewVersion(latestInstalled)
			return VersionResolution{
				VersionString: latestInstalled,
				IsSemVer:      err == nil,
				SemVer:        semverV,
				IsInstalled:   true,
			}, nil
		}

		for _, v := range sortedInstalledVersions {
			if strings.HasPrefix(v, request.UnparsedVersion) {
				// Since semver-compatible versions are sorted according to the semver spec
				// and are at the front of the list,
				// we can stop searching if the version prefix-matches the requested version.
				semverV, err := version.NewVersion(v)
				return VersionResolution{
					VersionString: v,
					IsSemVer:      err == nil,
					SemVer:        semverV,
					IsInstalled:   true,
				}, nil
			}
		}

		// If there is no match among installed versions, we check the released versions (despite the strategy being "latest installed").
		sortedReleasedVersions := logicallySortedVersions(releasedVersions)

		for _, v := range sortedReleasedVersions {
			if strings.HasPrefix(v, request.UnparsedVersion) {
				// Since semver-compatible versions are sorted according to the semver spec
				// and are at the front of the list,
				// we can stop searching if the version prefix-matches the requested version.
				semverV, err := version.NewVersion(v)
				return VersionResolution{
					VersionString: v,
					IsSemVer:      err == nil,
					SemVer:        semverV,
					IsInstalled:   false,
				}, nil
			}
		}

		return VersionResolution{}, &ErrNoMatchingVersion{AvailableVersions: releasedVersions, RequestedVersion: request.UnparsedVersion}
	case provider.ResolutionStrategyLatestReleased:
		sortedReleasedVersions := logicallySortedVersions(releasedVersions)
		if request.UnparsedVersion == "" {
			latestReleased := sortedReleasedVersions[0]
			if latestReleased == "" {
				return VersionResolution{}, &ErrNoMatchingVersion{
					AvailableVersions: releasedVersions,
					RequestedVersion:  "latest",
				}
			}
			isInstalled := slices.Contains(installedVersions, latestReleased)
			semverV, err := version.NewVersion(latestReleased)
			return VersionResolution{
				VersionString: latestReleased,
				IsSemVer:      err == nil,
				SemVer:        semverV,
				IsInstalled:   isInstalled,
			}, nil
		}
		for _, v := range sortedReleasedVersions {
			if strings.HasPrefix(v, request.UnparsedVersion) {
				// Since semver-compatible versions are sorted according to the semver spec
				// and are at the front of the list,
				// we can stop searching if the version prefix-matches the requested version.

				// Even though we search the released versions primarily,
				// it's still possible that the matching version is installed already.
				isInstalled := slices.Contains(installedVersions, v)

				semverV, err := version.NewVersion(v)
				return VersionResolution{
					VersionString: v,
					IsSemVer:      err == nil,
					SemVer:        semverV,
					IsInstalled:   isInstalled,
				}, nil
			}
		}
		return VersionResolution{}, &ErrNoMatchingVersion{AvailableVersions: releasedVersions, RequestedVersion: request.UnparsedVersion}
	}

	return VersionResolution{}, fmt.Errorf("unknown resolution strategy: %v", request.ResolutionStrategy)
}

// logicallySortedVersions reverse-sorts the given versions in a way that semver-compatible versions are sorted according to the semver spec,
// while non-semver versions are appended at the end in their own lexicographical order.
// This way, semver-compatible versions are prioritized over non-semver versions.
func logicallySortedVersions(versions []string) []string {
	var semverVersions version.Collection
	var nonSemverVersions []string
	for _, v := range versions {
		semverV, err := version.NewVersion(v)
		if err != nil {
			nonSemverVersions = append(nonSemverVersions, v)
			continue
		}
		semverVersions = append(semverVersions, semverV)
	}

	// semverVersions is of type version.Collection, which implements sort.Interface according to the semver spec.
	sort.Sort(sort.Reverse(semverVersions))
	// nonSemverVersions are only lexicographically sortable
	sort.Sort(sort.Reverse(sort.StringSlice(nonSemverVersions)))

	var sortedVersions []string
	for _, v := range semverVersions {
		sortedVersions = append(sortedVersions, v.Original())
	}

	sortedVersions = append(sortedVersions, nonSemverVersions...)
	return sortedVersions
}
