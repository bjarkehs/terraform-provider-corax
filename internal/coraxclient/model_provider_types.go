package coraxclient

// ModelProvider maps to components.schemas.ModelProvider
type ModelProvider struct {
	// Links map[string]HateoasLink `json:"_links,omitempty"`
	Name          string            `json:"name"`
	ProviderType  string            `json:"provider_type"`
	Configuration map[string]string `json:"configuration"` // Assuming string to string
	ID            string            `json:"id"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     *string           `json:"updated_at,omitempty"`
	CreatedBy     string            `json:"created_by"`
	UpdatedBy     *string           `json:"updated_by,omitempty"`
	// Deprecated fields: api_endpoint, api_key are omitted as they should be part of Configuration
}

// ModelProviderCreate maps to components.schemas.ModelProviderCreate
type ModelProviderCreate struct {
	Name          string            `json:"name"`
	ProviderType  string            `json:"provider_type"`
	Configuration map[string]string `json:"configuration"`
}

// ModelProviderUpdate maps to components.schemas.ModelProviderUpdate
// The API spec for ModelProviderUpdate includes `id` as required, which is unusual for an update payload.
// Typically, ID is in the path. Assuming `id` is not part of the request body for an update.
// The spec also makes all other fields required for PUT.
// For a more typical partial update, fields would be pointers.
// Sticking to the spec for now, which implies full replacement.
type ModelProviderUpdate struct {
	Name          string            `json:"name"`          // Required in API spec for PUT
	ProviderType  string            `json:"provider_type"` // Required in API spec for PUT
	Configuration map[string]string `json:"configuration"` // Required in API spec for PUT
	// ID string `json:"id"` // This is in the API spec for ModelProviderUpdate body, but usually not for PUT body. Omitting for now.
}
