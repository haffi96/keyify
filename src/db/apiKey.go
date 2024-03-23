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
)

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

func CreateApiKeyDatetimeRow(hashedKey string, workspaceId string, keyId string, req schemas.CreateKeyRequest, dbClient *dynamodb.Client) (schemas.ApiKeyDatetimeRow, error) {
	// Create ApiKeyDatetimeRow struct
	now := utils.TimeNow()
	apiKeyToAdd := schemas.ApiKeyDatetimeRow{
		WorkspaceIdApiId: fmt.Sprintf("workspaceId#%s-apiId#%s", workspaceId, req.ApiId),
		CreatedAtKeyId:   fmt.Sprintf("%s#apiKeyId#%s", now, keyId),
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
