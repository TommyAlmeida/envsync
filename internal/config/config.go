package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/tommyalmeida/envsync/pkg/schema"
)

type Config struct {
	Schema   schema.Schema     `yaml:"schema"`
	Defaults map[string]string `yaml:"defaults"`
	Rules    Rules             `yaml:"rules"`
}

type Rules struct {
	RequireAll     bool     `yaml:"require_all"`
	AllowExtra     bool     `yaml:"allow_extra"`
	IgnorePatterns []string `yaml:"ignore_patterns"`
}

func Load() (*Config, error) {
	var cfg Config

	cfg.Rules.AllowExtra = true
	cfg.Defaults = make(map[string]string)

	if viper.ConfigFileUsed() == "" {
		return &cfg, nil
	}

	configFile, err := os.ReadFile(viper.ConfigFileUsed())

	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(configFile, &cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.Schema.Variables == nil {
		cfg.Schema.Variables = make(map[string]schema.Variable)
	}

	slog.Info("Loaded config data", "cfg", cfg)

	return &cfg, nil
}
