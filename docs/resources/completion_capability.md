---
page_title: "corax_completion_capability Resource - terraform-provider-corax"
subcategory: ""
description: |-
  Manages a Corax Completion Capability. Completion capabilities define configurations for generating text completions, potentially with structured output.
---

# corax_completion_capability Resource

Manages a Corax Completion Capability. Completion capabilities define configurations for generating text completions, potentially with structured output.

## Example Usage

### Basic Text Completion

```terraform
resource "corax_completion_capability" "basic_text" {
  name               = "My Basic Text Completion"
  system_prompt      = "You are a helpful text completion model."
  completion_prompt  = "The quick brown fox jumps over the lazy "
  output_type        = "text"
}
```

### Completion with Structured Schema Output

```terraform
resource "corax_completion_capability" "structured_data_extractor" {
  name               = "User Data Extractor"
  system_prompt      = "Extract user information from the text and provide it in the specified JSON schema."
  completion_prompt  = "Process the following user data: {{UserData}}"
  variables          = ["UserData"]
  output_type        = "schema"

  schema_def = {
    name = jsonencode({
      type        = "string"
      description = "The full name of the user."
    })
    email = jsonencode({
      type        = "string"
      description = "The email address of the user."
    })
    age = jsonencode({
      type        = "integer"
      description = "The age of the user."
    })
    address = jsonencode({
      type = "object"
      description = "The user's address."
      properties = {
        street = { type = "string", description = "Street name and number." }
        city   = { type = "string", description = "City." }
        zip    = { type = "string", description = "Zip code." }
      }
    })
  }

  config {
    temperature = 0.2
  }
}
```

**Note on `schema_def`:** The values within the `schema_def` map should be JSON encoded strings representing the schema for each property. The provider will interpret these JSON strings.

## Argument Reference

The following arguments are supported:

- `name` - (String, Required) A user-defined name for the completion capability. Must be at least 1 character long.
- `system_prompt` - (String, Required) The system prompt that provides context or instructions to the completion model.
- `completion_prompt` - (String, Required) The main prompt for which a completion is generated. May include placeholders for variables (e.g., `{{MyVariable}}`).
- `output_type` - (String, Required) Defines the expected output format. Must be either `text` or `schema`.
- `is_public` - (Boolean, Optional) Indicates whether the capability is publicly accessible. Defaults to `false`.
- `model_id` - (String, Optional) The UUID of an existing [Model Deployment](./model_deployment.md) to use for this capability. If not provided, a default model for the 'completion' type may be used by the API.
- `project_id` - (String, Optional) The UUID of an existing [Project](./project.md) this capability belongs to.
- `variables` - (List of String, Optional) A list of variable names (strings) that can be interpolated into the `completion_prompt`.
- `schema_def` - (Map of String to Dynamic, Optional) Defines the structure of the output when `output_type` is `schema`. This is a map where keys are property names and values are **JSON encoded strings** defining the property's type, description, and potentially nested properties (for object types) or item types (for array types). This argument is **required** if `output_type` is `schema`.
- `config` - (Block, Optional) Configuration settings for the capability's behavior. See [Config Block](#config-block) below. (Identical to the config block in `corax_chat_capability`).

### Config Block

The `config` block supports the following (refer to [`corax_chat_capability`](./chat_capability.md#config-block) documentation for details on sub-blocks like `blob_config` and `data_retention`):

- `temperature` - (Number, Optional)
- `content_tracing` - (Boolean, Optional)
- `blob_config` - (Block, Optional)
- `data_retention` - (Block, Optional)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the completion capability (UUID).
- `type` - (String) The type of the capability, which will always be "completion".
- `created_by` - (String) User who created the capability.
- `updated_by` - (String) User who last updated the capability.
- `created_at` - (String) Creation timestamp.
- `updated_at` - (String) Last update timestamp.
- `archived_at` - (String, Nullable) Archival timestamp, if applicable.
- `owner` - (String) Owner of the capability.

## Import

Completion capabilities can be imported using their ID:

```sh
terraform import corax_completion_capability.my_completion_capability capability_id_here
```
