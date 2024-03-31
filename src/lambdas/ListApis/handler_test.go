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

func TestListApis(t *testing.T) {
	ctx := context.Background()
	d := ListApisDeps{
		DbClient:  db.GetMockDynamoClient(ctx),
		TableName: cfg.Config.ApiKeyTable,
	}

	// Create rootKey row for testing
	rootKey, _ := utils.GenerateApiKey("keyify_")
	hashedKey := utils.HashString(rootKey)
	workspaceId := fmt.Sprintf("workspace-%s", gofakeit.UUID())
	rootKeyReq := schemas.CreateRootKeyRequest{
		WorkspaceId: workspaceId,
	}
	db.CreateRootKeyRow(hashedKey, rootKeyReq, d.DbClient)

	// Create a 3 test APIs
	for i := 0; i < 3; i++ {
		apiId := fmt.Sprintf("api_%s", gofakeit.UUID())
		db.CreateApiRow(workspaceId, apiId, gofakeit.FirstName()+"'s API", d.DbClient)
	}

	// Test the handler
	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Authorization": "Bearer " + rootKey,
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

	assert.Equal(t, 3, keysList, "Expected 3 APIs in the list")
}
