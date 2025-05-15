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
	CreatedAt     types.String `tfsdk:"created_at"`    // Computed
	UpdatedAt     types.String `tfsdk:"updated_at"`    // Computed, Nullable
	CreatedBy     types.String `tfsdk:"created_by"`    // Computed
	UpdatedBy     types.String `tfsdk:"updated_by"`    // Computed, Nullable
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
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp of the model provider.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last update timestamp of the model provider.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User who created the model provider.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User who last updated the model provider.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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

// Helper to map TF model to API Create struct
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

// Helper to map TF model to API Update struct
// The API spec for ModelProviderUpdate implies all fields are required for PUT.
// This helper will construct a full object based on the plan.
func modelProviderResourceModelToAPIUpdate(ctx context.Context, plan ModelProviderResourceModel, diags *diag.Diagnostics) (*coraxclient.ModelProviderUpdate, error) {
	apiUpdate := &coraxclient.ModelProviderUpdate{
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

// Helper to map API response to TF model
func mapAPIModelProviderToResourceModel(ctx context.Context, apiProvider *coraxclient.ModelProvider, model *ModelProviderResourceModel, diags *diag.Diagnostics) {
	model.ID = types.StringValue(apiProvider.ID)
	model.Name = types.StringValue(apiProvider.Name)
	model.ProviderType = types.StringValue(apiProvider.ProviderType)

	configMap, mapDiags := types.MapValueFrom(ctx, types.StringType, apiProvider.Configuration)
	diags.Append(mapDiags...)
	model.Configuration = configMap

	model.CreatedAt = types.StringValue(apiProvider.CreatedAt)
	model.CreatedBy = types.StringValue(apiProvider.CreatedBy)

	if apiProvider.UpdatedAt != nil {
		model.UpdatedAt = types.StringValue(*apiProvider.UpdatedAt)
	} else {
		model.UpdatedAt = types.StringNull()
	}
	if apiProvider.UpdatedBy != nil {
		model.UpdatedBy = types.StringValue(*apiProvider.UpdatedBy)
	} else {
		model.UpdatedBy = types.StringNull()
	}
}

func (r *ModelProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ModelProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	tflog.Info(ctx, fmt.Sprintf("Model Provider %s created successfully with ID %s", plan.Name.ValueString(), plan.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ModelProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Model Provider %s", providerID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ModelProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ModelProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerID := plan.ID.ValueString() // ID comes from plan/state, not updatable
	tflog.Debug(ctx, fmt.Sprintf("Updating Model Provider with ID: %s", providerID))

	// API spec for ModelProviderUpdate implies all fields are required for PUT.
	// Construct the full payload from the plan.
	apiUpdatePayload, err := modelProviderResourceModelToAPIUpdate(ctx, plan, &resp.Diagnostics)
	if err != nil {
		return // Diagnostics already handled
	}
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProvider, err := r.client.UpdateModelProvider(ctx, providerID, *apiUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update model provider %s: %s", providerID, err))
		return
	}

	mapAPIModelProviderToResourceModel(ctx, updatedProvider, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Model Provider %s updated successfully", providerID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
