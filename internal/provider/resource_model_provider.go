// Copyright (c) Trifork

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-corax/internal/coraxclient"
)

const apiKeyConfigurationKey = "api_key"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ModelProviderResource{}
var _ resource.ResourceWithImportState = &ModelProviderResource{}

func NewModelProviderResource() resource.Resource {
	return &ModelProviderResource{}
}

// ModelProviderResource defines the resource implementation.
type ModelProviderResource struct {
	client *coraxclient.Client
}

// ModelProviderResourceModel describes the resource data model.
type ModelProviderResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ProviderType  types.String `tfsdk:"provider_type"`
	Configuration types.Map    `tfsdk:"configuration"` // Map of string to string, some values might be sensitive
}

func (r *ModelProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model_provider"
}

func (r *ModelProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Model Provider. Model Providers store configurations (like API keys and endpoints) for different LLM providers (e.g., Azure OpenAI, OpenAI, Bedrock).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the model provider (UUID).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A user-defined name for the model provider instance.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"provider_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of the model provider (e.g., 'azure_openai', 'openai', 'bedrock'). This should match a type known to the Corax API.",
				// TODO: Consider a validator if the list of types is fixed and small, or link to a data source for valid types.
			},
			"configuration": schema.MapAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Configuration key-value pairs for the model provider. Specific keys depend on the `provider_type`. For example, 'api_key', 'api_endpoint'. Some values may be sensitive.",
				Sensitive:           true, // Mark the whole map as sensitive as it often contains API keys.
			},
		},
	}
}

func (r *ModelProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coraxclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *coraxclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	r.client = client
}

// Helper to map TF model to API Create struct.
func modelProviderResourceModelToAPICreate(ctx context.Context, plan ModelProviderResourceModel, diags *diag.Diagnostics) (*coraxclient.ModelProviderCreate, error) {
	apiCreate := &coraxclient.ModelProviderCreate{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
	}

	configMap := make(map[string]string)
	respDiags := plan.Configuration.ElementsAs(ctx, &configMap, false)
	diags.Append(respDiags...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert configuration")
	}
	apiCreate.Configuration = configMap

	return apiCreate, nil
}

// Helper to map TF model to API Update struct.
// The API spec for ModelProviderUpdate implies all fields are required for PUT.
// This helper will construct a full object based on the plan.
func modelProviderResourceModelToAPIUpdate(ctx context.Context, plan ModelProviderResourceModel, diags *diag.Diagnostics) (*coraxclient.ModelProviderUpdate, error) {
	apiUpdate := &coraxclient.ModelProviderUpdate{
		ID:           plan.ID.ValueString(), // TODO: ID is currently required for update?
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
	}

	configMap := make(map[string]string)
	respDiags := plan.Configuration.ElementsAs(ctx, &configMap, false)
	diags.Append(respDiags...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert configuration for update")
	}
	apiUpdate.Configuration = configMap

	return apiUpdate, nil
}

// Helper to map API response to TF model.
func mapAPIModelProviderToResourceModel(ctx context.Context, apiProvider *coraxclient.ModelProvider, model *ModelProviderResourceModel, diags *diag.Diagnostics) {
	model.ID = types.StringValue(apiProvider.ID)
	model.Name = types.StringValue(apiProvider.Name)
	model.ProviderType = types.StringValue(apiProvider.ProviderType)

	configMap, mapDiags := types.MapValueFrom(ctx, types.StringType, apiProvider.Configuration)
	tflog.Debug(ctx, fmt.Sprintf("Mapping configuration: %v", configMap))
	diags.Append(mapDiags...)
	model.Configuration = configMap
}

