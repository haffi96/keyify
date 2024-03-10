package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	d := deps{
		ddbClient: db.GetMockDynamoClient(ctx),
		tableName: "ApiKeyTableDev",
	}

	resp, err := d.handler(ctx, events.APIGatewayProxyRequest{
		Body: `{"apiId":"api-1","name":"key-1","prefix":"key_","roles":["role-1","role-2"]}`,
	})

	assert.Equal(t, nil, err, "Handler returned an error")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp.Body), &result)
	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}
	assert.Equal(t, "apiId#api-1", result["apiId"], "Expected apiId to be api-1")

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
