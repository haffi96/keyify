package auth

import (
	"db"
	"errors"
	"log"
	"strings"
	"utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func VerifyAuthHeader(event events.APIGatewayProxyRequest, dbClient *dynamodb.Client) (string, error) {
	// Verify the request
	authHeader := event.Headers["Authorization"]
	log.Println("authHeader: ", authHeader)
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	rootKey := parts[1]

	log.Println("rootKey: ", rootKey)

	// Hash key for query
	hashedKey := utils.HashString(rootKey)

	log.Println("hashedKey: ", hashedKey)

	rootKeyRow, err := db.GetRootKeyRow(hashedKey, dbClient)

	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(rootKeyRow.WorkspaceId, "workspaceId#"), nil
}
