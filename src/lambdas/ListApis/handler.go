package main

import (
	"auth"
	"cfg"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"schemas"
	"src"
	"strings"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ListApisDeps src.Deps

func main() {
	d := ListApisDeps{
		DbClient:  db.GetDynamoClient(context.Background()),
		TableName: cfg.Config.ApiKeyTable,
	}
	lambda.Start(d.handler)
}

func (d *ListApisDeps) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request authentication
	workspaceId, err := auth.VerifyAuthHeader(event, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %s", err.Error())), nil
	}

	// Query the database for all APIs for the workspace
	apis, err := db.ListApis(workspaceId, d.DbClient)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error listing APIs: %s", err)), nil
	}

	// Return the API Id
	respBody := []schemas.GetApiResponse{}
	for _, api := range apis {
		respBody = append(respBody, schemas.GetApiResponse{
			ApiId:     strings.Split(api.ApiId, "#")[1],
			ApiName:   api.ApiName,
			CreatedAt: api.CreatedAt,
		})
	}

	// Construct the response
	respBodyJson, err := json.Marshal(respBody)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling response: %s", err)), nil
	}

	if len(respBody) == 0 {
		return utils.HttpErrorResponse(http.StatusNotFound, "No APIs found"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBodyJson),
	}, nil
}
