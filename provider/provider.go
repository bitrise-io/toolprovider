package provider

type ResolutionStrategy int

const (
	// TODO: ResolutionStrategyAuto?
	ResolutionStrategyStrict ResolutionStrategy = iota
	ResolutionStrategyLatestInstalled
	ResolutionStrategyLatestReleased
)

type ToolRequest struct {
	ToolName           string
	// UnparsedVersion is the version string as provided by the user.
	// It may or may not be a valid semantic version.
	UnparsedVersion    string
	ResolutionStrategy ResolutionStrategy
	// TODO: PostInstall script
}

type ToolInstallResult struct {
	ToolName           string
	IsAlreadyInstalled bool
	// ConcreteVersion is the version that was actually installed and we resolved to.
	// It may differ from the requested version if the requested version was not a concrete version.
	// This value may or may not be a valid semantic version.
	ConcreteVersion    string
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
