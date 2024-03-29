package schemas

// ApiKey struct to represent the data stored in DynamoDB
type ApiKeyIdRow struct {
	WorkspaceIdApiId string   `json:"workspaceIdapiId" dynamodbav:"pk"`
	KeyId            string   `json:"apiKeyId" dynamodbav:"sk"`
	HashedKey        string   `json:"-" dynamodbav:"hashedKey"` // Store hashed key securely
	ApiKey           string   `json:"key" dynamodbav:"-"`       // Not used, store hashed key instead
	Name             string   `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Prefix           string   `json:"prefix,omitempty" dynamodbav:"prefix,omitempty"`
	Roles            []string `json:"roles,omitempty" dynamodbav:"roles,omitempty"`
	CreatedAt        string   `json:"createdAt" dynamodbav:"createdAt"`
}

// ApiKey with datetime sk for sorting
type ApiKeyDatetimeRow struct {
	WorkspaceIdApiId string   `json:"workspaceIdApiId" dynamodbav:"pk"`
	CreatedAtKeyId   string   `json:"createdAtKeyId" dynamodbav:"sk"`
	HashedKey        string   `json:"-" dynamodbav:"hashedKey"` // Store hashed key securely
	ApiKey           string   `json:"key" dynamodbav:"-"`       // Not used, store hashed key instead
	Name             string   `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Prefix           string   `json:"prefix,omitempty" dynamodbav:"prefix,omitempty"`
	Roles            []string `json:"roles,omitempty" dynamodbav:"roles,omitempty"`
	CreatedAt        string   `json:"createdAt" dynamodbav:"createdAt"`
}

type HashedKeyRow struct {
	WorkspaceIdApiId string   `json:"workspaceIdApiId" dynamodbav:"pk"`
	HashedKey        string   `json:"-" dynamodbav:"sk"` // Store hashed key securely
	KeyId            string   `json:"apiKeyId" dynamodbav:"apiKeyId"`
	ApiKey           string   `json:"key" dynamodbav:"-"` // Not used, store hashed key instead
	Name             string   `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Prefix           string   `json:"prefix,omitempty" dynamodbav:"prefix,omitempty"`
	Roles            []string `json:"roles,omitempty" dynamodbav:"roles,omitempty"`
	CreatedAt        string   `json:"createdAt" dynamodbav:"createdAt"`
}

type VerifyHashedKeyInput struct {
	WorkspaceIdApiId string `json:"workspaceIdApiId" dynamodbav:"pk"`
	HashedKey        string `json:"-" dynamodbav:"sk"` // Store hashed key securely
}

type GetApiKeyInput struct {
	WorkspaceIdApiId string `json:"workspaceIdApiId" dynamodbav:"pk"`
	KeyId            string `json:"apiKeyId" dynamodbav:"sk"`
}

type ListApiKeysInput struct {
	WorkspaceIdApiId string `json:"workspaceIdApiId" dynamodbav:"pk"`
}

type RootKeyRow struct {
	RootKeyHash string   `json:"rootKeyHash" dynamodbav:"pk"`
	WorkspaceId string   `json:"workspaceId" dynamodbav:"sk"`
	Permissions []string `json:"permissions,omitempty" dynamodbav:"permissions,omitempty"`
	CreatedAt   string   `json:"createdAt" dynamodbav:"createdAt"`
}

type GetRootKeyInput struct {
	RootKeyHash string `json:"rootKeyHash" dynamodbav:"pk"`
}

type ApiRow struct {
	WorkspaceId string `json:"workspaceId" dynamodbav:"pk"`
	ApiName     string `json:"apiName" dynamodbav:"apiName"`
	ApiId       string `json:"apiId" dynamodbav:"sk"`
	CreatedAt   string `json:"createdAt" dynamodbav:"createdAt"`
}
