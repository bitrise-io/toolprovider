package main

import "testing"

func TestPrependPath(t *testing.T) {
	tests := []struct {
		name        string
		pathEnv     string
		newPath     string
		expected    string
	}{
		{
			name:        "empty path env",
			pathEnv:     "",
			newPath:     "/usr/local/bin",
			expected:    "/usr/local/bin",
		},
		{
			name:        "prepend to existing path",
			pathEnv:     "/usr/bin:/bin",
			newPath:     "/usr/local/bin",
			expected:    "/usr/local/bin:/usr/bin:/bin",
		},
		{
			name:        "remove duplicate and prepend",
			pathEnv:     "/usr/bin:/usr/local/bin:/bin",
			newPath:     "/usr/local/bin",
			expected:    "/usr/local/bin:/usr/bin:/bin",
		},
		{
			name:        "duplicate at end",
			pathEnv:     "/usr/bin:/bin:/usr/local/bin",
			newPath:     "/usr/local/bin",
			expected:    "/usr/local/bin:/usr/bin:/bin",
		},
		{
			name:        "single path duplicate",
			pathEnv:     "/usr/local/bin",
			newPath:     "/usr/local/bin",
			expected:    "/usr/local/bin",
		},
		{
			name:        "empty new path",
			pathEnv:     "/usr/bin:/bin",
			newPath:     "",
			expected:    ":/usr/bin:/bin",
		},
		{
			name:        "multiple duplicates",
			pathEnv:     "/usr/local/bin:/usr/bin:/usr/local/bin:/bin:/usr/local/bin",
			newPath:     "/usr/local/bin",
			expected:    "/usr/local/bin:/usr/bin:/bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prependPath(tt.pathEnv, tt.newPath)
			if result != tt.expected {
				t.Errorf("prependPath(%q, %q) = %q, want %q", tt.pathEnv, tt.newPath, result, tt.expected)
			}
		})
	}
}
