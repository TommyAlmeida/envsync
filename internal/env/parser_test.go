package env_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tommyalmeida/envsync/internal/env"
)

func TestParseFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expected    env.Vars
		expectError bool
	}{
		{
			name: "valid env file",
			content: `DATABASE_URL=postgres://localhost/test
                      PORT=3000
                      DEBUG=true
                      EMPTY_VAR=`,
			expected: env.Vars{
				"DATABASE_URL": "postgres://localhost/test",
				"PORT":         "3000",
				"DEBUG":        "true",
				"EMPTY_VAR":    "",
			},
			expectError: false,
		},
		{
			name: "env file with comments",
			content: `# This is a comment
                      DATABASE_URL=postgres://localhost/test
                      # Another comment
                      PORT=3000`,
			expected: env.Vars{
				"DATABASE_URL": "postgres://localhost/test",
				"PORT":         "3000",
			},
			expectError: false,
		},
		{
			name: "env file with quoted values",
			content: `MESSAGE="Hello World"
                      PATH="/usr/local/bin:/usr/bin"`,
			expected: env.Vars{
				"MESSAGE": "Hello World",
				"PATH":    "/usr/local/bin:/usr/bin",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(tempDir, tt.name+".env")
			err := os.WriteFile(filename, []byte(tt.content), 0644)

			if err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			result, err := env.ParseFile(filename)

			if tt.expectError && err == nil {
				t.Error("expected error, got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestParseFile_NonExistentFile(t *testing.T) {
	_, err := env.ParseFile("non-existent-file.env")
	if err == nil {
		t.Error("expected error for non-existent file, got none")
	}
}

func TestEnvVars_WriteToFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.env")

	envVars := env.Vars{
		"DATABASE_URL": "postgres://localhost/test",
		"PORT":         "3000",
		"DEBUG":        "true",
		"MESSAGE":      "Hello World",
	}

	err := envVars.WriteToFile(filename)

	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	result, err := env.ParseFile(filename)

	if err != nil {
		t.Fatalf("failed to read back file: %v", err)
	}

	if !reflect.DeepEqual(result, envVars) {
		t.Errorf("expected %+v, got %+v", envVars, result)
	}
}
