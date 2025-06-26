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
	UnparsedVersion    string
	ResolutionStrategy ResolutionStrategy
	// TODO: PostInstall script
}

type ToolInstallResult struct {
	ToolName           string
	IsAlreadyInstalled bool
	ConcreteVersion    string
}

type EnvironmentActivation struct {
	ContributedEnvVars map[string]string
	ContributedPaths   []string
}

// TODO: make it generic over the struct of the provider options
type ToolProvider interface {
	Bootstrap() error

	InstallTool(tool ToolRequest) (ToolInstallResult, error)

	ActivateEnv(result ToolInstallResult) (EnvironmentActivation, error)

	// TODO: IsInstalledNative(tool ToolRequest) (bool, error)
}
