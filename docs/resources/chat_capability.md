---
page_title: "corax_chat_capability Resource - terraform-provider-corax"
subcategory: ""
description: |-
  Manages a Corax Chat Capability. Chat capabilities define configurations for conversational AI models.
---

# corax_chat_capability Resource

Manages a Corax Chat Capability. Chat capabilities define configurations for conversational AI models.

## Example Usage

### Basic Chat Capability

```terraform
resource "corax_chat_capability" "basic_chat" {
  name          = "My Basic Assistant"
  system_prompt = "You are a friendly and helpful assistant."
}
```

### Chat Capability with Configuration

```terraform
resource "corax_chat_capability" "configured_chat" {
  name          = "My Configured Assistant"
  system_prompt = "You are a specialized assistant for coding questions."
  is_public     = true
  model_id      = "uuid-of-a-model-deployment" # Optional: Specify a model deployment
  project_id    = "uuid-of-a-project"        # Optional: Link to a project

  config {
    temperature     = 0.7
    content_tracing = true

    data_retention {
      timed {
        hours = 72 # Retain data for 72 hours
      }
      # Alternatively, for infinite retention:
      # infinite {
      #   enabled = true
      # }
    }

    blob_config {
      max_file_size_mb   = 15
      max_blobs          = 3
      allowed_mime_types = ["image/png", "image/jpeg", "application/pdf"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (String, Required) A user-defined name for the chat capability. Must be at least 1 character long.
- `system_prompt` - (String, Required) The system prompt that guides the behavior of the chat model.
- `is_public` - (Boolean, Optional) Indicates whether the capability is publicly accessible. Defaults to `false`.
- `model_id` - (String, Optional) The UUID of an existing [Model Deployment](./model_deployment.md) to use for this capability. If not provided, a default model for the 'chat' type may be used by the API.
- `project_id` - (String, Optional) The UUID of an existing [Project](./project.md) this capability belongs to.
- `config` - (Block, Optional) Configuration settings for the capability's behavior. See [Config Block](#config-block) below.

### Config Block

The `config` block supports the following:

- `temperature` - (Number, Optional) Controls randomness in response generation (0.0 to 1.0). Higher values make output more random.
- `content_tracing` - (Boolean, Optional) Whether content (prompts, completion data, variables) should be recorded in observability systems. Defaults to `true`. Automatically set to `false` by the API if `data_retention` is `timed`.
- `blob_config` - (Block, Optional) Configuration for handling file uploads (blobs) if the capability supports it. See [Blob Config Block](#blob-config-block) below.
- `data_retention` - (Block, Optional) Defines how long execution input and output data should be kept. Must define exactly one of `timed` or `infinite`. See [Data Retention Block](#data-retention-block) below.

### Blob Config Block

The `blob_config` block supports the following:

- `max_file_size_mb` - (Number, Optional) Maximum file size in megabytes for uploaded blobs. If not set, defaults to the API's defined value (e.g., 20MB).
- `max_blobs` - (Number, Optional) Maximum number of blobs that can be uploaded. If not set, defaults to the API's defined value (e.g., 10).
- `allowed_mime_types` - (List of String, Optional) List of allowed MIME types for uploaded blobs. If not set, defaults to the API's defined value (e.g., `["image/png", "image/jpeg"]`).

### Data Retention Block

The `data_retention` block must configure exactly one of the following nested blocks:

- `timed` - (Block, Optional) Retain data for a specific duration.
  - `hours` - (Number, Required) Duration in hours to retain data (minimum 1).
- `infinite` - (Block, Optional) Retain data indefinitely.
  - `enabled` - (Boolean, Required) Set to `true` to enable infinite data retention. This field must be explicitly set to `true`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the chat capability (UUID).
- `created_by` - (String) User who created the capability.
- `updated_by` - (String) User who last updated the capability.
- `created_at` - (String) Creation timestamp.
- `updated_at` - (String) Last update timestamp.
- `archived_at` - (String, Nullable) Archival timestamp, if applicable.

## Import

Chat capabilities can be imported using their ID:

```sh
terraform import corax_chat_capability.my_chat_capability capability_id_here
```
