package asdf_test

import (
	"testing"

	"github.com/bitrise-io/toolprovider"
	"github.com/bitrise-io/toolprovider/asdf"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestResolutionStrategies(t *testing.T) {
	tests := []struct {
		name               string
		requestedVersion   string
		strategy           toolprovider.ResolutionStrategy
		installedVersions  []string
		releasedVersions   []string
		expectedResolution asdf.VersionResolution
		expectedErr        error
	}{
		{
			name:             "Exact match with installed version, strict strategy",
			requestedVersion: "1.0.0",
			strategy:         toolprovider.ResolutionStrategyStrict,
			installedVersions: []string{
				"1.0.0",
				"1.0.1",
				"1.1.0",
			},
			releasedVersions: []string{
				"1.0.0",
				"1.0.1",
				"1.1.0",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "1.0.0",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("1.0.0")),
				IsInstalled:   true,
			},
		},
		{
			name:             "Exact version but not installed, strict strategy",
			requestedVersion: "1.0.0",
			strategy:         toolprovider.ResolutionStrategyStrict,
			installedVersions: []string{
				"1.0.1",
				"1.1.0",
			},
			releasedVersions: []string{
				"1.0.0",
				"1.0.1",
				"1.1.0",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "1.0.0",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("1.0.0")),
				IsInstalled:   false,
			},
		},
		{
			name:             "Exact version but not installed, latest installed strategy",
			requestedVersion: "1.0.0",
			strategy:         toolprovider.ResolutionStrategyLatestInstalled,
			installedVersions: []string{
				"1.0.1",
				"1.1.0",
			},
			releasedVersions: []string{
				"1.0.0",
				"1.0.1",
				"1.1.0",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "1.0.0",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("1.0.0")),
				IsInstalled:   false,
			},
		},
		{
			name:             "Partial match with installed version, latest installed strategy",
			requestedVersion: "20",
			strategy:         toolprovider.ResolutionStrategyLatestInstalled,
			installedVersions: []string{
				"18.6.3",
				"20.0.0",
				"20.1.0",
				"20.2.0",
				"21.0.0",
			},
			releasedVersions: []string{
				"18.6.3",
				"20.0.0",
				"20.1.0",
				"20.2.0",
				"20.5.0",
				"21.0.0",
				"22.0.0",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "20.2.0",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("20.2.0")),
				IsInstalled:   true,
			},
		},
		{
			name:             "Partial match with installed version, latest released strategy",
			requestedVersion: "20",
			strategy:         toolprovider.ResolutionStrategyLatestReleased,
			installedVersions: []string{
				"18.6.3",
				"20.0.0",
				"20.1.0",
				"20.2.0",
				"21.0.0",
			},
			releasedVersions: []string{
				"18.6.3",
				"20.0.0",
				"20.1.0",
				"20.2.0",
				"20.5.0",
				"21.0.0",
				"22.0.0",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "20.5.0",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("20.5.0")),
				IsInstalled:   false,
			},
		},
		{
			name:             "Possibly-semver versions and partial match should correctly report installed state",
			requestedVersion: "20",
			strategy:         toolprovider.ResolutionStrategyLatestInstalled,
			installedVersions: []string{
				"20.0.0",
				"20.1", // should be padded to 20.1.0
				"21.0.0",
			},
			releasedVersions: []string{
				"20.0.0",
				"20.1", // should be padded to 20.1.0
				"21.0.0",
				"22.0.0",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "20.1.0",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("20.1")),
				IsInstalled:   true,
			},
		},
		// TODO: non-semver versions (JDK, etc)
		// TODO: requesting a semver-compatible version when installed / released versions are non-semver
		// TODO: correct resolution among versions that only differ in pre-release tags and metadata
		// TODO: Go 1.20
		// TODO: matching a lower version (inverse of pessimistic operator)
		// TODO: ErrRequestedVersionNotSemVer
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			declaration := toolprovider.ToolRequest{
				ToolName:           "test-tool",
				UnparsedVersion:    tt.requestedVersion,
				ResolutionStrategy: tt.strategy,
			}

			resolvedV, err := asdf.ResolveVersion(
				declaration,
				tt.releasedVersions,
				tt.installedVersions,
			)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResolution, resolvedV)
			}
		})
	}
}
