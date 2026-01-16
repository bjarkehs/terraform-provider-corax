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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-corax/internal/coraxclient"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ModelDeploymentResource{}
var _ resource.ResourceWithImportState = &ModelDeploymentResource{}

func NewModelDeploymentResource() resource.Resource {
	return &ModelDeploymentResource{}
}

// ModelDeploymentResource defines the resource implementation.
type ModelDeploymentResource struct {
	client *coraxclient.Client
}

// ModelDeploymentResourceModel describes the resource data model.
type ModelDeploymentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`     // Nullable
	SupportedTasks types.List   `tfsdk:"supported_tasks"` // List of strings
	Configuration  types.Map    `tfsdk:"configuration"`   // Map of string to string
	IsActive       types.Bool   `tfsdk:"is_active"`
	ProviderID     types.String `tfsdk:"provider_id"`
}

func (r *ModelDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model_deployment"
}

func (r *ModelDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Model Deployment. Model Deployments link a specific model configuration from a Model Provider to be usable for certain tasks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the model deployment (UUID).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A user-defined name for the model deployment.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional description for the model deployment.",
			},
			"supported_tasks": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "A list of tasks this model deployment supports (e.g., 'chat', 'completion', 'embedding').",
				// TODO: Add validator for allowed enum values if strictly defined by API, or leave as free strings.
				// OpenAPI spec: items: {$ref: "#/components/schemas/CapabilityType"}
				// CapabilityType enum: ["chat", "completion", "embedding"]
			},
			"configuration": schema.MapAttribute{
				ElementType:         types.StringType, // Assuming string values for simplicity. API says object with additionalProperties.
				Required:            true,
				MarkdownDescription: "Configuration key-value pairs specific to the model deployment (e.g., model name, API version for Azure OpenAI).",
			},
			"is_active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Indicates whether the model deployment is active and usable. Defaults to true.",
			},
			"provider_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the Model Provider this deployment belongs to.",
				// TODO: Add validator for UUID format
			},
		},
	}
}

func (r *ModelDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func modelDeploymentResourceModelToAPICreate(ctx context.Context, plan ModelDeploymentResourceModel, diags *diag.Diagnostics) (*coraxclient.ModelDeploymentCreate, error) {
	apiCreate := &coraxclient.ModelDeploymentCreate{
		Name:       plan.Name.ValueString(),
		ProviderID: plan.ProviderID.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		desc := plan.Description.ValueString()
		apiCreate.Description = &desc
	}

	if !plan.IsActive.IsNull() && !plan.IsActive.IsUnknown() {
		isActive := plan.IsActive.ValueBool()
		apiCreate.IsActive = &isActive
	}

	respDiags := plan.SupportedTasks.ElementsAs(ctx, &apiCreate.SupportedTasks, false)
	diags.Append(respDiags...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert supported_tasks")
	}

	configMap := make(map[string]string)
	respDiags = plan.Configuration.ElementsAs(ctx, &configMap, false)
	diags.Append(respDiags...)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert configuration")
	}
	apiCreate.Configuration = configMap

	return apiCreate, nil
}

// Helper to map TF model to API Update struct.
// The API requires all fields for PUT (full replacement), so we always send the complete state.
func modelDeploymentResourceModelToAPIUpdate(ctx context.Context, plan ModelDeploymentResourceModel, state ModelDeploymentResourceModel, diags *diag.Diagnostics) (*coraxclient.ModelDeploymentUpdate, bool, error) {
	// Check if ProviderID is being changed (not allowed)
	if !plan.ProviderID.Equal(state.ProviderID) {
		diags.AddWarning("ProviderID Change", "ProviderID cannot be updated for a model deployment. This change will be ignored.")
	}

	// Check if any update is needed by comparing plan to state
	updateNeeded := !plan.Name.Equal(state.Name) ||
		!plan.Description.Equal(state.Description) ||
		!plan.IsActive.Equal(state.IsActive) ||
		!plan.SupportedTasks.Equal(state.SupportedTasks) ||
		!plan.Configuration.Equal(state.Configuration)

	if !updateNeeded {
		return nil, false, nil
	}

	// Build the full update payload from plan (API requires all fields)
	apiUpdate := &coraxclient.ModelDeploymentUpdate{
		Name:       plan.Name.ValueString(),
		ProviderID: state.ProviderID.ValueString(), // Use state, as ProviderID is not updatable
	}

	// Handle optional description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		desc := plan.Description.ValueString()
		apiUpdate.Description = &desc
	}

	// Handle optional is_active
	if !plan.IsActive.IsNull() && !plan.IsActive.IsUnknown() {
		isActive := plan.IsActive.ValueBool()
		apiUpdate.IsActive = &isActive
	}

	// Convert supported_tasks from TF list to string slice
	respDiags := plan.SupportedTasks.ElementsAs(ctx, &apiUpdate.SupportedTasks, false)
	diags.Append(respDiags...)
	if diags.HasError() {
		return nil, false, fmt.Errorf("failed to convert supported_tasks for update")
	}

	// Convert configuration from TF map to string map
	configMap := make(map[string]string)
	respDiags = plan.Configuration.ElementsAs(ctx, &configMap, false)
	diags.Append(respDiags...)
	if diags.HasError() {
		return nil, false, fmt.Errorf("failed to convert configuration for update")
	}
	apiUpdate.Configuration = configMap

	return apiUpdate, true, nil
}

// Helper to map API response to TF model.
func mapAPIModelDeploymentToResourceModel(ctx context.Context, apiDeployment *coraxclient.ModelDeployment, model *ModelDeploymentResourceModel, diags *diag.Diagnostics) {
	model.ID = types.StringValue(apiDeployment.ID)
	model.Name = types.StringValue(apiDeployment.Name)
	model.ProviderID = types.StringValue(apiDeployment.ProviderID)

	if apiDeployment.Description != nil {
		model.Description = types.StringValue(*apiDeployment.Description)
	} else {
		model.Description = types.StringNull()
	}
	if apiDeployment.IsActive != nil {
		model.IsActive = types.BoolValue(*apiDeployment.IsActive)
	} else {
		model.IsActive = types.BoolValue(true) // Default
	}

	supportedTasks, listDiags := types.ListValueFrom(ctx, types.StringType, apiDeployment.SupportedTasks)
	diags.Append(listDiags...)
	model.SupportedTasks = supportedTasks

	configMap, mapDiags := types.MapValueFrom(ctx, types.StringType, apiDeployment.Configuration)
	diags.Append(mapDiags...)
	model.Configuration = configMap
}

func (r *ModelDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiCreatePayload, err := modelDeploymentResourceModelToAPICreate(ctx, plan, &resp.Diagnostics)
	if err != nil {
		// Diagnostics already appended by helper
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating Model Deployment: %s", apiCreatePayload.Name))
	createdDeployment, err := r.client.CreateModelDeployment(ctx, *apiCreatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create model deployment, got error: %s", err))
		return
	}

	mapAPIModelDeploymentToResourceModel(ctx, createdDeployment, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Model Deployment %s created successfully with ID %s", plan.Name.ValueString(), plan.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading Model Deployment with ID: %s", deploymentID))

	apiDeployment, err := r.client.GetModelDeployment(ctx, deploymentID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Model Deployment %s not found, removing from state", deploymentID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read model deployment %s: %s", deploymentID, err))
		return
	}

	mapAPIModelDeploymentToResourceModel(ctx, apiDeployment, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Model Deployment %s", deploymentID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ModelDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating Model Deployment with ID: %s", deploymentID))

	apiUpdatePayload, updateNeeded, err := modelDeploymentResourceModelToAPIUpdate(ctx, plan, state, &resp.Diagnostics)
	if err != nil {
		// Diagnostics already appended
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if !updateNeeded {
		tflog.Debug(ctx, "No attribute changes detected for Model Deployment update.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Ensure state matches plan if no API call
		return
	}

	updatedDeployment, err := r.client.UpdateModelDeployment(ctx, deploymentID, *apiUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update model deployment %s: %s", deploymentID, err))
		return
	}

	mapAPIModelDeploymentToResourceModel(ctx, updatedDeployment, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Model Deployment %s updated successfully", deploymentID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting Model Deployment with ID: %s", deploymentID))

	err := r.client.DeleteModelDeployment(ctx, deploymentID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Model Deployment %s not found, already deleted", deploymentID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete model deployment %s: %s", deploymentID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Model Deployment %s deleted successfully", deploymentID))
}

func (r *ModelDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
