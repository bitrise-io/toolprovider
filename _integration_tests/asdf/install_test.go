package integration_tests

import (
	"testing"

	"github.com/bitrise-io/toolprovider"
	"github.com/bitrise-io/toolprovider/asdf"
	"github.com/stretchr/testify/require"
)

func TestAsdfInstall(t *testing.T) {
	// testEnv, err := createTestEnv(t, asdfInstallation{
	// 	flavor:  flavorAsdfRewrite,
	// 	version: "0.17.0",
	// 	plugins: []string{"nodejs"},
	// })
	testEnv, err := createTestEnv(t, asdfInstallation{
		flavor:  flavorAsdfClassic,
		version: "0.14.0",
		plugins: []string{"nodejs"},
	})
	require.NoError(t, err)

	provider := asdf.AsdfToolProvider{
		ExecEnv: asdf.ExecEnv{
			EnvVars:   testEnv.envVars,
			ShellInit: testEnv.shellInit,
		},
	}

	request := toolprovider.ToolRequest{
		ToolName:        "nodejs",
		UnparsedVersion: "18.16.0",
	}
	result, err := provider.InstallTool(request)
	require.NoError(t, err)
	require.Equal(t, "nodejs", result.ToolName)
	require.Equal(t, "18.16.0", result.ConcreteVersion)
	require.False(t, result.IsAlreadyInstalled)
}
