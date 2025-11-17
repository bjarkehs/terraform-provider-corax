// Copyright (c) Trifork

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-corax/internal/coraxclient" // TODO: Adjust this path if your module name is different
)

// Ensure CoraxProvider satisfies various provider interfaces.
var _ provider.Provider = &CoraxProvider{}
var _ provider.ProviderWithFunctions = &CoraxProvider{}
var _ provider.ProviderWithEphemeralResources = &CoraxProvider{}

// CoraxProvider defines the provider implementation.
type CoraxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CoraxProviderModel describes the provider data model.
type CoraxProviderModel struct {
	APIEndpoint types.String `tfsdk:"api_endpoint"`
	APIKey      types.String `tfsdk:"api_key"`
}

func (p *CoraxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "corax" // Updated TypeName
	resp.Version = p.version
}

func (p *CoraxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for Corax API.",
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint for the Corax API. Can also be set via CORAX_API_ENDPOINT environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API Key for the Corax API. Can also be set via CORAX_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *CoraxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CoraxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read configuration from environment variables if not set in config
	if data.APIEndpoint.IsNull() || data.APIEndpoint.ValueString() == "" {
		envEndpoint := os.Getenv("CORAX_API_ENDPOINT")
		if envEndpoint != "" {
			data.APIEndpoint = types.StringValue(envEndpoint)
			tflog.Debug(ctx, "Using CORAX_API_ENDPOINT from environment variable")
		}
	}

	if data.APIKey.IsNull() || data.APIKey.ValueString() == "" {
		envAPIKey := os.Getenv("CORAX_API_KEY")
		if envAPIKey != "" {
			data.APIKey = types.StringValue(envAPIKey)
			tflog.Debug(ctx, "Using CORAX_API_KEY from environment variable")
		}
	}

	// Validate required configuration
	if data.APIEndpoint.IsNull() || data.APIEndpoint.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing API Endpoint Configuration",
			"The provider cannot be configured without an API endpoint. "+
				"Set the api_endpoint attribute in the provider configuration or use the CORAX_API_ENDPOINT environment variable.",
		)
	}

	if data.APIKey.IsNull() || data.APIKey.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing API Key Configuration",
			"The provider cannot be configured without an API Key. "+
				"Set the api_key attribute in the provider configuration or use the CORAX_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Configuring Corax API client")
	tflog.Debug(ctx, "Corax API Endpoint: "+data.APIEndpoint.ValueString())
	// Do not log API key for security reasons, even at debug level.

	client, err := coraxclient.NewClient(data.APIEndpoint.ValueString(), data.APIKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Corax API client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Corax API client configured successfully")
}

func (p *CoraxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAPIKeyResource,
		NewProjectResource,
		NewChatCapabilityResource,             // Added Chat Capability
		NewCompletionCapabilityResource,       // Added Completion Capability
		NewModelDeploymentResource,            // Added Model Deployment
		NewModelProviderResource,              // Added Model Provider
		NewCapabilityTypeDefaultModelResource, // Added Capability Type Default Model
		// NewCollectionResource, // Removed as per new scope
		// NewDocumentResource,   // Removed as per new scope
		// NewEmbeddingsModelResource, // Removed as per new scope
	}
}

func (p *CoraxProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource { // Updated receiver to CoraxProvider
	return []func() ephemeral.EphemeralResource{}
}

func (p *CoraxProvider) DataSources(ctx context.Context) []func() datasource.DataSource { // Updated receiver to CoraxProvider
	return []func() datasource.DataSource{}
}

func (p *CoraxProvider) Functions(ctx context.Context) []func() function.Function { // Updated receiver to CoraxProvider
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CoraxProvider{ // Updated to CoraxProvider
			version: version,
		}
	}
}
