package main

import (
	"context"
	"encoding/json"
	"fmt"
	"schemas"
	"testing"
	"utils"

	"cfg"
	"db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestGetApiKeyHandler(t *testing.T) {
	ctx := context.Background()
	d := GetApiKeyDeps{
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

	// Create a test API key
	req := schemas.CreateKeyRequest{
		ApiId:  "api_" + gofakeit.UUID(),
		Name:   gofakeit.FirstName() + "'s test key",
		Prefix: "test_",
		Roles:  []string{"admin", "user"},
	}
	apiKeyId := utils.GenerateRandomId("key_")
	db.CreateApiKeyRow(utils.HashString("key_1234"), workspaceId, apiKeyId, req, d.DbClient)

	// Test the handler
	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Authorization": "Bearer " + rootKey,
		},
		QueryStringParameters: map[string]string{
			"apiId":    req.ApiId,
			"apiKeyId": apiKeyId,
		},
	})

	// Verify the response
	assert.Equal(t, nil, err, "Handler returned an error")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &result)
	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}

	assert.Equal(t, apiKeyId, result["apiKeyId"], "Expected keyId to be key_1234")
	assert.Equal(t, req.ApiId, result["apiId"], "Expected apiId to be api-1234")
	assert.Equal(t, req.Name, result["name"], "Expected name to be my test key")
	assert.Equal(t, "test_", result["prefix"], "Expected prefix to be test_")
	assert.Equal(t, []interface{}{"admin", "user"}, result["roles"], "Expected roles to be [admin, user]")

}
