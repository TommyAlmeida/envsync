package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/tommyalmeida/envsync/internal/adapter"
	"github.com/tommyalmeida/envsync/internal/env"
)

type SSMAdapter struct {
	client *ssm.Client
	region string
}

func NewSSMAdapter(adapterConfig adapter.Config) (adapter.Adapter, error) {
	region, ok := adapterConfig["region"].(string)

	if !ok || region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &SSMAdapter{
		client: ssm.NewFromConfig(cfg),
		region: region,
	}, nil
}

func (s *SSMAdapter) Name() string {
	return "aws-ssm"
}

func (s *SSMAdapter) Get(prefix string) (env.Vars, error) {
	ctx := context.TODO()
	vars := make(env.Vars)

	input := &ssm.GetParametersByPathInput{
		Path:      aws.String(s.normalizePath(prefix)),
		Recursive: aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}

	paginator := ssm.NewGetParametersByPathPaginator(s.client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get parameters: %w", err)
		}

		for _, param := range output.Parameters {
			if param.Name != nil && param.Value != nil {
				key := s.extractKeyFromPath(*param.Name, prefix)
				vars[key] = *param.Value
			}
		}
	}

	return vars, nil
}

func (s *SSMAdapter) Set(prefix string, vars env.Vars) error {
	ctx := context.TODO()

	for key, value := range vars {
		paramName := s.buildParameterPath(prefix, key)
		
		input := &ssm.PutParameterInput{
			Name:      aws.String(paramName),
			Value:     aws.String(value),
			Type:      types.ParameterTypeString,
			Overwrite: aws.Bool(true),
		}

		_, err := s.client.PutParameter(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to set parameter %s: %w", paramName, err)
		}
	}

	return nil
}

func (s *SSMAdapter) Delete(prefix string, keys []string) error {
	ctx := context.TODO()

	var paramNames []string
	for _, key := range keys {
		paramNames = append(paramNames, s.buildParameterPath(prefix, key))
	}

	if len(paramNames) == 0 {
		return nil
	}

	input := &ssm.DeleteParametersInput{
		Names: paramNames,
	}

	_, err := s.client.DeleteParameters(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete parameters: %w", err)
	}

	return nil
}

func (s *SSMAdapter) List(prefix string) ([]string, error) {
	ctx := context.TODO()
	var keys []string

	input := &ssm.GetParametersByPathInput{
		Path:      aws.String(s.normalizePath(prefix)),
		Recursive: aws.Bool(true),
	}

	paginator := ssm.NewGetParametersByPathPaginator(s.client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list parameters: %w", err)
		}

		for _, param := range output.Parameters {
			if param.Name != nil {
				key := s.extractKeyFromPath(*param.Name, prefix)
				keys = append(keys, key)
			}
		}
	}

	return keys, nil
}

func (s *SSMAdapter) normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}

func (s *SSMAdapter) buildParameterPath(prefix, key string) string {
	normalizedPrefix := s.normalizePath(prefix)
	return normalizedPrefix + key
}

func (s *SSMAdapter) extractKeyFromPath(fullPath, prefix string) string {
	normalizedPrefix := s.normalizePath(prefix)
	if strings.HasPrefix(fullPath, normalizedPrefix) {
		return strings.TrimPrefix(fullPath, normalizedPrefix)
	}
	return fullPath
}