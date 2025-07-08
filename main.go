package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/toolprovider/config"
	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/provider/asdf"
	"github.com/bitrise-io/toolprovider/provider/asdf/execenv"
)

func main() {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	bitriseModel, err := config.ParseBitriseYml(filepath.Join(workdir, "bitrise.yml"))
	if err != nil {
		panic(err)
	}
	toolDeclarations, err := config.ParseToolDeclarations(bitriseModel)
	if err != nil {
		panic(err)
	}

	if len(toolDeclarations) == 0 {
		fmt.Println("No tools to set up.")
		return
	}

	fmt.Println("Tools to set up:")

	for toolName, toolRequest := range toolDeclarations {
		resolutionStrategy := "strict" // default
		switch toolRequest.ResolutionStrategy {
		case provider.ResolutionStrategyLatestInstalled:
			resolutionStrategy = "closest_installed"
		case provider.ResolutionStrategyLatestReleased:
			resolutionStrategy = "closest_released"
		}

		fmt.Printf("- %s v%s (resolution: %s)\n",
			toolName,
			toolRequest.UnparsedVersion,
			resolutionStrategy)
	}

	fmt.Println()
	fmt.Println("Installing any missing tools...")

	asdfProvider := asdf.AsdfToolProvider{
		ExecEnv: execenv.ExecEnv{
			EnvVars:   convertEnvToMap(os.Environ()),
			ShellInit: "", // TODO
		},
	}
	var toolInstalls []provider.ToolInstallResult
	for toolName, toolRequest := range toolDeclarations {
		canonicalToolName := provider.GetCanonicalToolName(toolName)
		toolRequest.ToolName = canonicalToolName

		fmt.Printf("Installing %s v%s...\n", canonicalToolName, toolRequest.UnparsedVersion)
		result, err := asdfProvider.InstallTool(toolRequest)
		if err != nil {
			panic(err)
		}
		toolInstalls = append(toolInstalls, result)

		if result.IsAlreadyInstalled {
			fmt.Printf("%s v%s is already installed.\n", result.ToolName, result.ConcreteVersion)
		} else {
			fmt.Printf("Successfully installed %s v%s.\n", result.ToolName, result.ConcreteVersion)
		}
	}

	fmt.Println()
	fmt.Println("Activating environment with envman...")
	if os.Getenv("CI") == "" {
		_ = exec.Command("envman", "init").Run()
	}

	for _, install := range toolInstalls {
		activation, err := asdfProvider.ActivateEnv(install)
		if err != nil {
			panic(fmt.Errorf("activate tool %s: %w", install.ToolName, err))
		}
		err = extendEnvmanEnv(activation)
		if err != nil {
			panic(fmt.Errorf("extend envman env for %s: %w", install.ToolName, err))
		}
		fmt.Printf("Environment for %s activated.\n", install.ToolName)
	}
	fmt.Print("Environment setup complete!")

}

func convertEnvToMap(env []string) map[string]string {
	result := make(map[string]string)
	for _, envVar := range env {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func extendEnvmanEnv(activation provider.EnvironmentActivation) error {
	for k, v := range activation.ContributedEnvVars {
		cmd := exec.Command("envman", "add", "--key", k, "--value", v)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Println(string(out))
			return fmt.Errorf("add $%s to env: %w", k, err)
		}
	}

	if len(activation.ContributedPaths) > 0 {
		newPath := os.Getenv("PATH")
		for _, p := range activation.ContributedPaths {
			if !strings.Contains(newPath, p) {
				newPath = fmt.Sprintf("%s:%s", p, newPath)
			}
		}
		cmd := exec.Command("envman", "add", "--key", "PATH", "--value", newPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Println(string(out))
			return fmt.Errorf("update $PATH: %w", err)
		}
	}
	return nil
}