func (r *ModelProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ModelProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Store the planned configuration to preserve sensitive values like the full API key
	plannedConfiguration := plan.Configuration

	apiCreatePayload, err := modelProviderResourceModelToAPICreate(ctx, plan, &resp.Diagnostics)
	if err != nil {
		return // Diagnostics already handled
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating Model Provider: %s", apiCreatePayload.Name))
	createdProvider, err := r.client.CreateModelProvider(ctx, *apiCreatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create model provider, got error: %s", err))
		return
	}

	mapAPIModelProviderToResourceModel(ctx, createdProvider, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the planned configuration for "api_key" was set, ensure it's preserved
	// over any potentially truncated value returned by the API.
	if !plannedConfiguration.IsNull() && !plannedConfiguration.IsUnknown() {
		plannedConfigMap := make(map[string]string)
		diags := plannedConfiguration.ElementsAs(ctx, &plannedConfigMap, false)
		resp.Diagnostics.Append(diags...)

		if !resp.Diagnostics.HasError() {
			if fullAPIKey, ok := plannedConfigMap[apiKeyConfigurationKey]; ok && fullAPIKey != "" {
				currentConfigMap := make(map[string]string)
				// plan.Configuration might be null/unknown if API returned nothing for config
				if !plan.Configuration.IsNull() && !plan.Configuration.IsUnknown() {
					diags = plan.Configuration.ElementsAs(ctx, &currentConfigMap, false)
					resp.Diagnostics.Append(diags...)
				}

				if !resp.Diagnostics.HasError() {
					currentConfigMap[apiKeyConfigurationKey] = fullAPIKey // Overwrite with full key
					updatedConfigTFMap, mapDiags := types.MapValueFrom(ctx, types.StringType, currentConfigMap)
					resp.Diagnostics.Append(mapDiags...)
					if !resp.Diagnostics.HasError() {
						plan.Configuration = updatedConfigTFMap
					}
				}
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Model Provider %s created successfully with ID %s", plan.Name.ValueString(), plan.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ModelProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Store the prior state's configuration to preserve sensitive values like the full API key
	priorStateConfiguration := state.Configuration

	providerID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading Model Provider with ID: %s", providerID))

	apiProvider, err := r.client.GetModelProvider(ctx, providerID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Model Provider %s not found, removing from state", providerID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read model provider %s: %s", providerID, err))
		return
	}

	mapAPIModelProviderToResourceModel(ctx, apiProvider, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the prior state's configuration for "api_key" was set, ensure it's preserved
	// over any potentially truncated value returned by the API.
	if !priorStateConfiguration.IsNull() && !priorStateConfiguration.IsUnknown() {
		priorStateConfigMap := make(map[string]string)
		diags := priorStateConfiguration.ElementsAs(ctx, &priorStateConfigMap, false)
		resp.Diagnostics.Append(diags...)

		if !resp.Diagnostics.HasError() {
			if fullAPIKey, ok := priorStateConfigMap[apiKeyConfigurationKey]; ok && fullAPIKey != "" {
				currentConfigMap := make(map[string]string)
				// state.Configuration might be null/unknown if API returned nothing for config
				if !state.Configuration.IsNull() && !state.Configuration.IsUnknown() {
					diags = state.Configuration.ElementsAs(ctx, &currentConfigMap, false)
					resp.Diagnostics.Append(diags...)
				}

				if !resp.Diagnostics.HasError() {
					currentConfigMap[apiKeyConfigurationKey] = fullAPIKey // Overwrite with full key
					updatedConfigTFMap, mapDiags := types.MapValueFrom(ctx, types.StringType, currentConfigMap)
					resp.Diagnostics.Append(mapDiags...)
					if !resp.Diagnostics.HasError() {
						state.Configuration = updatedConfigTFMap
					}
				}
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Model Provider %s", providerID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ModelProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ModelProviderResourceModel // plan contains the configuration from the TF plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the configuration from the plan, as this is what we intend to set for sensitive data.
	// If the API modifies this (e.g. adds/removes keys, normalizes values),
	// using the planned configuration for the state can prevent "unexpected new value" errors
	// if those API modifications are not meant to be authoritative for the TF state.
	plannedConfiguration := plan.Configuration

	providerID := plan.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating Model Provider with ID: %s", providerID))

	apiUpdatePayload, err := modelProviderResourceModelToAPIUpdate(ctx, plan, &resp.Diagnostics)
	if err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProvider, err := r.client.UpdateModelProvider(ctx, providerID, *apiUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update model provider %s: %s", providerID, err))
		return
	}

	// Map the API response to a temporary model to get computed values
	var stateFromServer ModelProviderResourceModel
	mapAPIModelProviderToResourceModel(ctx, updatedProvider, &stateFromServer, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construct the final state:
	// Start with the original plan (name, provider_type, etc. as planned).
	// Update computed fields from the server's response.
	// Crucially, set Configuration to what was planned.
	finalState := plan
	finalState.Configuration = plannedConfiguration // Use the planned configuration
	// Name and ProviderType are taken from the 'plan' variable, which reflects the user's intent.
	// ID is not expected to change on update / is UseStateForUnknown or immutable.

	tflog.Info(ctx, fmt.Sprintf("Model Provider %s updated successfully", providerID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

func (r *ModelProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ModelProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting Model Provider with ID: %s", providerID))

	err := r.client.DeleteModelProvider(ctx, providerID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Model Provider %s not found, already deleted", providerID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete model provider %s: %s", providerID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Model Provider %s deleted successfully", providerID))
}

func (r *ModelProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
