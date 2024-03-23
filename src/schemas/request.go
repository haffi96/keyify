package schemas

// API request structure
type CreateApiRequest struct {
	Name string `json:"name"`
}

type CreateKeyRequest struct {
	ApiId  string   `json:"apiId"`
	Name   string   `json:"name,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

type VerifyKeyRequest struct {
	ApiId string `json:"apiId"`
	Key   string `json:"key"`
}
type CreateRootKeyRequest struct {
	WorkspaceId string   `json:"workspaceId"`
	Name        string   `json:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}
