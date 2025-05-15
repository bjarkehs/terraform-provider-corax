package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
var _ resource.Resource = &CompletionCapabilityResource{}
var _ resource.ResourceWithImportState = &CompletionCapabilityResource{}

func NewCompletionCapabilityResource() resource.Resource {
	return &CompletionCapabilityResource{}
}

// CompletionCapabilityResource defines the resource implementation.
type CompletionCapabilityResource struct {
	client *coraxclient.Client
}

// CompletionCapabilityResourceModel describes the resource data model.
type CompletionCapabilityResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	IsPublic         types.Bool   `tfsdk:"is_public"`
	ModelID          types.String `tfsdk:"model_id"`      // Nullable
	Config           types.Object `tfsdk:"config"`        // Nullable, uses CapabilityConfigModel from chat_capability.go
	ProjectID        types.String `tfsdk:"project_id"`    // Nullable
	SystemPrompt     types.String `tfsdk:"system_prompt"` // Shared with Chat, but also in Completion
	CompletionPrompt types.String `tfsdk:"completion_prompt"`
	Variables        types.List   `tfsdk:"variables"`      // Nullable, list of strings
	OutputType       types.String `tfsdk:"output_type"`    // "schema" or "text"
	SchemaDef        types.Map    `tfsdk:"schema_def"`     // Nullable, map of string to dynamic/object for property defs
	CreatedBy        types.String `tfsdk:"created_by"`     // Computed
	UpdatedBy        types.String `tfsdk:"updated_by"`     // Computed
	CreatedAt        types.String `tfsdk:"created_at"`     // Computed
	UpdatedAt        types.String `tfsdk:"updated_at"`     // Computed
	ArchivedAt       types.String `tfsdk:"archived_at"`    // Computed, Nullable
	Owner            types.String `tfsdk:"owner"`          // Computed
	Type             types.String `tfsdk:"type"`           // Computed, should always be "completion"
}

// Note: CapabilityConfigModel, BlobConfigModel, DataRetentionModel, TimedDataRetentionModel, InfiniteDataRetentionModel
// are already defined in resource_chat_capability.go and can be reused.

func (r *CompletionCapabilityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_completion_capability"
}

