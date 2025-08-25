package schema

import (
	"testing"
)

func TestSchema_ValidateVariable(t *testing.T) {
    schema := Schema{
        Variables: map[string]Variable{
            "REQUIRED_VAR": {
                Required: true,
                Type:     "string",
            },
            "NUMBER_VAR": {
                Required: false,
                Type:     "number",
            },
            "BOOLEAN_VAR": {
                Required: false,
                Type:     "boolean",
            },
            "EMAIL_VAR": {
                Required: false,
                Type:     "email",
            },
            "URL_VAR": {
                Required: false,
                Type:     "url",
            },
            "PATTERN_VAR": {
                Required: false,
                Type:     "string",
                Pattern:  "^[A-Z]+$",
            },
        },
    }

    tests := []struct {
        name     string
        variable string
        value    string
        wantErr  bool
        errMsg   string
    }{
        {
            name:     "required variable empty",
            variable: "REQUIRED_VAR",
            value:    "",
            wantErr:  true,
            errMsg:   "required variable is empty",
        },
        {
            name:     "required variable valid",
            variable: "REQUIRED_VAR",
            value:    "some value",
            wantErr:  false,
        },
        {
            name:     "valid number",
            variable: "NUMBER_VAR",
            value:    "123.45",
            wantErr:  false,
        },
        {
            name:     "invalid number",
            variable: "NUMBER_VAR",
            value:    "not-a-number",
            wantErr:  true,
            errMsg:   "not a valid number",
        },
        {
            name:     "valid boolean true",
            variable: "BOOLEAN_VAR",
            value:    "true",
            wantErr:  false,
        },
        {
            name:     "valid boolean 1",
            variable: "BOOLEAN_VAR",
            value:    "1",
            wantErr:  false,
        },
        {
            name:     "invalid boolean",
            variable: "BOOLEAN_VAR",
            value:    "maybe",
            wantErr:  true,
            errMsg:   "not a valid boolean",
        },
        {
            name:     "valid email",
            variable: "EMAIL_VAR",
            value:    "test@example.com",
            wantErr:  false,
        },
        {
            name:     "invalid email",
            variable: "EMAIL_VAR",
            value:    "invalid-email",
            wantErr:  true,
            errMsg:   "not a valid email",
        },
        {
            name:     "valid url",
            variable: "URL_VAR",
            value:    "https://example.com",
            wantErr:  false,
        },
        {
            name:     "invalid url",
            variable: "URL_VAR",
            value:    "not-a-url",
            wantErr:  true,
            errMsg:   "not a valid URL",
        },
        {
            name:     "valid pattern",
            variable: "PATTERN_VAR",
            value:    "HELLO",
            wantErr:  false,
        },
        {
            name:     "invalid pattern",
            variable: "PATTERN_VAR",
            value:    "hello",
            wantErr:  true,
            errMsg:   "value does not match pattern",
        },
        {
            name:     "unknown variable",
            variable: "UNKNOWN_VAR",
            value:    "value",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            errors := schema.ValidateVariable(tt.variable, tt.value)
            
            if tt.wantErr && len(errors) == 0 {
                t.Errorf("expected validation error, got none")
                return
            }
            
            if !tt.wantErr && len(errors) > 0 {
                t.Errorf("expected no validation error, got: %v", errors)
                return
            }
            
            if tt.wantErr && len(errors) > 0 {
                found := false
                for _, err := range errors {
                    if err.Variable == tt.variable {
                        found = true
                        if tt.errMsg != "" && !contains(err.Message, tt.errMsg) {
                            t.Errorf("expected error message to contain '%s', got '%s'", tt.errMsg, err.Message)
                        }
                        break
                    }
                }
                if !found {
                    t.Errorf("expected validation error for variable %s, but not found", tt.variable)
                }
            }
        })
    }
}

func contains(s, substr string) bool {
    if substr == "" {
        return true
    }

    if len(substr) > len(s) {
        return false
    }

    if s[:len(substr)] == substr {
        return true
    }

    if s[len(s)-len(substr):] == substr {
        return true
    }

    return containsMiddle(s, substr)
}

func containsMiddle(s, substr string) bool {
    for i := 1; i < len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}