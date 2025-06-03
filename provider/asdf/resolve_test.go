package asdf_test

import (
	"testing"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/asdf"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestStrictResolution(t *testing.T) {
	tests := []struct {
		name               string
		requestedVersion   string
		installedVersions  []string
		releasedVersions   []string
		expectedResolution asdf.VersionResolution
		expectedErr        error
	}{
		{
			name:             "Exact match with installed version",
			requestedVersion: "1.0.0",
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
			name:             "Exact version but not installed",
			requestedVersion: "1.0.0",
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
			name:             "Nonexistent version",
			requestedVersion: "2.0.0",
			installedVersions: []string{
				"1.0.1",
				"1.1.0",
			},
			releasedVersions: []string{
				"1.0.0",
				"1.0.1",
				"1.1.0",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
		{
			name:             "Old Golang versioning scheme",
			requestedVersion: "1.19",
			installedVersions: []string{
				"1.18",
				"1.18.3",
				"1.20",
			},
			releasedVersions: []string{
				"1.18",
				"1.18.1",
				"1.18.2",
				"1.18.3",
				"1.19",
				"1.19.1",
				"1.19.5",
				"1.20",
				"1.20.1",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "1.19",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("1.19")),
				IsInstalled:   false,
			},
		},
		{
			name:             "Non-semver tool, exact match with installed version",
			requestedVersion: "temurin-21.0.0+35.0.LTS",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   true,
			},
		},
		{
			name:             "Non-semver tool, exact match but not installed",
			requestedVersion: "temurin-21.0.0+35.0.LTS",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   false,
			},
		},
		{
			name:             "Non-semver tool, nonexistent version",
			requestedVersion: "temurin-21.0.0+39.0",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
		{
			name:             "Non-semver tool, requested version is semver",
			requestedVersion: "21.0.0",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
	}

	runVersionResolutionTests(t, tests, provider.ResolutionStrategyStrict)
}

func TestLatestInstalledResolution(t *testing.T) {
	tests := []struct {
		name               string
		requestedVersion   string
		installedVersions  []string
		releasedVersions   []string
		expectedResolution asdf.VersionResolution
		expectedErr        error
	}{
		{
			name:             "Exact version but not installed",
			requestedVersion: "1.0.0",
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
			name:             "Exact version and installed",
			requestedVersion: "1.0.0",
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
			name:             "Partial match with installed version",
			requestedVersion: "20",
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
			name:             "Possibly-semver versions and partial match should correctly report installed state",
			requestedVersion: "20",
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
		{
			name:             "Old Golang versioning scheme",
			requestedVersion: "1.19",
			installedVersions: []string{
				"1.18",
				"1.18.3",
				"1.19.5",
				"1.20",
			},
			releasedVersions: []string{
				"1.18",
				"1.18.1",
				"1.18.2",
				"1.18.3",
				"1.19",
				"1.19.1",
				"1.19.5",
				"1.20",
				"1.20.1",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "1.19.5",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("1.19.5")),
				IsInstalled:   true,
			},
		},
		{
			name:             "No partial match for installed version, fallback to released version match",
			requestedVersion: "20.3",
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
				"20.3.0",
				"20.5.0",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "20.3.0",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("20.3.0")),
				IsInstalled:   false,
			},
		},
		{
			name:             "Nonexistent version",
			requestedVersion: "2.0.0",
			installedVersions: []string{
				"1.0.1",
				"1.1.0",
			},
			releasedVersions: []string{
				"1.0.0",
				"1.0.1",
				"1.1.0",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
		{
			name:             "Non-semver tool, exact match with installed version",
			requestedVersion: "temurin-21.0.0+35.0.LTS",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   true,
			},
		},
		{
			name:             "Non-semver tool, exact match but not installed",
			requestedVersion: "temurin-21.0.0+35.0.LTS",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   false,
			},
		},
		{
			name:             "Non-semver tool, partial match with installed version",
			requestedVersion: "temurin-21",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   true,
			},
		},
		{
			name:             "Non-semver tool, partial match but not installed",
			requestedVersion: "temurin-21",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   false,
			},
		},
		{
			name:             "Non-semver tool, partial match with correct natural ordering",
			requestedVersion: "temurin-11.0.15",
			installedVersions: []string{
				"temurin-11.0.15+10",
				"temurin-11.0.15+100",
				"temurin-11.0.15+101",
			},
			releasedVersions: []string{
				"temurin-11.0.15+10",
				"temurin-11.0.15+100",
				"temurin-11.0.15+101",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-11.0.15+101",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   true,
			},
		},
		{
			name:             "Non-semver tool, nonexistent version",
			requestedVersion: "temurin-21.0.0+39.0",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
		{
			name:             "Non-semver tool, requested version is semver",
			requestedVersion: "21.0.0",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
	}

	runVersionResolutionTests(t, tests, provider.ResolutionStrategyLatestInstalled)
}

