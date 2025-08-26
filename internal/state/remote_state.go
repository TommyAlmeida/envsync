package state

import (
	"fmt"
	"os"
	"path/filepath"

	"maps"

	"github.com/tommyalmeida/envsync/internal/adapter"
	"github.com/tommyalmeida/envsync/internal/config"
	"github.com/tommyalmeida/envsync/internal/env"
)

type RemoteSyncResult struct {
	Added    []string `json:"added"`
	Updated  []string `json:"updated"`
	Deleted  []string `json:"deleted"`
	Skipped  []string `json:"skipped"`
	Source   string   `json:"source"`
	Target   string   `json:"target"`
}

type RemoteState struct {
	config  *config.Config
	adapter adapter.Adapter
}

func NewRemoteState(cfg *config.Config, adp adapter.Adapter) *RemoteState {
	return &RemoteState{
		config:  cfg,
		adapter: adp,
	}
}

func (rs *RemoteState) Pull(prefix string, targetFile string, dryRun bool) (RemoteSyncResult, error) {
	result := RemoteSyncResult{
		Source: fmt.Sprintf("%s:%s", rs.adapter.Name(), prefix),
		Target: targetFile,
	}

	remoteVars, err := rs.adapter.Get(prefix)

	if err != nil {
		return result, fmt.Errorf("failed to get remote variables: %w", err)
	}

	var localVars env.Vars
	if targetFile != "" {
		localVars, err = env.ParseFile(targetFile)

		if err != nil {
			localVars = make(env.Vars)
		}
	} else {
		localVars = make(env.Vars)
	}

	diff := env.CompareEnvs(remoteVars, localVars)

	newLocal := make(env.Vars)
	maps.Copy(newLocal, localVars)

	for _, key := range diff.Missing {
		newLocal[key] = remoteVars[key]
		result.Added = append(result.Added, key)
	}

	for key, diffVal := range diff.Different {
		newLocal[key] = diffVal.Source
		result.Updated = append(result.Updated, key)
	}

	if !dryRun && len(result.Added)+len(result.Updated) > 0 {
		if targetFile != "" {
			if err := rs.atomicWriteFile(targetFile, newLocal); err != nil {
				return result, fmt.Errorf("failed to write target file: %w", err)
			}
		}
	}

	return result, nil
}

func (rs *RemoteState) Push(sourceFile string, prefix string, dryRun bool) (RemoteSyncResult, error) {
	result := RemoteSyncResult{
		Source: sourceFile,
		Target: fmt.Sprintf("%s:%s", rs.adapter.Name(), prefix),
	}

	localVars, err := env.ParseFile(sourceFile)
	if err != nil {
		return result, fmt.Errorf("failed to parse source file: %w", err)
	}

	remoteVars, err := rs.adapter.Get(prefix)
	if err != nil {
		remoteVars = make(env.Vars)
	}

	diff := env.CompareEnvs(localVars, remoteVars)

	varsToSet := make(env.Vars)
	for _, key := range diff.Missing {
		varsToSet[key] = localVars[key]
		result.Added = append(result.Added, key)
	}

	for key, diffVal := range diff.Different {
		varsToSet[key] = diffVal.Source
		result.Updated = append(result.Updated, key)
	}

	if !dryRun {
		if len(varsToSet) > 0 {
			if err := rs.adapter.Set(prefix, varsToSet); err != nil {
				return result, fmt.Errorf("failed to set remote variables: %w", err)
			}
		}
	}

	return result, nil
}

func (rs *RemoteState) Diff(localFile string, prefix string) (env.DiffResult, error) {
	localVars, err := env.ParseFile(localFile)
	if err != nil {
		return env.DiffResult{}, fmt.Errorf("failed to parse local file: %w", err)
	}

	remoteVars, err := rs.adapter.Get(prefix)
	if err != nil {
		return env.DiffResult{}, fmt.Errorf("failed to get remote variables: %w", err)
	}

	return env.CompareEnvs(localVars, remoteVars), nil
}

func (rs *RemoteState) atomicWriteFile(targetFile string, vars env.Vars) error {
	dir := filepath.Dir(targetFile)
	tempFile, err := os.CreateTemp(dir, ".envsync-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	
	tempPath := tempFile.Name()
	tempFile.Close()

	defer func() {
		if _, err := os.Stat(tempPath); err == nil {
			os.Remove(tempPath)
		}
	}()

	if err := vars.WriteToFile(tempPath); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempPath, targetFile); err != nil {
		return fmt.Errorf("failed to move temp file: %w", err)
	}

	return nil
}
