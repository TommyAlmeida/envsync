package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/tommyalmeida/envsync/internal/config"
)

func TestLoad_DefaultConfig(t *testing.T) {
	viper.Reset()

	cfg, err := config.Load()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.True(t, cfg.Rules.AllowExtra)
	require.NotNil(t, cfg.Defaults)
	require.Empty(t, cfg.Schema.Variables)
}

func TestLoad_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yaml")

	// Fyi this doesn't really represent a complete config, just a setup to test the config loading
	content := `schema:
  variables:
    FOO:
      description: "Foo variable"
      default: "bar"
      required: true
      type: "string"
defaults:
  BAR: "very much bar here" 
rules:
  require_all: true
  allow_extra: false
  ignore_patterns:
    - "IGNORED_*"
`

	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	viper.SetConfigFile(filePath)

	cfg, err := config.Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Contains(t, cfg.Schema.Variables, "FOO")
	v := cfg.Schema.Variables["FOO"]
	require.Equal(t, "bar", v.Default)
	require.Equal(t, "Foo variable", v.Description)
	require.True(t, v.Required)
	require.Equal(t, "string", v.Type)

	require.Equal(t, "very much bar here", cfg.Defaults["BAR"])

	require.True(t, cfg.Rules.RequireAll)
	require.False(t, cfg.Rules.AllowExtra)
	require.Len(t, cfg.Rules.IgnorePatterns, 1)
	require.Equal(t, "IGNORED_*", cfg.Rules.IgnorePatterns[0])
}
