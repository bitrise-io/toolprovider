package provider

var toolAliasMap = map[string]string{
	"go":   "golang",
	"node": "nodejs",
}

func GetCanonicalToolName(toolName string) string {
	if canonicalName, exists := toolAliasMap[toolName]; exists {
		return canonicalName
	}
	return toolName
}
