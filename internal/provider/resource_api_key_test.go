package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	// testAccAPIKeyResourcePrefix is a prefix for the names of the
	// test API keys.
	testAccAPIKeyResourcePrefix = "tf-acc-test-apikey-"
)

// TestAccAPIKeyResource provides acceptance tests for the corax_api_key resource.
func TestAccAPIKeyResource(t *testing.T) {
	// API Key and Endpoint must be set via environment variables
	// CORAX_API_KEY and CORAX_API_ENDPOINT
	// Skip test if not set
	if os.Getenv("CORAX_API_KEY") == "" || os.Getenv("CORAX_API_ENDPOINT") == "" {
		t.Skip("CORAX_API_KEY and CORAX_API_ENDPOINT must be set for acceptance tests")
		return
	}

	// Generate a unique name for the API key resource for each test run
	// to prevent collisions and ensure idempotency.
	// rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	// For now, using a timestamp based name for simplicity, ensure it's valid for the API.
	// The API Key name from openapi.json has minLength: 1.
	// Let's use a fixed name for now and rely on cleanup.
	// A better approach would be to use random names and ensure proper cleanup.
	resourceName := "corax_api_key.test"
	apiKeyName := fmt.Sprintf("%s%d", testAccAPIKeyResourcePrefix, time.Now().UnixNano())
	//apiKeyNameUpdated := fmt.Sprintf("%s-updated-%d", testAccAPIKeyResourcePrefix, time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// ProviderFactories: testAccProviderFactories, // Assuming this is defined elsewhere or use internal.TestAccProtoV6ProviderFactories
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAPIKeyResourceConfig(apiKeyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", apiKeyName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "key"), // Key is sensitive, but should be in state
					resource.TestCheckResourceAttrSet(resourceName, "prefix"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
					// expires_at needs to be checked carefully due to time formatting
					// resource.TestCheckResourceAttr(resourceName, "expires_at", "some_expected_format_or_value"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The key attribute will not be present in the API response after creation,
				// so ImportStateVerifyIgnore will prevent a diff.
				ImportStateVerifyIgnore: []string{"key"},
			},
			// Update testing (API Key update is not supported, this step should confirm that behavior or be removed)
			// As per our resource_api_key.go, Update returns an error.
			// If we want to test that, we'd need a config that attempts an update and expect an error.
			// For now, since updates are not supported for name/expires_at (immutable post-creation via PUT),
			// we'll skip an explicit update test that changes these.
			// If other mutable fields existed and were supported by PUT, we'd test them.
			// The current API spec doesn't show a PUT for /v1/api-keys/{key_id}
			// so the resource correctly implements Update as not supported.

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAPIKeyResourceConfig(apiKeyName string) string {
	// Calculate expires_at for 1 hour from now in RFC3339 format
	expiresAt := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	return fmt.Sprintf(`
provider "corax" {
  # api_endpoint and api_key should be configured via environment variables
  # CORAX_API_ENDPOINT and CORAX_API_KEY for acceptance tests.
}

resource "corax_api_key" "test" {
  name       = "%s"
  expires_at = "%s"
}
`, apiKeyName, expiresAt)
}

// Note: The `testAccProtoV6ProviderFactories` variable is defined in `provider_test.go`
// and is available to this package. The `resource.TestCase` above uses it directly.
// No local definition of `testAccProtoV6ProviderFactories` is needed in this file.

// Helper to check if a string is a valid UUID.
func testAccCheckUUID(v string) error {
	if matched, _ := regexp.MatchString("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$", v); !matched {
		return fmt.Errorf("expected UUID, got %s", v)
	}
	return nil
}
