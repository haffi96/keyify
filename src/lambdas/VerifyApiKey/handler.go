package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"db"
	"schemas"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type deps struct {
	ddbClient *dynamodb.Client
}

func main() {
	d := deps{
		ddbClient: db.GetDynamoClient(context.Background()),
	}
	lambda.Start(d.handler)
}

func (d *deps) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse and validate request body
	var req schemas.VerifyKeyRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return utils.HttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err.Error()))
	}

	if req.ApiId == "" || req.Key == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, "Missing required fields: apiId and key")
	}

	// Hash the provided key
	hashedKey := utils.HashString(req.Key)

	key := schemas.VerifyHashedKeyInput{
		ApiId:     "apiId#" + req.ApiId,
		HashedKey: "hashedKey#" + hashedKey,
	}

	keyJson, err := attributevalue.MarshalMap(key)

	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling key: %s", err.Error()))
	}

	// GetItem operation
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("ApiKeyTableDev"),
		Key:       keyJson,
	}

	result, err := d.ddbClient.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error getting item from DynamoDB: %s", err.Error()))
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
