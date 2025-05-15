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

	"corax-terraform-provider/internal/coraxclient"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ChatCapabilityResource{}
var _ resource.ResourceWithImportState = &ChatCapabilityResource{}

func NewChatCapabilityResource() resource.Resource {
	return &ChatCapabilityResource{}
}

// ChatCapabilityResource defines the resource implementation.
type ChatCapabilityResource struct {
	client *coraxclient.Client
}

// ChatCapabilityResourceModel describes the resource data model.
type ChatCapabilityResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	IsPublic     types.Bool   `tfsdk:"is_public"`
	ModelID      types.String `tfsdk:"model_id"`   // Nullable
	Config       types.Object `tfsdk:"config"`     // Nullable
	ProjectID    types.String `tfsdk:"project_id"` // Nullable
	SystemPrompt types.String `tfsdk:"system_prompt"`
	// CollectionIDs types.List   `tfsdk:"collection_ids"` // Omitted for now as per decision to skip collection-related features
	CreatedBy  types.String `tfsdk:"created_by"`  // Computed
	UpdatedBy  types.String `tfsdk:"updated_by"`  // Computed
	CreatedAt  types.String `tfsdk:"created_at"`  // Computed
	UpdatedAt  types.String `tfsdk:"updated_at"`  // Computed
	ArchivedAt types.String `tfsdk:"archived_at"` // Computed, Nullable
	Owner      types.String `tfsdk:"owner"`       // Computed
	Type       types.String `tfsdk:"type"`        // Computed, should always be "chat"
}

func (r *ChatCapabilityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chat_capability"
}

