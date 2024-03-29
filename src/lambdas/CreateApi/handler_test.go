package main

import (
	"cfg"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"schemas"
	"testing"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestCreateApiHandler(t *testing.T) {
	ctx := context.Background()
	d := CreateApiDeps{
		DbClient:  db.GetMockDynamoClient(ctx),
		TableName: cfg.Config.ApiKeyTable,
	}

	// Create rootKey row for testing
	rootKey, _ := utils.GenerateApiKey("apikeyservice_")
	hashedKey := utils.HashString(rootKey)
	rootKeyReq := schemas.CreateRootKeyRequest{
		WorkspaceId: fmt.Sprintf("workspace-%s", gofakeit.UUID()),
	}
	db.CreateRootKeyRow(hashedKey, rootKeyReq, d.DbClient)

	apiName := "my-test-api"

	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Authorization": "Bearer " + rootKey,
		},
		Body: fmt.Sprintf(`{"name":"%s"}`, apiName),
	})

	assert.Equal(t, nil, err, "Handler returned an error")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &result)
	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}
	if _, ok := result["apiId"]; !ok {
		t.Errorf("Expected 'apiId' in response body, but not found")
	}
}