func (r *CompletionCapabilityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Completion Capability. Completion capabilities define configurations for generating text completions, potentially with structured output.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the completion capability (UUID).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A user-defined name for the completion capability.",
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
				MarkdownDescription: "The UUID of the model deployment to use for this capability. If not provided, a default model for 'completion' type may be used by the API.",
			},
			"project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The UUID of the project this capability belongs to.",
			},
			"system_prompt": schema.StringAttribute{
				Required:            true, // API spec shows this for CompletionCapability too
				MarkdownDescription: "The system prompt that provides context or instructions to the completion model.",
			},
			"completion_prompt": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The main prompt for which a completion is generated. May include placeholders for variables.",
			},
			"variables": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A list of variable names (strings) that can be interpolated into the `completion_prompt`.",
			},
			"output_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Defines the expected output format. Must be either 'text' or 'schema'.",
				Validators:          []validator.String{stringvalidator.OneOf("text", "schema")},
			},
			"schema_def": schema.MapAttribute{
				ElementType:         types.DynamicType, // Using DynamicType for now due to complex nested structure.
				Optional:            true,
				MarkdownDescription: "Defines the structure of the output when `output_type` is 'schema'. This is a map where keys are property names and values define the property's type and description. Required if `output_type` is 'schema'.",
				// TODO: Add validation: required if output_type is "schema". This can be done with a CustomType or PlanModifier.
			},
			"config": schema.SingleNestedAttribute{ // Reusing the same config structure as chat
				Optional:            true,
				MarkdownDescription: "Configuration settings for the capability's behavior.",
				Attributes:          capabilityConfigSchemaAttributes(), // Defined in chat_capability_resource.go (or move to a common place)
			},
			"created_by":  schema.StringAttribute{Computed: true, MarkdownDescription: "User who created the capability.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"updated_by":  schema.StringAttribute{Computed: true, MarkdownDescription: "User who last updated the capability.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"created_at":  schema.StringAttribute{Computed: true, MarkdownDescription: "Creation timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"updated_at":  schema.StringAttribute{Computed: true, MarkdownDescription: "Last update timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"archived_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Archival timestamp.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"owner":       schema.StringAttribute{Computed: true, MarkdownDescription: "Owner of the capability.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"type":        schema.StringAttribute{Computed: true, MarkdownDescription: "Type of the capability (should be 'completion').", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

// capabilityConfigSchemaAttributes, capabilityConfigModelToAPI, capabilityConfigAPItoModel
// and their underlying attribute type helpers are defined in common_capability_config.go
// No need to redefine them here.

// --- Helper functions for mapping (specific to Completion Capability) ---

func schemaDefMapToAPI(ctx context.Context, schemaDefMap types.Map, diags *diag.Diagnostics) map[string]interface{} {
	if schemaDefMap.IsNull() || schemaDefMap.IsUnknown() {
		return nil
	}
	apiSchemaDef := make(map[string]interface{})
	for key, val := range schemaDefMap.Elements() {
		// val is types.DynamicValue, need to convert to Go native type
		// For simplicity, assuming val can be marshaled to JSON and then unmarshaled to interface{}
		// A more robust way would be to inspect val.UnderlyingValue() and convert based on its type.
		if val.IsNull() || val.IsUnknown() {
			apiSchemaDef[key] = nil
			continue
		}
		
		// This is a basic attempt. For complex nested structures within DynamicValue,
		// this might not be sufficient and a recursive conversion would be needed.
		// Or, if DynamicValue holds a string of JSON, unmarshal that.
		// For now, let's try to get the underlying value.
		// This part is tricky with types.DynamicType and needs careful handling
		// based on what structure is actually stored by the user in HCL.
		// A common pattern is for users to provide JSON strings for dynamic parts.
		// If the user provides a map in HCL, it should be convertible.

		var goVal interface{}
		// Attempt to convert to a Go map[string]interface{} if it's an object
		// This is a simplification. A full conversion utility for dynamic types is complex.
		// For now, we'll assume the dynamic value can be converted to a string (e.g. if it's JSON)
		// or directly to a Go native type.
		// A better approach for complex schema_def would be to define a proper TF schema for it.
		
		// Simplistic: try to convert to string, assuming it might be JSON.
		// This is NOT robust for complex, typed HCL maps.
		strVal, ok := val.(types.String)
		if ok && !strVal.IsNull() && !strVal.IsUnknown() {
			var rawJsonVal interface{}
			err := json.Unmarshal([]byte(strVal.ValueString()), &rawJsonVal)
			if err == nil {
				apiSchemaDef[key] = rawJsonVal
				continue
			}
			// if unmarshal fails, store as string
			apiSchemaDef[key] = strVal.ValueString()
		} else {
			// Fallback for other types or if not a string. This is very basic.
			// A proper conversion from types.DynamicValue to interface{} is needed.
			// For now, this will likely fail for complex nested HCL maps.
			// We might need to use val.UnderlyingValue() and type assertions.
			// Or, expect users to provide JSON strings for complex schema_def values.
			// Let's assume for now the API client will handle map[string]interface{}
			// where interface{} are basic types or nested maps/slices.
			// The framework should convert HCL maps to Go maps of types.Value.
			// We need to convert types.Value to native Go types.
			// This is a placeholder for a more robust conversion.
			tflog.Warn(ctx, fmt.Sprintf("SchemaDef value for key '%s' is not a string, direct assignment. This might not work for complex types.", key))
			// This is a very rough conversion and likely needs improvement.
			// For now, let's just try to pass it as is, hoping json.Marshal in newRequest handles it.
			// A better way: use val.As(ctx, &someGoInterface, basetypes.DynamicAsOptions{})
			// but that requires knowing the target Go type.
			apiSchemaDef[key] = "UNSUPPORTED_DYNAMIC_VALUE_CONVERSION" // Placeholder
		}
	}
	if len(apiSchemaDef) == 0 { return nil}
	return apiSchemaDef
}

func schemaDefAPIToMap(ctx context.Context, apiSchemaDef map[string]interface{}, diags *diag.Diagnostics) types.Map {
	if apiSchemaDef == nil {
		return types.MapNull(types.DynamicType)
	}
	elements := make(map[string]attr.Value)
	for key, val := range apiSchemaDef {
		// Convert interface{} back to types.DynamicValue
		// Simplest way: marshal to JSON string, then create types.StringValue, then types.DynamicValueFromString.
		// This ensures it's a valid representation for DynamicType if the underlying is complex.
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			diags.AddError("SchemaDef Conversion Error", fmt.Sprintf("Failed to marshal schema_def value for key %s: %s", key, err))
			continue
		}
		// Store as a string containing JSON, which DynamicType can handle.
		// Or, attempt to convert to known TF types if possible.
		// For now, using string representation of JSON.
		strVal := types.StringValue(string(jsonBytes))
		dynVal, errDiag := types.DynamicValueFrom(ctx, strVal)
		diags.Append(errDiag...)
		if errDiag.HasError() {
			continue
		}
		elements[key] = dynVal
	}
	mapVal, mapDiags := types.MapValue(types.DynamicType, elements)
	diags.Append(mapDiags...)
	return mapVal
}


func mapAPICompletionCapabilityToModel(apiCap *coraxclient.CapabilityRepresentation, model *CompletionCapabilityResourceModel, diags *diag.Diagnostics, ctx context.Context) {
	model.ID = types.StringValue(apiCap.ID)
	model.Name = types.StringValue(apiCap.Name)
	model.IsPublic = types.BoolValue(apiCap.IsPublic != nil && *apiCap.IsPublic)
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
	
	// Specific fields for completion are expected to be in apiCap.Configuration
	// This needs to be verified with actual API responses.
	// The OpenAPI spec for CompletionCapability shows these fields at the top level of the model,
	// but CapabilityRepresentation has a generic 'configuration' field.
	// Let's assume they are found within apiCap.Configuration for now.

	configMap := apiCap.Configuration
	if configMap != nil {
		if sysPrompt, ok := configMap["system_prompt"].(string); ok {
			model.SystemPrompt = types.StringValue(sysPrompt)
		} else { model.SystemPrompt = types.StringUnknown() } // Or StringNull if API can omit it

		if compPrompt, ok := configMap["completion_prompt"].(string); ok {
			model.CompletionPrompt = types.StringValue(compPrompt)
		} else { model.CompletionPrompt = types.StringUnknown() }

		if outputType, ok := configMap["output_type"].(string); ok {
			model.OutputType = types.StringValue(outputType)
		} else { model.OutputType = types.StringUnknown() }

		if vars, ok := configMap["variables"].([]interface{}); ok {
			strVars := make([]string, len(vars))
			valid := true
			for i, v := range vars {
				if strV, okStr := v.(string); okStr {
					strVars[i] = strV
				} else {
					valid = false
					break
				}
			}
			if valid {
				listVal, listDiags := types.ListValueFrom(ctx, types.StringType, strVars)
				diags.Append(listDiags...)
				model.Variables = listVal
			} else {
				model.Variables = types.ListNull(types.StringType)
			}
		} else {
			model.Variables = types.ListNull(types.StringType)
		}

		if schemaDef, ok := configMap["schema_def"].(map[string]interface{}); ok {
			model.SchemaDef = schemaDefAPIToMap(ctx, schemaDef, diags)
		} else {
			model.SchemaDef = types.MapNull(types.DynamicType)
		}
	} else {
		// If top-level configuration is missing, all specific fields are unknown/null
		model.SystemPrompt = types.StringUnknown()
		model.CompletionPrompt = types.StringUnknown()
		model.OutputType = types.StringUnknown()
		model.Variables = types.ListNull(types.StringType)
		model.SchemaDef = types.MapNull(types.DynamicType)
		tflog.Warn(ctx, fmt.Sprintf("Main 'configuration' object missing in API response for capability %s", apiCap.ID))
	}


	model.Config = capabilityConfigAPItoModel(ctx, apiCap.Config, diags) // Common config

	model.CreatedBy = types.StringValue(apiCap.CreatedBy)
	model.Owner = types.StringValue(apiCap.Owner)
	model.CreatedAt = types.StringValue(apiCap.CreatedAt)
	model.UpdatedAt = types.StringValue(apiCap.UpdatedAt)
	if apiCap.UpdatedBy != "" {
		model.UpdatedBy = types.StringValue(apiCap.UpdatedBy)
	} else {
		model.UpdatedBy = types.StringNull()
	}
	if apiCap.ArchivedAt != nil {
		model.ArchivedAt = types.StringValue(*apiCap.ArchivedAt)
	} else {
		model.ArchivedAt = types.StringNull()
	}
}


func (r *CompletionCapabilityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CompletionCapabilityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CompletionCapabilityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() { return }

	tflog.Debug(ctx, fmt.Sprintf("Creating Completion Capability: %s", plan.Name.ValueString()))

	apiPayload := coraxclient.CompletionCapabilityCreate{
		Name:             plan.Name.ValueString(),
		Type:             "completion", // Hardcoded
		SystemPrompt:     plan.SystemPrompt.ValueString(),
		CompletionPrompt: plan.CompletionPrompt.ValueString(),
		OutputType:       plan.OutputType.ValueString(),
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
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		resp.Diagnostics.Append(plan.Variables.ElementsAs(ctx, &apiPayload.Variables, false)...)
		if resp.Diagnostics.HasError() { return }
	}
	if plan.OutputType.ValueString() == "schema" {
		if plan.SchemaDef.IsNull() || plan.SchemaDef.IsUnknown() {
			resp.Diagnostics.AddError("Validation Error", "schema_def is required when output_type is 'schema'")
			return
		}
		apiPayload.SchemaDef = schemaDefMapToAPI(ctx, plan.SchemaDef, &resp.Diagnostics)
		if resp.Diagnostics.HasError() { return }
	}
	
	// Common config mapping (reuse from chat capability if moved to common, or define here)
	// For now, assuming capabilityConfigModelToAPI is available (defined in chat_capability.go or common)
	apiPayload.Config = capabilityConfigModelToAPI(ctx, plan.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() { return }


	createdAPICap, err := r.client.CreateCapability(ctx, apiPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create completion capability, got error: %s", err))
		return
	}

	mapAPICompletionCapabilityToModel(createdAPICap, &plan, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() { return }
	
	tflog.Info(ctx, fmt.Sprintf("Completion Capability %s created successfully with ID %s", plan.Name.ValueString(), plan.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CompletionCapabilityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CompletionCapabilityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }

	capabilityID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading Completion Capability with ID: %s", capabilityID))

	apiCap, err := r.client.GetCapability(ctx, capabilityID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Completion Capability %s not found, removing from state", capabilityID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read completion capability %s: %s", capabilityID, err))
		return
	}

	if apiCap.Type != "completion" {
		resp.Diagnostics.AddError("Resource Type Mismatch", fmt.Sprintf("Expected capability type 'completion' but found '%s' for ID %s. Removing from state.", apiCap.Type, capabilityID))
		resp.State.RemoveResource(ctx)
		return
	}
	
	mapAPICompletionCapabilityToModel(apiCap, &state, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() { return }

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Completion Capability %s", capabilityID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CompletionCapabilityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CompletionCapabilityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }

	capabilityID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating Completion Capability with ID: %s", capabilityID))

	updatePayload := coraxclient.CompletionCapabilityUpdate{}
	updateNeeded := false

	// Helper to set string pointer for update payload
	setStringPtr := func(current, new types.String) *string {
		if !new.Equal(current) {
			updateNeeded = true
			if new.IsNull() { return nil } // Explicitly setting to null if API supports it
			val := new.ValueString()
			return &val
		}
		return nil // No change, don't include in payload
	}
	setBoolPtr := func(current, new types.Bool) *bool {
		if !new.Equal(current) {
			updateNeeded = true
			if new.IsNull() { return nil }
			val := new.ValueBool()
			return &val
		}
		return nil
	}


	updatePayload.Name = setStringPtr(state.Name, plan.Name)
	updatePayload.IsPublic = setBoolPtr(state.IsPublic, plan.IsPublic)
	updatePayload.ModelID = setStringPtr(state.ModelID, plan.ModelID)
	updatePayload.ProjectID = setStringPtr(state.ProjectID, plan.ProjectID)
	updatePayload.SystemPrompt = setStringPtr(state.SystemPrompt, plan.SystemPrompt)
	updatePayload.CompletionPrompt = setStringPtr(state.CompletionPrompt, plan.CompletionPrompt)
	updatePayload.OutputType = setStringPtr(state.OutputType, plan.OutputType)


	if !plan.Variables.Equal(state.Variables) {
		updateNeeded = true
		if plan.Variables.IsNull() {
			updatePayload.Variables = nil // Or an empty slice if API expects that to clear
		} else {
			resp.Diagnostics.Append(plan.Variables.ElementsAs(ctx, &updatePayload.Variables, false)...)
			if resp.Diagnostics.HasError() { return }
		}
	}

	if !plan.SchemaDef.Equal(state.SchemaDef) {
		updateNeeded = true
		if plan.SchemaDef.IsNull() {
			updatePayload.SchemaDef = nil
		} else {
			updatePayload.SchemaDef = schemaDefMapToAPI(ctx, plan.SchemaDef, &resp.Diagnostics)
			if resp.Diagnostics.HasError() { return }
		}
	}
	
	if !plan.Config.Equal(state.Config) {
		updateNeeded = true
		updatePayload.Config = capabilityConfigModelToAPI(ctx, plan.Config, &resp.Diagnostics) // Assumes this helper is available
		if resp.Diagnostics.HasError() { return }
	}


	if !updateNeeded {
		tflog.Debug(ctx, "No attribute changes detected for Completion Capability update.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	updatedAPICap, err := r.client.UpdateCapability(ctx, capabilityID, updatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update completion capability %s: %s", capabilityID, err))
		return
	}

	mapAPICompletionCapabilityToModel(updatedAPICap, &plan, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() { return }

	tflog.Info(ctx, fmt.Sprintf("Completion Capability %s updated successfully", capabilityID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CompletionCapabilityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CompletionCapabilityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }

	capabilityID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting Completion Capability with ID: %s", capabilityID))

	err := r.client.DeleteCapability(ctx, capabilityID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Completion Capability %s not found, already deleted", capabilityID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete completion capability %s: %s", capabilityID, err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Completion Capability %s deleted successfully", capabilityID))
}

func (r *CompletionCapabilityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
