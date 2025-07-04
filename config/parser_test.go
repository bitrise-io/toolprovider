package config_test

import (
	"testing"

	"github.com/bitrise-io/toolprovider/provider"
	"github.com/bitrise-io/toolprovider/config"
	"github.com/stretchr/testify/assert"
)

func TestParseBitriseYml(t *testing.T) {
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
					ToolName:        "golang",
					UnparsedVersion: "1.16.3",
				},
				"nodejs": {
					ToolName:           "nodejs",
					UnparsedVersion:    "20",
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
				if actual.ToolName != expected.ToolName || actual.UnparsedVersion != expected.UnparsedVersion {
					t.Fatalf("%s: expected %v, got %v", key, expected, actual)
				}
			}
		})
	}
}
