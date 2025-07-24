package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bitrise-io/bitrise/v2/bitrise"
	"github.com/bitrise-io/bitrise/v2/models"
	"github.com/bitrise-io/toolprovider/provider"
)

const keyExperimental = "experimental"
const keyToolDeclarations = "tools"
const keyToolConfig = "tool_config"
const latestSyntaxPattern = `(.*):latest$`
const installedSyntaxPattern = `(.*):installed$`

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
	for toolName, toolData := range toolBlock {
		// TODO: string or int
		var versionString string
		var pluginIdentifier *string

		switch v := toolData.(type) {
		case string:
			// If it's a string, it should be a version string only.
			versionString = strings.TrimSpace(v)
		case map[string]any:
			// If it's a map, it should contain version and optionally plugin fields.
			ver, ok := v["version"].(string)
			if !ok {
				if v["version"] == nil {
					// User aims to provide the shortest form, skipping version field, so we decide later what to do with no version.
					ver = ""
				} else {
					return nil, fmt.Errorf("parse bitrise.yml: meta.%s.%s.%s.version is not a string", keyExperimental, keyToolDeclarations, toolName)
				}
			}
			versionString = strings.TrimSpace(ver)
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

func defaultToolConfig() ToolConfig {
	return ToolConfig{
		Provider: "asdf",
	}
}

func ParseToolConfig(bitriseYml models.BitriseDataModel) (ToolConfig, error) {
	if bitriseYml.Meta == nil {
		return ToolConfig{}, fmt.Errorf("parse bitrise.yml: meta block is not defined")
	}

	metaBlock := bitriseYml.Meta[keyExperimental].(map[string]any)
	if metaBlock == nil {
		return ToolConfig{}, fmt.Errorf("parse bitrise.yml: meta.%s block is not defined", keyExperimental)
	}

	toolConfig := defaultToolConfig()

	toolConfigBlock, exists := metaBlock[keyToolConfig]
	if !exists {
		return toolConfig, nil // No explicit tool config, return default
	}

	toolConfigMap, ok := toolConfigBlock.(map[string]any)
	if !ok {
		return ToolConfig{}, fmt.Errorf("parse bitrise.yml: meta.%s.%s block is not a map", keyExperimental, keyToolConfig)
	}
	if toolConfigMap == nil {
		return toolConfig, nil
	}

	for key, value := range toolConfigMap {
		switch key {
		case "provider":
			toolConfig.Provider = value.(string)
		}
	}

	return toolConfig, nil
}
