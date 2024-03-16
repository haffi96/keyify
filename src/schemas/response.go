package schemas

// API response structure
type CreateKeyResponse struct {
	ApiId string `json:"apiId"`
	KeyId string `json:"keyId"`
	Key   string `json:"key"`
}

type GetApiKeyResponse struct {
	ApiId  string   `json:"apiId"`
	KeyId  string   `json:"apiKeyId"`
	Name   string   `json:"name,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

type CreateRootKeyResponse struct {
	RootKey string `json:"rootKey"`
}
