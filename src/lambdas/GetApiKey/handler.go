package main

import (
	"auth"
	"cfg"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

type GetApiKeyDeps src.Deps

func main() {
	d := GetApiKeyDeps{
		DbClient:  db.GetDynamoClient(context.Background()),
		TableName: cfg.Config.ApiKeyTable,
	}
	lambda.Start(d.handler)
}

func (d *GetApiKeyDeps) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request authentication
	workspaceId, err := auth.VerifyAuthHeader(request, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %s", err.Error())), nil
	}
	log.Printf("workspaceId: %s", workspaceId)

	// Parse API ID and key ID from request parameters
	apiId := request.QueryStringParameters["apiId"]
	keyId := request.QueryStringParameters["apiKeyId"]

	if apiId == "" || keyId == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, "Missing required query parameters: apiId and apiKeyId"), nil
	}

	// Construct DynamoDB key
	key := schemas.GetApiKeyInput{
		WorkspaceIdApiId: "workspaceId#" + workspaceId + "-" + "apiId#" + apiId,
		KeyId:            "apiKeyId#" + keyId,
	}

	keyJson, err := attributevalue.MarshalMap(key)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling key: %v", err)), nil
	}

	// Get item from DynamoDB
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
		Key:       keyJson,
	}

	result, err := d.DbClient.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error getting item from DynamoDB: %v", err)), nil
	}

	// Check if key exists and return 404 if not found
	if result.Item == nil {
		return utils.HttpErrorResponse(http.StatusNotFound, "API key not found"), nil
	}

	// Extract and return relevant data (excluding sensitive fields)
	apiKeyData := schemas.ApiKeyIdRow{}
	err = attributevalue.UnmarshalMap(result.Item, &apiKeyData)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling item: %v", err)), nil
	}

	respBody := schemas.GetApiKeyResponse{
		ApiId:  apiId,
		KeyId:  keyId,
		Name:   apiKeyData.Name,
		Prefix: apiKeyData.Prefix,
		Roles:  apiKeyData.Roles,
	}
	respBodyJSON, err := json.Marshal(respBody)

	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling response: %v", err)), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBodyJSON),
	}, nil

}
