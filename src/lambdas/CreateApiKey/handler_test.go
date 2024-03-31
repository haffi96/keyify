package main

import (
	"context"
	"encoding/json"
	"fmt"
	"schemas"
	"strings"
	"testing"
	"utils"

	"cfg"
	"db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestCreateApiKeyHandler(t *testing.T) {
	ctx := context.Background()
	d := CreateApiKeyDeps{
		DbClient:  db.GetMockDynamoClient(ctx),
		TableName: cfg.Config.ApiKeyTable,
	}

	// Create rootKey row for testing
	rootKey, _ := utils.GenerateApiKey("keyify_")
	hashedKey := utils.HashString(rootKey)
	rootKeyReq := schemas.CreateRootKeyRequest{
		WorkspaceId: fmt.Sprintf("workspace-%s", gofakeit.UUID()),
	}
	db.CreateRootKeyRow(hashedKey, rootKeyReq, d.DbClient)

	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Authorization": "Bearer " + rootKey,
		},
		Body: `{"apiId":"api-1","name":"key-1","prefix":"key_","roles":["role-1","role-2"]}`,
	})

	assert.Equal(t, nil, err, "Handler returned an error")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &result)
	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}
	assert.Equal(t, "api-1", result["apiId"], "Expected apiId to be api-1")

	// Assert prefix of the key
	key := result["key"].(string)
	keyParts := strings.SplitN(key, "_", 2)
	keyPrefix := keyParts[0] + "_"
	assert.Equal(t, "key_", keyPrefix, "Expected key prefix to be 'key_'")

	if _, ok := result["keyId"]; !ok {
		t.Errorf("Expected 'keyId' in response body, but not found")
	}
	if _, ok := result["key"]; !ok {
		t.Errorf("Expected 'key' in response body, but not found")
	}
}
