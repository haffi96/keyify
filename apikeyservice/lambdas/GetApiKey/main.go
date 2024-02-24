package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type apiKeyIdRow struct {
	ApiId     string   `json:"apiId" dynamodbav:"pk"`
	KeyId     string   `json:"apiKeyId" dynamodbav:"sk"`
	HashedKey string   `json:"-" dynamodbav:"hashedKey"` // Store hashed key securely
	ApiKey    string   `json:"key" dynamodbav:"-"`       // Not used, store hashed key instead
	Name      string   `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Prefix    string   `json:"prefix,omitempty" dynamodbav:"prefix,omitempty"`
	Roles     []string `json:"roles,omitempty" dynamodbav:"roles,omitempty"`
}

type GetApiKeyResponse struct {
	ApiId  string   `json:"apiId"`
	KeyId  string   `json:"apiKeyId"`
	Name   string   `json:"name,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse API ID and key ID from request parameters
	apiId := request.QueryStringParameters["apiId"]
	keyId := request.QueryStringParameters["apiKeyId"]

	if apiId == "" || keyId == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing required query parameters: apiId and apiKeyId",
		}, nil
	}

	// Configure AWS SDK client
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-2"), // Replace with your region
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error loading SDK config: %v", err),
		}, nil
	}

	svc := dynamodb.NewFromConfig(cfg)

	// Construct DynamoDB key
	key := apiKeyIdRow{
		ApiId: "apiId#" + apiId,
		KeyId: "apiKeyId#" + keyId,
	}

	keyJson, err := attributevalue.MarshalMap(key)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error marshalling key: %v", err),
		}, nil
	}

	// Get item from DynamoDB
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("ApiKeyTableDev"), // Replace with your table name
		Key:       keyJson,
	}

	result, err := svc.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error getting item from DynamoDB: %v", err),
		}, nil
	}

	// Check if key exists and return 404 if not found
	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "API key not found",
		}, nil
	}

	// Extract and return relevant data (excluding sensitive fields)
	apiKeyData := apiKeyIdRow{}
	err = attributevalue.UnmarshalMap(result.Item, &apiKeyData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	respBody := GetApiKeyResponse{
		ApiId:  apiId,
		KeyId:  keyId,
		Name:   apiKeyData.Name,
		Prefix: apiKeyData.Prefix,
		Roles:  apiKeyData.Roles,
	}
	respBodyJSON, err := attributevalue.MarshalMap(respBody)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error marshalling response: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("%v", respBodyJSON),
	}, nil

}

func main() {
	lambda.Start(handler)
}
