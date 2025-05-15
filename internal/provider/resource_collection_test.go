package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

)

const (
	testAccCollectionResourcePrefix = "tf-acc-test-collection-"
)

// TestAccCollectionResource provides acceptance tests for the corax_collection resource.
func TestAccCollectionResource(t *testing.T) {
	if os.Getenv("CORAX_API_KEY") == "" || os.Getenv("CORAX_API_ENDPOINT") == "" {
		t.Skip("CORAX_API_KEY and CORAX_API_ENDPOINT must be set for acceptance tests")
		return
	}

	// Pre-requisite: Project
	projectRName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectName := fmt.Sprintf("%s%s-for-collection", testAccProjectResourcePrefix, projectRName)
	projectResourceFullName := "corax_project.test_for_collection"

	// Collection details
	collectionRName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	collectionName := fmt.Sprintf("%s%s", testAccCollectionResourcePrefix, collectionRName)
	collectionNameUpdated := fmt.Sprintf("%s-updated", collectionName)
	collectionDesc := "Test collection description"
	collectionDescUpdated := "Test collection description updated"
	collectionResourceFullName := "corax_collection.test"

	// Embeddings Model ID (assuming a valid one exists or can be faked for tests if API allows)
	// For now, we'll test with and without it. If a real one is needed, this needs adjustment.
	// testEmbeddingsModelID := "your-test-embeddings-model-id" // Replace if needed

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// --- Step 1: Create Collection (minimal) ---
			{
				Config: testAccCollectionResourceConfigBasic(projectName, collectionName, collectionDesc),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Project checks
					resource.TestCheckResourceAttr(projectResourceFullName, "name", projectName),
					// Collection checks
					resource.TestCheckResourceAttr(collectionResourceFullName, "name", collectionName),
					resource.TestCheckResourceAttr(collectionResourceFullName, "description", collectionDesc),
					resource.TestCheckResourceAttrPair(collectionResourceFullName, "project_id", projectResourceFullName, "id"),
					resource.TestCheckResourceAttrSet(collectionResourceFullName, "id"),
					resource.TestMatchResourceAttr(collectionResourceFullName, "id", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					resource.TestCheckResourceAttrSet(collectionResourceFullName, "created_at"),
					resource.TestCheckResourceAttrSet(collectionResourceFullName, "created_by"),
					resource.TestCheckResourceAttr(collectionResourceFullName, "document_count", "0"),
					resource.TestCheckResourceAttr(collectionResourceFullName, "size_bytes", "0"),
					resource.TestCheckResourceAttrSet(collectionResourceFullName, "status"), // e.g., "ready" or "creating" initially
					resource.TestCheckResourceAttr(collectionResourceFullName, "embeddings_model_id", ""), // Expect null/empty if not set
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.%", "0"), // Expect empty map
				),
			},
			// --- Step 2: Update Collection (name, description) ---
			{
				Config: testAccCollectionResourceConfigBasic(projectName, collectionNameUpdated, collectionDescUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collectionResourceFullName, "name", collectionNameUpdated),
					resource.TestCheckResourceAttr(collectionResourceFullName, "description", collectionDescUpdated),
				),
			},
			// --- Step 3: Update Collection (add metadata_schema) ---
			{
				Config: testAccCollectionResourceConfigWithMetadata(projectName, collectionNameUpdated, collectionDescUpdated, map[string]string{"tag": "string", "year": "number"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.%", "2"),
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.tag", "string"),
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.year", "number"),
				),
			},
			// --- Step 4: Update Collection (change metadata_schema) ---
			{
				Config: testAccCollectionResourceConfigWithMetadata(projectName, collectionNameUpdated, collectionDescUpdated, map[string]string{"category": "string", "rating": "number", "is_featured": "boolean"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.%", "3"),
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.category", "string"),
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.rating", "number"),
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.is_featured", "boolean"),
				),
			},
			// --- Step 5: Update Collection (clear metadata_schema by omitting it - depends on API behavior and provider implementation) ---
			// This test assumes that omitting metadata_schema in config and having `Optional: true, Computed: true` with UseStateForUnknown
			// might not clear it. Explicit null might be needed, or specific handling in Update.
			// For now, we'll test setting it to an empty map.
			// If the API/provider supports clearing by omitting, that's another test case.
			{
				Config: testAccCollectionResourceConfigWithMetadata(projectName, collectionNameUpdated, collectionDescUpdated, map[string]string{}), // Empty map
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collectionResourceFullName, "metadata_schema.%", "0"),
				),
			},
			// --- Step 6: Update Collection (add embeddings_model_id) ---
			// This requires a valid or mockable embeddings_model_id.
			// For now, this step is commented out.
			/*
			{
				Config: testAccCollectionResourceConfigWithEmbeddingsModel(projectName, collectionNameUpdated, collectionDescUpdated, testEmbeddingsModelID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(collectionResourceFullName, "embeddings_model_id", testEmbeddingsModelID),
				),
			},
			*/
			// --- Step 7: ImportState testing ---
			{
				ResourceName:      collectionResourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
				// Potentially ignore computed fields that might change or not be perfectly stable for verify, e.g. size_bytes if documents were added outside TF.
				// For now, assume all schema fields are verifiable.
			},
			// Delete testing automatically occurs in TestCase for collection.
			// The prerequisite project will also be deleted.
		},
	})
}

