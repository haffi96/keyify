package main

import (
	"auth"
	"cfg"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"schemas"
	"src"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CreateApiDeps src.Deps

func main() {
	d := CreateApiDeps{
		DbClient:  db.GetDynamoClient(context.Background()),
		TableName: cfg.Config.ApiKeyTable,
	}
	lambda.Start(d.handler)
}

func (d *CreateApiDeps) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request authentication
	workspaceId, err := auth.VerifyAuthHeader(event, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %s", err.Error())), nil
	}
	log.Printf("workspaceId: %s", workspaceId)

	// Parse and validate request body
	var req schemas.CreateApiRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error())), nil
	}

	// Parse and validate request body
	if req.Name == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, "Missing required field: name"), nil
	}

	// Generate random API Id
	apiId := utils.GenerateRandomId("api_")

	// Create the API
	apiToAdd, err := db.CreateApiRow(workspaceId, apiId, req.Name, d.DbClient)

	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error creating API: %s", err.Error())), nil
	}

	// Return the API Id
	resp := schemas.CreateApiResponse{
		ApiId: apiToAdd.ApiId,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling response: %s", err.Error())), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respJSON),
	}, nil
}
