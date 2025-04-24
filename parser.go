package toolprovider

import (
	"fmt"

	"github.com/bitrise-io/bitrise/v2/bitrise"
	"github.com/bitrise-io/bitrise/v2/models"
)

const keyExperimental = "experimental"
const keyToolDeclarations = "tools"

func ParseBitriseYml(path string) (models.BitriseDataModel, error) {
	model, _, err := bitrise.ReadBitriseConfig(path, bitrise.ValidationTypeMinimal)
	if err != nil {
		return models.BitriseDataModel{}, fmt.Errorf("parse bitrise.yml: %v", err)
	}

	return model, nil
}

func ParseToolDeclarations(bitriseYml models.BitriseDataModel) (map[string]ToolRequest, error) {
	if bitriseYml.Meta == nil {
		return nil, fmt.Errorf("parse bitrise.yml: meta block is not defined")
	}

	metaBlock := bitriseYml.Meta[keyExperimental].(map[string]any)
	if metaBlock == nil {
		return nil, fmt.Errorf("parse bitrise.yml: meta.%s block is not defined", keyExperimental)
	}

	toolBlock := metaBlock[keyToolDeclarations].(map[string]any)
	if toolBlock == nil {
		return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s block is not defined", keyExperimental, keyToolDeclarations)
	}

	toolDeclarations := make(map[string]ToolRequest)
	for toolName, toolData := range toolBlock {
		toolDataMap, ok := toolData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s block is not a map", keyExperimental, keyToolDeclarations, toolName)
		}

		// TODO: string or int
		version, ok := toolDataMap["version"].(string)
		if !ok {
			return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s.version is not a string", keyExperimental, keyToolDeclarations, toolName)
		}

		var resolutionStrategy ResolutionStrategy
		if toolDataMap["resolution_strategy"] != nil {
			resolutionStrategyString, ok := toolDataMap["resolution_strategy"].(string)
			if !ok {
				return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s.resolution_strategy is not a string", keyExperimental, keyToolDeclarations, toolName)
			}
			switch resolutionStrategyString {
			case "":
				resolutionStrategy = ResolutionStrategyStrict
			case "strict":
				resolutionStrategy = ResolutionStrategyStrict
			case "closest_installed":
				resolutionStrategy = ResolutionStrategyLatestInstalled
			case "closest_released":
				resolutionStrategy = ResolutionStrategyLatestReleased
			}
		}

		toolDeclarations[toolName] = ToolRequest{
			ToolName:           toolName,
			UnparsedVersion:    version,
			ResolutionStrategy: resolutionStrategy,
		}
	}

	return toolDeclarations, nil

}
