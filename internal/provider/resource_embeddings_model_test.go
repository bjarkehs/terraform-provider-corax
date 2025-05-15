package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	// "github.com/hashicorp/terraform-plugin-framework/types" // Not directly used in this basic test structure
)

const (
	testAccEmbeddingsModelResourcePrefix = "tf-acc-test-em-"
)

// TestAccEmbeddingsModelResource provides acceptance tests for the corax_embeddings_model resource.
func TestAccEmbeddingsModelResource(t *testing.T) {
	if os.Getenv("CORAX_API_KEY") == "" || os.Getenv("CORAX_API_ENDPOINT") == "" {
		t.Skip("CORAX_API_KEY and CORAX_API_ENDPOINT must be set for acceptance tests")
		return
	}

	// For some providers, a real API key for a model might be needed for full validation.
	// For now, we'll test with placeholder values or assume the API can validate without making external calls
	// or that it has a "custom" provider type that doesn't require a live key for schema testing.
	// If the API attempts to validate the model by calling the provider (e.g. OpenAI), these tests might fail
	// without valid (even if dummy) credentials for that provider.

	modelRName1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	modelName1 := fmt.Sprintf("%s%s-custom", testAccEmbeddingsModelResourcePrefix, modelRName1)
	modelDesc1 := "Test custom embeddings model"
	modelProvider1 := "custom" // Assuming 'custom' is a valid provider type that might not need real keys
	modelNameOnProvider1 := "my-custom-model-v1"
	modelDimensions1 := int64(384)
	modelMaxTokens1 := int64(512)

	modelName1Updated := fmt.Sprintf("%s-updated", modelName1)
	modelDesc1Updated := "Test custom embeddings model - updated description"
	// apiBaseUrl1 := "http://localhost:8080/v1/embeddings" // Example, if needed

	modelResourceFullName := "corax_embeddings_model.test_custom"

	// --- Test case for a model that might require an API key (e.g., OpenAI) ---
	// modelRName2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	// modelName2 := fmt.Sprintf("%s%s-openai", testAccEmbeddingsModelResourcePrefix, modelRName2)
	// modelProvider2 := "openai"
	// modelNameOnProvider2 := "text-embedding-ada-002" // A known OpenAI model
	// modelDimensions2 := int64(1536)
	// modelApiKey2 := "sk-dummykeyfortest" // Placeholder
	// modelResourceOpenAIFullName := "corax_embeddings_model.test_openai"


	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// --- Step 1: Create Custom Embeddings Model (minimal) ---
			{
				Config: testAccEmbeddingsModelResourceConfigCustomBasic(modelName1, modelProvider1, modelNameOnProvider1, modelDimensions1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(modelResourceFullName, "id"),
					resource.TestMatchResourceAttr(modelResourceFullName, "id", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					resource.TestCheckResourceAttr(modelResourceFullName, "name", modelName1),
					resource.TestCheckResourceAttr(modelResourceFullName, "description", ""), // Optional, so empty if not set
					resource.TestCheckResourceAttr(modelResourceFullName, "model_provider", modelProvider1),
					resource.TestCheckResourceAttr(modelResourceFullName, "model_name_on_provider", modelNameOnProvider1),
					resource.TestCheckResourceAttr(modelResourceFullName, "dimensions", fmt.Sprintf("%d", modelDimensions1)),
					resource.TestCheckResourceAttrSet(modelResourceFullName, "max_tokens"), // Should be computed or have a default from API
					resource.TestCheckResourceAttr(modelResourceFullName, "api_key", ""), // Optional, sensitive, empty if not set
					resource.TestCheckResourceAttr(modelResourceFullName, "api_base_url", ""), // Optional, empty if not set
					resource.TestCheckResourceAttrSet(modelResourceFullName, "status"),
					resource.TestCheckResourceAttrSet(modelResourceFullName, "is_default"),
					resource.TestCheckResourceAttrSet(modelResourceFullName, "created_at"),
				),
			},
			// --- Step 2: Update Custom Embeddings Model (name, description, api_base_url) ---
			{
				Config: testAccEmbeddingsModelResourceConfigCustomFull(modelName1Updated, modelDesc1Updated, modelProvider1, modelNameOnProvider1, modelDimensions1, modelMaxTokens1, "", "http://custom-embeddings.local/v1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelResourceFullName, "name", modelName1Updated),
					resource.TestCheckResourceAttr(modelResourceFullName, "description", modelDesc1Updated),
					resource.TestCheckResourceAttr(modelResourceFullName, "api_base_url", "http://custom-embeddings.local/v1"),
					resource.TestCheckResourceAttrSet(modelResourceFullName, "updated_at"),
				),
			},
			// --- Step 3: Import Custom Embeddings Model ---
			{
				ResourceName:      modelResourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
				// ImportStateVerifyIgnore: []string{"api_key"}, // API key is sensitive and not returned by GET
			},
			// --- Step 4: Test ForceNew on immutable fields (e.g., dimensions) ---
			// This step changes 'dimensions', which should trigger a replacement.
			{
				Config: testAccEmbeddingsModelResourceConfigCustomBasic(modelName1Updated, modelProvider1, modelNameOnProvider1, modelDimensions1+10), // Changed dimensions
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(modelResourceFullName, "name", modelName1Updated), // Name should persist if not changed
					resource.TestCheckResourceAttr(modelResourceFullName, "dimensions", fmt.Sprintf("%d", modelDimensions1+10)), // New dimension
					// ID should change due to ForceNew
				),
			},

			// --- Test Case for a model requiring API Key (if API supports validation or storing it) ---
			// This test might require specific API behavior or a mock.
			/*
			{
				Config: testAccEmbeddingsModelResourceConfigWithKey(modelName2, modelProvider2, modelNameOnProvider2, modelDimensions2, modelApiKey2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(modelResourceOpenAIFullName, "id"),
					resource.TestCheckResourceAttr(modelResourceOpenAIFullName, "name", modelName2),
					resource.TestCheckResourceAttr(modelResourceOpenAIFullName, "model_provider", modelProvider2),
					resource.TestCheckResourceAttr(modelResourceOpenAIFullName, "model_name_on_provider", modelNameOnProvider2),
					resource.TestCheckResourceAttr(modelResourceOpenAIFullName, "dimensions", fmt.Sprintf("%d", modelDimensions2)),
					// resource.TestCheckResourceAttr(modelResourceOpenAIFullName, "api_key", modelApiKey2), // This will fail if API doesn't return it
					resource.TestCheckResourceAttrSet(modelResourceOpenAIFullName, "status"),
				),
			},
			// Update API Key
			{
				Config: testAccEmbeddingsModelResourceConfigWithKey(modelName2, modelProvider2, modelNameOnProvider2, modelDimensions2, modelApiKey2+"-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttr(modelResourceOpenAIFullName, "api_key", modelApiKey2+"-updated"), // This will fail
					resource.TestCheckResourceAttrSet(modelResourceOpenAIFullName, "updated_at"),
				),
			},
			*/

			// Delete testing automatically occurs in TestCase.
		},
	})
}

