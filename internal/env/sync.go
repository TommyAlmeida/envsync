package env

import (
	"fmt"

	"maps"

	"github.com/tommyalmeida/envsync/internal/config"
)

type SyncResult struct {
    Added    []string `json:"added"`
    Skipped  []string `json:"skipped"`
    FilePath string   `json:"file_path"`
}

type Syncer struct {
    config *config.Config
}

func NewSyncer(cfg *config.Config) *Syncer {
    return &Syncer{config: cfg}
}

func (s *Syncer) Sync(source, target EnvVars, targetFile string, dryRun bool) (SyncResult, error) {
    result := SyncResult{
        FilePath: targetFile,
    }
    
    diff := CompareEnvs(source, target)
    
    newTarget := make(EnvVars)
    maps.Copy(newTarget, target)
    
    for _, key := range diff.Missing {
        sourceValue := source[key]
        defaultValue := s.getDefaultValue(key, sourceValue)
        
        newTarget[key] = defaultValue
        result.Added = append(result.Added, key)
    }
    
    if !dryRun && len(result.Added) > 0 {
        if err := newTarget.WriteToFile(targetFile); err != nil {
            return result, fmt.Errorf("failed to write target file: %w", err)
        }
    }
    
    return result, nil
}

func (s *Syncer) getDefaultValue(key, originalValue string) string {
    if defaultVal, exists := s.config.Defaults[key]; exists {
        return defaultVal
    }

    if variable, exists := s.config.Schema.Variables[key]; exists && variable.Default != "" {
        return variable.Default
    }

    return originalValue
}