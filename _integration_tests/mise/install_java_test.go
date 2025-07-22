package mise

import (
	"testing"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/mise"
	"github.com/stretchr/testify/require"
)

func TestMiseInstallJavaVersion(t *testing.T) {
	tests := []struct {
		name               string
		requestedVersion   string
		resolutionStrategy provider.ResolutionStrategy
		expectedVersion    string
	}{
		{
			name:               "OpenJDK major version only",
			requestedVersion:   "21",
			resolutionStrategy: provider.ResolutionStrategyStrict,
			expectedVersion:    "21.0.2",
		},
		{
			name:               "OpenJDK major version only, latest released",
			requestedVersion:   "17",
			resolutionStrategy: provider.ResolutionStrategyLatestReleased,
			expectedVersion:    "17.0.2",
		},
		{
			name:               "Temurin major version only",
			requestedVersion:   "temurin-21",
			resolutionStrategy: provider.ResolutionStrategyLatestReleased,
			expectedVersion:    "temurin-21.0.7+6.0.LTS",
		},
		{
			name:               "Temurin exact version",
			requestedVersion:   "temurin-17.0.8+101",
			resolutionStrategy: provider.ResolutionStrategyStrict,
			expectedVersion:    "temurin-17.0.8+101",
		},
	}

	for _, tt := range tests {
		miseProvider := mise.MiseToolProvider{}
		t.Run(tt.name, func(t *testing.T) {
			request := provider.ToolRequest{
				ToolName:           "java",
				UnparsedVersion:    tt.requestedVersion,
				ResolutionStrategy: tt.resolutionStrategy,
			}
			result, err := miseProvider.InstallTool(request)
			require.NoError(t, err)
			require.Equal(t, "java", result.ToolName)
			require.Equal(t, tt.expectedVersion, result.ConcreteVersion)
			require.False(t, result.IsAlreadyInstalled)
		})
	}
}
