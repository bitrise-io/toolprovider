package asdf

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/bitrise-io/bitrise/v2/log"
	"github.com/bitrise-io/toolprovider/provider"
	"github.com/hashicorp/go-version"
)

type ErrNoMatchingVersion struct {
	AvailableVersions []string
}

func (e ErrNoMatchingVersion) Error() string {
	if len(e.AvailableVersions) == 0 {
		return "no matching version found"
	}

	versionList := ""
	for _, v := range e.AvailableVersions {
		versionList += fmt.Sprintf("- %s\n", v)
	}
	return fmt.Sprintf("no matching version found, available versions: \n%s", versionList)
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
		if request.ResolutionStrategy != provider.ResolutionStrategyStrict {
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

	if request.ResolutionStrategy == provider.ResolutionStrategyStrict {
		if slices.Contains(releasedVersions, request.UnparsedVersion) {
			requestedSemVer, err := version.NewVersion(request.UnparsedVersion)
			return VersionResolution{
				VersionString: request.UnparsedVersion,
				IsSemVer:      err == nil,
				SemVer:        requestedSemVer,
				IsInstalled:   slices.Contains(installedVersions, request.UnparsedVersion),
			}, nil
		}
		return VersionResolution{}, ErrNoMatchingVersion{AvailableVersions: releasedVersions}
	}

	switch request.ResolutionStrategy {
	case provider.ResolutionStrategyLatestInstalled:
		// Installed versions are checked first because strategy is "latest installed"
		sortedInstalledVersions := logicallySortedVersions(installedVersions)
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

		return VersionResolution{}, ErrNoMatchingVersion{AvailableVersions: releasedVersions}
	case provider.ResolutionStrategyLatestReleased:
		sortedReleasedVersions := logicallySortedVersions(releasedVersions)
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
		return VersionResolution{}, ErrNoMatchingVersion{AvailableVersions: releasedVersions}
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
