package schemas

// API response structure
type CreateApiResponse struct {
	ApiId string `json:"apiId"`
}

type GetApiResponse struct {
	ApiId     string `json:"apiId"`
	ApiName   string `json:"apiName"`
	CreatedAt string `json:"createdAt"`
}

type CreateKeyResponse struct {
	ApiId string `json:"apiId"`
	KeyId string `json:"keyId"`
	Key   string `json:"key"`
}

type GetApiKeyResponse struct {
	ApiId     string   `json:"apiId"`
	KeyId     string   `json:"apiKeyId"`
	Name      string   `json:"name,omitempty"`
	Prefix    string   `json:"prefix,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	CreatedAt string   `json:"createdAt"`
}

type CreateRootKeyResponse struct {
	RootKey string `json:"rootKey"`
}
