package main

import (
	"context"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"utils"

	"cfg"
	"schemas"
	"src"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CreateApiKeyDeps src.Deps

func main() {
	d := CreateApiKeyDeps{
		DbClient:  db.GetDynamoClient(context.Background()),
		TableName: cfg.Config.ApiKeyTable,
	}
	lambda.Start(d.handler)
}

func (d *CreateApiKeyDeps) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse and validate request body
	var req schemas.CreateKeyRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error())), nil
	}

	// Parse and validate request body
	if req.ApiId == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Missing required field: %s", req.ApiId)), nil
	}
	// ... Add more validations as needed

	// Generate random API key
	apiKey, err := utils.GenerateApiKey(req.Prefix)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error generating api key: %s", err.Error())), nil
	}

	// Hash the API key
	hashedKey := utils.HashString(apiKey)

	// Create key ID
	keyId := utils.GenerateKeyId()

	// Create ApiKey struct
	apiKeyToAdd, err := db.CreateApiKeyRow(hashedKey, keyId, req, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error adding key to DynamoDB: %s", err.Error())), nil
	}

	// Create HashedKey struct
	err = db.CreateHashedKeyRow(hashedKey, keyId, req, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error adding key to DynamoDB: %s", err.Error())), nil
	}

	respBody := schemas.CreateKeyResponse{
		ApiId: strings.TrimPrefix(apiKeyToAdd.ApiId, "apiId#"),
		KeyId: strings.TrimPrefix(apiKeyToAdd.KeyId, "apiKeyId#"),
		Key:   apiKey,
	}

	// Marshal the response body
	respBodyJSON, err := json.Marshal(respBody)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling response: %s", err.Error())), nil
	}

	// Return a success response with masked key
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBodyJSON),
	}, nil
}
