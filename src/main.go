package src

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type Deps struct {
	DbClient  *dynamodb.Client
	TableName string
}
