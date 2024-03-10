package main

import (
	"context"
	"db"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mathRand "math/rand"
	"net/http"
	"strings"
	"utils"

	"schemas"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type deps struct {
	ddbClient *dynamodb.Client
	tableName string
}

func main() {
	d := deps{
		ddbClient: db.GetDynamoClient(context.Background()),
		tableName: "ApiKeyTableDev",
	}
	lambda.Start(d.handler)
}

func (d *deps) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse and validate request body
	var req schemas.CreateKeyRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error()))
	}

	// Parse and validate request body
	if req.ApiId == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Missing required field: %s", req.ApiId))
	}
	// ... Add more validations as needed

	// Generate random API key
	apiKeyBytes, err := utils.GenerateRandomBytes(36)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error generating random bytes: %s", err.Error()))
	}
	apiKey := base64.URLEncoding.EncodeToString(apiKeyBytes)

	// Hash the API key
	hashedKey := utils.HashString(apiKey)

	// Create key ID
	keyId, err := generateKeyId()
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error generating key ID: %s", err.Error()))
	}

	// Create ApiKey struct
	apiKeyToAdd := schemas.ApiKeyIdRow{
		ApiId:     "apiId#" + req.ApiId,
		KeyId:     "apiKeyId#" + keyId,
		HashedKey: hashedKey,
		Name:      req.Name,
		Prefix:    req.Prefix,
		Roles:     req.Roles,
	}
	apiKeyToAddJSON, err := attributevalue.MarshalMap(apiKeyToAdd)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling key to JSON: %s", err.Error()))
	}
	// Put the key in the DynamoDB table
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      apiKeyToAddJSON,
	}

	// Put item for apiKeyId lookup
	_, err = d.ddbClient.PutItem(context.TODO(), putInput)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error adding key to DynamoDB: %s", err.Error()))
	}

	// Create HashedKey struct
	hashedKeyToAdd := schemas.HashedKeyRow{
		ApiId:     "apiId#" + req.ApiId,
		HashedKey: "hashedKey#" + hashedKey,
		KeyId:     keyId,
		Name:      req.Name,
		Prefix:    req.Prefix,
		Roles:     req.Roles,
	}
	hashedKeyToAddJSON, err := attributevalue.MarshalMap(hashedKeyToAdd)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling key to JSON: %s", err.Error()))
	}
	// Put the key in the DynamoDB table
	putHashedKeyInput := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      hashedKeyToAddJSON,
	}
	// Put item for HashedKey lookup
	_, err = d.ddbClient.PutItem(context.TODO(), putHashedKeyInput)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error adding key to DynamoDB: %s", err.Error()))
	}

	respBody := schemas.CreateKeyResponse{
		ApiId: apiKeyToAdd.ApiId,
		KeyId: strings.TrimPrefix(apiKeyToAdd.KeyId, "apiKeyId#"),
		Key:   req.Prefix + apiKey,
	}

	// Marshal the response body
	respBodyJSON, err := json.Marshal(respBody)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling response: %s", err.Error()))
	}

	// Return a success response with masked key
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBodyJSON),
	}, nil
}

func generateKeyId() (string, error) {
	// Prefix to append to the random string
	prefix := "key_"

	// Length of the random string (excluding prefix)
	randomStringLength := 16

	// Characters to use for the random string
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Generate random string directly
	b := make([]byte, randomStringLength)
	for i := range b {
		b[i] = chars[mathRand.Intn(len(chars))] // Select a random character
	}
	randomString := string(b)

	// Combine the prefix and random string
	keyId := prefix + randomString

	return keyId, nil
}
