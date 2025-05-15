# Corax API Key Resource (`corax_api_key`)

Manages a Corax API Key. API keys are used to authenticate requests to the Corax API.

## Example Usage

```hcl
provider "corax" {
  # Configure API endpoint and key, potentially via environment variables
  # CORAX_API_ENDPOINT and CORAX_API_KEY
}

resource "corax_api_key" "my_app_key" {
  name       = "my-application-key"
  expires_at = "2025-12-31T23:59:59Z" # RFC3339 format
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required|String) The name of the API key. This is a friendly name for identifying the key.
- `expires_at` - (Required|String) The expiration date and time for the API key, specified in RFC3339 format (e.g., `YYYY-MM-DDTHH:mm:ssZ`). After this time, the key will no longer be valid.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the API key.
- `key` - (String|Sensitive) The API key secret. This value is only available upon creation of the resource and will not be returned by subsequent reads. Store this value securely.
- `prefix` - (String) The prefix of the API key (e.g., `corax-sk-`).
- `created_by` - (String) The identifier of the user or entity that created the API key.
- `created_at` - (String) The timestamp (RFC3339 format) when the API key was created.
- `updated_at` - (String) The timestamp (RFC3339 format) when the API key was last updated. This may be null if the key has not been updated.
- `is_active` - (Boolean) Indicates whether the API key is currently active and can be used for authentication.
- `last_used_at` - (String) The timestamp (RFC3339 format) when the API key was last used. This may be null if the key has not been used yet.
- `usage_count` - (Number) The number of times the API key has been used.

## Import

Corax API Keys can be imported using their `id`. For example:

```shell
terraform import corax_api_key.my_app_key your_api_key_id
```

**Note:** When an API key is imported, the `key` attribute (the secret itself) cannot be retrieved from the API. Terraform will still manage the lifecycle of the imported API key, but the `key` attribute will be unknown in the state. If you need the key value, it must be obtained when the key is first created.
