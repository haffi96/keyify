package main

import (
	"cfg"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"schemas"
	"testing"
	"time"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func createApiKeyRows(workspaceId string, apiId string, dbClient *dynamodb.Client) {
	// Create a test API key
	req := schemas.CreateKeyRequest{
		ApiId:  apiId,
		Name:   gofakeit.FirstName() + "'s test key",
		Prefix: "test_",
		Roles:  []string{"admin", "user"},
	}
	apiKeyId := utils.GenerateRandomId("key_")
	hashedKey := utils.HashString(gofakeit.Word())
	db.CreateApiKeyRow(hashedKey, workspaceId, apiKeyId, req, dbClient)
	db.CreateApiKeyDatetimeRow(hashedKey, workspaceId, apiKeyId, req, dbClient)
	db.CreateHashedKeyRow(hashedKey, workspaceId, apiKeyId, req, dbClient)
}

func TestListApiKeysHandler(t *testing.T) {
	ctx := context.Background()
	d := ListApiKeysDeps{
		DbClient:  db.GetMockDynamoClient(ctx),
		TableName: cfg.Config.ApiKeyTable,
	}

	// Create rootKey row for testing
	rootKey, _ := utils.GenerateApiKey("apikeyservice_")
	hashedKey := utils.HashString(rootKey)
	workspaceId := fmt.Sprintf("workspace-%s", gofakeit.UUID())
	rootKeyReq := schemas.CreateRootKeyRequest{
		WorkspaceId: workspaceId,
	}
	db.CreateRootKeyRow(hashedKey, rootKeyReq, d.DbClient)

	// Create a 3 test API keys
	apiId := fmt.Sprintf("api_%s", gofakeit.UUID())
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 100)
		createApiKeyRows(workspaceId, apiId, d.DbClient)
	}

	// Test the handler
	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Authorization": "Bearer " + rootKey,
		},
		QueryStringParameters: map[string]string{
			"apiId": apiId,
		},
	})

	// Verify the response
	assert.Equal(t, nil, err, "Handler returned an error")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &result)

	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}

	keysList := len(result)

	assert.Equal(t, 3, keysList, "Expected 3 keys in the list")
}
