package env

import (
	"reflect"
	"testing"
)

func TestCompareEnvs(t *testing.T) {
    tests := []struct {
        name     string
        source   EnvVars
        target   EnvVars
        expected DiffResult
    }{
        {
            name: "identical environments",
            source: EnvVars{
                "VAR1": "value1",
                "VAR2": "value2",
            },
            target: EnvVars{
                "VAR1": "value1",
                "VAR2": "value2",
            },
            expected: DiffResult{
                Missing:   nil,
                Extra:     nil,
                Different: map[string]Diff{},
                Same:      []string{"VAR1", "VAR2"},
            },
        },
        {
            name: "missing variables",
            source: EnvVars{
                "VAR1": "value1",
                "VAR2": "value2",
            },
            target: EnvVars{
                "VAR1": "value1",
            },
            expected: DiffResult{
                Missing:   []string{"VAR2"},
                Extra:     nil,
                Different: map[string]Diff{},
                Same:      []string{"VAR1"},
            },
        },
        {
            name: "extra variables",
            source: EnvVars{
                "VAR1": "value1",
            },
            target: EnvVars{
                "VAR1": "value1",
                "VAR2": "value2",
            },
            expected: DiffResult{
                Missing:   nil,
                Extra:     []string{"VAR2"},
                Different: map[string]Diff{},
                Same:      []string{"VAR1"},
            },
        },
        {
            name: "different values",
            source: EnvVars{
                "VAR1": "value1",
                "VAR2": "old_value",
            },
            target: EnvVars{
                "VAR1": "value1",
                "VAR2": "new_value",
            },
            expected: DiffResult{
                Missing: nil,
                Extra:   nil,
                Different: map[string]Diff{
                    "VAR2": {Source: "old_value", Target: "new_value"},
                },
                Same: []string{"VAR1"},
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CompareEnvs(tt.source, tt.target)
            
            if !reflect.DeepEqual(result.Missing, tt.expected.Missing) {
                t.Errorf("expected Missing=%v, got Missing=%v", tt.expected.Missing, result.Missing)
            }
            
            if !reflect.DeepEqual(result.Extra, tt.expected.Extra) {
                t.Errorf("expected Extra=%v, got Extra=%v", tt.expected.Extra, result.Extra)
            }
            
            if !reflect.DeepEqual(result.Different, tt.expected.Different) {
                t.Errorf("expected Different=%v, got Different=%v", tt.expected.Different, result.Different)
            }
            
            if len(result.Same) != len(tt.expected.Same) {
                t.Errorf("expected Same length=%d, got length=%d", len(tt.expected.Same), len(result.Same))
            }
        })
    }
}