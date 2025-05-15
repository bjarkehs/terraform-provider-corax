package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest" // For random strings if needed
)

const testAccModelDeploymentProviderIDEnvVar = "CORAX_TEST_MODEL_PROVIDER_ID"

func TestAccModelDeploymentResource_basic(t *testing.T) {
	if os.Getenv("CORAX_API_ENDPOINT") == "" || os.Getenv("CORAX_API_KEY") == "" {
		t.Skip("Skipping acceptance test: CORAX_API_ENDPOINT or CORAX_API_KEY not set")
	}
	testProviderID := os.Getenv(testAccModelDeploymentProviderIDEnvVar)
	if testProviderID == "" {
		t.Skipf("Skipping acceptance test: %s must be set", testAccModelDeploymentProviderIDEnvVar)
	}

	resourceName := "corax_model_deployment.test"
	deploymentName := "tf-acc-test-deployment-basic"
	// deploymentNameUpdated := deploymentName + "-updated" // Not used yet, but for update test

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccModelDeploymentResourceBasicConfig(deploymentName, testProviderID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", deploymentName),
					resource.TestCheckResourceAttr(resourceName, "provider_id", testProviderID),
					resource.TestCheckResourceAttr(resourceName, "supported_tasks.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "supported_tasks.0", "chat"),
					resource.TestCheckResourceAttr(resourceName, "supported_tasks.1", "completion"),
					resource.TestCheckResourceAttr(resourceName, "configuration.model_name", "gpt-3.5-turbo"),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
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
				// ImportStateVerifyIgnore: []string{"some_attribute_that_might_not_be_in_api_response"},
			},
			// Update and Read testing
			{
				Config: testAccModelDeploymentResourceUpdatedConfig(deploymentName+"-updated", testProviderID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", deploymentName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "is_active", "false"),
					resource.TestCheckResourceAttr(resourceName, "supported_tasks.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "supported_tasks.0", "embedding"),
					resource.TestCheckResourceAttr(resourceName, "configuration.model_name", "text-embedding-ada-002"),
					resource.TestCheckResourceAttr(resourceName, "configuration.api_version", "2023-05-15"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccModelDeploymentResourceBasicConfig(name, providerID string) string {
	return fmt.Sprintf(`
provider "corax" {
  # api_endpoint = "..." 
  # api_key      = "..."
}

resource "corax_model_deployment" "test" {
  name            = "%s"
  provider_id     = "%s"
  supported_tasks = ["chat", "completion"]
  configuration = {
    model_name = "gpt-3.5-turbo"
  }
  is_active       = true
  description     = "Basic test deployment"
}
`, name, providerID)
}

func testAccModelDeploymentResourceUpdatedConfig(name, providerID string) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_model_deployment" "test" {
  name            = "%s"
  provider_id     = "%s"
  supported_tasks = ["embedding"] # Changed
  configuration = {
    model_name   = "text-embedding-ada-002" # Changed
    api_version  = "2023-05-15"             # Added
  }
  is_active       = false # Changed
  description     = "Updated description" # Changed
}
`, name, providerID)
}

// testAccPreCheck is defined in provider_test.go
