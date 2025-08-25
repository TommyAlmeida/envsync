package env

import (
	"testing"

	"github.com/tommyalmeida/envsync/pkg/schema"
)

func TestValidator_Validate(t *testing.T) {
    testSchema := schema.Schema{
        Variables: map[string]schema.Variable{
            "REQUIRED_VAR": {
                Required: true,
                Type:     "string",
            },
            "OPTIONAL_VAR": {
                Required: false,
                Type:     "string",
            },
            "NUMBER_VAR": {
                Required: true,
                Type:     "number",
            },
        },
    }
    
    validator := NewValidator(testSchema)
    
    tests := []struct {
        name     string
        envVars  EnvVars
        expected ValidationResult
    }{
        {
            name: "valid env vars",
            envVars: EnvVars{
                "REQUIRED_VAR": "value",
                "NUMBER_VAR":   "123",
            },
            expected: ValidationResult{
                Valid:   true,
                Errors:  []schema.ValidationError{},
                Missing: nil,
                Extra:   nil,
            },
        },
        {
            name: "missing required variable",
            envVars: EnvVars{
                "NUMBER_VAR": "123",
            },
            expected: ValidationResult{
                Valid:   false,
                Missing: []string{"REQUIRED_VAR"},
            },
        },
        {
            name: "invalid type",
            envVars: EnvVars{
                "REQUIRED_VAR": "value",
                "NUMBER_VAR":   "not-a-number",
            },
            expected: ValidationResult{
                Valid: false,
            },
        },
        {
            name: "extra variables",
            envVars: EnvVars{
                "REQUIRED_VAR": "value",
                "NUMBER_VAR":   "123",
                "EXTRA_VAR":    "extra",
            },
            expected: ValidationResult{
                Valid: true,
                Extra: []string{"EXTRA_VAR"},
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validator.Validate(tt.envVars)
            
            if result.Valid != tt.expected.Valid {
                t.Errorf("expected Valid=%v, got Valid=%v", tt.expected.Valid, result.Valid)
            }
            
            if tt.expected.Missing != nil {
                if len(result.Missing) != len(tt.expected.Missing) {
                    t.Errorf("expected %d missing variables, got %d", len(tt.expected.Missing), len(result.Missing))
                }
            }
            
            if tt.expected.Extra != nil {
                if len(result.Extra) != len(tt.expected.Extra) {
                    t.Errorf("expected %d extra variables, got %d", len(tt.expected.Extra), len(result.Extra))
                }
            }
        })
    }
}