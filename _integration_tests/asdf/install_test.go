package integration_tests

import (
	"testing"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/asdf"
	"github.com/stretchr/testify/require"
)

func TestAsdfInstallClassic(t *testing.T) {
	testEnv, err := createTestEnv(t, asdfInstallation{
		flavor:  flavorAsdfClassic,
		version: "0.14.0",
		plugins: []string{"nodejs"},
	})
	require.NoError(t, err)

	asdfProvider := asdf.AsdfToolProvider{
		ExecEnv: asdf.ExecEnv{
			EnvVars:   testEnv.envVars,
			ShellInit: testEnv.shellInit,
		},
	}

	request := provider.ToolRequest{
		ToolName:        "nodejs",
		UnparsedVersion: "18.16.0",
	}
	result, err := asdfProvider.InstallTool(request)
	require.NoError(t, err)
	require.Equal(t, "nodejs", result.ToolName)
	require.Equal(t, "18.16.0", result.ConcreteVersion)
	require.False(t, result.IsAlreadyInstalled)
}

func TestAsdfInstallRewrite(t *testing.T) {
	testEnv, err := createTestEnv(t, asdfInstallation{
		flavor:  flavorAsdfRewrite,
		version: "0.18.0",
		plugins: []string{"nodejs"},
	})
	require.NoError(t, err)

	asdfProvider := asdf.AsdfToolProvider{
		ExecEnv: asdf.ExecEnv{
			EnvVars:   testEnv.envVars,
			ShellInit: testEnv.shellInit,
		},
	}

	request := provider.ToolRequest{
		ToolName:        "nodejs",
		UnparsedVersion: "18.16.0",
	}
	result, err := asdfProvider.InstallTool(request)
	require.NoError(t, err)
	require.Equal(t, "nodejs", result.ToolName)
	require.Equal(t, "18.16.0", result.ConcreteVersion)
	require.False(t, result.IsAlreadyInstalled)
}
