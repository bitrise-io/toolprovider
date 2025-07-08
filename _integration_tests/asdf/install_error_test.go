package integration_tests

import (
	"testing"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/asdf"
	"github.com/bitrise-io/toolprovider/provider/asdf/execenv"
	"github.com/stretchr/testify/require"
)

func TestNoMatchingVersionError(t *testing.T) {
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
	request := provider.ToolRequest{
		ToolName:        "nodejs",
		UnparsedVersion: "22",
		ResolutionStrategy: provider.ResolutionStrategyStrict,
	}
	_, err = asdfProvider.InstallTool(request)
	require.Error(t, err)

	var installErr provider.ToolInstallError
	require.ErrorAs(t, err, &installErr)
	require.Equal(t, "nodejs", installErr.ToolName)
	require.Equal(t, "22", installErr.RequestedVersion)
	require.Contains(t, installErr.Error(), "No exact match found for 22")
	require.Contains(t, installErr.Recommendation, "22:latest")
	require.Contains(t, installErr.Recommendation, "22:installed")
}
