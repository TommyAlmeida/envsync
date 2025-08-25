package env

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

type EnvVars map[string]string

func ParseFile(filename string) (EnvVars, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("expected a file but got a directory: %s", filename)
	}

	vars, err := godotenv.Read(filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read env file %s: %w", filename, err)
	}

	return EnvVars(vars), nil
}

func (e EnvVars) Keys() []string {
	if e == nil {
		return nil
	}

	keys := make([]string, 0, len(e))

	for k := range e {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (e EnvVars) WriteToFile(filename string) error {
	if e == nil {
		return fmt.Errorf("EnvVars is nil")
	}

	var lines []string
	for _, key := range e.Keys() {
		if strings.TrimSpace(key) == "" {
			continue
		}

		value := e[key]

		if strings.ContainsAny(value, " \t\n\r") || strings.Contains(value, `"`) {
			value = `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
		}

		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	content := strings.Join(lines, "\n") + "\n"

	return os.WriteFile(filename, []byte(content), 0600)
}