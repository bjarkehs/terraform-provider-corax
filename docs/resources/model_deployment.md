---
page_title: "corax_model_deployment Resource - terraform-provider-corax"
subcategory: ""
description: |-
  Manages a Corax Model Deployment. Model Deployments link a specific model configuration from a Model Provider to be usable for certain tasks (e.g., chat, completion).
---

# corax_model_deployment Resource

Manages a Corax Model Deployment. Model Deployments make a specific model from a Model Provider available for use with certain capabilities.

## Example Usage

```terraform
resource "corax_model_deployment" "my_gpt35_turbo" {
  name        = "GPT-3.5 Turbo Deployment"
  description = "Deployment for general chat and completion tasks using GPT-3.5 Turbo."
  provider_id = "uuid-of-your-corax-model-provider" # Replace with your Model Provider ID

  supported_tasks = ["chat", "completion"]

  configuration = {
    "model_name"  = "gpt-3.5-turbo"
    // Example for Azure OpenAI, other providers might have different keys:
    // "api_version" = "2023-05-15"
    // "deployment_name" = "my-azure-deployment-name"
  }

  is_active = true
}

resource "corax_model_deployment" "my_embedding_model" {
  name        = "Text Embedding Ada v2"
  description = "Deployment for text embedding tasks."
  provider_id = "uuid-of-your-corax-model-provider" # Replace with your Model Provider ID

  supported_tasks = ["embedding"]

  configuration = {
    "model_name" = "text-embedding-ada-002"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (String, Required) A user-defined name for the model deployment. Must be at least 1 character long.
- `provider_id` - (String, Required) The UUID of the [Model Provider](./model_provider.md) this deployment belongs to.
- `supported_tasks` - (List of String, Required) A list of tasks this model deployment supports. Allowed values are `chat`, `completion`, `embedding`.
- `configuration` - (Map of String to String, Required) Configuration key-value pairs specific to the model deployment. These are passed to the underlying model provider. Common keys include `model_name`. For Azure OpenAI, you might also need `api_version` and `deployment_name`. Refer to your Model Provider's requirements.
- `description` - (String, Optional) An optional description for the model deployment. Defaults to null if not provided.
- `is_active` - (Boolean, Optional) Indicates whether the model deployment is active and usable. Defaults to `true`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the model deployment (UUID).
- `created_at` - (String) Creation timestamp of the model deployment.
- `updated_at` - (String, Nullable) Last update timestamp of the model deployment.
- `created_by` - (String) User who created the model deployment.
- `updated_by` - (String, Nullable) User who last updated the model deployment.

## Import

Model Deployments can be imported using their ID:

```sh
terraform import corax_model_deployment.my_gpt35_turbo deployment_id_here
```
