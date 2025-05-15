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

// --- Collection Methods ---

// CreateCollection creates a new collection.
// Corresponds to POST /v1/collections
func (c *Client) CreateCollection(ctx context.Context, collectionData CollectionCreate) (*Collection, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collections", collectionData)
	if err != nil {
		return nil, err
	}

	var createdCollection Collection
	if err := c.doRequest(req, &createdCollection); err != nil {
		return nil, err
	}
	return &createdCollection, nil
}

// GetCollection retrieves a specific collection by its ID.
// Corresponds to GET /v1/collections/{collection_id}
func (c *Client) GetCollection(ctx context.Context, collectionID string) (*Collection, error) {
	if strings.TrimSpace(collectionID) == "" {
		return nil, fmt.Errorf("collectionID cannot be empty")
	}
	path := fmt.Sprintf("/v1/collections/%s", collectionID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var collection Collection
	if err := c.doRequest(req, &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

// UpdateCollection updates a specific collection by its ID.
// Corresponds to PUT /v1/collections/{collection_id}
func (c *Client) UpdateCollection(ctx context.Context, collectionID string, collectionData CollectionUpdate) (*Collection, error) {
	if strings.TrimSpace(collectionID) == "" {
		return nil, fmt.Errorf("collectionID cannot be empty")
	}
	path := fmt.Sprintf("/v1/collections/%s", collectionID)
	req, err := c.newRequest(ctx, http.MethodPut, path, collectionData)
	if err != nil {
		return nil, err
	}

	var updatedCollection Collection
	if err := c.doRequest(req, &updatedCollection); err != nil {
		return nil, err
	}
	return &updatedCollection, nil
}

// DeleteCollection deletes a specific collection by its ID.
// Corresponds to DELETE /v1/collections/{collection_id}
// Expects a 204 No Content on success.
func (c *Client) DeleteCollection(ctx context.Context, collectionID string) error {
	if strings.TrimSpace(collectionID) == "" {
		return fmt.Errorf("collectionID cannot be empty")
	}
	path := fmt.Sprintf("/v1/collections/%s", collectionID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, nil) // No body expected on 204
}

// --- Document Methods ---

// UpsertDocument creates or updates a document within a collection.
// Corresponds to PUT /v1/collections/{collection_id}/documents/{document_id}
// The API returns the created/updated document.
func (c *Client) UpsertDocument(ctx context.Context, collectionID string, documentID string, documentData DocumentUpdate) (*Document, error) {
	if strings.TrimSpace(collectionID) == "" {
		return nil, fmt.Errorf("collectionID cannot be empty")
	}
	if strings.TrimSpace(documentID) == "" {
		return nil, fmt.Errorf("documentID cannot be empty")
	}

	path := fmt.Sprintf("/v1/collections/%s/documents/%s", collectionID, documentID)
	req, err := c.newRequest(ctx, http.MethodPut, path, documentData)
	if err != nil {
		return nil, err
	}

	var upsertedDocument Document
	if err := c.doRequest(req, &upsertedDocument); err != nil {
		return nil, err
	}
	// The API response for a document doesn't include collection_id, add it for context.
	upsertedDocument.CollectionID = collectionID
	return &upsertedDocument, nil
}

// GetDocument retrieves a specific document by its collection ID and document ID.
// Corresponds to GET /v1/collections/{collection_id}/documents/{document_id}
func (c *Client) GetDocument(ctx context.Context, collectionID string, documentID string) (*Document, error) {
	if strings.TrimSpace(collectionID) == "" {
		return nil, fmt.Errorf("collectionID cannot be empty")
	}
	if strings.TrimSpace(documentID) == "" {
		return nil, fmt.Errorf("documentID cannot be empty")
	}

	path := fmt.Sprintf("/v1/collections/%s/documents/%s", collectionID, documentID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var document Document
	if err := c.doRequest(req, &document); err != nil {
		return nil, err
	}
	// The API response for a document doesn't include collection_id, add it for context.
	document.CollectionID = collectionID
	return &document, nil
}

// DeleteDocument deletes a specific document by its collection ID and document ID.
// Corresponds to DELETE /v1/collections/{collection_id}/documents/{document_id}
// Expects a 204 No Content on success.
func (c *Client) DeleteDocument(ctx context.Context, collectionID string, documentID string) error {
	if strings.TrimSpace(collectionID) == "" {
		return fmt.Errorf("collectionID cannot be empty")
	}
	if strings.TrimSpace(documentID) == "" {
		return fmt.Errorf("documentID cannot be empty")
	}

	path := fmt.Sprintf("/v1/collections/%s/documents/%s", collectionID, documentID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, nil) // No body expected on 204
}

// --- Embeddings Model Methods ---

// CreateEmbeddingsModel creates a new embeddings model.
// Corresponds to POST /v1/embeddings-models
func (c *Client) CreateEmbeddingsModel(ctx context.Context, modelData EmbeddingsModelCreate) (*EmbeddingsModel, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/embeddings-models", modelData)
	if err != nil {
		return nil, err
	}

	var createdModel EmbeddingsModel
	if err := c.doRequest(req, &createdModel); err != nil {
		return nil, err
	}
	return &createdModel, nil
}

// GetEmbeddingsModel retrieves a specific embeddings model by its ID.
// Corresponds to GET /v1/embeddings-models/{model_id}
func (c *Client) GetEmbeddingsModel(ctx context.Context, modelID string) (*EmbeddingsModel, error) {
	if strings.TrimSpace(modelID) == "" {
		return nil, fmt.Errorf("modelID cannot be empty")
	}
	path := fmt.Sprintf("/v1/embeddings-models/%s", modelID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var model EmbeddingsModel
	if err := c.doRequest(req, &model); err != nil {
		return nil, err
	}
	return &model, nil
}

// UpdateEmbeddingsModel updates a specific embeddings model by its ID.
// Corresponds to PUT /v1/embeddings-models/{model_id}
func (c *Client) UpdateEmbeddingsModel(ctx context.Context, modelID string, modelData EmbeddingsModelUpdate) (*EmbeddingsModel, error) {
	if strings.TrimSpace(modelID) == "" {
		return nil, fmt.Errorf("modelID cannot be empty")
	}
	path := fmt.Sprintf("/v1/embeddings-models/%s", modelID)
	req, err := c.newRequest(ctx, http.MethodPut, path, modelData)
	if err != nil {
		return nil, err
	}

	var updatedModel EmbeddingsModel
	if err := c.doRequest(req, &updatedModel); err != nil {
		return nil, err
	}
	return &updatedModel, nil
}

// DeleteEmbeddingsModel deletes a specific embeddings model by its ID.
// Corresponds to DELETE /v1/embeddings-models/{model_id}
// Expects a 204 No Content on success.
func (c *Client) DeleteEmbeddingsModel(ctx context.Context, modelID string) error {
	if strings.TrimSpace(modelID) == "" {
		return fmt.Errorf("modelID cannot be empty")
	}
	path := fmt.Sprintf("/v1/embeddings-models/%s", modelID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.doRequest(req, nil) // No body expected on 204
}
