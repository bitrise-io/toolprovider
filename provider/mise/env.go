package mise

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/bitrise-io/toolprovider/provider"
)

type envOutput map[string]string

// envVarsForTool returns the env vars required for the given tool version to be available and work correctly in
// a shell environment. This includes $PATH additions and other env vars, such as $JAVA_HOME, $GOROOT, etc.
func (m *MiseToolProvider) envVarsForTool(installResult provider.ToolInstallResult) (envOutput, error) {
	// Note: --quiet hides warnings and other plain text lines that would break JSON parsing.
	envCmd := exec.Command("mise", "env", "--quiet", "--json", fmt.Sprintf("%s@%s", installResult.ToolName, installResult.ConcreteVersion))
	data, err := envCmd.CombinedOutput()
	if err != nil {
		return envOutput{}, fmt.Errorf("mise env %s@%s: %w", installResult.ToolName, installResult.ConcreteVersion, err)
	}

	var env envOutput
	err = json.Unmarshal(data, &env)
	if err != nil {
		return envOutput{}, fmt.Errorf("parse mise env output: %w\n%s", err, string(data))
	}

	return env, nil
}

func processEnvOutput(envs envOutput) provider.EnvironmentActivation {
	// `mise env` returns tool-specific envs, as well as a new $PATH with the tool-specific dirs prepended.
	envsWithoutPath := maps.Clone(envs)
	delete(envsWithoutPath, "PATH")

	var pathsAddedByMise []string
	pathEnv, exists := envs["PATH"]
	if exists && pathEnv != "" {
		misePaths := strings.Split(pathEnv, ":")
		processPathEnv := os.Getenv("PATH")
		processPaths := strings.Split(processPathEnv, ":")

		// Track paths we've already added to avoid duplicates
		addedPaths := make(map[string]bool)
		for _, p := range misePaths {
			if p != "" && !slices.Contains(processPaths, p) && !addedPaths[p] {
				pathsAddedByMise = append(pathsAddedByMise, p)
				addedPaths[p] = true
			}
		}
	}

	return provider.EnvironmentActivation{
		ContributedEnvVars: envsWithoutPath,
		ContributedPaths:   pathsAddedByMise,
	}
}
