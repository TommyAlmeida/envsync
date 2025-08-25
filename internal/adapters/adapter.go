package adapters

import "github.com/tommyalmeida/envsync/internal/config"



type Adapter interface {
    Sync(config *config.Config) error         // Sync environment variables to the provider
    Validate(config *config.Config) error     // Validate configuration for the provider
    Name() string                             // Return the adapter name (e.g., "aws", "gcp")
}