package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-corax/internal/coraxclient" // TODO: Adjust if your module name is different
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &APIKeyResource{}
var _ resource.ResourceWithImportState = &APIKeyResource{}

func NewAPIKeyResource() resource.Resource {
	return &APIKeyResource{}
}

// APIKeyResource defines the resource implementation.
type APIKeyResource struct {
	client *coraxclient.Client
}

// APIKeyResourceModel describes the resource data model.
type APIKeyResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ExpiresAt  types.String `tfsdk:"expires_at"`
	Key        types.String `tfsdk:"key"`
	Prefix     types.String `tfsdk:"prefix"`
	IsActive   types.Bool   `tfsdk:"is_active"`
	LastUsedAt types.String `tfsdk:"last_used_at"`
	UsageCount types.Int64  `tfsdk:"usage_count"`
}

func (r *APIKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *APIKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax API Key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the API key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the API key.",
			},
			"expires_at": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The expiration date and time for the API key (RFC3339 format).",
				// TODO: Add validation for RFC3339 format if possible, or handle in Create/Update.
			},
			"key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The API key secret. This is only available upon creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // Or potentially a modifier to ensure it's only set on create
				},
			},
			"prefix": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The prefix of the API key.",
			},
			"is_active": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates whether the API key is active.",
			},
			"last_used_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the API key was last used.",
			},
			"usage_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of times the API key has been used.",
			},
		},
	}
}

func (r *APIKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*coraxclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *coraxclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *APIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data APIKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating API Key with name: %s, expires_at: %s", data.Name.ValueString(), data.ExpiresAt.ValueString()))

	apiKeyInput := coraxclient.ApiKeyCreate{
		Name:      data.Name.ValueString(),
		ExpiresAt: data.ExpiresAt.ValueString(),
	}

	createdAPIKey, err := r.client.CreateAPIKey(ctx, apiKeyInput)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create API key, got error: %s", err))
		return
	}

	data.ID = types.StringValue(createdAPIKey.ID)
	data.Key = types.StringValue(createdAPIKey.Key) // This is sensitive and usually only returned on create
	data.Prefix = types.StringValue(createdAPIKey.Prefix)
	data.IsActive = types.BoolValue(createdAPIKey.IsActive)
	if createdAPIKey.LastUsedAt != nil && *createdAPIKey.LastUsedAt != "" {
		data.LastUsedAt = types.StringValue(*createdAPIKey.LastUsedAt)
	} else {
		data.LastUsedAt = types.StringNull()
	}
	data.UsageCount = types.Int64Value(int64(createdAPIKey.UsageCount))
	data.Name = types.StringValue(createdAPIKey.Name) // Re-set name in case API modifies it
	if createdAPIKey.ExpiresAt != nil {
		data.ExpiresAt = types.StringValue(*createdAPIKey.ExpiresAt) // Re-set expires_at
	} else {
		data.ExpiresAt = types.StringNull() // Should not happen based on schema (required)
	}

	tflog.Info(ctx, fmt.Sprintf("API Key created successfully with ID: %s", createdAPIKey.ID))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data APIKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiKeyID := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading API Key with ID: %s", apiKeyID))

	apiKey, err := r.client.GetAPIKey(ctx, apiKeyID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("API Key with ID %s not found, removing from state", apiKeyID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API key %s, got error: %s", apiKeyID, err))
		return
	}

	data.Name = types.StringValue(apiKey.Name)
	if apiKey.ExpiresAt != nil {
		data.ExpiresAt = types.StringValue(*apiKey.ExpiresAt)
	} else {
		data.ExpiresAt = types.StringNull() // Should not happen
		}
		data.Prefix = types.StringValue(apiKey.Prefix)
		data.IsActive = types.BoolValue(apiKey.IsActive)
	if apiKey.LastUsedAt != nil && *apiKey.LastUsedAt != "" {
		data.LastUsedAt = types.StringValue(*apiKey.LastUsedAt)
	} else {
		data.LastUsedAt = types.StringNull()
	}
	data.UsageCount = types.Int64Value(int64(apiKey.UsageCount))
	// Note: The 'key' field is typically not returned by a GET request for security reasons.
	// It should remain as it was set during creation (or import). data.Key is already populated from state.

	tflog.Debug(ctx, fmt.Sprintf("Successfully read API Key with ID: %s", apiKeyID))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Updating API Keys is not supported. Please create a new API Key and delete the old one if changes are needed.",
	)
}

func (r *APIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data APIKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiKeyID := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting API Key with ID: %s", apiKeyID))

	err := r.client.DeleteAPIKey(ctx, apiKeyID)
	if err != nil {
		// If the key is already gone, that's a successful delete from Terraform's perspective.
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("API Key with ID %s already deleted, removing from state", apiKeyID))
			resp.State.RemoveResource(ctx) // Ensure it's removed if not already
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete API key %s, got error: %s", apiKeyID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("API Key with ID %s deleted successfully", apiKeyID))
	// If the delete is successful, the resource is automatically removed from state by Terraform.
}

func (r *APIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
