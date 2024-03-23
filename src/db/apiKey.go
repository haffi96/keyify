package db

import (
	"cfg"
	"context"
	"fmt"
	"schemas"
	"utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CreateApiKeyRow creates a new row in the ApiKey table
func CreateApiKeyRow(hashedKey string, workspaceId string, keyId string, req schemas.CreateKeyRequest, dbClient *dynamodb.Client) (schemas.ApiKeyIdRow, error) {
	// Create ApiKey struct
	apiKeyToAdd := schemas.ApiKeyIdRow{
		WorkspaceIdApiId: fmt.Sprintf("workspaceId#%s-apiId#%s", workspaceId, req.ApiId),
		KeyId:            "apiKeyId#" + keyId,
		HashedKey:        hashedKey,
		Name:             req.Name,
		Prefix:           req.Prefix,
		Roles:            req.Roles,
		CreatedAt:        utils.TimeNow(),
	}
	apiKeyToAddJSON, err := attributevalue.MarshalMap(apiKeyToAdd)
	if err != nil {
		return schemas.ApiKeyIdRow{}, err
	}
	// Put the key in the DynamoDB table
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.Config.ApiKeyTable),
		Item:      apiKeyToAddJSON,
	}

	// Put item for apiKeyId lookup
	_, err = dbClient.PutItem(context.TODO(), putInput)
	if err != nil {
		return schemas.ApiKeyIdRow{}, err
	}

	return apiKeyToAdd, nil
}

// CreateApiKeyDatetimeRow creates a new row in the ApiKey table for datetime sorting
func CreateApiKeyDatetimeRow(hashedKey string, workspaceId string, keyId string, req schemas.CreateKeyRequest, dbClient *dynamodb.Client) (schemas.ApiKeyDatetimeRow, error) {
	// Create ApiKeyDatetimeRow struct
	now := utils.TimeNow()
	apiKeyToAdd := schemas.ApiKeyDatetimeRow{
		WorkspaceIdApiId: fmt.Sprintf("workspaceId#%s-apiId#%s", workspaceId, req.ApiId),
		CreatedAtKeyId:   fmt.Sprintf("createdAt#%s-apiKeyId#%s", now, keyId),
		HashedKey:        hashedKey,
		Name:             req.Name,
		Prefix:           req.Prefix,
		Roles:            req.Roles,
		CreatedAt:        now,
	}
	apiKeyToAddJSON, err := attributevalue.MarshalMap(apiKeyToAdd)
	if err != nil {
		return schemas.ApiKeyDatetimeRow{}, err
	}
	// Put the key in the DynamoDB table
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.Config.ApiKeyTable),
		Item:      apiKeyToAddJSON,
	}

	// Put item for apiKeyId lookup
	_, err = dbClient.PutItem(context.TODO(), putInput)
	if err != nil {
		return schemas.ApiKeyDatetimeRow{}, err
	}

	return apiKeyToAdd, nil
}

// CreateHashedKeyRow creates a new row in the ApiKey table for quick hashed key lookup
func CreateHashedKeyRow(hashedKey string, workspaceId string, keyId string, req schemas.CreateKeyRequest, dbClient *dynamodb.Client) error {
	// Create HashedKey struct
	hashedKeyToAdd := schemas.HashedKeyRow{
		WorkspaceIdApiId: fmt.Sprintf("workspaceId#%s-apiId#%s", workspaceId, req.ApiId),
		HashedKey:        "hashedKey#" + hashedKey,
		KeyId:            keyId,
		Name:             req.Name,
		Prefix:           req.Prefix,
		Roles:            req.Roles,
		CreatedAt:        utils.TimeNow(),
	}
	hashedKeyToAddJSON, err := attributevalue.MarshalMap(hashedKeyToAdd)
	if err != nil {
		return err
	}
	// Put the key in the DynamoDB table
	putHashedKeyInput := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.Config.ApiKeyTable),
		Item:      hashedKeyToAddJSON,
	}
	// Put item for HashedKey lookup
	_, err = dbClient.PutItem(context.TODO(), putHashedKeyInput)
	if err != nil {
		return err
	}

	return nil
}

// GetApiKey retrieves an API key from the database by pk
func GetApiKey(workspaceId string, apiId string, keyId string, dbClient *dynamodb.Client) (schemas.ApiKeyIdRow, error) {
	// Construct DynamoDB key
	key := schemas.GetApiKeyInput{
		WorkspaceIdApiId: "workspaceId#" + workspaceId + "-" + "apiId#" + apiId,
		KeyId:            "apiKeyId#" + keyId,
	}

	keyJson, err := attributevalue.MarshalMap(key)
	if err != nil {
		return schemas.ApiKeyIdRow{}, err
	}

	// Get item from DynamoDB
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(cfg.Config.ApiKeyTable),
		Key:       keyJson,
	}

	result, err := dbClient.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return schemas.ApiKeyIdRow{}, err
	}

	// Check if key exists and return 404 if not found
	if result.Item == nil {
		return schemas.ApiKeyIdRow{}, nil
	}

	// Extract and return relevant data (excluding sensitive fields)
	apiKeyData := schemas.ApiKeyIdRow{}
	err = attributevalue.UnmarshalMap(result.Item, &apiKeyData)
	if err != nil {
		return schemas.ApiKeyIdRow{}, err
	}

	return apiKeyData, nil
}

// ListApiKeys retrieves all API keys for a given workspace and API, sorted by createdAt.
//
// Most recent keys are returned first.
func ListApiKeys(workspaceId string, apiId string, dbClient *dynamodb.Client) ([]schemas.ApiKeyIdRow, error) {
	// Construct DynamoDB key
	key := schemas.ListApiKeysInput{
		WorkspaceIdApiId: "workspaceId#" + workspaceId + "-" + "apiId#" + apiId,
	}

	// Get item from DynamoDB
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(cfg.Config.ApiKeyTable),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: key.WorkspaceIdApiId},
			":sk": &types.AttributeValueMemberS{Value: "createdAt"},
		},
		ScanIndexForward: aws.Bool(false),
	}

	// Execute the query
	result, err := dbClient.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, err
	}

	// Unmarshal the results
	keys := []schemas.ApiKeyIdRow{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
