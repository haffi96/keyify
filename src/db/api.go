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

func CreateApiRow(workspaceId string, apiId string, dbClient *dynamodb.Client) (schemas.ApiRow, error) {
	// Create ApiKey struct
	apiKeyToAdd := schemas.ApiRow{
		WorkspaceId: fmt.Sprintf("workspaceId#%s", workspaceId),
		ApiId:       "apiId#" + apiId,
		CreatedAt:   utils.TimeNow(),
	}

	apiKeyToAddJSON, err := attributevalue.MarshalMap(apiKeyToAdd)
	if err != nil {
		return schemas.ApiRow{}, err
	}

	// Put the key in the DynamoDB table
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.Config.ApiKeyTable),
		Item:      apiKeyToAddJSON,
	}

	// Put item for apiKeyId lookup
	_, err = dbClient.PutItem(context.TODO(), putInput)
	if err != nil {
		return schemas.ApiRow{}, err
	}

	return apiKeyToAdd, nil
}
