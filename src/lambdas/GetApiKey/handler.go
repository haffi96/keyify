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

func (d *deps) handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse API ID and key ID from request parameters
	apiId := request.QueryStringParameters["apiId"]
	keyId := request.QueryStringParameters["apiKeyId"]

	if apiId == "" || keyId == "" {
		return utils.HttpErrorResponse(http.StatusBadRequest, "Missing required query parameters: apiId and apiKeyId")
	}

	// Construct DynamoDB key
	key := schemas.GetApiKeyInput{
		ApiId: "apiId#" + apiId,
		KeyId: "apiKeyId#" + keyId,
	}

	keyJson, err := attributevalue.MarshalMap(key)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling key: %v", err))
	}

	// Get item from DynamoDB
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("ApiKeyTableDev"), // Replace with your table name
		Key:       keyJson,
	}

	result, err := d.ddbClient.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error getting item from DynamoDB: %v", err))
	}

	// Check if key exists and return 404 if not found
	if result.Item == nil {
		return utils.HttpErrorResponse(http.StatusNotFound, "API key not found")
	}

	// Extract and return relevant data (excluding sensitive fields)
	apiKeyData := schemas.ApiKeyIdRow{}
	err = attributevalue.UnmarshalMap(result.Item, &apiKeyData)
	if err != nil {
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling item: %v", err))
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
		return utils.HttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error marshalling response: %v", err))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBodyJSON),
	}, nil

}
