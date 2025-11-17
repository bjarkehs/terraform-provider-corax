// Copyright (c) Trifork


package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"corax": providerserver.NewProtocol6WithError(New("test")()), // Changed "scaffolding" to "corax"
}

// testAccProtoV6ProviderFactoriesWithEcho includes the echo provider alongside the corax provider.
// It allows for testing assertions on data returned by an ephemeral resource during Open.
// The echoprovider is used to arrange tests by echoing ephemeral data into the Terraform state.
// This lets the data be referenced in test assertions with state checks.
//
//nolint:unused // retained for potential future acceptance tests involving echo provider
var testAccProtoV6ProviderFactoriesWithEcho = map[string]func() (tfprotov6.ProviderServer, error){
	"corax": providerserver.NewProtocol6WithError(New("test")()), // Changed "scaffolding" to "corax"
	"echo":  echoprovider.NewProviderServer(),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CORAX_API_ENDPOINT"); v == "" {
		t.Fatal("CORAX_API_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("CORAX_API_KEY"); v == "" {
		t.Fatal("CORAX_API_KEY must be set for acceptance tests")
	}
}
