package schemas

// ApiKey struct to represent the data stored in DynamoDB
type ApiKeyIdRow struct {
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

type VerifyHashedKeyInput struct {
	ApiId     string `json:"apiId" dynamodbav:"pk"`
	HashedKey string `json:"-" dynamodbav:"sk"` // Store hashed key securely
}

type GetApiKeyInput struct {
	ApiId string `json:"apiId" dynamodbav:"pk"`
	KeyId string `json:"apiKeyId" dynamodbav:"sk"`
}
