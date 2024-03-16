package auth

import (
	"db"
	"errors"
	"strings"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func VerifyAuthHeader(event events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (string, error) {
	// Verify the request
	authHeader := event.Headers["Authorization"]
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	rootKey := parts[1]

	// Hash key for query
	hashedKey := utils.HashString(rootKey)

	rootKeyRow, err := db.GetRootKeyRow(hashedKey, dbClient)

	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(rootKeyRow.WorkspaceId, "workspaceId#"), nil
}
