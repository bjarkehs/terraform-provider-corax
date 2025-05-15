package coraxclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	apiKeyHeader   = "X-API-Key"
)

// Client manages communication with the Corax API.
type Client struct {
	// HTTP client used to communicate with the API.
	httpClient *http.Client

	// Base URL for API requests. Must include scheme and host.
	BaseURL *url.URL

	// API key for authentication.
	APIKey string

	// UserAgent for client
	UserAgent string
}

// NewClient returns a new Corax API client.
func NewClient(baseURLStr string, apiKey string) (*Client, error) {
	if strings.TrimSpace(baseURLStr) == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("apiKey cannot be empty")
	}

	parsedBaseURL, err := url.ParseRequestURI(baseURLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}
	if parsedBaseURL.Scheme == "" || parsedBaseURL.Host == "" {
		return nil, fmt.Errorf("baseURL must include scheme and host")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		BaseURL:   parsedBaseURL,
		APIKey:    apiKey,
		UserAgent: "corax-terraform-provider/0.0.1", // TODO: Make version dynamic
	}, nil
}

// APIError represents an error response from the Corax API.
type APIError struct {
	StatusCode int
	Message    string
	Body       []byte
	// TODO: Could include a more structured error, e.g. from HTTPValidationError schema
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error: status %d, message: %s", e.StatusCode, e.Message)
}

// ErrNotFound is returned when a resource is not found (HTTP 404).
var ErrNotFound = &APIError{StatusCode: http.StatusNotFound, Message: "resource not found"}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	relURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path: %w", err)
	}

	fullURL := c.BaseURL.ResolveReference(relURL)

	var reqBody io.ReadWriter
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(apiKeyHeader, c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) doRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Body:       respBodyBytes,
		}
		// Try to unmarshal into a standard error structure if available
		// For now, just use a generic message or the body itself if it's short.
		if len(respBodyBytes) > 0 && len(respBodyBytes) < 512 { // Arbitrary limit for error message
			apiErr.Message = string(respBodyBytes)
		} else {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return apiErr
	}

	if v != nil {
		if err := json.Unmarshal(respBodyBytes, v); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w, body: %s", err, string(respBodyBytes))
		}
	}

	return nil
}

// CreateAPIKey creates a new API key.
// Corresponds to POST /v1/api-keys
func (c *Client) CreateAPIKey(ctx context.Context, apiKeyData ApiKeyCreate) (*ApiKey, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/api-keys", apiKeyData)
	if err != nil {
		return nil, err
	}

	var createdAPIKey ApiKey
	if err := c.doRequest(req, &createdAPIKey); err != nil {
		return nil, err
	}
	return &createdAPIKey, nil
}

