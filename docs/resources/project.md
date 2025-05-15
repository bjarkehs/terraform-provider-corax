# Corax Project Resource (`corax_project`)

Manages a Corax Project. Projects are used as containers to organize knowledge collections and capabilities within the Corax platform.

## Example Usage

```hcl
provider "corax" {
  # Configure API endpoint and key, potentially via environment variables
  # CORAX_API_ENDPOINT and CORAX_API_KEY
}

resource "corax_project" "my_research_project" {
  name        = "My Research Project Q3"
  description = "Project to store all research data and models for Q3."
  is_public   = false
}

resource "corax_project" "public_demos" {
  name        = "Public Demos"
  description = "A project for publicly accessible demonstrations."
  is_public   = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required|String) The name of the project. Must be at least 1 character long.
- `description` - (Optional|String) An optional description for the project.
- `is_public` - (Optional|Boolean) Indicates whether the project is public. If not specified, defaults to `false` (private).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier (UUID) for the project.
- `created_by` - (String) The identifier of the user or entity that created the project.
- `updated_by` - (String) The identifier of the user or entity that last updated the project. This may be null.
- `created_at` - (String) The timestamp (RFC3339 format) when the project was created.
- `updated_at` - (String) The timestamp (RFC3339 format) when the project was last updated. This may be null if the project has not been updated.
- `owner` - (String) The owner of the project.
- `collection_count` - (Number) The number of knowledge collections currently associated with this project.
- `capability_count` - (Number) The number of capabilities (e.g., models, functions) currently associated with this project.

## Import

Corax Projects can be imported using their `id`. For example:

```shell
terraform import corax_project.my_research_project your_project_id
```
