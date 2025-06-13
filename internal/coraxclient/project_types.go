package coraxclient

// ProjectCreate represents the request body for creating a project.
// Based on openapi.json components.schemas.ProjectCreate.
type ProjectCreate struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty"` // API defaults to false if not provided
}

// ProjectUpdate represents the request body for updating a project.
// Based on openapi.json components.schemas.ProjectUpdate.
type ProjectUpdate struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsPublic    bool    `json:"is_public"`
}

// Project represents the project details.
// Based on openapi.json components.schemas.Project.
type Project struct {
	// Links       map[string]HateoasLink `json:"_links,omitempty"` // HateoasLink not defined yet
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	IsPublic        bool    `json:"is_public"`
	CreatedBy       string  `json:"created_by"`
	UpdatedBy       *string `json:"updated_by,omitempty"` // Can be null
	CreatedAt       string  `json:"created_at"`           // Expected format: date-time
	UpdatedAt       *string `json:"updated_at,omitempty"` // Can be null; Expected format: date-time
	Owner           string  `json:"owner"`
	CollectionCount int     `json:"collection_count"`
	CapabilityCount int     `json:"capability_count"`
}

// Note: HateoasLink definition is still pending from api_key_types.go
// if it becomes necessary for client operations.
