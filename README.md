# Terraform Provider for Corax

A Terraform provider for managing [Corax](https://corax.ai) resources. Corax is an AI/LLM capability management platform that provides a unified interface for configuring and deploying language model capabilities.

## Features

- **Projects** - Organize capabilities and resources
- **Model Providers** - Configure LLM providers (Azure OpenAI, OpenAI, Bedrock)
- **Model Deployments** - Link model configurations to specific tasks
- **Chat Capabilities** - Configure conversational AI endpoints
- **Completion Capabilities** - Configure text completion with optional structured output
- **API Keys** - Manage API keys for accessing Corax services

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (only for development)

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    corax = {
      source = "registry.terraform.io/trifork/corax"
    }
  }
}
```

## Authentication

The provider requires an API endpoint and API key to authenticate with the Corax API. You can configure these in two ways:

### Option 1: Environment Variables (Recommended)

```bash
export CORAX_API_ENDPOINT="https://api.corax.ai"
export CORAX_API_KEY="your-api-key"
```

Then configure the provider without explicit credentials:

```hcl
provider "corax" {}
```

### Option 2: Provider Configuration

```hcl
provider "corax" {
  api_endpoint = "https://api.corax.ai"
  api_key      = var.corax_api_key  # Use a variable to avoid hardcoding secrets
}
```

## Usage

### Complete Example: Azure OpenAI Setup

This example demonstrates a complete setup with Azure OpenAI, including a project, model provider, model deployment, and both chat and completion capabilities.

```hcl
terraform {
  required_providers {
    corax = {
      source = "registry.terraform.io/trifork/corax"
    }
  }
}

provider "corax" {}

# Create a project to organize resources
resource "corax_project" "my_project" {
  name        = "my-ai-project"
  description = "Production AI capabilities"
  is_public   = false
}

# Configure an Azure OpenAI provider
resource "corax_model_provider" "azure_openai" {
  name          = "azure-openai-prod"
  provider_type = "azure_openai"

  configuration = {
    api_key      = var.azure_openai_api_key
    api_endpoint = "https://my-instance.openai.azure.com"
  }
}

# Create a model deployment for GPT-4
resource "corax_model_deployment" "gpt4" {
  name        = "gpt-4-deployment"
  description = "GPT-4 model for chat and completion tasks"
  provider_id = corax_model_provider.azure_openai.id
  is_active   = true

  supported_tasks = ["chat", "completion"]

  configuration = {
    model_name  = "gpt-4"
    api_version = "2024-02-15-preview"
  }
}

# Set GPT-4 as the default model for chat capabilities
resource "corax_capability_type_default_model" "chat_default" {
  capability_type             = "chat"
  default_model_deployment_id = corax_model_deployment.gpt4.id
}

# Create a chat capability
resource "corax_chat_capability" "customer_support" {
  name          = "customer-support-chat"
  project_id    = corax_project.my_project.id
  model_id      = corax_model_deployment.gpt4.id
  system_prompt = "You are a helpful customer support assistant. Be friendly, professional, and concise."
  is_public     = false

  config {
    temperature = 0.7

    data_retention {
      type  = "timed"
      hours = 720  # 30 days
    }
  }
}

# Create a completion capability with structured output
resource "corax_completion_capability" "sentiment_analyzer" {
  name              = "sentiment-analyzer"
  semantic_id       = "analyze-sentiment"
  project_id        = corax_project.my_project.id
  model_id          = corax_model_deployment.gpt4.id
  system_prompt     = "You are a sentiment analysis system. Analyze the given text and return structured sentiment data."
  completion_prompt = "Analyze the sentiment of the following text: {{text}}"
  output_type       = "schema"
  variables         = ["text"]
  is_public         = false

  schema_def = jsonencode({
    type = "object"
    properties = {
      sentiment = {
        type = "string"
        enum = ["positive", "negative", "neutral"]
      }
      confidence = {
        type    = "number"
        minimum = 0
        maximum = 1
      }
      explanation = {
        type = "string"
      }
    }
    required = ["sentiment", "confidence"]
  })

  config {
    temperature     = 0.3
    content_tracing = true

    data_retention {
      type = "infinite"
    }
  }
}

# Create an API key for accessing the capabilities
resource "corax_api_key" "app_key" {
  name       = "production-app-key"
  expires_at = "2026-12-31T23:59:59Z"
}

# Output the API key (only available on creation)
output "api_key" {
  value     = corax_api_key.app_key.key
  sensitive = true
}
```

### Example: OpenAI Setup

```hcl
resource "corax_model_provider" "openai" {
  name          = "openai-prod"
  provider_type = "openai"

  configuration = {
    api_key = var.openai_api_key
  }
}

resource "corax_model_deployment" "gpt4_turbo" {
  name        = "gpt-4-turbo"
  provider_id = corax_model_provider.openai.id
  is_active   = true

  supported_tasks = ["chat", "completion"]

  configuration = {
    model_name = "gpt-4-turbo"
  }
}
```

### Example: Chat Capability with File Uploads

```hcl
resource "corax_chat_capability" "document_analyzer" {
  name          = "document-analyzer"
  project_id    = corax_project.my_project.id
  model_id      = corax_model_deployment.gpt4.id
  system_prompt = "You are a document analysis assistant. Analyze uploaded documents and answer questions about them."

  config {
    temperature = 0.5

    blob_config {
      max_blobs        = 5
      max_file_size_mb = 10
      allowed_mime_types = [
        "application/pdf",
        "text/plain",
        "image/png",
        "image/jpeg"
      ]
    }

    data_retention {
      type  = "timed"
      hours = 168  # 7 days
    }
  }
}
```

### Example: Completion Capability with Custom Parameters

```hcl
resource "corax_completion_capability" "summarizer" {
  name              = "text-summarizer"
  project_id        = corax_project.my_project.id
  system_prompt     = "You are a text summarization system."
  completion_prompt = "Summarize the following {{content_type}} in {{length}} sentences: {{content}}"
  output_type       = "text"
  variables         = ["content", "content_type", "length"]

  config {
    temperature = 0.4

    custom_parameters = {
      max_tokens = 500
      top_p      = 0.9
    }

    data_retention {
      type = "infinite"
    }
  }
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `corax_project` | Manages projects for organizing capabilities |
| `corax_model_provider` | Configures LLM providers (Azure OpenAI, OpenAI, Bedrock) |
| `corax_model_deployment` | Links model configurations to specific tasks |
| `corax_chat_capability` | Configures conversational AI capabilities |
| `corax_completion_capability` | Configures text completion capabilities |
| `corax_capability_type_default_model` | Sets default models for capability types |
| `corax_api_key` | Manages API keys for Corax access |

For detailed schema information, see the [documentation](./docs/).

## Developing the Provider

### Building

```shell
go install
```

### Generating Documentation

```shell
make generate
```

### Running Tests

Acceptance tests create real resources and may incur costs.

```shell
make testacc
```

### Adding Dependencies

```shell
go get github.com/author/dependency
go mod tidy
```

## License

See [LICENSE](./LICENSE) for details.