func TestLatestReleasedResolution(t *testing.T) {
	tests := []struct {
		name               string
		requestedVersion   string
		installedVersions  []string
		releasedVersions   []string
		expectedResolution asdf.VersionResolution
		expectedErr        error
	}{
		{
			name:             "Partial match with installed latest version",
			requestedVersion: "20",
			installedVersions: []string{
				"18.6.3",
				"20.0.0",
				"20.1.0",
				"20.2.0",
				"20.5.0",
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
				IsInstalled:   true,
			},
		},
		{
			name:             "Partial match with both installed and non-installed versions",
			requestedVersion: "20",
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
			name:             "Partial match with released version only",
			requestedVersion: "20",
			installedVersions: []string{
				"18.6.3",
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
			name:             "Exact version matches installed version",
			requestedVersion: "18.6.3",
			installedVersions: []string{
				"18.6.3",
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
				VersionString: "18.6.3",
				IsSemVer:      true,
				SemVer:        version.Must(version.NewVersion("18.6.3")),
				IsInstalled:   true,
			},
		},
		{
			name:             "Nonexistent version",
			requestedVersion: "2.0.0",
			installedVersions: []string{
				"1.0.1",
				"1.1.0",
			},
			releasedVersions: []string{
				"1.0.0",
				"1.0.1",
				"1.1.0",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
		{
			name:             "Non-semver tool, exact match with both installed and released versions",
			requestedVersion: "temurin-21.0.0+35.0.LTS",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+33.0.LTS",
				"temurin-21.0.0+35.0.LTS",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+33.0.LTS",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   true,
			},
		},
		{
			name:             "Non-semver tool, exact match but not installed",
			requestedVersion: "temurin-21.0.0+35.0.LTS",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+33.0.LTS",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+33.0.LTS",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   false,
			},
		},
		{
			name:             "Non-semver tool, partial match with installed version",
			requestedVersion: "temurin-21",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+33.0.LTS",
				"temurin-21.0.0+35.0.LTS",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+33.0.LTS",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   true,
			},
		},
		{
			name:             "Non-semver tool, partial match but not installed",
			requestedVersion: "temurin-21",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+33.0.LTS",
				"temurin-21.0.0+35.0.LTS",
				"temurin-23.0.0+35.0.LTS",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-21.0.0+35.0.LTS",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   false,
			},
		},
		{
			name:             "Non-semver tool, partial match with correct natural ordering",
			requestedVersion: "temurin-11.0.15",
			installedVersions: []string{
				"temurin-11.0.15+10",
				"temurin-11.0.15+100",
				"temurin-11.0.15+101",
			},
			releasedVersions: []string{
				"temurin-11.0.15+10",
				"temurin-11.0.15+100",
				"temurin-11.0.15+101",
			},
			expectedResolution: asdf.VersionResolution{
				VersionString: "temurin-11.0.15+101",
				IsSemVer:      false,
				SemVer:        nil,
				IsInstalled:   true,
			},
		},
		{
			name:             "Non-semver tool, nonexistent version",
			requestedVersion: "temurin-21.0.0+39.0",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
		{
			name:             "Non-semver tool, requested version is semver",
			requestedVersion: "21.0.0",
			installedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
			},
			releasedVersions: []string{
				"openjdk-21",
				"oracle-21",
				"temurin-11.0.15+10",
				"temurin-17.0.4+101",
				"temurin-21.0.0+35.0.LTS",
			},
			expectedErr: asdf.ErrNoMatchingVersion,
		},
	}

	runVersionResolutionTests(t, tests, provider.ResolutionStrategyLatestReleased)
}

func runVersionResolutionTests(
	t *testing.T,
	tests []struct {
		name               string
		requestedVersion   string
		installedVersions  []string
		releasedVersions   []string
		expectedResolution asdf.VersionResolution
		expectedErr        error
	},
	strategy provider.ResolutionStrategy,
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			declaration := provider.ToolRequest{
				ToolName:           "test-tool",
				UnparsedVersion:    tt.requestedVersion,
				ResolutionStrategy: strategy,
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

// TODO: correct resolution among versions that only differ in pre-release tags and metadata
