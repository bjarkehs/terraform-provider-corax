package coraxclient

// CollectionCreate represents the request body for creating a collection.
// Based on openapi.json components.schemas.CollectionCreate
type CollectionCreate struct {
	Name              string            `json:"name"`
	Description       *string           `json:"description,omitempty"`
	ProjectID         string            `json:"project_id"`
	EmbeddingsModelID *string           `json:"embeddings_model_id,omitempty"`
	MetadataSchema    map[string]string `json:"metadata_schema,omitempty"` // API allows null, omitempty handles if not provided
}

// CollectionUpdate represents the request body for updating a collection.
// Based on openapi.json components.schemas.CollectionUpdate
type CollectionUpdate struct {
	Name              *string            `json:"name,omitempty"`
	Description       *string            `json:"description,omitempty"` // To clear, send null or empty string as per API behavior
	EmbeddingsModelID *string            `json:"embeddings_model_id,omitempty"`
	MetadataSchema    *map[string]string `json:"metadata_schema,omitempty"` // Pointer to allow sending explicit null to clear the schema if API supports it
}

// Collection represents the collection details.
// Based on openapi.json components.schemas.Collection
type Collection struct {
	// Links       map[string]HateoasLink `json:"_links,omitempty"` // HateoasLink not defined yet
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Description       *string           `json:"description,omitempty"` // Nullable
	ProjectID         string            `json:"project_id"`
	EmbeddingsModelID *string           `json:"embeddings_model_id,omitempty"` // Nullable
	MetadataSchema    map[string]string `json:"metadata_schema,omitempty"`     // Nullable
	CreatedBy         string            `json:"created_by"`
	UpdatedBy         *string           `json:"updated_by,omitempty"` // Nullable
	CreatedAt         string            `json:"created_at"`           // Expected format: date-time
	UpdatedAt         *string           `json:"updated_at,omitempty"` // Nullable; Expected format: date-time
	DocumentCount     int64             `json:"document_count"`       // Using int64 to match Terraform model
	SizeBytes         int64             `json:"size_bytes"`           // Using int64 to match Terraform model
	Status            string            `json:"status"`
}

// Note: HateoasLink definition is still pending from api_key_types.go
// if it becomes necessary for client operations.
