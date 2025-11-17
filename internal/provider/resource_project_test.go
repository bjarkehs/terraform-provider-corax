// Copyright (c) Trifork

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest" // For random strings
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"regexp"
	"testing"
)

const (
	testAccProjectResourcePrefix = "tf-acc-test-project-"
)

// TestAccProjectResource provides acceptance tests for the corax_project resource.
func TestAccProjectResource(t *testing.T) {
	if os.Getenv("CORAX_API_KEY") == "" || os.Getenv("CORAX_API_ENDPOINT") == "" {
		t.Skip("CORAX_API_KEY and CORAX_API_ENDPOINT must be set for acceptance tests")
		return
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectName := fmt.Sprintf("%s%s", testAccProjectResourcePrefix, rName)
	projectNameUpdated := fmt.Sprintf("%s-updated", projectName)
	projectDesc := "Test project description"
	projectDescUpdated := "Test project description updated"
	resourceFullName := "corax_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories, // From provider_test.go
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigBasic(projectName, projectDesc),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", projectName),
					resource.TestCheckResourceAttr(resourceFullName, "description", projectDesc),
					resource.TestCheckResourceAttr(resourceFullName, "is_public", "false"), // Default
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceFullName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceFullName, "owner"),
					resource.TestCheckResourceAttr(resourceFullName, "collection_count", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "capability_count", "0"),
					// Check for UUID format for ID
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
				),
			},
			// Update and Read testing (Name and Description)
			{
				Config: testAccProjectResourceConfigBasic(projectNameUpdated, projectDescUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", projectNameUpdated),
					resource.TestCheckResourceAttr(resourceFullName, "description", projectDescUpdated),
					resource.TestCheckResourceAttr(resourceFullName, "is_public", "false"), // Should remain default
				),
			},
			// Update and Read testing (IsPublic)
			{
				Config: testAccProjectResourceConfigWithPublic(projectNameUpdated, projectDescUpdated, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", projectNameUpdated),
					resource.TestCheckResourceAttr(resourceFullName, "description", projectDescUpdated),
					resource.TestCheckResourceAttr(resourceFullName, "is_public", "true"),
				),
			},
			// Update and Read testing (IsPublic back to false)
			{
				Config: testAccProjectResourceConfigWithPublic(projectNameUpdated, projectDescUpdated, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", projectNameUpdated),
					resource.TestCheckResourceAttr(resourceFullName, "description", projectDescUpdated),
					resource.TestCheckResourceAttr(resourceFullName, "is_public", "false"),
				),
			},
			// Update and Read testing (Clear description)
			{
				Config: testAccProjectResourceConfigNoDescription(projectNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", projectNameUpdated),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""), // Expect empty or null
					resource.TestCheckResourceAttr(resourceFullName, "is_public", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
				// No attributes to ignore for project typically
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfigBasic(projectName, description string) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_project" "test" {
  name        = "%s"
  description = "%s"
}
`, projectName, description)
}

func testAccProjectResourceConfigWithPublic(projectName, description string, isPublic bool) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_project" "test" {
  name        = "%s"
  description = "%s"
  is_public   = %t
}
`, projectName, description, isPublic)
}

func testAccProjectResourceConfigNoDescription(projectName string) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_project" "test" {
  name        = "%s"
  # description is omitted
}
`, projectName)
}
