package config

import (
	"fmt"
	"regexp"

	"github.com/bitrise-io/bitrise/v2/bitrise"
	"github.com/bitrise-io/bitrise/v2/models"
	"github.com/bitrise-io/toolprovider/provider"
)

const keyExperimental = "experimental"
const keyToolDeclarations = "tools"
const latestSyntaxPattern = `(.+):latest$`
const installedSyntaxPattern = `(.+):installed$`

func ParseBitriseYml(path string) (models.BitriseDataModel, error) {
	model, _, err := bitrise.ReadBitriseConfig(path, bitrise.ValidationTypeMinimal)
	if err != nil {
		return models.BitriseDataModel{}, fmt.Errorf("parse bitrise.yml: %v", err)
	}

	return model, nil
}

func ParseToolDeclarations(bitriseYml models.BitriseDataModel) (map[string]provider.ToolRequest, error) {
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

	latestSyntaxPattern, err := regexp.Compile(latestSyntaxPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex pattern: %v", err)
	}
	preinstalledSyntaxPattern, err := regexp.Compile(installedSyntaxPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex pattern: %v", err)
	}

	toolDeclarations := make(map[string]provider.ToolRequest)
	for toolName, toolVersion := range toolBlock {
		// TODO: string or int
		var versionString string
		var pluginIdentifier *string

		switch v := toolVersion.(type) {
		case string:
			versionString = v
		case map[string]any:
			ver, ok := v["version"].(string)
			if !ok {
				return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s.version is not a string", keyExperimental, keyToolDeclarations, toolName)
			}
			versionString = ver
			if pluginVal, ok := v["plugin"]; ok && pluginVal != nil {
				pluginStr, ok := pluginVal.(string)
				if !ok {
					return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s.plugin is not a string", keyExperimental, keyToolDeclarations, toolName)
				}
				pluginIdentifier = &pluginStr
			}
		default:
			return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s is not a string or map", keyExperimental, keyToolDeclarations, toolName)
		}

		var resolutionStrategy provider.ResolutionStrategy
		var plainVersion string
		if latestSyntaxPattern.MatchString(versionString) {
			resolutionStrategy = provider.ResolutionStrategyLatestReleased
			matches := latestSyntaxPattern.FindStringSubmatch(versionString)
			if len(matches) > 1 {
				plainVersion = matches[1]
			} else {
				return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s.version does not match latest syntax: %s", keyExperimental, keyToolDeclarations, toolName, versionString)
			}
		} else if preinstalledSyntaxPattern.MatchString(versionString) {
			resolutionStrategy = provider.ResolutionStrategyLatestInstalled
			matches := preinstalledSyntaxPattern.FindStringSubmatch(versionString)
			if len(matches) > 1 {
				plainVersion = matches[1]
			} else {
				return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s.version does not match preinstalled syntax: %s", keyExperimental, keyToolDeclarations, toolName, versionString)
			}
		} else {
			resolutionStrategy = provider.ResolutionStrategyStrict
			plainVersion = versionString
		}

		toolDeclarations[toolName] = provider.ToolRequest{
			ToolName:           toolName,
			UnparsedVersion:    plainVersion,
			ResolutionStrategy: resolutionStrategy,
			PluginIdentifier:   pluginIdentifier,
		}
	}

	return toolDeclarations, nil
}
