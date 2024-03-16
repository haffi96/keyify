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
	d := VerifyKeyDeps{
		DbClient:  db.GetMockDynamoClient(ctx),
		TableName: cfg.Config.ApiKeyTable,
	}

	// Create a test API hasked key row
	req := schemas.CreateKeyRequest{
		ApiId:  "api" + gofakeit.UUID(),
		Name:   gofakeit.FirstName() + "'s test key",
		Prefix: "test_",
		Roles:  []string{"admin", "user"},
	}
	apiKeyId := utils.GenerateKeyId()
	apiKey, _ := utils.GenerateApiKey(req.Prefix)
	db.CreateHashedKeyRow(utils.HashString(apiKey), apiKeyId, req, d.DbClient)

	// Test the handler
	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Body: fmt.Sprintf(`{"apiId": "%s", "key": "%s"}`, req.ApiId, apiKey),
	})

	// Verify the response
	assert.Equal(t, nil, err, "Handler returned an error")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &result)
	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}

	assert.Equal(t, true, result["isValidKey"], "Expected api key to be valid")
}
