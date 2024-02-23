package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ApiKey struct {
	ApiKeyId string `json:"api_key_id" dynamodbav:"api_key_id"`
}

func HandleRequest(ctx context.Context) (string, error) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-2"),
	)
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	key := ApiKey{
		ApiKeyId: "5cb67f61-4577-4f22-9fc2-ed91486dced5",
	}

	k, err := attributevalue.MarshalMap(key)

	if err != nil {
		return "", fmt.Errorf("failed to marshal Record, %v", err)
	}

	// Build the request with its input parameters
	resp, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("ApiKeyTableDev"),
		Key:       k,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get item: %v", err)
	}

	return fmt.Sprintf("Item: %v", resp.Item), nil
}

func main() {
	lambda.Start(HandleRequest)
}
