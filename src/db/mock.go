package db

import (
	"context"
	"fmt"

	"cfg"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func DeleteAllDbItems(ctx context.Context, dbClient *dynamodb.Client) {
	// Scan all items
	items, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(cfg.Config.ApiKeyTable),
	})

	if err != nil {
		fmt.Printf("Error scanning table: %v", err)
	}

	// Delete all items
	for _, item := range items.Items {
		_, _ = dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(cfg.Config.ApiKeyTable),
			Key:       map[string]types.AttributeValue{"pk": item["pk"], "sk": item["sk"]},
		})
	}
}
