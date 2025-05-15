package coraxclient

// DocumentCreate represents the data for creating a single document.
// This is based on the items in the DocumentIngest schema.
type DocumentCreate struct {
	ID          *string                `json:"id,omitempty"` // Optional client-provided ID
	TextContent *string                `json:"text_content,omitempty"`
	JsonContent map[string]interface{} `json:"json_content,omitempty"` // For arbitrary JSON object
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// DocumentUpdate represents the request body for updating a document.
// Based on openapi.json components.schemas.DocumentUpdate
type DocumentUpdate struct {
	TextContent *string                `json:"text_content,omitempty"`
	JsonContent map[string]interface{} `json:"json_content,omitempty"` // For arbitrary JSON object
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Document represents the document details as returned by the API.
// Based on openapi.json components.schemas.Document
type Document struct {
	// Links       map[string]HateoasLink `json:"_links,omitempty"`
	ID               string                 `json:"id"`
	CollectionID     string                 `json:"collection_id"` // Not in API response schema for GET Document, but useful context
	Content          interface{}            `json:"content"`       // oneOf: [string, object]
	TextContent      *string                `json:"text_content,omitempty"`
	JsonContent      map[string]interface{} `json:"json_content,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	TokenCount       int                    `json:"token_count"`
	ChunkCount       int                    `json:"chunk_count"`
	EmbeddingsStatus string                 `json:"embeddings_status"`
	CreatedBy        string                 `json:"created_by"`
	UpdatedBy        *string                `json:"updated_by,omitempty"`
	CreatedAt        string                 `json:"created_at"` // Expected format: date-time
	UpdatedAt        *string                `json:"updated_at,omitempty"` // Expected format: date-time
}

// Note: For JsonContent in DocumentCreate and DocumentUpdate, if the user provides a JSON string,
// the Terraform provider will need to unmarshal it into map[string]interface{} before sending.
// Alternatively, the client methods could accept a raw json.RawMessage or string and handle it.
// For now, map[string]interface{} is used, assuming the provider prepares this structure.

// The API for creating documents is POST /v1/collections/{collection_id}/documents
// and it takes a DocumentIngest which is an array of DocumentCreate-like objects.
// For a single document resource, we'll likely use the PUT endpoint for create/update.
// PUT /v1/collections/{collection_id}/documents/{document_id}
// The body for this PUT is DocumentUpdate. This implies that for creation with a client-specified ID,
// we'd use PUT. If the ID is API-generated, a POST to the batch endpoint (with a single item array)
// might be an option, or the API might have a dedicated single POST not explicitly shown for DocumentCreate.
// Given the PUT takes DocumentUpdate, and DocumentUpdate doesn't have an ID field,
// the ID is purely in the path. This is standard.
// The DocumentCreate struct above includes an optional ID, which might be used if we were to
// implement a client method that uses the batch POST /v1/collections/{collection_id}/documents.
// For the single document resource, the ID will be in the path for PUT.
