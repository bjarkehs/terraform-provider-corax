package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCompletionCapabilityResource_basic(t *testing.T) {
	if os.Getenv("CORAX_API_ENDPOINT") == "" || os.Getenv("CORAX_API_KEY") == "" {
		t.Skip("Skipping acceptance test: CORAX_API_ENDPOINT or CORAX_API_KEY not set")
	}
	
	resourceName := "corax_completion_capability.test_basic"
	capabilityName := "tf-acc-test-completion-basic"
	systemPrompt := "You are a text completion model."
	completionPrompt := "Once upon a time, "

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCompletionCapabilityResourceBasicConfig(capabilityName, systemPrompt, completionPrompt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName),
					resource.TestCheckResourceAttr(resourceName, "system_prompt", systemPrompt),
					resource.TestCheckResourceAttr(resourceName, "completion_prompt", completionPrompt),
					resource.TestCheckResourceAttr(resourceName, "output_type", "text"), // Default if not specified, or should be required? Schema says required.
					resource.TestCheckResourceAttr(resourceName, "type", "completion"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCompletionCapabilityResourceBasicConfig(capabilityName+"-upd", systemPrompt+" upd", completionPrompt+"there was a..."),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "system_prompt", systemPrompt+" upd"),
					resource.TestCheckResourceAttr(resourceName, "completion_prompt", completionPrompt+"there was a..."),
				),
			},
		},
	})
}

func TestAccCompletionCapabilityResource_withSchemaOutput(t *testing.T) {
	if os.Getenv("CORAX_API_ENDPOINT") == "" || os.Getenv("CORAX_API_KEY") == "" {
		t.Skip("Skipping acceptance test: CORAX_API_ENDPOINT or CORAX_API_KEY not set")
	}

	resourceName := "corax_completion_capability.test_schema"
	capabilityName := "tf-acc-test-completion-schema"
	systemPrompt := "Extract structured data."
	completionPrompt := "User: John Doe, Age: 30, City: New York."
	
	// Note: The schema_def uses jsonencode for simplicity in HCL.
	// The provider's DynamicType handling for schema_def needs to correctly parse this.
	// Our current schemaDefMapToAPI is basic and might need users to provide JSON strings.

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCompletionCapabilityResourceSchemaOutputConfig(capabilityName, systemPrompt, completionPrompt),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName),
					resource.TestCheckResourceAttr(resourceName, "output_type", "schema"),
					resource.TestCheckResourceAttr(resourceName, "variables.#", "0"), // No variables in this example
					// Checking schema_def is tricky with DynamicType.
					// We'd ideally check specific keys if the provider could parse it into a structured way.
					// For now, just check that it's set.
					resource.TestCheckResourceAttrSet(resourceName, "schema_def.name.type"), 
					resource.TestCheckResourceAttr(resourceName, "schema_def.name.type", "\"string\""), // Note: these will be JSON strings
					resource.TestCheckResourceAttr(resourceName, "schema_def.name.description", "\"The name of the user\""),
					resource.TestCheckResourceAttr(resourceName, "schema_def.age.type", "\"integer\""),
				),
			},
			// Update: Change output_type to text and remove schema_def
			{
				Config: testAccCompletionCapabilityResourceBasicConfig(capabilityName, systemPrompt, completionPrompt), // Re-use basic config
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName),
					resource.TestCheckResourceAttr(resourceName, "output_type", "text"),
					resource.TestCheckResourceAttr(resourceName, "schema_def.%", "0"), // Should be empty
				),
			},
		},
	})
}


func testAccCompletionCapabilityResourceBasicConfig(name, sysPrompt, compPrompt string) string {
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_completion_capability" "test_basic" {
  name               = "%s"
  system_prompt      = "%s"
  completion_prompt  = "%s"
  output_type        = "text"
}
`, name, sysPrompt, compPrompt)
}

func testAccCompletionCapabilityResourceSchemaOutputConfig(name, sysPrompt, compPrompt string) string {
	// Using jsonencode for schema_def values for easier HCL representation.
	// The provider needs to handle these stringified JSON values if DynamicType is used this way.
	// A more robust solution would be a well-defined schema for schema_def itself.
	return fmt.Sprintf(`
provider "corax" {}

resource "corax_completion_capability" "test_schema" {
  name               = "%s"
  system_prompt      = "%s"
  completion_prompt  = "%s"
  output_type        = "schema"
  
  variables = ["User", "Age", "City"] # Example variables

  schema_def = {
    name = jsonencode({
      type        = "string"
      description = "The name of the user"
    })
    age = jsonencode({
      type        = "integer"
      description = "The age of the user"
    })
    city = jsonencode({
      type        = "string"
      description = "The city where the user lives"
    })
    is_student = jsonencode({
      type = "boolean"
      description = "Is the user a student"
    })
    details = jsonencode({
        type = "object"
        description = "Further details"
        properties = {
            hobby = { type = "string", description = "User's hobby"}
            occupation = { type = "string", description = "User's occupation"}
        }
    })
  }

  config {
    temperature = 0.5
  }
}
`, name, sysPrompt, compPrompt)
}
