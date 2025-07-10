package provider

type ToolPlugin struct {
	PluginName string
	PluginURL  string
}

var toolPluginMap = map[string]ToolPlugin{
	"flutter": {PluginName: "flutter", PluginURL: "https://github.com/asdf-community/asdf-flutter.git"},
	"tuist":   {PluginName: "tuist", PluginURL: "https://github.com/tuist/asdf-tuist.git"},
	"kotlin":  {PluginName: "kotlin", PluginURL: "https://github.com/asdf-community/asdf-kotlin.git"},
	"java":    {PluginName: "java", PluginURL: "https://github.com/halcyon/asdf-java.git"},
	"ruby":    {PluginName: "ruby", PluginURL: "https://github.com/asdf-vm/asdf-ruby.git"},
	"nodejs":  {PluginName: "nodejs", PluginURL: "https://github.com/asdf-vm/asdf-nodejs.git"},
}

func GetToolPlugin(toolName string) *ToolPlugin {
	canonicalName := GetCanonicalToolName(toolName)
	if toolPlugin, exists := toolPluginMap[canonicalName]; exists {
		return &toolPlugin
	}
	return nil
}
