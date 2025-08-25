package schema

import (
	"fmt"
	"regexp"
)

type Schema struct {
	Variables map[string]Variable `yaml:"variables"`
}

type Variable struct {
	Required    bool   `yaml:"required"`
	Type        string `yaml:"type"` // string, number, boolean, url, email
	Pattern     string `yaml:"pattern"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
}

type ValidationError struct {
	Variable string `json:"variable"`
	Message  string `json:"message"`
}

func (s Schema) ValidateVariable(name, value string) []ValidationError {
	var errors []ValidationError

	variable, exists := s.Variables[name]

	if !exists {
		return errors
	}

	if variable.Required && value == "" {
		errors = append(errors, ValidationError{
			Variable: name,
			Message:  "required variable is empty",
		})

		return errors
	}

	if value == "" {
		return errors
	}

	if err := s.validateType(variable.Type, value); err != nil {
		errors = append(errors, ValidationError{
			Variable: name,
			Message:  fmt.Sprintf("type validation failed: %v", err),
		})
	}

	if variable.Pattern != "" {
		if matched, err := regexp.MatchString(variable.Pattern, value); err != nil {
			errors = append(errors, ValidationError{
				Variable: name,
				Message:  fmt.Sprintf("pattern validation error: %v", err),
			})
		} else if !matched {
			errors = append(errors, ValidationError{
				Variable: name,
				Message:  fmt.Sprintf("value does not match pattern: %s", variable.Pattern),
			})
		}
	}

	return errors
}

func (s Schema) validateType(varType, value string) error {
	switch varType {
	case "string", "":
		return nil
	case "number":
		if matched, _ := regexp.MatchString(`^\d+(\.\d+)?$`, value); !matched {
			return fmt.Errorf("not a valid number")
		}
	case "boolean":
		if matched, _ := regexp.MatchString(`^(true|false|1|0|yes|no|on|off)$`, value); !matched {
			return fmt.Errorf("not a valid boolean")
		}
	case "url":
		if matched, _ := regexp.MatchString(`^https?://`, value); !matched {
			return fmt.Errorf("not a valid URL")
		}
	case "email":
		if matched, _ := regexp.MatchString(`^[^\s@]+@[^\s@]+\.[^\s@]+$`, value); !matched {
			return fmt.Errorf("not a valid email")
		}
	default:
		return fmt.Errorf("unknown type: %s", varType)
	}

	return nil
}
