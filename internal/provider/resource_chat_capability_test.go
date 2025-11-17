// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChatCapabilityResource_basic(t *testing.T) {
	// Skip if CORAX_API_ENDPOINT or CORAX_API_KEY is not set, as these are needed for acceptance tests
	if os.Getenv("CORAX_API_ENDPOINT") == "" || os.Getenv("CORAX_API_KEY") == "" {
		t.Skip("Skipping acceptance test: CORAX_API_ENDPOINT or CORAX_API_KEY not set")
	}

	resourceName := "corax_chat_capability.test"
	capabilityName := "tf-acc-test-chat-cap-basic"
	systemPrompt := "You are a helpful assistant."

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccChatCapabilityResourceBasicConfig(capabilityName, systemPrompt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName),
					resource.TestCheckResourceAttr(resourceName, "system_prompt", systemPrompt),
					resource.TestCheckResourceAttr(resourceName, "type", "chat"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"), // Default
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// ImportStateVerifyIgnore: []string{"config"}, // Config might have defaults applied by API not set in HCL
			},
			// Update and Read testing (e.g., change name and system_prompt)
			{
				Config: testAccChatCapabilityResourceBasicConfig(capabilityName+"-updated", systemPrompt+" Updated."),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "system_prompt", systemPrompt+" Updated."),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccChatCapabilityResource_withConfig(t *testing.T) {
	if os.Getenv("CORAX_API_ENDPOINT") == "" || os.Getenv("CORAX_API_KEY") == "" {
		t.Skip("Skipping acceptance test: CORAX_API_ENDPOINT or CORAX_API_KEY not set")
	}

	resourceName := "corax_chat_capability.test_with_config"
	capabilityName := "tf-acc-test-chat-cap-config"
	systemPrompt := "You are a configured assistant."

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChatCapabilityResourceWithConfig(capabilityName, systemPrompt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName),
					resource.TestCheckResourceAttr(resourceName, "system_prompt", systemPrompt),
					resource.TestCheckResourceAttr(resourceName, "config.temperature", "0.7"),
					resource.TestCheckResourceAttr(resourceName, "config.content_tracing", "true"),
					resource.TestCheckResourceAttr(resourceName, "config.data_retention.timed.hours", "24"),
					resource.TestCheckResourceAttr(resourceName, "config.blob_config.max_file_size_mb", "10"),
					resource.TestCheckResourceAttr(resourceName, "config.blob_config.max_blobs", "5"),
					resource.TestCheckResourceAttr(resourceName, "config.blob_config.allowed_mime_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.blob_config.allowed_mime_types.0", "image/jpeg"),
				),
			},
			// Update config
			{
				Config: testAccChatCapabilityResourceWithUpdatedConfig(capabilityName, systemPrompt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName), // Name unchanged
					resource.TestCheckResourceAttr(resourceName, "config.temperature", "0.8"),
					resource.TestCheckResourceAttr(resourceName, "config.content_tracing", "false"),
					resource.TestCheckResourceAttr(resourceName, "config.data_retention.infinite.enabled", "true"),
				),
			},
		},
	})
}

func testAccChatCapabilityResourceBasicConfig(name, systemPrompt string) string {
	return fmt.Sprintf(`
provider "corax" {
  # api_endpoint = "..." # Ensure these are set via ENV for tests
  # api_key      = "..."
}

resource "corax_chat_capability" "test" {
  name          = "%s"
  system_prompt = "%s"
}
`, name, systemPrompt)
}

func testAccChatCapabilityResourceWithConfig(name, systemPrompt string) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_chat_capability" "test_with_config" {
  name          = "%s"
  system_prompt = "%s"
  is_public     = true
  
  config {
    temperature     = 0.7
    content_tracing = true
    data_retention {
      type  = "timed"
      hours = 24
    }
    blob_config {
      max_file_size_mb = 10
      max_blobs          = 5
      allowed_mime_types = ["image/jpeg"]
    }
  }
}
`, name, systemPrompt)
}

func testAccChatCapabilityResourceWithUpdatedConfig(name, systemPrompt string) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_chat_capability" "test_with_config" {
  name          = "%s"
  system_prompt = "%s" // System prompt could also be updated here if desired
  is_public     = true 
  
  config {
    temperature     = 0.8
    content_tracing = false // Updated
    data_retention {
      type = "infinite" // Changed from timed to infinite
    }
    // blob_config removed, should revert to API defaults or be null if API allows removal
  }
}
`, name, systemPrompt)
}

// testAccPreCheck validates necessary provider configuration for acceptance tests.
// It's typically defined in provider_test.go but can be here if specific to this resource.
// func testAccPreCheck(t *testing.T) {
// 	if v := os.Getenv("CORAX_API_ENDPOINT"); v == "" {
// 		t.Fatal("CORAX_API_ENDPOINT must be set for acceptance tests")
// 	}
// 	if v := os.Getenv("CORAX_API_KEY"); v == "" {
// 		t.Fatal("CORAX_API_KEY must be set for acceptance tests")
// 	}
// }
