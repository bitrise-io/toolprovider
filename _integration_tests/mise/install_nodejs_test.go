package mise

import (
	"testing"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/mise"
	"github.com/stretchr/testify/require"
)

func TestMiseInstallNodeVersion(t *testing.T) {
	tests := []struct {
		name               string
		requestedVersion   string
		resolutionStrategy provider.ResolutionStrategy
		expectedVersion    string
	}{
		{"Install specific version", "18.16.0", provider.ResolutionStrategyStrict, "18.16.0"},
		{"Install partial major version", "20", provider.ResolutionStrategyLatestInstalled, "20.19.3"},
		{"Install partial major.minor version", "18.10", provider.ResolutionStrategyLatestReleased, "18.10.0"},
	}

	for _, tt := range tests {
		miseProvider := mise.MiseToolProvider{}
		t.Run(tt.name, func(t *testing.T) {
			request := provider.ToolRequest{
				ToolName:           "nodejs",
				UnparsedVersion:    tt.requestedVersion,
				ResolutionStrategy: tt.resolutionStrategy,
			}
			result, err := miseProvider.InstallTool(request)
			require.NoError(t, err)
			require.Equal(t, "nodejs", result.ToolName)
			require.Equal(t, tt.expectedVersion, result.ConcreteVersion)
			require.False(t, result.IsAlreadyInstalled)
		})
	}
}
