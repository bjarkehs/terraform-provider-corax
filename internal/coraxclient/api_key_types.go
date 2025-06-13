package coraxclient

// ApiKeyCreate represents the request body for creating an API key.
// Based on openapi.json components.schemas.ApiKeyCreate.
type ApiKeyCreate struct {
	Name      string `json:"name"`
	ExpiresAt string `json:"expires_at"` // Expected format: date-time
}

// ApiKey represents the API key details.
// Based on openapi.json components.schemas.ApiKey.
type ApiKey struct {
	// Links       map[string]HateoasLink `json:"_links,omitempty"` // HateoasLink not defined yet
	ID         string  `json:"id"`
	Prefix     string  `json:"prefix,omitempty"`
	Key        string  `json:"key"` // This is sensitive and usually only returned on create
	Name       string  `json:"name"`
	ExpiresAt  *string `json:"expires_at"` // Pointer to handle null
	CreatedBy  string  `json:"created_by"`
	CreatedAt  string  `json:"created_at"` // Expected format: date-time
	UpdatedAt  *string `json:"updated_at"` // Pointer to handle null; Expected format: date-time
	IsActive   bool    `json:"is_active,omitempty"`
	LastUsedAt *string `json:"last_used_at"` // Pointer to handle null; Expected format: date-time
	UsageCount int     `json:"usage_count,omitempty"`
}

// TODO: Define HateoasLink if it's needed for client operations,
// based on openapi.json components.schemas.HateoasLink
// type HateoasLink struct {
// 	Href string `json:"href"`
// 	Type string `json:"type,omitempty"` // Corresponds to HTTPMethod
// }
