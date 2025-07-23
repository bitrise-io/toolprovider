package mise

import (
	"fmt"

	"github.com/bitrise-io/toolprovider/provider"
)

type MiseToolProvider struct {
}

func (m *MiseToolProvider) ID() string {
	return "mise"
}

func (m *MiseToolProvider) Bootstrap() error {
	// TODO
	return nil
}

func (m *MiseToolProvider) InstallTool(tool provider.ToolRequest) (provider.ToolInstallResult, error) {
	isAlreadyInstalled, err := isAlreadyInstalled(tool, m.resolveToLatestInstalled)
	if err != nil {
		return provider.ToolInstallResult{}, err
	}

	err = m.installToolVersion(tool)
	if err != nil {
		return provider.ToolInstallResult{}, err
	}

	concreteVersion, err := m.resolveToConcreteVersionAfterInstall(tool)
	if err != nil {
		return provider.ToolInstallResult{}, fmt.Errorf("resolve exact version after install: %w", err)
	}

	return provider.ToolInstallResult{
		ToolName:           tool.ToolName,
		IsAlreadyInstalled: isAlreadyInstalled,
		ConcreteVersion:    concreteVersion,
	}, nil
}

func (m *MiseToolProvider) ActivateEnv(result provider.ToolInstallResult) (provider.EnvironmentActivation, error) {
	envs, err := m.envVarsForTool(result)
	if err != nil {
		return provider.EnvironmentActivation{}, fmt.Errorf("get mise env: %w", err)
	}

	return processEnvOutput(envs), nil
}
