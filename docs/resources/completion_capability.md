---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "corax_completion_capability Resource - corax"
subcategory: ""
description: |-
  Manages a Corax Completion Capability. Completion capabilities define configurations for generating text completions, potentially with structured output.
---

# corax_completion_capability (Resource)

Manages a Corax Completion Capability. Completion capabilities define configurations for generating text completions, potentially with structured output.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `completion_prompt` (String) The main prompt for which a completion is generated. May include placeholders for variables.
- `name` (String) A user-defined name for the completion capability.
- `output_type` (String) Defines the expected output format. Must be either 'text' or 'schema'.
- `system_prompt` (String) The system prompt that provides context or instructions to the completion model.

### Optional

- `config` (Attributes) Configuration settings for the capability's behavior. (see [below for nested schema](#nestedatt--config))
- `is_public` (Boolean) Indicates whether the capability is publicly accessible. Defaults to false.
- `model_id` (String) The UUID of the model deployment to use for this capability. If not provided, a default model for 'completion' type may be used by the API.
- `project_id` (String) The UUID of the project this capability belongs to.
- `schema_def` (Dynamic) Defines the structure of the output when `output_type` is 'schema'. This can be an HCL map or a JSON string. Required if `output_type` is 'schema'.
- `variables` (Set of String) A set of variable names (strings) that can be interpolated into the `completion_prompt`. Order is not significant.

### Read-Only

- `id` (String) The unique identifier for the completion capability (UUID).
- `owner` (String) Owner of the capability.
- `type` (String) Type of the capability (should be 'completion').

<a id="nestedatt--config"></a>
### Nested Schema for `config`

Optional:

- `blob_config` (Attributes) Configuration for handling file uploads (blobs) if the capability supports it. (see [below for nested schema](#nestedatt--config--blob_config))
- `content_tracing` (Boolean) Whether content (prompts, completion data, variables) should be recorded in observability systems. Automatically set to false by the API for timed data retention.
- `data_retention` (Attributes) Defines how long execution input and output data should be kept. Configure with 'type' and optionally 'hours'. (see [below for nested schema](#nestedatt--config--data_retention))
- `temperature` (Number) Controls randomness in response generation (0.0 to 1.0). Higher values make output more random.

<a id="nestedatt--config--blob_config"></a>
### Nested Schema for `config.blob_config`

Optional:

- `allowed_mime_types` (List of String) List of allowed MIME types for uploaded blobs.
- `max_blobs` (Number) Maximum number of blobs that can be uploaded.
- `max_file_size_mb` (Number) Maximum file size in megabytes for uploaded blobs.


<a id="nestedatt--config--data_retention"></a>
### Nested Schema for `config.data_retention`

Required:

- `type` (String) Type of data retention. Must be 'timed' or 'infinite'.

Optional:

- `hours` (Number) Duration in hours to retain data. Required if type is 'timed'. Must not be set if type is 'infinite'. Minimum 1.
