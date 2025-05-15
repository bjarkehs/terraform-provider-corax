package coraxclient

// EmbeddingsModelCreate represents the request body for creating an embeddings model.
// Based on openapi.json components.schemas.EmbeddingsModelCreate
type EmbeddingsModelCreate struct {
	Name                string  `json:"name"`
	Description         *string `json:"description,omitempty"`
	ModelProvider       string  `json:"model_provider"`
	ModelNameOnProvider string  `json:"model_name_on_provider"`
	Dimensions          int     `json:"dimensions"` // API uses integer, TF model uses int64
	ApiKey              *string `json:"api_key,omitempty"`    // Sensitive
	ApiBaseUrl          *string `json:"api_base_url,omitempty"`
	MaxTokens           *int    `json:"max_tokens,omitempty"` // API uses integer
}

// EmbeddingsModelUpdate represents the request body for updating an embeddings model.
// Based on openapi.json components.schemas.EmbeddingsModelUpdate
type EmbeddingsModelUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ApiKey      *string `json:"api_key,omitempty"` // Sensitive
	ApiBaseUrl  *string `json:"api_base_url,omitempty"`
	// Note: model_provider, model_name_on_provider, dimensions, max_tokens are not updatable via this schema.
	// Status and is_default are likely not directly updatable by user.
}

// EmbeddingsModel represents the embeddings model details.
// Based on openapi.json components.schemas.EmbeddingsModel
type EmbeddingsModel struct {
	// Links       map[string]HateoasLink `json:"_links,omitempty"`
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	Description         *string `json:"description,omitempty"`
	ModelProvider       string  `json:"model_provider"`
	ModelNameOnProvider string  `json:"model_name_on_provider"`
	Dimensions          int     `json:"dimensions"`    // API uses integer
	MaxTokens           *int    `json:"max_tokens,omitempty"` // API uses integer
	ApiKey              *string `json:"api_key,omitempty"`    // Only present in response if just set? Usually not returned.
	ApiBaseUrl          *string `json:"api_base_url,omitempty"`
	Status              string  `json:"status"`
	IsDefault           bool    `json:"is_default"`
	CreatedBy           string  `json:"created_by"`
	UpdatedBy           *string `json:"updated_by,omitempty"`
	CreatedAt           string  `json:"created_at"` // Expected format: date-time
	UpdatedAt           *string `json:"updated_at,omitempty"` // Expected format: date-time
}
