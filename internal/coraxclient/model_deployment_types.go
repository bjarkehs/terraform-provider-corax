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
// Note: The API spec for ModelDeploymentUpdate is identical to ModelDeploymentCreate.
// All fields are required in the API spec for PUT, which is unusual for an update.
// Typically, updates are partial. If the API truly requires all fields for PUT,
// the provider's Update method will need to send all current state values.
// For now, defining struct as partial (pointers) to align with typical update patterns.
// If API enforces full replacement, this struct and update logic will need adjustment.
type ModelDeploymentUpdate struct {
	Name           *string           `json:"name,omitempty"`
	Description    *string           `json:"description,omitempty"` // Allow clearing description
	SupportedTasks []string          `json:"supported_tasks,omitempty"`
	Configuration  map[string]string `json:"configuration,omitempty"`
	IsActive       *bool             `json:"is_active,omitempty"`
	ProviderID     *string           `json:"provider_id,omitempty"` // ProviderID might not be updatable, check API behavior
}
