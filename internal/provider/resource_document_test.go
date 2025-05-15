package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testAccDocumentResourcePrefix = "tf-acc-test-doc-"
)

// TestAccDocumentResource provides acceptance tests for the corax_document resource.
func TestAccDocumentResource(t *testing.T) {
	if os.Getenv("CORAX_API_KEY") == "" || os.Getenv("CORAX_API_ENDPOINT") == "" {
		t.Skip("CORAX_API_KEY and CORAX_API_ENDPOINT must be set for acceptance tests")
		return
	}

	// Pre-requisites
	projectRName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectName := fmt.Sprintf("%s%s-for-doc", testAccProjectResourcePrefix, projectRName)
	projectResourceFullName := "corax_project.test_for_doc"

	collectionRName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	collectionName := fmt.Sprintf("%s%s-for-doc", testAccCollectionResourcePrefix, collectionRName)
	collectionResourceFullName := "corax_collection.test_for_doc"

	// Document details
	docRName1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	// docID1UserProvided := fmt.Sprintf("doc-%s", docRName1) // Test with user-provided ID
	docTextContent1 := "This is the first test document with plain text content."
	docTextContent1Updated := "This is the updated text content for the first test document."
	docResourceFullName := "corax_document.test_text"

	docRName2 := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	// docID2UserProvided := fmt.Sprintf("doc-%s", docRName2)
	docJsonContent2 := `{"title": "JSON Document", "version": 1, "tags": ["test", "json"]}`
	docJsonContent2Updated := `{"title": "JSON Document Updated", "version": 2, "tags": ["test", "json", "updated"], "status": "final"}`
	docResourceJsonFullName := "corax_document.test_json"


	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// --- Test Text Document ---
			// Step 1: Create text document (API generates ID)
			{
				Config: testAccDocumentResourceConfigText(projectName, collectionName, docTextContent1, map[string]interface{}{"source": "text_test", "priority": 1.0}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(docResourceFullName, "id"), // API generated
					resource.TestMatchResourceAttr(docResourceFullName, "id", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					resource.TestCheckResourceAttrPair(docResourceFullName, "collection_id", collectionResourceFullName, "id"),
					resource.TestCheckResourceAttr(docResourceFullName, "text_content", docTextContent1),
					resource.TestCheckResourceAttr(docResourceFullName, "json_content", ""), // Should be null/empty
					resource.TestCheckResourceAttr(docResourceFullName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(docResourceFullName, "metadata.source", "text_test"),
					resource.TestCheckResourceAttr(docResourceFullName, "metadata.priority", "1"), // Note: float becomes string in TF state for dynamic
					resource.TestCheckResourceAttrSet(docResourceFullName, "created_at"),
					resource.TestCheckResourceAttrSet(docResourceFullName, "token_count"),
					resource.TestCheckResourceAttrSet(docResourceFullName, "chunk_count"),
					resource.TestCheckResourceAttrSet(docResourceFullName, "embeddings_status"),
				),
			},
			// Step 2: Update text document (content and metadata)
			{
				Config: testAccDocumentResourceConfigText(projectName, collectionName, docTextContent1Updated, map[string]interface{}{"source": "text_test_updated", "status": "reviewed", "priority": 2.0}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(docResourceFullName, "text_content", docTextContent1Updated),
					resource.TestCheckResourceAttr(docResourceFullName, "metadata.%", "3"),
					resource.TestCheckResourceAttr(docResourceFullName, "metadata.source", "text_test_updated"),
					resource.TestCheckResourceAttr(docResourceFullName, "metadata.status", "reviewed"),
					resource.TestCheckResourceAttr(docResourceFullName, "metadata.priority", "2"),
					resource.TestCheckResourceAttrSet(docResourceFullName, "updated_at"),
				),
			},
			// Step 3: Import text document
			{
				ResourceName:      docResourceFullName,
				ImportState:       true,
				ImportStateIdFunc: testAccCoraxDocumentImportStateIdFunc(docResourceFullName),
				ImportStateVerify: true,
				// ImportStateVerifyIgnore: []string{"metadata.priority"}, // if float precision causes issues
			},

			// --- Test JSON Document ---
			// Step 4: Create JSON document (user provided ID)
			{
				Config: testAccDocumentResourceConfigJsonWithID(projectName, collectionName, fmt.Sprintf("json-doc-%s", docRName2), docJsonContent2, map[string]interface{}{"type": "json_payload", "version": 1.0}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(docResourceJsonFullName, "id", fmt.Sprintf("json-doc-%s", docRName2)), // User provided
					resource.TestCheckResourceAttrPair(docResourceJsonFullName, "collection_id", collectionResourceFullName, "id"),
					resource.TestCheckResourceAttr(docResourceJsonFullName, "json_content", docJsonContent2),
					resource.TestCheckResourceAttr(docResourceJsonFullName, "text_content", ""), // Should be null/empty
					resource.TestCheckResourceAttr(docResourceJsonFullName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(docResourceJsonFullName, "metadata.type", "json_payload"),
					resource.TestCheckResourceAttr(docResourceJsonFullName, "metadata.version", "1"),
				),
			},
			// Step 5: Update JSON document (content and metadata)
			{
				Config: testAccDocumentResourceConfigJsonWithID(projectName, collectionName, fmt.Sprintf("json-doc-%s", docRName2), docJsonContent2Updated, map[string]interface{}{"type": "json_payload_v2", "processed": true}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(docResourceJsonFullName, "json_content", docJsonContent2Updated),
					resource.TestCheckResourceAttr(docResourceJsonFullName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(docResourceJsonFullName, "metadata.type", "json_payload_v2"),
					resource.TestCheckResourceAttr(docResourceJsonFullName, "metadata.processed", "true"),
				),
			},
			// Step 6: Import JSON document
			{
				ResourceName:      docResourceJsonFullName,
				ImportState:       true,
				ImportStateIdFunc: testAccCoraxDocumentImportStateIdFunc(docResourceJsonFullName),
				ImportStateVerify: true,
			},
			// --- Test Validation: ExactlyOneOf text_content or json_content ---
			// This would require a config that's expected to fail.
			// {
			//  Config: testAccDocumentResourceConfigInvalidContent(projectName, collectionName),
			//  ExpectError: regexp.MustCompile(".*one of `text_content` or `json_content` must be provided.*"), // Adjust regex
			// },
			// {
			//  Config: testAccDocumentResourceConfigBothContents(projectName, collectionName),
			//  ExpectError: regexp.MustCompile(".*ExactlyOneOf.*"), // Adjust regex
			// },
		},
	})
}

func testAccCoraxDocumentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["collection_id"], rs.Primary.Attributes["id"]), nil
	}
}


