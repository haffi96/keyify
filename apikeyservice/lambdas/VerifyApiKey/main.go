package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type HashedKeyRow struct {
	ApiId     string `json:"apiId" dynamodbav:"pk"`
	HashedKey string `json:"-" dynamodbav:"sk"` // Store hashed key securely
}

type VerifyKeyRequest struct {
	ApiId string `json:"apiId"`
	Key   string `json:"key"`
}

type VerifyKeyResponse struct {
	Valid bool `json:"isValidKey"`
}

func main() {
	lambda.Start(handler)
}

func errorResponse(statusCode int, message string, args ...interface{}) (events.APIGatewayProxyResponse, error) {
	body := fmt.Sprintf(message, args...)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}, nil
}

func hashKey(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-2"),
	)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error loading SDK config: %s", err.Error())
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Parse and validate request body
	var req VerifyKeyRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return errorResponse(http.StatusBadRequest, "Invalid request body: %s", err.Error())
	}

	if req.ApiId == "" || req.Key == "" {
		return errorResponse(http.StatusBadRequest, "Missing required fields: apiId and key")
	}

	// Hash the provided key
	hashedKey := hashKey(req.Key)

	key := HashedKeyRow{
		ApiId:     "apiId#" + req.ApiId,
		HashedKey: "hashedKey#" + hashedKey,
	}

	keyJson, err := attributevalue.MarshalMap(key)

	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error marshalling key: %s", err.Error())
	}

	// GetItem operation
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("ApiKeyTableDev"),
		Key:       keyJson,
	}

	result, err := svc.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error getting item from DynamoDB: %s", err.Error())
	}

	// Check if the key exists and has valid data
	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"isValidKey": false}`,
		}, nil
	}

	// Additional checks based on your requirements (e.g., key expiration, roles)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"isValidKey": true}`,
	}, nil
}
