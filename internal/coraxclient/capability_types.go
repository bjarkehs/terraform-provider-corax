// Copyright (c) HashiCorp, Inc.

package coraxclient

// --- Common Capability Structures ---

// CapabilityConfig maps to components.schemas.CapabilityConfig.
type CapabilityConfig struct {
	Temperature      *float64               `json:"temperature,omitempty"`
	BlobConfig       *BlobConfig            `json:"blob_config,omitempty"`
	DataRetention    *DataRetention         `json:"data_retention,omitempty"` // Polymorphic
	ContentTracing   *bool                  `json:"content_tracing,omitempty"`
	CustomParameters map[string]interface{} `json:"custom_parameters,omitempty"`
}

// BlobConfig maps to components.schemas.BlobConfig.
type BlobConfig struct {
	MaxFileSizeMB    *int     `json:"max_file_size_mb,omitempty"`
	MaxBlobs         *int     `json:"max_blobs,omitempty"`
	AllowedMimeTypes []string `json:"allowed_mime_types,omitempty"`
}

// DataRetention is a wrapper for polymorphic retention types.
type DataRetention struct {
	Type  string `json:"type"`            // "timed" or "infinite"
	Hours *int   `json:"hours,omitempty"` // For TimedDataRetention
}

// --- Chat Capability Specific Structures ---

// ChatCapabilityCreate maps to components.schemas.ChatCapabilityCreate.
type ChatCapabilityCreate struct {
	Name         string            `json:"name"`
	IsPublic     *bool             `json:"is_public,omitempty"`
	Type         string            `json:"type"` // Should always be "chat"
	ModelID      *string           `json:"model_id,omitempty"`
	Config       *CapabilityConfig `json:"config,omitempty"`
	ProjectID    *string           `json:"project_id,omitempty"`
	SystemPrompt string            `json:"system_prompt"`
	// CollectionIDs []string       `json:"collection_ids,omitempty"` // Omitted for now
}

// ChatCapabilityUpdate maps to components.schemas.ChatCapabilityUpdate.
type ChatCapabilityUpdate struct {
	Name         *string           `json:"name,omitempty"` // Note: API spec says name is required here, but usually updates are partial.
	IsPublic     *bool             `json:"is_public,omitempty"`
	Type         *string           `json:"type,omitempty"` // Should always be "chat" if sent
	ModelID      *string           `json:"model_id,omitempty"`
	Config       *CapabilityConfig `json:"config,omitempty"`
	ProjectID    *string           `json:"project_id,omitempty"`
	SystemPrompt *string           `json:"system_prompt,omitempty"`
	// CollectionIDs []string       `json:"collection_ids,omitempty"` // Omitted for now
}

// CapabilityRepresentation maps to components.schemas.CapabilityRepresentation
// This is a generic structure returned by the API for GET requests.
// We will use its fields to populate specific chat or completion capability models.
type CapabilityRepresentation struct {
	// Links map[string]HateoasLink `json:"_links,omitempty"`
	Name          string                 `json:"name"`
	IsPublic      *bool                  `json:"is_public"` // API default false
	Type          string                 `json:"type"`      // "chat" or "completion"
	ModelID       *string                `json:"model_id"`
	Config        *CapabilityConfig      `json:"config"` // API returns the resolved config
	ProjectID     *string                `json:"project_id"`
	ID            string                 `json:"id"`
	SemanticID    string                 `json:"semantic_id"`
	CreatedBy     string                 `json:"created_by"`
	UpdatedBy     string                 `json:"updated_by"` // API shows this as non-nullable, but might be null in practice
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	ArchivedAt    *string                `json:"archived_at"`
	Owner         string                 `json:"owner"`
	Input         map[string]interface{} `json:"input"`         // For CapabilityRepresentation
	Output        map[string]interface{} `json:"output"`        // For CapabilityRepresentation
	Configuration map[string]interface{} `json:"configuration"` // For CapabilityRepresentation

	// Chat-specific fields from ChatCapability (if type is "chat")
	// These are not directly in CapabilityRepresentation but are part of the underlying ChatCapability
	// that CapabilityRepresentation might represent.
	// The API GET /v1/capabilities/{id} returns CapabilityRepresentation.
	// The actual chat/completion specific fields like system_prompt are part of the 'configuration' field in CapabilityRepresentation.
	// Let's adjust based on how the API actually returns these.
	// The openapi spec shows `system_prompt` directly in `ChatCapability` schema,
	// and `CapabilityRepresentation` has a `configuration` field of type object.
	// It's likely that `system_prompt` etc. are nested within this `configuration` field in the response.
	// For now, we'll assume the specific fields are directly accessible or need to be extracted from `configuration`.
	// The `POST` body for creating a capability is `anyOf[CompletionCapabilityCreate, ChatCapabilityCreate]`.
	// The `GET` response is `CapabilityRepresentation`.
	// The `PUT` body is `anyOf[CompletionCapabilityUpdate, ChatCapabilityUpdate]`.

	// Let's assume for now that when we GET a capability, if it's a chat type,
	// the specific chat fields are available, perhaps through the 'configuration' map.
	// The client methods will need to handle this.
	// For the purpose of defining request/response structs for the client,
	// we'll use the specific Create/Update structs for POST/PUT.
	// For GET, we'll get CapabilityRepresentation and then potentially map its 'configuration'
	// or other fields to our more specific Terraform models.
}

