package env_test

import (
	"testing"

	"github.com/tommyalmeida/envsync/internal/config"
	"github.com/tommyalmeida/envsync/pkg/schema"

	"github.com/tommyalmeida/envsync/internal/env"
)

func TestSyncer_Sync(t *testing.T) {
	cfg := &config.Config{
		Schema: schema.Schema{
			Variables: map[string]schema.Variable{
				"VAR1": {Default: "default1"},
				"VAR2": {Default: "default2"},
			},
		},
		Defaults: map[string]string{
			"VAR3": "config_default",
		},
	}

	syncer := env.NewSyncer(cfg)
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		source   env.Vars
		target   env.Vars
		dryRun   bool
		expected env.SyncResult
	}{
		{
			name: "sync missing variables",
			source: env.Vars{
				"VAR1": "source_value1",
				"VAR2": "source_value2",
				"VAR3": "source_value3",
			},
			target: env.Vars{
				"VAR1": "target_value1",
			},
			dryRun: true,
			expected: env.SyncResult{
				Added: []string{"VAR2", "VAR3"},
			},
		},
		{
			name: "no sync needed",
			source: env.Vars{
				"VAR1": "value1",
			},
			target: env.Vars{
				"VAR1": "value1",
			},
			dryRun: true,
			expected: env.SyncResult{
				Added: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetFile := tempDir + "/" + tt.name + ".env"

			result, err := syncer.Sync(tt.source, tt.target, targetFile, tt.dryRun)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result.Added) != len(tt.expected.Added) {
				t.Errorf("expected %d added variables, got %d", len(tt.expected.Added), len(result.Added))
			}
		})
	}
}
