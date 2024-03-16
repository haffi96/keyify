package main

import (
	"auth"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cfg"
	"db"
	"schemas"
	"src"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type VerifyKeyDeps src.Deps

func main() {
	d := VerifyKeyDeps{
		DbClient:  db.GetDynamoClient(context.Background()),
		TableName: cfg.Config.ApiKeyTable,
	}
	lambda.Start(d.handler)
}

func (d *VerifyKeyDeps) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request authentication
	workspaceId, err := auth.VerifyAuthHeader(event, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %s", err.Error())), nil
	}
	log.Printf("workspaceId: %s", workspaceId)

	// Parse and validate request body
	var req schemas.VerifyKeyRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error())), nil
	}

	if req.ApiId == "" || req.Key == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, "Missing required fields: apiId and key"), nil
	}

	// Hash the provided key
	hashedKey := utils.HashString(req.Key)

	key := schemas.VerifyHashedKeyInput{
		WorkspaceIdApiId: "workspaceId#" + workspaceId + "-" + "apiId#" + req.ApiId,
		HashedKey:        "hashedKey#" + hashedKey,
	}

	keyJson, err := attributevalue.MarshalMap(key)

	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling key: %s", err.Error())), nil
	}

	// GetItem operation
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
		Key:       keyJson,
	}

	result, err := d.DbClient.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error getting item from DynamoDB: %s", err.Error())), nil
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
