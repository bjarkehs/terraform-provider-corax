// Copyright (c) Trifork

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					// Check that schema_def is set and contains the expected JSON string.
					// The new schemaDefAPIToMap converts the whole map into a single JSON string stored in types.Dynamic.
					resource.TestCheckResourceAttr(resourceName, "schema_def", `{`+
						`"age":{"description":"The age of the user","type":"integer"},`+
						`"city":{"description":"The city where the user lives","type":"string"},`+
						`"details":{"description":"Further details","properties":{"hobby":{"description":"User's hobby","type":"string"},"occupation":{"description":"User's occupation","type":"string"}},"type":"object"},`+
						`"is_student":{"description":"Is the user a student","type":"boolean"},`+
						`"name":{"description":"The name of the user","type":"string"}`+
						`}`), // Note: Order of keys in JSON might vary, this could be fragile. A custom check function might be better.
				),
			},
			// Update: Change output_type to text and remove schema_def
			{
				Config: testAccCompletionCapabilityResourceBasicConfig(capabilityName, systemPrompt, completionPrompt), // Re-use basic config
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", capabilityName),
					resource.TestCheckResourceAttr(resourceName, "output_type", "text"),
					// When output_type is "text", schema_def should be null.
					resource.TestCheckResourceAttr(resourceName, "schema_def", ""), // types.DynamicNull() often results in an empty string representation for TestCheckResourceAttr
					// A more robust check might be resource.TestCheckResourceAttrIsNull if supported and reliable for DynamicType
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
