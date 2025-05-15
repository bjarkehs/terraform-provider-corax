# Corax Collection Resource (`corax_collection`)

Manages a Corax Knowledge Collection. Collections are used to store and organize documents along with their vector embeddings, enabling semantic search and other AI-powered functionalities. Each collection belongs to a project.

## Example Usage

```hcl
provider "corax" {
  # Configure API endpoint and key
}

resource "corax_project" "main_project" {
  name        = "Main Research Project"
  description = "Primary project for all research activities."
}

resource "corax_collection" "research_papers_q1" {
  project_id  = corax_project.main_project.id
  name        = "Research Papers Q1 2024"
  description = "Collection of all research papers published or analyzed in Q1 2024."
  # embeddings_model_id = "your-embeddings-model-uuid" # Optional: Specify a custom embeddings model

  metadata_schema = {
    "author"       = "string"
    "publish_date" = "string" # Consider a date type if API supports richer types
    "journal"      = "string"
    "keywords"     = "array"  # Assuming 'array' implies an array of strings or simple types
    "page_count"   = "number"
  }
}

resource "corax_collection" "product_docs" {
  project_id = corax_project.main_project.id
  name       = "Product Documentation"
  # No description
  # No custom embeddings model (uses default)
  # No metadata schema (or a very simple one if needed)
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required|String) The name of the collection. Must be at least 1 character long.
- `project_id` - (Required|String) The UUID of the project this collection belongs to.
- `description` - (Optional|String) An optional description for the collection.
- `embeddings_model_id` - (Optional|String) The UUID of a specific embeddings model to use for this collection. If not provided, the API's default embeddings model will be used.
- `metadata_schema` - (Optional|Map of String to String) A map defining the schema for custom metadata that can be attached to documents within this collection.
  - Keys are the names of the metadata properties.
  - Values are strings indicating the data type of the property. Supported types are: `"string"`, `"number"`, `"boolean"`, `"array"`, `"object"`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier (UUID) for the collection.
- `created_by` - (String) The identifier of the user or entity that created the collection.
- `updated_by` - (String) The identifier of the user or entity that last updated the collection. This may be null.
- `created_at` - (String) The timestamp (RFC3339 format) when the collection was created.
- `updated_at` - (String) The timestamp (RFC3339 format) when the collection was last updated. This may be null.
- `document_count` - (Number) The number of documents currently stored in the collection.
- `size_bytes` - (Number) The total size (in bytes) of the documents and their embeddings stored in the collection.
- `status` - (String) The current operational status of the collection (e.g., `ready`, `indexing`, `failed`).

## Import

Corax Collections can be imported using their `id`. For example:

```shell
terraform import corax_collection.research_papers_q1 your_collection_id
```
