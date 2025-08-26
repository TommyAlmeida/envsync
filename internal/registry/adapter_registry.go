package registry

import (
	"fmt"

	"github.com/tommyalmeida/envsync/internal/adapter"
	"github.com/tommyalmeida/envsync/internal/adapter/aws"
)

// Due to circular dependecies the registry of the default adapters its on this package

var registry = make(map[string]adapter.AdapterFactory)

func init() {
	Register("aws-ssm", aws.NewSSMAdapter)
}

func Register(name string, factory adapter.AdapterFactory) {
	registry[name] = factory
}

func Create(name string, config adapter.Config) (adapter.Adapter, error) {
	factory, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("unknown adapter: %s", name)
	}

	return factory(config)
}

func ListAvailable() []string {
	var adapters []string
	for name := range registry {
		adapters = append(adapters, name)
	}
	return adapters
}

