package provider

import "fmt"

type ResolutionStrategy int

const (
	// TODO: ResolutionStrategyAuto?
	ResolutionStrategyStrict ResolutionStrategy = iota
	ResolutionStrategyLatestInstalled
	ResolutionStrategyLatestReleased
)

type ToolRequest struct {
	ToolName string
	// UnparsedVersion is the version string as provided by the user.
	// It may or may not be a valid semantic version.
	UnparsedVersion    string
	ResolutionStrategy ResolutionStrategy
	// PluginIdentifier is an optional identifier for the tool plugin.
	PluginIdentifier *string
	// TODO: PostInstall script
}

type ToolInstallResult struct {
	ToolName           string
	IsAlreadyInstalled bool
	// ConcreteVersion is the version that was actually installed and we resolved to.
	// It may differ from the requested version if the requested version was not a concrete version.
	// This value may or may not be a valid semantic version.
	ConcreteVersion string
}

type ToolInstallError struct {
	ToolName         string
	RequestedVersion string

	// Optional fields
	RawOutput      string
	Cause          string
	Recommendation string
}

func (e ToolInstallError) Error() string {
	msg := fmt.Sprintf("Error: failed to install %s %s", e.ToolName, e.RequestedVersion)

	if e.Cause != "" {
		msg += "\nCause: " + e.Cause
	}

	if e.Recommendation != "" {
		msg += "\nRecommendation: " + e.Recommendation
	}

	if e.RawOutput != "" {
		msg += "\nAdditional info: " + e.RawOutput
	}

	return msg
}

type EnvironmentActivation struct {
	ContributedEnvVars map[string]string
	ContributedPaths   []string
}

type ToolProvider interface {
	Bootstrap() error

	InstallTool(tool ToolRequest) (ToolInstallResult, error)

	ActivateEnv(result ToolInstallResult) (EnvironmentActivation, error)

	// TODO: IsInstalledNative(tool ToolRequest) (bool, error)
}
