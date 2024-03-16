package db

import (
	"cfg"
	"context"
	"fmt"
	"schemas"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func CreateHashedKeyRow(hashedKey string, workspaceId string, keyId string, req schemas.CreateKeyRequest, dbClient *dynamodb.Client) error {
	// Create HashedKey struct
	hashedKeyToAdd := schemas.HashedKeyRow{
		WorkspaceIdApiId: fmt.Sprintf("workspaceId#%s-apiId#%s", workspaceId, req.ApiId),
		HashedKey:        "hashedKey#" + hashedKey,
		KeyId:            keyId,
		Name:             req.Name,
		Prefix:           req.Prefix,
		Roles:            req.Roles,
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

func CreateRootKeyRow(hashedRootKey string, req schemas.CreateRootKeyRequest, dbClient *dynamodb.Client) error {
	// Create RootKey struct
	rootKeyToAdd := schemas.RootKeyRow{
		RootKeyHash: "rootKeyHash#" + hashedRootKey,
		WorkspaceId: "workspaceId#" + req.WorkspaceId,
		Permissions: req.Permissions,
	}
	rootKeyToAddJSON, err := attributevalue.MarshalMap(rootKeyToAdd)
	if err != nil {
		return err
	}
	// Put the key in the DynamoDB table
	putRootKeyInput := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.Config.ApiKeyTable),
		Item:      rootKeyToAddJSON,
	}
	// Put item for RootKey lookup
	_, err = dbClient.PutItem(context.TODO(), putRootKeyInput)
	if err != nil {
		return err
	}

	return nil
}

func GetRootKeyRow(rootKeyHash string, dbClient *dynamodb.Client) (schemas.RootKeyRow, error) {

	// Get item from DynamoDB
	getRootKeyInput := &dynamodb.QueryInput{
		TableName:              aws.String(cfg.Config.ApiKeyTable),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "rootKeyHash#" + rootKeyHash},
		},
	}
	getRootKeyOutput, err := dbClient.Query(context.TODO(), getRootKeyInput)
	if err != nil {
		return schemas.RootKeyRow{}, err
	}

	rootKey := schemas.RootKeyRow{}
	err = attributevalue.UnmarshalMap(getRootKeyOutput.Items[0], &rootKey)
	if err != nil {
		return schemas.RootKeyRow{}, err
	}

	return rootKey, nil
}