// --- Completion Capability Specific Structures ---

// CompletionCapabilityCreate maps to components.schemas.CompletionCapabilityCreate.
type CompletionCapabilityCreate struct {
	Name             string                 `json:"name"`
	IsPublic         *bool                  `json:"is_public,omitempty"`
	Type             string                 `json:"type"` // Should always be "completion"
	SemanticID       *string                `json:"semantic_id,omitempty"`
	ModelID          *string                `json:"model_id,omitempty"`
	Config           *CapabilityConfig      `json:"config,omitempty"`
	ProjectID        *string                `json:"project_id,omitempty"`
	SystemPrompt     string                 `json:"system_prompt"`
	CompletionPrompt string                 `json:"completion_prompt"`
	Variables        []string               `json:"variables,omitempty"`
	OutputType       string                 `json:"output_type"`          // "schema" or "text"
	SchemaDef        map[string]interface{} `json:"schema_def,omitempty"` // Used if output_type is "schema"
}

// CompletionCapabilityUpdate maps to components.schemas.CompletionCapabilityUpdate.
type CompletionCapabilityUpdate struct {
	Name             *string                `json:"name,omitempty"`
	IsPublic         *bool                  `json:"is_public,omitempty"`
	Type             *string                `json:"type,omitempty"` // Should always be "completion" if sent
	SemanticID       *string                `json:"semantic_id,omitempty"`
	ModelID          *string                `json:"model_id,omitempty"`
	Config           *CapabilityConfig      `json:"config,omitempty"`
	ProjectID        *string                `json:"project_id,omitempty"`
	SystemPrompt     *string                `json:"system_prompt,omitempty"`
	CompletionPrompt *string                `json:"completion_prompt,omitempty"`
	Variables        []string               `json:"variables,omitempty"` // To clear, send empty list? To leave unchanged, omit.
	OutputType       *string                `json:"output_type,omitempty"`
	SchemaDef        map[string]interface{} `json:"schema_def,omitempty"`
}

// --- Capability Type Specific Structures ---

// DefaultModelDeploymentUpdate maps to components.schemas.DefaultModelDeploymentUpdate.
type DefaultModelDeploymentUpdate struct {
	DefaultModelDeploymentID string `json:"default_model_deployment_id"`
}

// CapabilityTypeRepresentation maps to components.schemas.CapabilityTypeRepresentation.
type CapabilityTypeRepresentation struct {
	// Links map[string]HateoasLink `json:"_links,omitempty"`
	ID                       string  `json:"id"`   // This is the capability_type string like "chat"
	Name                     string  `json:"name"` // Display name like "Chat"
	DefaultModelDeploymentID *string `json:"default_model_deployment_id,omitempty"`
	// Embedded map[string]ModelDeployment `json:"_embedded,omitempty"` // Assuming ModelDeployment is defined elsewhere
}

// CapabilityTypesRepresentation maps to components.schemas.CapabilityTypesRepresentation
// Used for GET /v1/capability-types.
type CapabilityTypesRepresentation struct {
	// Links   map[string]HateoasLink         `json:"_links,omitempty"`
	Embedded []CapabilityTypeRepresentation `json:"_embedded"`
}