// Config for basic collection with a prerequisite project
func testAccCollectionResourceConfigBasic(projectName, collectionName, description string) string {
	return fmt.Sprintf(`
resource "corax_project" "test_for_collection" {
  name = "%s"
}

resource "corax_collection" "test" {
  project_id  = corax_project.test_for_collection.id
  name        = "%s"
  description = "%s"
}
`, projectName, collectionName, description)
}

// Config with metadata_schema
func testAccCollectionResourceConfigWithMetadata(projectName, collectionName, description string, metadataSchema map[string]string) string {
	metadataSchemaStrings := []string{}
	for k, v := range metadataSchema {
		metadataSchemaStrings = append(metadataSchemaStrings, fmt.Sprintf(`    "%s" = "%s"`, k, v))
	}
	metadataSchemaHCL := ""
	if len(metadataSchemaStrings) > 0 {
		// Corrected: strings.Join returns a string, no need for types.StringValue here.
		metadataSchemaHCL = fmt.Sprintf("{\n%s\n  }", strings.Join(metadataSchemaStrings, ",\n"))
	} else {
		metadataSchemaHCL = "{}"
	}


	return fmt.Sprintf(`
resource "corax_project" "test_for_collection" {
  name = "%s"
}

resource "corax_collection" "test" {
  project_id      = corax_project.test_for_collection.id
  name            = "%s"
  description     = "%s"
  metadata_schema = %s
}
`, projectName, collectionName, description, metadataSchemaHCL)
}

// Config with embeddings_model_id
func testAccCollectionResourceConfigWithEmbeddingsModel(projectName, collectionName, description, embeddingsModelID string) string {
	return fmt.Sprintf(`
resource "corax_project" "test_for_collection" {
  name = "%s"
}

resource "corax_collection" "test" {
  project_id          = corax_project.test_for_collection.id
  name                = "%s"
  description         = "%s"
  embeddings_model_id = "%s"
}
`, projectName, collectionName, description, embeddingsModelID)
}

// Config to test clearing description (by omitting it)
func testAccCollectionResourceConfigOmitDescription(projectName, collectionName string) string {
	return fmt.Sprintf(`
resource "corax_project" "test_for_collection" {
  name = "%s"
}

resource "corax_collection" "test" {
  project_id  = corax_project.test_for_collection.id
  name        = "%s"
  # description is omitted
}
`, projectName, collectionName)
}
