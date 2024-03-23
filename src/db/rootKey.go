package db

import (
	"cfg"
	"context"
	"errors"
	"schemas"
	"utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateRootKeyRow(hashedRootKey string, req schemas.CreateRootKeyRequest, dbClient *dynamodb.Client) error {
	// Create RootKey struct
	rootKeyToAdd := schemas.RootKeyRow{
		RootKeyHash: "rootKeyHash#" + hashedRootKey,
		WorkspaceId: "workspaceId#" + req.WorkspaceId,
		Permissions: req.Permissions,
		CreatedAt:   utils.TimeNow(),
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

	if len(getRootKeyOutput.Items) == 0 {
		return schemas.RootKeyRow{}, errors.New("root key invalid")
	}

	rootKey := schemas.RootKeyRow{}
	err = attributevalue.UnmarshalMap(getRootKeyOutput.Items[0], &rootKey)
	if err != nil {
		return schemas.RootKeyRow{}, err
	}

	return rootKey, nil
}
