package adapter

import "github.com/tommyalmeida/envsync/internal/env"

type Adapter interface {
	Name() string
	Get(prefix string) (env.Vars, error)
	Set(prefix string, vars env.Vars) error
	Delete(prefix string, keys []string) error
	List(prefix string) ([]string, error)
}

type Config map[string]interface{}

type AdapterFactory func(config Config) (Adapter, error)