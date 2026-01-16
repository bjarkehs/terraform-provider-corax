// Copyright (c) Trifork

package coraxclient

// ModelDeployment maps to components.schemas.ModelDeployment.
type ModelDeployment struct {
	// Links map[string]HateoasLink `json:"_links,omitempty"` // Assuming HateoasLink is defined elsewhere or not strictly needed for TF state
	Name           string            `json:"name"`
	Description    *string           `json:"description,omitempty"`
	SupportedTasks []string          `json:"supported_tasks"`     // Enum: "chat", "completion", "embedding"
	Configuration  map[string]string `json:"configuration"`       // Assuming string to string for simplicity based on TF schema choice
	IsActive       *bool             `json:"is_active,omitempty"` // API default true
	ProviderID     string            `json:"provider_id"`
	ID             string            `json:"id"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      *string           `json:"updated_at,omitempty"`
	CreatedBy      string            `json:"created_by"`
	UpdatedBy      *string           `json:"updated_by,omitempty"`
	// Deprecated fields from OpenAPI spec are omitted: api_version, model_name, deployment_name
}

// ModelDeploymentCreate maps to components.schemas.ModelDeploymentCreate.
type ModelDeploymentCreate struct {
	Name           string            `json:"name"`
	Description    *string           `json:"description,omitempty"`
	SupportedTasks []string          `json:"supported_tasks"`
	Configuration  map[string]string `json:"configuration"`
	IsActive       *bool             `json:"is_active,omitempty"`
	ProviderID     string            `json:"provider_id"`
}

// ModelDeploymentUpdate maps to components.schemas.ModelDeploymentUpdate
// Note: The API spec for ModelDeploymentUpdate requires all fields for PUT.
// The API performs a full replacement, not a partial update.
// The provider's Update method must send all current state values.
type ModelDeploymentUpdate struct {
	Name           string            `json:"name"`
	Description    *string           `json:"description,omitempty"` // Optional, can be null
	SupportedTasks []string          `json:"supported_tasks"`
	Configuration  map[string]string `json:"configuration"`
	IsActive       *bool             `json:"is_active,omitempty"` // Optional, defaults to true
	ProviderID     string            `json:"provider_id"`
}