// Config for a "custom" embeddings model with minimal fields
func testAccEmbeddingsModelResourceConfigCustomBasic(name, provider, nameOnProvider string, dimensions int64) string {
	return fmt.Sprintf(`
resource "corax_embeddings_model" "test_custom" {
  name                   = "%s"
  model_provider         = "%s"
  model_name_on_provider = "%s"
  dimensions             = %d
}
`, name, provider, nameOnProvider, dimensions)
}

// Config for a "custom" embeddings model with all optional fields
func testAccEmbeddingsModelResourceConfigCustomFull(name, description, provider, nameOnProvider string, dimensions, maxTokens int64, apiKey, apiBaseUrl string) string {
	descStr := ""
	if description != "" {
		descStr = fmt.Sprintf(`description = "%s"`, description)
	}
	maxTokensStr := ""
	if maxTokens > 0 {
		maxTokensStr = fmt.Sprintf(`max_tokens = %d`, maxTokens)
	}
	apiKeyStr := ""
	if apiKey != "" {
		apiKeyStr = fmt.Sprintf(`api_key = "%s"`, apiKey)
	}
	apiBaseUrlStr := ""
	if apiBaseUrl != "" {
		apiBaseUrlStr = fmt.Sprintf(`api_base_url = "%s"`, apiBaseUrl)
	}

	return fmt.Sprintf(`
resource "corax_embeddings_model" "test_custom" {
  name                   = "%s"
  %s
  model_provider         = "%s"
  model_name_on_provider = "%s"
  dimensions             = %d
  %s
  %s
  %s
}
`, name, descStr, provider, nameOnProvider, dimensions, maxTokensStr, apiKeyStr, apiBaseUrlStr)
}


// Config for a model that might require an API key (e.g., OpenAI)
func testAccEmbeddingsModelResourceConfigWithKey(name, provider, nameOnProvider string, dimensions int64, apiKey string) string {
	return fmt.Sprintf(`
resource "corax_embeddings_model" "test_openai" {
  name                   = "%s"
  model_provider         = "%s"
  model_name_on_provider = "%s"
  dimensions             = %d
  api_key                = "%s"
  // max_tokens can be omitted to test API default/derivation
}
`, name, provider, nameOnProvider, dimensions, apiKey)
}
