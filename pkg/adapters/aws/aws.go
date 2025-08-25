package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/tommyalmeida/envsync/internal/adapters"
	cfg "github.com/tommyalmeida/envsync/internal/config"
)

type AWSAdapter struct {
    client *ssm.Client
}

func New(adapterConfig map[string]string) adapters.Adapter {
    region, ok := adapterConfig["region"]

    if !ok {
        log.Fatal("AWS adapter requires 'region' in config")
    }

    cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))

    if err != nil {
        log.Fatalf("Failed to load AWS config: %v", err)
    }

    return &AWSAdapter{client: ssm.NewFromConfig(cfg)}
}

func (a *AWSAdapter) Sync(config *cfg.Config) error {
    log.Printf("Syncing to AWS with config: %+v", config)
    return nil
}

func (a *AWSAdapter) Validate(config *cfg.Config) error {
    return nil
}

func (a *AWSAdapter) Name() string {
    return "aws"
}

func init() {
    adapters.Register("aws", New)
}