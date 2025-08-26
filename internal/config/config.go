package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/tommyalmeida/envsync/pkg/schema"
)

type Config struct {
	Schema   schema.Schema     `yaml:"schema"`
	Defaults map[string]string `yaml:"defaults"`
	Rules    Rules             `yaml:"rules"`
	Adapter  AdapterConfig            `yaml:"adapter"`
}

type Rules struct {
	RequireAll     bool     `yaml:"require_all"`
	AllowExtra     bool     `yaml:"allow_extra"`
	IgnorePatterns []string `yaml:"ignore_patterns"`
}

type AdapterConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]any `yaml:"config"`
}

func Load() (*Config, error) {
	var cfg Config

	cfg.Rules.AllowExtra = true
	cfg.Defaults = make(map[string]string)
	cfg.Adapter.Config = make(map[string]any)

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

	return &cfg, nil
}