func formatMetadataToHCL(metadata map[string]interface{}) string {
	if len(metadata) == 0 {
		return "{}"
	}
	var parts []string
	for k, v := range metadata {
		// Basic formatting, might need improvement for complex types or quotes
		switch val := v.(type) {
		case string:
			parts = append(parts, fmt.Sprintf(`    "%s" = "%s"`, k, strings.ReplaceAll(val, "\"", "\\\"")))
		case int, int64, float64, bool:
			parts = append(parts, fmt.Sprintf(`    "%s" = %v`, k, val))
		default:
			// For maps or slices, marshal to JSON string as a fallback for HCL representation
			jsonVal, err := json.Marshal(v)
			if err == nil {
				parts = append(parts, fmt.Sprintf(`    "%s" = %s`, k, string(jsonVal))) // Use raw JSON
			} else {
				// Fallback for unhandled types, may cause HCL errors
				parts = append(parts, fmt.Sprintf(`    "%s" = "%v"`, k, v))
			}
		}
	}
	return fmt.Sprintf("{\n%s\n  }", strings.Join(parts, ",\n"))
}


// Config for text document with API-generated ID
func testAccDocumentResourceConfigText(projectName, collectionName, textContent string, metadata map[string]interface{}) string {
	metadataHCL := formatMetadataToHCL(metadata)
	return fmt.Sprintf(`
resource "corax_project" "test_for_doc" {
  name = "%s"
}

resource "corax_collection" "test_for_doc" {
  project_id = corax_project.test_for_doc.id
  name       = "%s"
}

resource "corax_document" "test_text" {
  collection_id = corax_collection.test_for_doc.id
  text_content  = "%s"
  metadata      = %s
}
`, projectName, collectionName, textContent, metadataHCL)
}

// Config for JSON document with user-provided ID
func testAccDocumentResourceConfigJsonWithID(projectName, collectionName, docID, jsonContent string, metadata map[string]interface{}) string {
	metadataHCL := formatMetadataToHCL(metadata)
	// Escape backticks in jsonContent for HCL heredoc or quoted string
	escapedJsonContent := strings.ReplaceAll(jsonContent, "`", "$${'`'}") // For heredoc
	escapedJsonContent = strings.ReplaceAll(escapedJsonContent, "\\", "\\\\")
	escapedJsonContent = strings.ReplaceAll(escapedJsonContent, "\"", "\\\"")


	return fmt.Sprintf(`
resource "corax_project" "test_for_doc" {
  name = "%s"
}

resource "corax_collection" "test_for_doc" {
  project_id = corax_project.test_for_doc.id
  name       = "%s"
}

resource "corax_document" "test_json" {
  id            = "%s"
  collection_id = corax_collection.test_for_doc.id
  json_content  = "%s"
  metadata      = %s
}
`, projectName, collectionName, docID, escapedJsonContent, metadataHCL)
}

/*
// Config that should fail validation (neither text_content nor json_content)
func testAccDocumentResourceConfigInvalidContent(projectName, collectionName string) string {
	return fmt.Sprintf(`
resource "corax_project" "test_for_doc" {
  name = "%s"
}
resource "corax_collection" "test_for_doc" {
  project_id = corax_project.test_for_doc.id
  name       = "%s"
}
resource "corax_document" "test_invalid" {
  collection_id = corax_collection.test_for_doc.id
  metadata = {
    key = "value"
  }
}
`, projectName, collectionName)
}

// Config that should fail validation (both text_content and json_content)
func testAccDocumentResourceConfigBothContents(projectName, collectionName string) string {
	return fmt.Sprintf(`
resource "corax_project" "test_for_doc" {
  name = "%s"
}
resource "corax_collection" "test_for_doc" {
  project_id = corax_project.test_for_doc.id
  name       = "%s"
}
resource "corax_document" "test_both" {
  collection_id = corax_collection.test_for_doc.id
  text_content  = "some text"
  json_content  = "{\"key\":\"value\"}"
}
`, projectName, collectionName)
}
*/
