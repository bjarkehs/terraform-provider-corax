package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModelProviderResource_basic(t *testing.T) {
	if os.Getenv("CORAX_API_ENDPOINT") == "" || os.Getenv("CORAX_API_KEY") == "" {
		t.Skip("Skipping acceptance test: CORAX_API_ENDPOINT or CORAX_API_KEY not set")
	}

	resourceName := "corax_model_provider.test"
	// providerName := "tf-acc-test-provider-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	providerName := "tf-acc-test-provider-basic" // Using a fixed name for simplicity in this example
	providerType := "azure_openai"               // Example type, ensure this is valid in your Corax API

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccModelProviderResourceBasicConfig(providerName, providerType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", providerName),
					resource.TestCheckResourceAttr(resourceName, "provider_type", providerType),
					resource.TestCheckResourceAttr(resourceName, "configuration.api_key", "test-api-key"),
					resource.TestCheckResourceAttr(resourceName, "configuration.api_endpoint", "https://example-azure.openai.com/"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Sensitive attributes like configuration.api_key might not be fully verifiable on import if not returned by GET
				// ImportStateVerifyIgnore: []string{"configuration.api_key"},
			},
			// Update and Read testing
			{
				Config: testAccModelProviderResourceUpdatedConfig(providerName+"-updated", providerType), // Name update
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", providerName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "provider_type", providerType), // Type usually not updatable
					resource.TestCheckResourceAttr(resourceName, "configuration.api_key", "updated-test-api-key"),
					resource.TestCheckResourceAttr(resourceName, "configuration.api_endpoint", "https://updated-example-azure.openai.com/"),
					resource.TestCheckResourceAttr(resourceName, "configuration.custom_header", "X-My-Header: test_value"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccModelProviderResourceBasicConfig(name, providerType string) string {
	return fmt.Sprintf(`
provider "corax" {
  # api_endpoint = "..." 
  # api_key      = "..."
}

resource "corax_model_provider" "test" {
  name           = "%s"
  provider_type  = "%s"
  configuration = {
    api_key      = "test-api-key"
    api_endpoint = "https://example-azure.openai.com/"
  }
}
`, name, providerType)
}

func testAccModelProviderResourceUpdatedConfig(name, providerType string) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_model_provider" "test" {
  name           = "%s"
  provider_type  = "%s" # Provider type is often immutable
  configuration = {
    api_key       = "updated-test-api-key" # Updated
    api_endpoint  = "https://updated-example-azure.openai.com/" # Updated
    custom_header = "X-My-Header: test_value" # Added
  }
}
`, name, providerType)
}

// testAccPreCheck is defined in provider_test.go
