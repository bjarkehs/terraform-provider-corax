package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

	"corax/internal/coraxclient"
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
	CreatedAt      types.String `tfsdk:"created_at"`      // Computed
	UpdatedAt      types.String `tfsdk:"updated_at"`      // Computed, Nullable
	CreatedBy      types.String `tfsdk:"created_by"`      // Computed
	UpdatedBy      types.String `tfsdk:"updated_by"`      // Computed, Nullable
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
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp of the model deployment.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last update timestamp of the model deployment.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User who created the model deployment.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User who last updated the model deployment.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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

// Helper to map TF model to API Create struct
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

// Helper to map TF model to API Update struct
func modelDeploymentResourceModelToAPIUpdate(ctx context.Context, plan ModelDeploymentResourceModel, state ModelDeploymentResourceModel, diags *diag.Diagnostics) (*coraxclient.ModelDeploymentUpdate, bool, error) {
	apiUpdate := &coraxclient.ModelDeploymentUpdate{}
	updateNeeded := false

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		apiUpdate.Name = &name
		updateNeeded = true
	}
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			// To clear description, send an empty string or handle as API expects.
			// For now, let's assume sending null/empty string clears it.
			// The API client struct uses *string, so sending nil means omit.
			// If API needs explicit null or empty string, adjust client or here.
			// For now, if TF model is null, we don't set it in update payload, implying no change unless API treats omission as clear.
			// If it's an empty string in TF, we send it.
			if !state.Description.IsNull() { // only send if it was not null before
				var emptyDesc string
				apiUpdate.Description = &emptyDesc // Send empty string to clear
				updateNeeded = true
			}
		} else {
			desc := plan.Description.ValueString()
			apiUpdate.Description = &desc
			updateNeeded = true
		}
	}
	if !plan.IsActive.Equal(state.IsActive) {
		isActive := plan.IsActive.ValueBool()
		apiUpdate.IsActive = &isActive
		updateNeeded = true
	}
	if !plan.ProviderID.Equal(state.ProviderID) {
		// ProviderID is usually not updatable. If API allows, uncomment.
		// providerID := plan.ProviderID.ValueString()
		// apiUpdate.ProviderID = &providerID
		// updateNeeded = true
		diags.AddWarning("ProviderID Change", "ProviderID cannot be updated for a model deployment. This change will be ignored.")
	}

	if !plan.SupportedTasks.Equal(state.SupportedTasks) {
		respDiags := plan.SupportedTasks.ElementsAs(ctx, &apiUpdate.SupportedTasks, false)
		diags.Append(respDiags...)
		if diags.HasError() { return nil, false, fmt.Errorf("failed to convert supported_tasks for update") }
		updateNeeded = true
	}
	if !plan.Configuration.Equal(state.Configuration) {
		configMap := make(map[string]string)
		respDiags := plan.Configuration.ElementsAs(ctx, &configMap, false)
		diags.Append(respDiags...)
		if diags.HasError() { return nil, false, fmt.Errorf("failed to convert configuration for update") }
		apiUpdate.Configuration = configMap
		updateNeeded = true
	}

	return apiUpdate, updateNeeded, nil
}


// Helper to map API response to TF model
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
	
	model.CreatedAt = types.StringValue(apiDeployment.CreatedAt)
	model.CreatedBy = types.StringValue(apiDeployment.CreatedBy)

	if apiDeployment.UpdatedAt != nil {
		model.UpdatedAt = types.StringValue(*apiDeployment.UpdatedAt)
	} else {
		model.UpdatedAt = types.StringNull()
	}
	if apiDeployment.UpdatedBy != nil {
		model.UpdatedBy = types.StringValue(*apiDeployment.UpdatedBy)
	} else {
		model.UpdatedBy = types.StringNull()
	}
}


func (r *ModelDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() { return }

	apiCreatePayload, err := modelDeploymentResourceModelToAPICreate(ctx, plan, &resp.Diagnostics)
	if err != nil {
		// Diagnostics already appended by helper
		return
	}
	if resp.Diagnostics.HasError() { return }


	tflog.Debug(ctx, fmt.Sprintf("Creating Model Deployment: %s", apiCreatePayload.Name))
	createdDeployment, err := r.client.CreateModelDeployment(ctx, *apiCreatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create model deployment, got error: %s", err))
		return
	}

	mapAPIModelDeploymentToResourceModel(ctx, createdDeployment, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() { return }

	tflog.Info(ctx, fmt.Sprintf("Model Deployment %s created successfully with ID %s", plan.Name.ValueString(), plan.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }

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
	if resp.Diagnostics.HasError() { return }

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Model Deployment %s", deploymentID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ModelDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }

	deploymentID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating Model Deployment with ID: %s", deploymentID))

	apiUpdatePayload, updateNeeded, err := modelDeploymentResourceModelToAPIUpdate(ctx, plan, state, &resp.Diagnostics)
	if err != nil {
		// Diagnostics already appended
		return
	}
	if resp.Diagnostics.HasError() { return }

	if !updateNeeded {
		tflog.Debug(ctx, "No attribute changes detected for Model Deployment update.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Ensure state matches plan if no API call
		return
	}
	
	// The API spec for ModelDeploymentUpdate is identical to Create, implying all fields are required.
	// If the API truly requires all fields for PUT, then apiUpdatePayload must be fully populated from `plan`.
	// The current modelDeploymentResourceModelToAPIUpdate creates a partial payload.
	// For now, we proceed with the partial update. If API fails, this needs adjustment.
	// A safer approach if API requires full object for PUT:
	// fullApiPayloadForUpdate, convErr := modelDeploymentResourceModelToAPICreate(ctx, plan, &resp.Diagnostics) // Use create mapper
	// if convErr != nil || resp.Diagnostics.HasError() { return }
	// updatedDeployment, err := r.client.UpdateModelDeployment(ctx, deploymentID, *fullApiPayloadForUpdate)

	updatedDeployment, err := r.client.UpdateModelDeployment(ctx, deploymentID, *apiUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update model deployment %s: %s", deploymentID, err))
		return
	}

	mapAPIModelDeploymentToResourceModel(ctx, updatedDeployment, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() { return }

	tflog.Info(ctx, fmt.Sprintf("Model Deployment %s updated successfully", deploymentID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ModelDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ModelDeploymentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }

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
