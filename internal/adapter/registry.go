package adapter

import (
	"fmt"
)

var registry = make(map[string]AdapterFactory)

func Register(name string, factory AdapterFactory) {
	registry[name] = factory
}

func Create(name string, config Config) (Adapter, error) {
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