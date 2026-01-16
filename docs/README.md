# Corax Terraform Provider

Terraform provider for managing AI/LLM capabilities via the Corax API.

## Provider Configuration

```hcl
provider "corax" {
  api_endpoint = "https://api.corax.ai"  # or CORAX_API_ENDPOINT env var
  api_key      = "your-api-key"          # or CORAX_API_KEY env var
}
```

## Resources

### corax_project
Organizes capabilities and resources.

```hcl
resource "corax_project" "example" {
  name        = "my-project"
  description = "Optional description"
  is_public   = false
}
```

### corax_model_provider
Configures LLM providers (azure_openai, openai, bedrock).

```hcl
resource "corax_model_provider" "azure" {
  name          = "azure-openai"
  provider_type = "azure_openai"
  configuration = {
    api_key      = var.azure_api_key
    api_endpoint = "https://my-instance.openai.azure.com"
  }
}
```

### corax_model_deployment
Links a model provider to specific tasks.

```hcl
resource "corax_model_deployment" "gpt4" {
  name            = "gpt-4"
  provider_id     = corax_model_provider.azure.id
  supported_tasks = ["chat", "completion"]
  is_active       = true
  configuration   = {
    model_name  = "gpt-4"
    api_version = "2024-02-15-preview"
  }
}
```

### corax_capability_type_default_model
Sets the default model for a capability type.

```hcl
resource "corax_capability_type_default_model" "chat_default" {
  capability_type             = "chat"  # chat, completion, or embedding
  default_model_deployment_id = corax_model_deployment.gpt4.id
}
```

### corax_chat_capability
Configures conversational AI endpoints.

```hcl
resource "corax_chat_capability" "support" {
  name          = "customer-support"
  project_id    = corax_project.example.id
  model_id      = corax_model_deployment.gpt4.id
  system_prompt = "You are a helpful assistant."
  is_public     = false

  config {
    temperature = 0.7
    data_retention {
      type  = "timed"
      hours = 720
    }
    blob_config {
      max_blobs          = 5
      max_file_size_mb   = 10
      allowed_mime_types = ["application/pdf", "image/png"]
    }
  }
}
```

### corax_completion_capability
Configures text completion with optional structured output.

```hcl
resource "corax_completion_capability" "analyzer" {
  name              = "sentiment-analyzer"
  semantic_id       = "analyze-sentiment"
  project_id        = corax_project.example.id
  model_id          = corax_model_deployment.gpt4.id
  system_prompt     = "You are a sentiment analysis system."
  completion_prompt = "Analyze: {{text}}"
  output_type       = "schema"  # or "text"
  variables         = ["text"]

  schema_def = jsonencode({
    type = "object"
    properties = {
      sentiment  = { type = "string", enum = ["positive", "negative", "neutral"] }
      confidence = { type = "number" }
    }
    required = ["sentiment", "confidence"]
  })

  config {
    temperature     = 0.3
    content_tracing = true
    custom_parameters = { max_tokens = 500 }
    data_retention { type = "infinite" }
  }
}
```

### corax_api_key
Manages API keys for accessing Corax.

```hcl
resource "corax_api_key" "app" {
  name       = "production-key"
  expires_at = "2026-12-31T23:59:59Z"
}

output "api_key" {
  value     = corax_api_key.app.key  # only available on creation
  sensitive = true
}
```

## Quick Reference

| Resource | Required Fields |
|----------|-----------------|
| `corax_project` | `name` |
| `corax_model_provider` | `name`, `provider_type`, `configuration` |
| `corax_model_deployment` | `name`, `provider_id`, `supported_tasks`, `configuration` |
| `corax_capability_type_default_model` | `capability_type`, `default_model_deployment_id` |
| `corax_chat_capability` | `name`, `system_prompt` |
| `corax_completion_capability` | `name`, `system_prompt`, `completion_prompt`, `output_type` |
| `corax_api_key` | `name`, `expires_at` |
