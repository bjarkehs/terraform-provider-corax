---
page_title: "corax_model_provider Resource - terraform-provider-corax"
subcategory: ""
description: |-
  Manages a Corax Model Provider. Model Providers store configurations (like API keys and endpoints) for different LLM providers (e.g., Azure OpenAI, OpenAI, Bedrock).
---

# corax_model_provider Resource

Manages a Corax Model Provider. Model Providers are used to configure access to various underlying Large Language Model (LLM) providers such as Azure OpenAI, OpenAI, Anthropic, Bedrock, etc. This resource stores the necessary credentials and settings to connect to these external services.

## Example Usage

```terraform
resource "corax_model_provider" "my_azure_openai" {
  name          = "My Azure OpenAI Service"
  provider_type = "azure_openai" # This type must be supported by your Corax API instance

  configuration = {
    "api_key"      = "your-azure-openai-api-key"
    "api_endpoint" = "https://your-resource-name.openai.azure.com/"
    // Other provider-specific keys like "api_version" might be needed
    // "api_version" = "2023-07-01-preview"
  }
}

resource "corax_model_provider" "my_openai" {
  name          = "OpenAI GPT-4"
  provider_type = "openai" # This type must be supported by your Corax API instance

  configuration = {
    "api_key" = "sk-your-openai-api-key"
    // "organization_id" = "your-org-id" # Optional for OpenAI
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (String, Required) A user-defined name for this model provider configuration (e.g., "Primary Azure OpenAI", "OpenAI GPT-4").
- `provider_type` - (String, Required) The type of the model provider. This string must match one of the types supported by your Corax API instance (e.g., `azure_openai`, `openai`, `bedrock`, `anthropic_claude`). You can typically get a list of available types via the `/v1/model-provider-types` API endpoint.
- `configuration` - (Map of String to String, Required, Sensitive) A map of key-value pairs specific to the `provider_type`. This is where you'll put API keys, base URLs, and other necessary credentials or settings.
  - **Important:** Since this map often contains sensitive values like API keys, the entire map is marked as sensitive.
  - The exact keys required depend on the `provider_type`. Refer to the Corax API documentation or the `/v1/model-provider-types/{provider_type}/model-provider-configuration` endpoint for the expected configuration schema for a given type.
  - Common examples:
    - For `azure_openai`: `api_key`, `api_endpoint`, `api_version`.
    - For `openai`: `api_key`, `organization_id` (optional).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the model provider (UUID).
- `created_at` - (String) Creation timestamp of the model provider.
- `updated_at` - (String, Nullable) Last update timestamp of the model provider.
- `created_by` - (String) User who created the model provider.
- `updated_by` - (String, Nullable) User who last updated the model provider.

## Import

Model Providers can be imported using their ID:

```sh
terraform import corax_model_provider.my_azure_openai provider_id_here
```
