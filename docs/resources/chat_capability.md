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

- `name` - (Required, String) A user-defined name for the chat capability. Must be at least 1 character long.
- `system_prompt` - (Required, String) The system prompt that guides the behavior of the chat model.
- `is_public` - (Optional, Boolean) Indicates whether the capability is publicly accessible. Defaults to `false`.
- `model_id` - (Optional, String) The UUID of the model deployment to use for this capability. If not provided, a default model for 'chat' type may be used by the API.
- `project_id` - (Optional, String) The UUID of the project this capability belongs to. If not provided, it might be associated with a default or no project.
- `config` - (Optional, Block) Configuration settings for the capability's behavior. See [Config Block](#config-block) below.

### Config Block

The `config` block supports the following:

- `temperature` - (Optional, Number) Controls randomness in response generation (0.0 to 1.0). Higher values make output more random.
- `content_tracing` - (Optional, Boolean) Whether content (prompts, completion data, variables) should be recorded in observability systems. Defaults to `true`. Automatically set to `false` by the API if `data_retention` is `timed`.
- `blob_config` - (Optional, Block) Configuration for handling file uploads (blobs) if the capability supports it. See [Blob Config Block](#blob-config-block) below.
- `data_retention` - (Optional, Block) Defines how long execution input and output data should be kept. Must define exactly one of `timed` or `infinite`. See [Data Retention Block](#data-retention-block) below.

### Blob Config Block

The `blob_config` block supports the following:

- `max_file_size_mb` - (Optional, Number) Maximum file size in megabytes for uploaded blobs. Defaults to API defined value (e.g., 20MB).
- `max_blobs` - (Optional, Number) Maximum number of blobs that can be uploaded. Defaults to API defined value (e.g., 10).
- `allowed_mime_types` - (Optional, List of String) List of allowed MIME types for uploaded blobs. Defaults to API defined value (e.g., `["image/png", "image/jpeg"]`).

### Data Retention Block

The `data_retention` block must configure exactly one of the following nested blocks:

- `timed` - (Block) Retain data for a specific duration.
  - `hours` - (Required, Number) Duration in hours to retain data (minimum 1).
- `infinite` - (Block) Retain data indefinitely.
  - `enabled` - (Required, Boolean) Set to `true` to enable infinite data retention. This field must be explicitly set to `true`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the chat capability (UUID).
- `type` - (String) The type of the capability, which will always be "chat".
- `created_by` - (String) User who created the capability.
- `updated_by` - (String) User who last updated the capability.
- `created_at` - (String) Creation timestamp.
- `updated_at` - (String) Last update timestamp.
- `archived_at` - (String) Archival timestamp, if applicable.
- `owner` - (String) Owner of the capability.

## Import

Chat capabilities can be imported using their ID:

```sh
terraform import corax_chat_capability.my_chat_capability capability_id_here
```
