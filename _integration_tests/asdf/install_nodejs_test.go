package integration_tests

import (
	"testing"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/asdf"
	"github.com/bitrise-io/toolprovider/provider/asdf/execenv"
	"github.com/stretchr/testify/require"
)

func TestAsdfInstallNodeVersion(t *testing.T) {
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
		testEnv, err := createTestEnv(t, asdfInstallation{
			flavor:  flavorAsdfClassic,
			version: "0.14.0",
			plugins: []string{"nodejs"},
		})
		require.NoError(t, err)

		asdfProvider := asdf.AsdfToolProvider{
			ExecEnv: execenv.ExecEnv{
				EnvVars:   testEnv.envVars,
				ShellInit: testEnv.shellInit,
			},
		}
		t.Run(tt.name, func(t *testing.T) {
			request := provider.ToolRequest{
				ToolName:           "nodejs",
				UnparsedVersion:    tt.requestedVersion,
				ResolutionStrategy: tt.resolutionStrategy,
			}
			result, err := asdfProvider.InstallTool(request)
			require.NoError(t, err)
			require.Equal(t, "nodejs", result.ToolName)
			require.Equal(t, tt.expectedVersion, result.ConcreteVersion)
			require.False(t, result.IsAlreadyInstalled)
		})
	}
}
