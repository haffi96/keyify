package main

import (
	"context"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"utils"

	"cfg"
	"schemas"
	"src"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CreateRootKeyDeps src.Deps

func main() {
	d := CreateRootKeyDeps{
		DbClient:  db.GetDynamoClient(context.Background()),
		TableName: cfg.Config.ApiKeyTable,
	}
	lambda.Start(d.handler)
}

func (d *CreateRootKeyDeps) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse and validate request body
	var req schemas.CreateRootKeyRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error())), nil
	}

	// Parse and validate request body
	if req.WorkspaceId == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Missing required field: %s", req.WorkspaceId)), nil
	}
	// ... Add more validations as needed

	// Generate random API key
	rootKey, err := utils.GenerateApiKey("apikeyservice_")
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error generating api key: %s", err.Error())), nil
	}

	// Hash the Root key
	hashedKey := utils.HashString(rootKey)

	// Create ApiKey struct
	err = db.CreateRootKeyRow(hashedKey, req, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error adding key to DynamoDB: %s", err.Error())), nil
	}
	respBody := schemas.CreateRootKeyResponse{
		RootKey: rootKey,
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
