package db

import (
	"context"

	"cfg"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetMockDynamoClient(ctx context.Context) *dynamodb.Client {
	// Configure AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Config.AwsRegion),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: "http://localhost:4566"}, nil
				})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		panic("Failed to load aws sdk config, " + err.Error())
	}

	client := dynamodb.NewFromConfig(cfg)

	return client
}