func (r *ChatCapabilityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Chat Capability. Chat capabilities define configurations for conversational AI models.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the chat capability (UUID).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A user-defined name for the chat capability.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"is_public": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether the capability is publicly accessible. Defaults to false.",
			},
			"model_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The UUID of the model deployment to use for this capability. If not provided, a default model for 'chat' type may be used by the API.",
				// TODO: Add validator for UUID format
			},
			"project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The UUID of the project this capability belongs to. If not provided, it might be associated with a default or no project.",
				// TODO: Add validator for UUID format
			},
			"system_prompt": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The system prompt that guides the behavior of the chat model.",
			},
			// "collection_ids": schema.ListAttribute{ // Omitted for now
			// 	ElementType:         types.StringType,
			// 	Optional:            true,
			// 	MarkdownDescription: "A list of collection UUIDs to be used for retrieval augmentation (RAG) by this chat capability.",
			// },
			"config": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration settings for the capability's behavior.",
				Attributes:          capabilityConfigSchemaAttributes(), // Use shared schema attributes
			},
			"created_by":  schema.StringAttribute{Computed: true, MarkdownDescription: "User who created the capability.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"updated_by":  schema.StringAttribute{Computed: true, MarkdownDescription: "User who last updated the capability.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"created_at":  schema.StringAttribute{Computed: true, MarkdownDescription: "Creation timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"updated_at":  schema.StringAttribute{Computed: true, MarkdownDescription: "Last update timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"archived_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Archival timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"owner":       schema.StringAttribute{Computed: true, MarkdownDescription: "Owner of the capability.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"type":        schema.StringAttribute{Computed: true, MarkdownDescription: "Type of the capability (should be 'chat').", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

func (r *ChatCapabilityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Helper functions for mapping (capabilityConfigModelToAPI, capabilityConfigAPItoModel are now in common_capability_config.go)

func mapAPICapabilityToChatModel(apiCap *coraxclient.CapabilityRepresentation, model *ChatCapabilityResourceModel, diags *diag.Diagnostics, ctx context.Context) {
	model.ID = types.StringValue(apiCap.ID)
	model.Name = types.StringValue(apiCap.Name)
	model.IsPublic = types.BoolValue(apiCap.IsPublic != nil && *apiCap.IsPublic) // API default is false
	model.Type = types.StringValue(apiCap.Type)

	if apiCap.ModelID != nil {
		model.ModelID = types.StringValue(*apiCap.ModelID)
	} else {
		model.ModelID = types.StringNull()
	}
	if apiCap.ProjectID != nil {
		model.ProjectID = types.StringValue(*apiCap.ProjectID)
	} else {
		model.ProjectID = types.StringNull()
	}

	// SystemPrompt is likely in apiCap.Configuration for chat capabilities
	// This needs to be confirmed based on actual API response structure.
	// Assuming it's directly in `configuration` map for now.
	if sysPrompt, ok := apiCap.Configuration["system_prompt"].(string); ok {
		model.SystemPrompt = types.StringValue(sysPrompt)
	} else {
		// This might indicate an issue if system_prompt is expected for chat type
		// Or it might be truly optional in some API views. For TF, it's required in schema.
		// If it's missing on read for an existing resource, it's problematic.
		// For now, if not found, make it null/unknown and let TF diff handle it.
		model.SystemPrompt = types.StringUnknown()
		tflog.Warn(ctx, fmt.Sprintf("System prompt not found in API response configuration for capability %s", apiCap.ID))
	}

	model.Config = capabilityConfigAPItoModel(ctx, apiCap.Config, diags)

	model.CreatedBy = types.StringValue(apiCap.CreatedBy)
	model.Owner = types.StringValue(apiCap.Owner)
	model.CreatedAt = types.StringValue(apiCap.CreatedAt)
	model.UpdatedAt = types.StringValue(apiCap.UpdatedAt) // API spec says non-null, but can be same as CreatedAt

	if apiCap.UpdatedBy != "" { // API spec says string, not *string
		model.UpdatedBy = types.StringValue(apiCap.UpdatedBy)
	} else {
		model.UpdatedBy = types.StringNull() // Or types.StringValue(apiCap.CreatedBy) if updatedby is never null
	}
	if apiCap.ArchivedAt != nil {
		model.ArchivedAt = types.StringValue(*apiCap.ArchivedAt)
	} else {
		model.ArchivedAt = types.StringNull()
	}
}

func (r *ChatCapabilityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ChatCapabilityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating Chat Capability: %s", plan.Name.ValueString()))

	apiPayload := coraxclient.ChatCapabilityCreate{
		Name:         plan.Name.ValueString(),
		Type:         "chat", // Hardcoded for this resource
		SystemPrompt: plan.SystemPrompt.ValueString(),
	}

	if !plan.IsPublic.IsNull() && !plan.IsPublic.IsUnknown() {
		isPublic := plan.IsPublic.ValueBool()
		apiPayload.IsPublic = &isPublic
	}
	if !plan.ModelID.IsNull() && !plan.ModelID.IsUnknown() {
		modelID := plan.ModelID.ValueString()
		apiPayload.ModelID = &modelID
	}
	if !plan.ProjectID.IsNull() && !plan.ProjectID.IsUnknown() {
		projectID := plan.ProjectID.ValueString()
		apiPayload.ProjectID = &projectID
	}

	apiPayload.Config = capabilityConfigModelToAPI(ctx, plan.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createdAPICap, err := r.client.CreateCapability(ctx, apiPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create chat capability, got error: %s", err))
		return
	}

	mapAPICapabilityToChatModel(createdAPICap, &plan, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Chat Capability %s created successfully with ID %s", plan.Name.ValueString(), plan.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ChatCapabilityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ChatCapabilityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	capabilityID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading Chat Capability with ID: %s", capabilityID))

	apiCap, err := r.client.GetCapability(ctx, capabilityID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Chat Capability %s not found, removing from state", capabilityID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read chat capability %s: %s", capabilityID, err))
		return
	}

	if apiCap.Type != "chat" {
		resp.Diagnostics.AddError("Resource Type Mismatch", fmt.Sprintf("Expected capability type 'chat' but found '%s' for ID %s. Removing from state.", apiCap.Type, capabilityID))
		resp.State.RemoveResource(ctx)
		return
	}

	//currentConfig := state.Config // Preserve potentially more detailed config from state if API is lossy

	mapAPICapabilityToChatModel(apiCap, &state, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	// If API returns a less detailed config, try to merge or prefer state if certain fields are not returned by GET
	// For now, mapAPICapabilityToChatModel will overwrite. If specific config fields are write-only,
	// they would need to be handled by preserving from `currentConfig`.
	// Example: if apiCap.Config is nil but currentConfig was not, we might want to keep currentConfig.
	// This depends on API behavior for GET /capabilities/{id} regarding the 'config' field.
	// The current mapping helper `capabilityConfigAPItoModel` handles nil apiConfig.

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Chat Capability %s", capabilityID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ChatCapabilityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ChatCapabilityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	capabilityID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating Chat Capability with ID: %s", capabilityID))

	updatePayload := coraxclient.ChatCapabilityUpdate{}
	updateNeeded := false

	if !plan.Name.Equal(state.Name) {
		val := plan.Name.ValueString()
		updatePayload.Name = &val
		updateNeeded = true
	}
	if !plan.IsPublic.Equal(state.IsPublic) {
		val := plan.IsPublic.ValueBool()
		updatePayload.IsPublic = &val
		updateNeeded = true
	}
	if !plan.ModelID.Equal(state.ModelID) {
		if plan.ModelID.IsNull() {
			var nilStr *string
			updatePayload.ModelID = nilStr
		} else {
			val := plan.ModelID.ValueString()
			updatePayload.ModelID = &val
		}
		updateNeeded = true
	}
	if !plan.ProjectID.Equal(state.ProjectID) {
		if plan.ProjectID.IsNull() {
			var nilStr *string
			updatePayload.ProjectID = nilStr
		} else {
			val := plan.ProjectID.ValueString()
			updatePayload.ProjectID = &val
		}
		updateNeeded = true
	}
	if !plan.SystemPrompt.Equal(state.SystemPrompt) {
		val := plan.SystemPrompt.ValueString()
		updatePayload.SystemPrompt = &val
		updateNeeded = true
	}

	// Config update
	if !plan.Config.Equal(state.Config) {
		updatePayload.Config = capabilityConfigModelToAPI(ctx, plan.Config, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		updateNeeded = true // even if config becomes null
	}

	if !updateNeeded {
		tflog.Debug(ctx, "No attribute changes detected for Chat Capability update.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Ensure state matches plan
		return
	}

	// Ensure type is "chat" if sending, though API might ignore it on PUT
	// chatType := "chat"
	// updatePayload.Type = &chatType // API schema for Update doesn't show type as updatable for specific types.

	updatedAPICap, err := r.client.UpdateCapability(ctx, capabilityID, updatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update chat capability %s: %s", capabilityID, err))
		return
	}

	mapAPICapabilityToChatModel(updatedAPICap, &plan, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Chat Capability %s updated successfully", capabilityID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ChatCapabilityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ChatCapabilityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	capabilityID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting Chat Capability with ID: %s", capabilityID))

	err := r.client.DeleteCapability(ctx, capabilityID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Chat Capability %s not found, already deleted", capabilityID))
			resp.State.RemoveResource(ctx) // Remove from state if not found
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete chat capability %s: %s", capabilityID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Chat Capability %s deleted successfully", capabilityID))
}

func (r *ChatCapabilityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	// After ID is set, Read will be called. Read needs to verify the type is "chat".
}
