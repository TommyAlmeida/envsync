package adapters

import (
	"fmt"
	"log"
)

type AdapterFactory func(config map[string]string) Adapter

var registry = make(map[string]AdapterFactory)

func Register(name string, factory AdapterFactory) {
    if _, exists := registry[name]; exists {
        log.Printf("Warning: Overwriting adapter registration for %s", name)
    }

    registry[name] = factory
}

func GetAdapter(name string, config map[string]string) (Adapter, error) {
    factory, ok := registry[name]

    if !ok {
        return nil, fmt.Errorf("adapter %s not registered", name)
    }
		
    return factory(config), nil
}