---
page_title: "corax_capability_type_default_model Resource - terraform-provider-corax"
subcategory: ""
description: |-
  Manages the default Model Deployment for a specific Capability Type (e.g., 'chat', 'completion', 'embedding') within the Corax API.
---

# corax_capability_type_default_model Resource

Manages the default Model Deployment for a specific Capability Type (e.g., `chat`, `completion`, `embedding`). When a Capability is created without an explicit `model_id`, the Corax API will use the default model deployment associated with its type.

## Example Usage

```terraform
resource "corax_model_provider" "my_provider" {
  name          = "Test Azure OpenAI"
  provider_type = "azure_openai" # Ensure this type is supported
  configuration = {
    "api_key"      = "your-azure-api-key"
    "api_endpoint" = "https://your-instance.openai.azure.com/"
  }
}

resource "corax_model_deployment" "chat_model" {
  name            = "GPT-3.5 Turbo for Chat"
  provider_id     = corax_model_provider.my_provider.id
  supported_tasks = ["chat"]
  configuration = {
    "model_name"    = "gpt-35-turbo" # Or your Azure deployment name for this model
    "api_version"   = "2023-07-01-preview"
  }
}

resource "corax_model_deployment" "completion_model" {
  name            = "Davinci for Completions"
  provider_id     = corax_model_provider.my_provider.id
  supported_tasks = ["completion"]
  configuration = {
    "model_name"    = "text-davinci-003" # Or your Azure deployment name
    "api_version"   = "2023-07-01-preview"
  }
}

// Set the default model for 'chat' capabilities
resource "corax_capability_type_default_model" "chat_default" {
  capability_type              = "chat"
  default_model_deployment_id  = corax_model_deployment.chat_model.id
}

// Set the default model for 'completion' capabilities
resource "corax_capability_type_default_model" "completion_default" {
  capability_type              = "completion"
  default_model_deployment_id  = corax_model_deployment.completion_model.id
}
```

## Argument Reference

The following arguments are supported:

- `capability_type` - (String, Required, Forces new resource) The type of the capability for which to set the default model. Allowed values are `chat`, `completion`, `embedding`. Changing this attribute will result in the destruction and recreation of the resource, as it identifies a different capability type's default setting.
- `default_model_deployment_id` - (String, Required) The UUID of an existing [Model Deployment](./model_deployment.md) to set as the default for this `capability_type`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `name` - (String) The display name of the capability type (e.g., "Chat", "Completion"). This is read-only from the API.

**Note on ID:** The ID of this resource is the `capability_type` itself.

## Import

This resource can be imported using the `capability_type`:

```sh
terraform import corax_capability_type_default_model.chat_default chat
```
