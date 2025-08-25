package schema_test

import (
	"strings"
	"testing"

	"github.com/tommyalmeida/envsync/pkg/schema"
)

func TestSchema_ValidateVariable(t *testing.T) {
	t.Parallel()

	schema := schema.Schema{
		Variables: map[string]schema.Variable{
			"REQUIRED_VAR": {Required: true, Type: "string"},
			"NUMBER_VAR":   {Required: false, Type: "number"},
			"BOOLEAN_VAR":  {Required: false, Type: "boolean"},
			"EMAIL_VAR":    {Required: false, Type: "email"},
			"URL_VAR":      {Required: false, Type: "url"},
			"PATTERN_VAR":  {Required: false, Type: "string", Pattern: "^[A-Z]+$"},
		},
	}

	tests := []struct {
		name     string
		variable string
		value    string
		wantErr  bool
		errMsg   string
	}{
		{"required variable empty", "REQUIRED_VAR", "", true, "required variable is empty"},
		{"required variable valid", "REQUIRED_VAR", "some value", false, ""},
		{"valid number", "NUMBER_VAR", "123.45", false, ""},
		{"invalid number", "NUMBER_VAR", "not-a-number", true, "not a valid number"},
		{"valid boolean true", "BOOLEAN_VAR", "true", false, ""},
		{"valid boolean 1", "BOOLEAN_VAR", "1", false, ""},
		{"invalid boolean", "BOOLEAN_VAR", "maybe", true, "not a valid boolean"},
		{"valid email", "EMAIL_VAR", "test@example.com", false, ""},
		{"invalid email", "EMAIL_VAR", "invalid-email", true, "not a valid email"},
		{"valid url", "URL_VAR", "https://example.com", false, ""},
		{"invalid url", "URL_VAR", "not-a-url", true, "not a valid URL"},
		{"valid pattern", "PATTERN_VAR", "HELLO", false, ""},
		{"invalid pattern", "PATTERN_VAR", "hello", true, "value does not match pattern"},
		{"unknown variable", "UNKNOWN_VAR", "value", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			errors := schema.ValidateVariable(tt.variable, tt.value)

			if tt.wantErr {
				if len(errors) == 0 {
					t.Errorf("expected validation error, got none")
					return
				}
				if !errorContains(errors, tt.variable, tt.errMsg) {
					t.Errorf("expected error message to contain '%s'", tt.errMsg)
				}
			} else if len(errors) > 0 {
				t.Errorf("expected no validation error, got: %v", errors)
			}
		})
	}
}

func errorContains(errors []schema.ValidationError, variable, substr string) bool {
	for _, err := range errors {
		if err.Variable == variable && strings.Contains(err.Message, substr) {
			return true
		}
	}
	return false
}
