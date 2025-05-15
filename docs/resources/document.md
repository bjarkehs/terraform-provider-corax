# Corax Document Resource (`corax_document`)

Manages a Corax Document within a specific Corax Collection. Documents are the fundamental units of information storage and typically contain text or structured JSON data, along with associated metadata. Their content is processed for embeddings to enable semantic search and other AI functionalities.

## Example Usage

```hcl
provider "corax" {
  # Configure API endpoint and key
}

resource "corax_project" "main_project" {
  name = "My Main Project"
}

resource "corax_collection" "product_specs" {
  project_id = corax_project.main_project.id
  name       = "Product Specifications"
  metadata_schema = {
    "version" = "number"
    "status"  = "string"
  }
}

// Document with text content (API generates ID)
resource "corax_document" "spec_v1_intro" {
  collection_id = corax_collection.product_specs.id
  text_content  = "This document contains the introduction to version 1 of the product specification."
  metadata = {
    "version" = 1
    "status"  = "draft"
    "author"  = "jane.doe@example.com"
  }
}

// Document with JSON content and user-provided ID
resource "corax_document" "spec_v1_features" {
  collection_id = corax_collection.product_specs.id
  id            = "spec-v1-features-doc-uuid" // Optional: user-defined UUID or other unique ID string
  json_content = jsonencode({
    title    = "Product Features - V1"
    sections = [
      {
        heading = "Core Functionality"
        details = "Describes the main features..."
      },
      {
        heading = "Advanced Options"
        details = "Details on advanced configuration..."
      }
    ]
  })
  metadata = {
    "version" = 1
    "status"  = "final"
    "reviewed_by" = ["john.doe@example.com", "alice.smith@example.com"]
  }
}
```

## Argument Reference

The following arguments are supported:

- `collection_id` - (Required|String) The UUID of the Corax Collection this document belongs to. Changing this forces a new resource.
- `id` - (Optional|String|ForceNew) The unique identifier (e.g., UUID) for the document. If not provided, the API will generate one. If you provide an ID and then change it in the configuration, it will force the creation of a new document resource.
- `text_content` - (Optional|String) The plain text content of the document. You must provide either `text_content` or `json_content`.
- `json_content` - (Optional|String) The JSON content of the document, provided as a valid JSON string. You must provide either `text_content` or `json_content`. Use `jsonencode()` for complex objects.
- `metadata` - (Optional|Map of Dynamic) A map of metadata to associate with the document. Keys must be strings. Values can be strings, numbers, booleans, or JSON-encodable nested maps/lists. This metadata should ideally conform to the `metadata_schema` defined on the parent collection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The unique identifier for the document.
- `content` - (String) The raw content of the document as returned by the API. This can be a string (if `text_content` was primary) or a JSON string representation of an object (if `json_content` was primary).
- `token_count` - (Number) The number of tokens calculated for the document's content.
- `chunk_count` - (Number) The number of chunks the document was divided into for the purpose of generating embeddings.
- `embeddings_status` - (String) The current status of the embeddings generation for this document (e.g., `completed`, `pending`, `failed`).
- `created_by` - (String) The identifier of the user or entity that created the document.
- `updated_by` - (String) The identifier of the user or entity that last updated the document. This may be null.
- `created_at` - (String) The timestamp (RFC3339 format) when the document was created.
- `updated_at` - (String) The timestamp (RFC3339 format) when the document was last updated. This may be null.

## Import

Corax Documents can be imported using a composite ID of the format `collection_id/document_id`.

```shell
terraform import corax_document.spec_v1_intro your_collection_id/your_document_id
```
