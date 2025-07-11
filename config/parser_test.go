package config_test

import (
	"testing"

	"github.com/bitrise-io/toolprovider/config"
	"github.com/bitrise-io/toolprovider/provider"
	"github.com/stretchr/testify/assert"
)

var flutterPlugin = "flutter::https://github.com/asdf-community/asdf-flutter.git"

func TestParseTools(t *testing.T) {
	tests := []struct {
		name     string
		ymlPath  string
		expected map[string]provider.ToolRequest
	}{
		{
			name:    "Valid YML",
			ymlPath: "testdata/valid.bitrise.yml",
			expected: map[string]provider.ToolRequest{
				"golang": {
					ToolName:           "golang",
					UnparsedVersion:    "1.16.3",
					ResolutionStrategy: provider.ResolutionStrategyStrict,
				},
				"nodejs": {
					ToolName:           "nodejs",
					UnparsedVersion:    "20",
					ResolutionStrategy: provider.ResolutionStrategyLatestInstalled,
				},
				"ruby": {
					ToolName:           "ruby",
					UnparsedVersion:    "3.2",
					ResolutionStrategy: provider.ResolutionStrategyLatestReleased,
				},
				"flutter": {
					ToolName:           "flutter",
					UnparsedVersion:    "3.32.5-stable",
					ResolutionStrategy: provider.ResolutionStrategyStrict,
					PluginIdentifier:   &flutterPlugin,
				},
				"python": {
					ToolName:           "python",
					UnparsedVersion:    "3.13",
					ResolutionStrategy: provider.ResolutionStrategyLatestInstalled,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bitriseYml, err := config.ParseBitriseYml(tt.ymlPath)
			assert.NoError(t, err)
			assert.NotNil(t, bitriseYml)

			toolDeclarations, err := config.ParseToolDeclarations(bitriseYml)
			assert.NoError(t, err)

			if len(toolDeclarations) != len(tt.expected) {
				t.Fatalf("expected %d tool declarations, got %d", len(tt.expected), len(toolDeclarations))
			}

			for key, expected := range tt.expected {
				actual, exists := toolDeclarations[key]
				if !exists {
					t.Fatalf("expected tool declaration for %s not found", key)
				}
				if actual.ToolName != expected.ToolName || actual.UnparsedVersion != expected.UnparsedVersion || actual.ResolutionStrategy != expected.ResolutionStrategy {
					t.Fatalf("%s: expected %v, got %v", key, expected, actual)
				}
				if expected.PluginIdentifier != nil && actual.PluginIdentifier == nil {
					t.Fatalf("%s: expected plugin identifier %s, got nil", key, *expected.PluginIdentifier)
				}
				if expected.PluginIdentifier == nil && actual.PluginIdentifier != nil {
					t.Fatalf("%s: expected no plugin identifier, got %s", key, *actual.PluginIdentifier)
				}
				if expected.PluginIdentifier != nil && actual.PluginIdentifier != nil && *expected.PluginIdentifier != *actual.PluginIdentifier {
					t.Fatalf("%s: expected plugin identifier %s, got %s", key, *expected.PluginIdentifier, *actual.PluginIdentifier)
				}
			}
		})
	}
}

func TestParseToolConfig(t *testing.T) {
	tests := []struct {
		name     string
		ymlPath  string
		expected config.ToolConfig
	}{
		{
			name:    "No explicit config",
			ymlPath: "testdata/valid.bitrise.yml",
			expected: config.ToolConfig{
				Provider: "asdf",
			},
		},
		{
			name:    "Custom tool config",
			ymlPath: "testdata/custom_config.bitrise.yml",
			expected: config.ToolConfig{
				Provider: "asdf",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bitriseYml, err := config.ParseBitriseYml(tt.ymlPath)
			assert.NoError(t, err)
			assert.NotNil(t, bitriseYml)

			toolConfig, err := config.ParseToolConfig(bitriseYml)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, toolConfig)
		})
	}
}
