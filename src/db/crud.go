package db

import (
	"cfg"
	"context"
	"schemas"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func CreateApiKeyRow(hashedKey string, keyId string, req schemas.CreateKeyRequest, dbClient *dynamodb.Client) (schemas.ApiKeyIdRow, error) {
	// Create ApiKey struct
	apiKeyToAdd := schemas.ApiKeyIdRow{
		ApiId:     "apiId#" + req.ApiId,
		KeyId:     "apiKeyId#" + keyId,
		HashedKey: hashedKey,
		Name:      req.Name,
		Prefix:    req.Prefix,
		Roles:     req.Roles,
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

func CreateHashedKeyRow(hashedKey string, keyId string, req schemas.CreateKeyRequest, dbClient *dynamodb.Client) error {
	// Create HashedKey struct
	hashedKeyToAdd := schemas.HashedKeyRow{
		ApiId:     "apiId#" + req.ApiId,
		HashedKey: "hashedKey#" + hashedKey,
		KeyId:     keyId,
		Name:      req.Name,
		Prefix:    req.Prefix,
		Roles:     req.Roles,
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
