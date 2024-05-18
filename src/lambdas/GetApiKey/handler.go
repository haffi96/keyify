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

	apiKeyData, err := db.GetApiKey(workspaceId, apiId, keyId, d.DbClient)

	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error getting API key: %v", err)), nil
	}

	respBody := schemas.GetApiKeyResponse{
		ApiId:     apiId,
		KeyId:     keyId,
		Name:      apiKeyData.Name,
		Prefix:    apiKeyData.Prefix,
		Roles:     apiKeyData.Roles,
		CreatedAt: apiKeyData.CreatedAt,
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