// GetAPIKey retrieves a specific API key by its ID.
// Corresponds to GET /v1/api-keys/{key_id}
func (c *Client) GetAPIKey(ctx context.Context, keyID string) (*ApiKey, error) {
	if strings.TrimSpace(keyID) == "" {
		return nil, fmt.Errorf("keyID cannot be empty")
	}
	path := fmt.Sprintf("/v1/api-keys/%s", keyID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var apiKey ApiKey
	if err := c.doRequest(req, &apiKey); err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// DeleteAPIKey deletes a specific API key by its ID.
// Corresponds to DELETE /v1/api-keys/{key_id}
// The OpenAPI spec indicates a 200 response with an empty JSON object {} on success.
func (c *Client) DeleteAPIKey(ctx context.Context, keyID string) error {
	if strings.TrimSpace(keyID) == "" {
		return fmt.Errorf("keyID cannot be empty")
	}
	path := fmt.Sprintf("/v1/api-keys/%s", keyID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	// Expecting an empty JSON object {} on successful delete, so pass a dummy struct or nil for v.
	// If the API returns no body on 200/204, doRequest will handle it.
	return c.doRequest(req, nil)
}

// --- Project Methods ---

// CreateProject creates a new project.
// Corresponds to POST /v1/projects
func (c *Client) CreateProject(ctx context.Context, projectData ProjectCreate) (*Project, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/projects", projectData)
	if err != nil {
		return nil, err
	}

	var createdProject Project
	if err := c.doRequest(req, &createdProject); err != nil {
		return nil, err
	}
	return &createdProject, nil
}

// GetProject retrieves a specific project by its ID.
// Corresponds to GET /v1/projects/{project_id}
func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	if strings.TrimSpace(projectID) == "" {
		return nil, fmt.Errorf("projectID cannot be empty")
	}
	path := fmt.Sprintf("/v1/projects/%s", projectID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := c.doRequest(req, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject updates a specific project by its ID.
// Corresponds to PUT /v1/projects/{project_id}
func (c *Client) UpdateProject(ctx context.Context, projectID string, projectData ProjectUpdate) (*Project, error) {
	if strings.TrimSpace(projectID) == "" {
		return nil, fmt.Errorf("projectID cannot be empty")
	}
	path := fmt.Sprintf("/v1/projects/%s", projectID)
	req, err := c.newRequest(ctx, http.MethodPut, path, projectData)
	if err != nil {
		return nil, err
	}

	var updatedProject Project
	if err := c.doRequest(req, &updatedProject); err != nil {
		return nil, err
	}
	return &updatedProject, nil
}

// DeleteProject deletes a specific project by its ID.
// Corresponds to DELETE /v1/projects/{project_id}
// Expects a 204 No Content on success.
func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	if strings.TrimSpace(projectID) == "" {
		return fmt.Errorf("projectID cannot be empty")
	}
	path := fmt.Sprintf("/v1/projects/%s", projectID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, nil) // No body expected on 204
}

// --- Collection Methods --- (REMOVED)
// --- Document Methods --- (REMOVED)
// --- Embeddings Model Methods --- (REMOVED)

// --- Capability Methods ---

// CreateCapability creates a new capability.
// The payload should be either ChatCapabilityCreate or CompletionCapabilityCreate.
// Corresponds to POST /v1/capabilities
func (c *Client) CreateCapability(ctx context.Context, capabilityData interface{}) (*CapabilityRepresentation, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/capabilities", capabilityData)
	if err != nil {
		return nil, err
	}

	var createdCapability CapabilityRepresentation
	if err := c.doRequest(req, &createdCapability); err != nil {
		return nil, err
	}
	return &createdCapability, nil
}

// GetCapability retrieves a specific capability by its ID.
// Corresponds to GET /v1/capabilities/{capability_id}
func (c *Client) GetCapability(ctx context.Context, capabilityID string) (*CapabilityRepresentation, error) {
	if strings.TrimSpace(capabilityID) == "" {
		return nil, fmt.Errorf("capabilityID cannot be empty")
	}
	path := fmt.Sprintf("/v1/capabilities/%s", capabilityID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var capability CapabilityRepresentation
	if err := c.doRequest(req, &capability); err != nil {
		return nil, err
	}
	return &capability, nil
}

// UpdateCapability updates a specific capability by its ID.
// The payload should be either ChatCapabilityUpdate or CompletionCapabilityUpdate.
// Corresponds to PUT /v1/capabilities/{capability_id}
func (c *Client) UpdateCapability(ctx context.Context, capabilityID string, capabilityData interface{}) (*CapabilityRepresentation, error) {
	if strings.TrimSpace(capabilityID) == "" {
		return nil, fmt.Errorf("capabilityID cannot be empty")
	}
	path := fmt.Sprintf("/v1/capabilities/%s", capabilityID)
	req, err := c.newRequest(ctx, http.MethodPut, path, capabilityData)
	if err != nil {
		return nil, err
	}

	var updatedCapability CapabilityRepresentation
	if err := c.doRequest(req, &updatedCapability); err != nil {
		return nil, err
	}
	return &updatedCapability, nil
}

// DeleteCapability deletes a specific capability by its ID.
// Corresponds to DELETE /v1/capabilities/{capability_id}
// Expects a 204 No Content on success.
func (c *Client) DeleteCapability(ctx context.Context, capabilityID string) error {
	if strings.TrimSpace(capabilityID) == "" {
		return fmt.Errorf("capabilityID cannot be empty")
	}
	path := fmt.Sprintf("/v1/capabilities/%s", capabilityID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, nil) // No body expected on 204
}

// --- ModelDeployment Methods ---

// CreateModelDeployment creates a new model deployment.
// Corresponds to POST /v1/model-deployments
func (c *Client) CreateModelDeployment(ctx context.Context, deploymentData ModelDeploymentCreate) (*ModelDeployment, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/model-deployments", deploymentData)
	if err != nil {
		return nil, err
	}

	var createdDeployment ModelDeployment
	if err := c.doRequest(req, &createdDeployment); err != nil {
		return nil, err
	}
	return &createdDeployment, nil
}

// GetModelDeployment retrieves a specific model deployment by its ID.
// Corresponds to GET /v1/model-deployments/{deployment_id}
func (c *Client) GetModelDeployment(ctx context.Context, deploymentID string) (*ModelDeployment, error) {
	if strings.TrimSpace(deploymentID) == "" {
		return nil, fmt.Errorf("deploymentID cannot be empty")
	}
	path := fmt.Sprintf("/v1/model-deployments/%s", deploymentID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var deployment ModelDeployment
	if err := c.doRequest(req, &deployment); err != nil {
		return nil, err
	}
	return &deployment, nil
}

// UpdateModelDeployment updates a specific model deployment by its ID.
// Corresponds to PUT /v1/model-deployments/{deployment_id}
func (c *Client) UpdateModelDeployment(ctx context.Context, deploymentID string, deploymentData ModelDeploymentUpdate) (*ModelDeployment, error) {
	if strings.TrimSpace(deploymentID) == "" {
		return nil, fmt.Errorf("deploymentID cannot be empty")
	}
	path := fmt.Sprintf("/v1/model-deployments/%s", deploymentID)
	req, err := c.newRequest(ctx, http.MethodPut, path, deploymentData)
	if err != nil {
		return nil, err
	}

	var updatedDeployment ModelDeployment
	if err := c.doRequest(req, &updatedDeployment); err != nil {
		return nil, err
	}
	return &updatedDeployment, nil
}

// DeleteModelDeployment deletes a specific model deployment by its ID.
// Corresponds to DELETE /v1/model-deployments/{deployment_id}
// Expects a 204 No Content on success.
func (c *Client) DeleteModelDeployment(ctx context.Context, deploymentID string) error {
	if strings.TrimSpace(deploymentID) == "" {
		return fmt.Errorf("deploymentID cannot be empty")
	}
	path := fmt.Sprintf("/v1/model-deployments/%s", deploymentID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, nil) // No body expected on 204
}

// --- ModelProvider Methods ---

// CreateModelProvider creates a new model provider.
// Corresponds to POST /v1/model-providers
func (c *Client) CreateModelProvider(ctx context.Context, providerData ModelProviderCreate) (*ModelProvider, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/model-providers", providerData)
	if err != nil {
		return nil, err
	}

	var createdProvider ModelProvider
	if err := c.doRequest(req, &createdProvider); err != nil {
		return nil, err
	}
	return &createdProvider, nil
}

// GetModelProvider retrieves a specific model provider by its ID.
// Corresponds to GET /v1/model-providers/{provider_id}
func (c *Client) GetModelProvider(ctx context.Context, providerID string) (*ModelProvider, error) {
	if strings.TrimSpace(providerID) == "" {
		return nil, fmt.Errorf("providerID cannot be empty")
	}
	path := fmt.Sprintf("/v1/model-providers/%s", providerID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var provider ModelProvider
	if err := c.doRequest(req, &provider); err != nil {
		return nil, err
	}
	return &provider, nil
}

// UpdateModelProvider updates a specific model provider by its ID.
// Corresponds to PUT /v1/model-providers/{provider_id}
func (c *Client) UpdateModelProvider(ctx context.Context, providerID string, providerData ModelProviderUpdate) (*ModelProvider, error) {
	if strings.TrimSpace(providerID) == "" {
		return nil, fmt.Errorf("providerID cannot be empty")
	}
	path := fmt.Sprintf("/v1/model-providers/%s", providerID)
	req, err := c.newRequest(ctx, http.MethodPut, path, providerData)
	if err != nil {
		return nil, err
	}

	var updatedProvider ModelProvider
	if err := c.doRequest(req, &updatedProvider); err != nil {
		return nil, err
	}
	return &updatedProvider, nil
}

// DeleteModelProvider deletes a specific model provider by its ID.
// Corresponds to DELETE /v1/model-providers/{provider_id}
// Expects a 204 No Content on success.
func (c *Client) DeleteModelProvider(ctx context.Context, providerID string) error {
	if strings.TrimSpace(providerID) == "" {
		return fmt.Errorf("providerID cannot be empty")
	}
	path := fmt.Sprintf("/v1/model-providers/%s", providerID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, nil) // No body expected on 204
}

// --- CapabilityType Methods ---

// GetCapabilityType retrieves a specific capability type definition.
// Corresponds to GET /v1/capability-types/{capability_type}
func (c *Client) GetCapabilityType(ctx context.Context, capabilityType string) (*CapabilityTypeRepresentation, error) {
	if strings.TrimSpace(capabilityType) == "" {
		return nil, fmt.Errorf("capabilityType cannot be empty")
	}
	path := fmt.Sprintf("/v1/capability-types/%s", capabilityType)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var capTypeRep CapabilityTypeRepresentation
	if err := c.doRequest(req, &capTypeRep); err != nil {
		return nil, err
	}
	return &capTypeRep, nil
}

// SetCapabilityTypeDefaultModel sets the default model deployment for a capability type.
// Corresponds to PUT /v1/capability-types/{capability_type}
func (c *Client) SetCapabilityTypeDefaultModel(ctx context.Context, capabilityType string, data DefaultModelDeploymentUpdate) (*CapabilityTypeRepresentation, error) {
	if strings.TrimSpace(capabilityType) == "" {
		return nil, fmt.Errorf("capabilityType cannot be empty")
	}
	path := fmt.Sprintf("/v1/capability-types/%s", capabilityType)
	req, err := c.newRequest(ctx, http.MethodPut, path, data)
	if err != nil {
		return nil, err
	}

	var capTypeRep CapabilityTypeRepresentation
	if err := c.doRequest(req, &capTypeRep); err != nil {
		return nil, err
	}
	return &capTypeRep, nil
}

// ListCapabilityTypes retrieves all capability type definitions.
// Corresponds to GET /v1/capability-types
func (c *Client) ListCapabilityTypes(ctx context.Context) (*CapabilityTypesRepresentation, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/capability-types", nil)
	if err != nil {
		return nil, err
	}

	var capTypesRep CapabilityTypesRepresentation
	if err := c.doRequest(req, &capTypesRep); err != nil {
		return nil, err
	}
	return &capTypesRep, nil
}
