package env

import (
	"fmt"
	"log"
	"sort"

	"github.com/tommyalmeida/envsync/pkg/schema"
)

type ValidationResult struct {
	Valid   bool                     `json:"valid"`
	Errors  []schema.ValidationError `json:"errors,omitempty"`
	Missing []string                 `json:"missing,omitempty"`
	Extra   []string                 `json:"extra,omitempty"`
}

type Validator struct {
	schema schema.Schema
	debug  bool
}

func NewValidator(s schema.Schema) *Validator {
	return &Validator{schema: s, debug: false}
}

func (v *Validator) Validate(envVars Vars) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Errors: []schema.ValidationError{},
	}

	if v.debug {
		log.Printf("DEBUG: Schema variables: %v\n", v.getSchemaKeys())
		log.Printf("DEBUG: Env variables: %v\n", envVars.Keys())
	}

	for name, variable := range v.schema.Variables {
		value, exists := envVars[name]
		if !exists {
			if variable.Required {
				result.Missing = append(result.Missing, name)
				result.Valid = false
			}
			continue
		}

		if errors := v.schema.ValidateVariable(name, value); len(errors) > 0 {
			result.Errors = append(result.Errors, errors...)
			result.Valid = false
		}
	}

	for name := range envVars {
		if _, exists := v.schema.Variables[name]; !exists {
			result.Extra = append(result.Extra, name)
		}
	}

	sort.Strings(result.Missing)
	sort.Strings(result.Extra)

	return result
}

func (v *Validator) getSchemaKeys() []string {
	keys := make([]string, 0, len(v.schema.Variables))

	for k := range v.schema.Variables {
		keys = append(keys, fmt.Sprintf("'%s'", k))
	}

	sort.Strings(keys)
	return keys
}
