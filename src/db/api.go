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

func CreateApiRow(workspaceId string, apiId string, apiName string, dbClient *dynamodb.Client) (schemas.ApiRow, error) {
	// Create ApiKey struct
	apiKeyToAdd := schemas.ApiRow{
		WorkspaceId: fmt.Sprintf("workspaceId#%s", workspaceId),
		ApiName:     apiName,
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

func ListApis(workspaceId string, dbClient *dynamodb.Client) ([]schemas.ApiRow, error) {
	// Query the table
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(cfg.Config.ApiKeyTable),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("workspaceId#%s", workspaceId)},
		},
	}

	// Get the items
	output, err := dbClient.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, err
	}

	// Unmarshal the items
	var apis []schemas.ApiRow
	err = attributevalue.UnmarshalListOfMaps(output.Items, &apis)
	if err != nil {
		return nil, err
	}

	return apis, nil
}
