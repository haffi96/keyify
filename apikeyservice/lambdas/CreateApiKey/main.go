package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ApiKey struct to represent the data stored in DynamoDB
type apiKeyIdRow struct {
	ApiId     string   `json:"apiId" dynamodbav:"pk"`
	KeyId     string   `json:"apiKeyId" dynamodbav:"sk"`
	HashedKey string   `json:"-" dynamodbav:"hashedKey"` // Store hashed key securely
	ApiKey    string   `json:"key" dynamodbav:"-"`       // Not used, store hashed key instead
	Name      string   `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Prefix    string   `json:"prefix,omitempty" dynamodbav:"prefix,omitempty"`
	Roles     []string `json:"roles,omitempty" dynamodbav:"roles,omitempty"`
}

type HashedKeyRow struct {
	ApiId     string   `json:"apiId" dynamodbav:"pk"`
	HashedKey string   `json:"-" dynamodbav:"sk"` // Store hashed key securely
	KeyId     string   `json:"apiKeyId" dynamodbav:"apiKeyId"`
	ApiKey    string   `json:"key" dynamodbav:"-"` // Not used, store hashed key instead
	Name      string   `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Prefix    string   `json:"prefix,omitempty" dynamodbav:"prefix,omitempty"`
	Roles     []string `json:"roles,omitempty" dynamodbav:"roles,omitempty"`
}

// API request structure
type CreateKeyRequest struct {
	ApiId  string   `json:"apiId"`
	Name   string   `json:"name,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

// API response structure
type CreateKeyResponse struct {
	KeyId string `json:"keyId"`
	Key   string `json:"key"` // **Do not return the actual API key**
}

func main() {
	lambda.Start(handler)
}

func errorResponse(statusCode int, message string, args ...interface{}) (events.APIGatewayProxyResponse, error) {
	body := fmt.Sprintf(message, args...)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func hashString(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func generateKeyId() (string, error) {
	// Prefix to append to the random string
	prefix := "key_"

	// Length of the random string (excluding prefix)
	randomStringLength := 16

	// Characters to use for the random string
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Generate random string directly
	b := make([]byte, randomStringLength)
	for i := range b {
		b[i] = chars[mathRand.Intn(len(chars))] // Select a random character
	}
	randomString := string(b)

	// Combine the prefix and random string
	keyId := prefix + randomString

	return keyId, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-2"),
	)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error loading SDK config: %s", err.Error())
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Parse and validate request body
	var req CreateKeyRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return errorResponse(http.StatusBadRequest, "Invalid request body: %s", err.Error())
	}

	// Parse and validate request body
	if req.ApiId == "" {
		return errorResponse(http.StatusBadRequest, "Missing required field: apiId")
	}
	// ... Add more validations as needed

	// Generate random API key
	apiKeyBytes, err := generateRandomBytes(36)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error generating random bytes: %s", err.Error())
	}
	apiKey := base64.URLEncoding.EncodeToString(apiKeyBytes)

	// Hash the API key
	hashedKey := hashString(apiKey)

	// Create key ID
	keyId, err := generateKeyId()
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error generating key ID: %s", err.Error())
	}

	// Create ApiKey struct
	apiKeyToAdd := apiKeyIdRow{
		ApiId:     "apiId#" + req.ApiId,
		KeyId:     "apiKeyId#" + keyId,
		HashedKey: hashedKey,
		Name:      req.Name,
		Prefix:    req.Prefix,
		Roles:     req.Roles,
	}
	apiKeyToAddJSON, err := attributevalue.MarshalMap(apiKeyToAdd)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error marshalling key to JSON: %s", err.Error())
	}
	// Put the key in the DynamoDB table
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("ApiKeyTableDev"),
		Item:      apiKeyToAddJSON,
	}

	// Put item for apiKeyId lookup
	_, err = svc.PutItem(context.TODO(), putInput)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error adding key to DynamoDB: %s", err.Error())
	}

	// Create HashedKey struct
	hashedKeyToAdd := HashedKeyRow{
		ApiId:     "apiId#" + req.ApiId,
		HashedKey: "hashedKey#" + hashedKey,
		KeyId:     keyId,
		Name:      req.Name,
		Prefix:    req.Prefix,
		Roles:     req.Roles,
	}
	hashedKeyToAddJSON, err := attributevalue.MarshalMap(hashedKeyToAdd)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error marshalling key to JSON: %s", err.Error())
	}
	// Put the key in the DynamoDB table
	putHashedKeyInput := &dynamodb.PutItemInput{
		TableName: aws.String("ApiKeyTableDev"),
		Item:      hashedKeyToAddJSON,
	}
	// Put item for HashedKey lookup
	_, err = svc.PutItem(context.TODO(), putHashedKeyInput)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Error adding key to DynamoDB: %s", err.Error())
	}

	// Return a success response with masked key
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"keyId": "%s"}`, strings.TrimPrefix(apiKeyToAdd.KeyId, "apiKeyId#")),
	}, nil
}
