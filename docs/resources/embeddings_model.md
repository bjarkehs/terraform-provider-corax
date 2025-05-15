# Corax Embeddings Model Resource (`corax_embeddings_model`)

Manages a Corax Embeddings Model configuration. Embeddings models are responsible for converting text or other data into numerical vector representations (embeddings), which are then used for semantic search, clustering, and other AI-driven tasks within Corax.

You can configure models from various providers (like OpenAI, Cohere, Sentence-Transformers) or specify custom/self-hosted models.

## Example Usage

### Custom Model (e.g., self-hosted Sentence Transformer)

```hcl
provider "corax" {
  # Configure API endpoint and key
}

resource "corax_embeddings_model" "my_mini_lm" {
  name                   = "My Local MiniLM-L6-v2"
  description            = "A self-hosted instance of all-MiniLM-L6-v2."
  model_provider         = "custom" // Or "sentence-transformers" if API has specific handling
  model_name_on_provider = "all-MiniLM-L6-v2" // Model identifier for your custom setup
  dimensions             = 384
  max_tokens             = 256 // Max input tokens for this specific model deployment
  api_base_url           = "http://my-embeddings-service.internal:8000/embed"
  # api_key is not needed if the custom service doesn't require it
}
```

### OpenAI Model (e.g., text-embedding-ada-002)

```hcl
resource "corax_embeddings_model" "openai_ada" {
  name                   = "OpenAI Ada v2"
  description            = "OpenAI's text-embedding-ada-002 model."
  model_provider         = "openai"
  model_name_on_provider = "text-embedding-ada-002"
  dimensions             = 1536
  # max_tokens might be defaulted by the API for known OpenAI models
  api_key                = var.openai_api_key # Sensitive, use a variable
}

variable "openai_api_key" {
  type        = string
  description = "API key for OpenAI."
  sensitive   = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required|String) A user-defined name for this embeddings model configuration.
- `model_provider` - (Required|String|ForceNew) The provider of the embeddings model. Examples: `openai`, `cohere`, `sentence-transformers`, `custom`. Changing this forces a new resource.
- `model_name_on_provider` - (Required|String|ForceNew) The specific model name as recognized by the provider (e.g., `text-embedding-ada-002`, `embed-english-v2.0`, `all-MiniLM-L6-v2`). Changing this forces a new resource.
- `dimensions` - (Required|Number|ForceNew) The number of dimensions (the length of the vector) that the embeddings model outputs. Changing this forces a new resource.
- `description` - (Optional|String) An optional description for the embeddings model configuration.
- `max_tokens` - (Optional|Number|ForceNew) The maximum number of input tokens the model can handle in a single request. If not provided, the API might use a default or derive it based on the model. Changing this forces a new resource.
- `api_key` - (Optional|String|Sensitive) The API key required for accessing the model if it's a third-party proprietary model (e.g., OpenAI API key). This is a sensitive field.
- `api_base_url` - (Optional|String) The base URL for the API of a custom or self-hosted embeddings model. This is typically used when `model_provider` is `custom` or a similar type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier (UUID) for the embeddings model configuration.
- `status` - (String) The operational status of the embeddings model configuration (e.g., `active`, `error`, `pending_validation`). This status might reflect whether Corax can successfully communicate with or use the configured model.
- `is_default` - (Boolean) Indicates if this is the default embeddings model for the Corax instance.
- `created_by` - (String) The identifier of the user or entity that created the embeddings model configuration.
- `updated_by` - (String) The identifier of the user or entity that last updated the configuration. This may be null.
- `created_at` - (String) The timestamp (RFC3339 format) when the configuration was created.
- `updated_at` - (String) The timestamp (RFC3339 format) when the configuration was last updated. This may be null.

## Import

Corax Embeddings Model configurations can be imported using their `id`.

```shell
terraform import corax_embeddings_model.my_mini_lm your_embeddings_model_id
```
