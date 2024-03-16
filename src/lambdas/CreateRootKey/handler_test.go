package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"cfg"
	"db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestCreateRootKeyHandler(t *testing.T) {
	ctx := context.Background()
	d := CreateRootKeyDeps{
		DbClient:  db.GetMockDynamoClient(ctx),
		TableName: cfg.Config.ApiKeyTable,
	}

	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Body: `{"workspaceId":"tenant_123","name":"rooyKey-1","permissions":["api.api-123.read_api","api.api-123.write_api"]}`,
	})

	assert.Equal(t, nil, err, "Handler returned an error")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &result)
	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}

	// Assert prefix of the key
	key := result["rootKey"].(string)
	keyParts := strings.SplitN(key, "_", 2)
	keyPrefix := keyParts[0] + "_"
	assert.Equal(t, "apikeyservice_", keyPrefix, "Expected key prefix to be 'apikeyservice_'")

	if _, ok := result["rootKey"]; !ok {
		t.Errorf("Expected 'keyId' in response body, but not found")
	}
}
