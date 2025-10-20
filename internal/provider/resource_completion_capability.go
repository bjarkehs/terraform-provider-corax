// Copyright (c) HashiCorp, Inc.

//nolint:staticcheck // using json.Marshal on framework types for normalization
package provider

import (
	"context"
	"encoding/json"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-corax/internal/coraxclient"
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
	ID               types.String  `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	SemanticID       types.String  `tfsdk:"semantic_id"` // Optional
	IsPublic         types.Bool    `tfsdk:"is_public"`
	ModelID          types.String  `tfsdk:"model_id"`      // Nullable
	Config           types.Object  `tfsdk:"config"`        // Nullable, uses CapabilityConfigModel from chat_capability.go
	ProjectID        types.String  `tfsdk:"project_id"`    // Nullable
	SystemPrompt     types.String  `tfsdk:"system_prompt"` // Shared with Chat, but also in Completion
	CompletionPrompt types.String  `tfsdk:"completion_prompt"`
	Variables        types.Set     `tfsdk:"variables"`   // Nullable, set of strings
	OutputType       types.String  `tfsdk:"output_type"` // "schema" or "text"
	SchemaDef        types.Dynamic `tfsdk:"schema_def"`  // Nullable, for structured output definition
	Owner            types.String  `tfsdk:"owner"`       // Computed
	Type             types.String  `tfsdk:"type"`        // Computed, should always be "completion"
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
			"semantic_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A semantic identifier for the completion capability that can be used for referencing.",
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
			"variables": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A set of variable names (strings) that can be interpolated into the `completion_prompt`. Order is not significant.",
			},
			"output_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Defines the expected output format. Must be either 'text' or 'schema'.",
				Validators:          []validator.String{stringvalidator.OneOf("text", "schema")},
			},
			"schema_def": schema.DynamicAttribute{
				Optional:            true,
				MarkdownDescription: "Defines the structure of the output when `output_type` is 'schema'. This can be an HCL map or a JSON string. Required if `output_type` is 'schema', must be null or omitted if `output_type` is 'text'.",
				PlanModifiers: []planmodifier.Dynamic{
					normalizeSchemaDef(),
				},
			},
			"config": schema.SingleNestedAttribute{ // Reusing the same config structure as chat
				Optional:            true,
				MarkdownDescription: "Configuration settings for the capability's behavior.",
				Attributes:          capabilityConfigSchemaAttributes(), // Defined in chat_capability_resource.go (or move to a common place)
			},
			"owner": schema.StringAttribute{Computed: true, MarkdownDescription: "Owner of the capability.", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"type":  schema.StringAttribute{Computed: true, MarkdownDescription: "Type of the capability (should be 'completion').", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

// normalizeSchemaDefDynamicModifier is a plan modifier that normalizes a JSON string
// stored in a types.DynamicValue by unmarshalling and re-marshalling it,
// which sorts object keys alphabetically.
type normalizeSchemaDefDynamicModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m normalizeSchemaDefDynamicModifier) Description(ctx context.Context) string {
	return "Normalizes the schema_def attribute to a canonical JSON string representation by sorting object keys. This applies to both JSON string inputs and HCL map/object inputs."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m normalizeSchemaDefDynamicModifier) MarkdownDescription(ctx context.Context) string {
	return "Normalizes the `schema_def` attribute to a canonical JSON string. If `schema_def` is provided as a JSON string, it is parsed and re-serialized. If provided as an HCL map or object, it is converted to a map, then serialized to JSON. This process ensures that object keys within the resulting JSON string are alphabetically sorted. This helps prevent inconsistencies and ensures a canonical form in the plan, regardless of the input format (JSON string or HCL map/object)."
}

// PlanModifyDynamic implements the plan modification logic.
func (m normalizeSchemaDefDynamicModifier) PlanModifyDynamic(ctx context.Context, req planmodifier.DynamicRequest, resp *planmodifier.DynamicResponse) {
	// If the planned value is null or unknown, don't modify it
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	underlyingVal := req.PlanValue.UnderlyingValue()
	var data map[string]interface{}

	switch val := underlyingVal.(type) {
	case types.String:
		if val.IsNull() || val.IsUnknown() {
			return
		}
		jsonStr := val.ValueString()
		if jsonStr == "" { // An empty string is not valid JSON for a map.
			// Consider if an empty string should be an error or treated as null.
			// For now, returning, assuming schema validation will catch if it's problematic.
			return
		}
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			// Not valid JSON, so we don't normalize. Let schema validation catch it.
			// Or add a warning: resp.Diagnostics.AddAttributeWarning(req.Path, "Non-JSON String for schema_def", "...")
			return
		}
	case types.Object:
		if val.IsNull() || val.IsUnknown() {
			return
		}

		jsonBytes, err := json.Marshal(val)
		if err != nil {
			resp.Diagnostics.AddAttributeError(req.Path, "SchemaDef HCL Object Marshal Error", fmt.Sprintf("Failed to marshal HCL object for schema_def to JSON: %s", err.Error()))
			return
		}

		err = json.Unmarshal(jsonBytes, &data)
		if err != nil {
			resp.Diagnostics.AddAttributeError(req.Path, "SchemaDef HCL Object Unmarshal Error", fmt.Sprintf("Failed to marshal HCL object for schema_def to JSON: %s", err.Error()))
		}
	case types.Map:
		if val.IsNull() || val.IsUnknown() {
			return
		}
		// For types.Map, it's often easier to marshal to JSON and then unmarshal to map[string]interface{}
		// or directly convert if elements are compatible.
		// Here, we'll use the marshal/unmarshal approach for simplicity and consistency with how it might be handled elsewhere.
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			resp.Diagnostics.AddAttributeError(req.Path, "SchemaDef HCL Map Marshal Error", fmt.Sprintf("Failed to marshal HCL map for schema_def to JSON: %s", err.Error()))
			return
		}
		err = json.Unmarshal(jsonBytes, &data)
		if err != nil {
			resp.Diagnostics.AddAttributeError(req.Path, "SchemaDef HCL Map Unmarshal Error", fmt.Sprintf("Failed to unmarshal intermediate JSON for schema_def HCL map: %s", err.Error()))
			return
		}
	default:
		// Not a type we're handling for normalization (e.g., already a different dynamic type, or some other TF type)
		return
	}

	// If data is nil (e.g. from an empty JSON string or empty HCL map that resulted in nil map),
	// we might want to set the plan to types.DynamicNull() or types.StringNull() depending on desired behavior.
	// For now, if data is nil, Marshal will produce "null" string.
	if data == nil {
		return
	}

	// Marshal it back to get the canonical (sorted keys) version.
	normalizedBytes, err := json.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Failed to Normalize schema_def", fmt.Sprintf("Error re-marshalling schema_def to JSON: %s", err))
		return
	}

	normalizedStringValue := types.StringValue(string(normalizedBytes))
	resp.PlanValue = types.DynamicValue(normalizedStringValue)
}

// Ensure the implementation satisfies the interface.
var _ planmodifier.Dynamic = normalizeSchemaDefDynamicModifier{}

// Helper function to create the modifier.
func normalizeSchemaDef() planmodifier.Dynamic {
	return normalizeSchemaDefDynamicModifier{}
}

// capabilityConfigSchemaAttributes, capabilityConfigModelToAPI, capabilityConfigAPItoModel
// and their underlying attribute type helpers are defined in common_capability_config.go
// No need to redefine them here.

// --- Helper functions for mapping (specific to Completion Capability) ---

func schemaDefMapToAPI(ctx context.Context, schemaDef types.Dynamic, diags *diag.Diagnostics) map[string]interface{} {
	if schemaDef.IsNull() || schemaDef.IsUnknown() {
		return nil
	}

	underlyingVal := schemaDef.UnderlyingValue()
	var goMap map[string]interface{}

	switch val := underlyingVal.(type) {
	case types.String:
		if val.IsNull() || val.IsUnknown() {
			return nil // Or an empty map, depending on desired behavior for empty/null JSON string
		}
		err := json.Unmarshal([]byte(val.ValueString()), &goMap)
		if err != nil {
			diags.AddError("SchemaDef JSON String Error", fmt.Sprintf("schema_def was provided as a string, but it's not valid JSON for a map: %s. Content: %s", err.Error(), val.ValueString()))
			return nil
		}
		return goMap
	case types.Object:
		// Use the As method of types.Object to convert to map[string]interface{}
		// Ensure the target type for As is compatible with how objects are structured.
		// map[string]interface{} is a common target.
		convDiags := val.As(ctx, &goMap, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		diags.Append(convDiags...)
		if convDiags.HasError() {
			return nil
		}
		return goMap
	case types.Map:
		// For types.Map, marshal to JSON and then unmarshal to map[string]interface{}
		if val.IsNull() || val.IsUnknown() {
			return nil
		}
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			diags.AddError("SchemaDef Map Marshal Error", fmt.Sprintf("Failed to marshal HCL map for schema_def to JSON: %s", err.Error()))
			return nil
		}
		err = json.Unmarshal(jsonBytes, &goMap)
		if err != nil {
			diags.AddError("SchemaDef Map Unmarshal Error", fmt.Sprintf("Failed to unmarshal intermediate JSON for schema_def map: %s", err.Error()))
			return nil
		}
		return goMap
	default:
		// If it's not a string, object, or map that we can convert, it's an unsupported type for schema_def.
		diags.AddError("SchemaDef Type Error",
			fmt.Sprintf("schema_def has an unsupported underlying type: %T. "+
				"It should be an HCL map/object or a valid JSON string representing such a structure.", underlyingVal))
		return nil
	}
}

func schemaDefAPIToMap(apiSchemaDef map[string]interface{}, diags *diag.Diagnostics) types.Dynamic {
	if apiSchemaDef == nil {
		return types.DynamicNull()
	}

	jsonBytes, err := json.Marshal(apiSchemaDef)
	if err != nil {
		diags.AddError("SchemaDef API Conversion Error", fmt.Sprintf("Failed to marshal schema_def from API to JSON: %s", err))
		return types.DynamicNull()
	}

	strVal := types.StringValue(string(jsonBytes))
	// types.DynamicValue(attr.Value) constructor returns types.Dynamic, not (types.Dynamic, diag.Diagnostics)
	dynVal := types.DynamicValue(strVal)
	return dynVal
}

func mapAPICompletionCapabilityToModel(apiCap *coraxclient.CapabilityRepresentation, model *CompletionCapabilityResourceModel, diags *diag.Diagnostics, ctx context.Context) {
	model.ID = types.StringValue(apiCap.ID)
	model.SemanticID = types.StringValue(apiCap.SemanticID)
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

	// Populate SystemPrompt and CompletionPrompt from apiCap.Configuration
	if apiCap.Configuration != nil {
		if sysPrompt, ok := apiCap.Configuration["system_prompt"].(string); ok {
			model.SystemPrompt = types.StringValue(sysPrompt)
		} else {
			// If key is missing or not a string, treat as unknown.
			// Per schema, system_prompt is required, so Unknown is appropriate if not found/convertible.
			model.SystemPrompt = types.StringUnknown()
		}

		if compPrompt, ok := apiCap.Configuration["completion_prompt"].(string); ok {
			model.CompletionPrompt = types.StringValue(compPrompt)
		} else {
			// Per schema, completion_prompt is required.
			model.CompletionPrompt = types.StringUnknown()
		}
	} else {
		// apiCap.Configuration map itself is nil
		model.SystemPrompt = types.StringUnknown()
		model.CompletionPrompt = types.StringUnknown()
		tflog.Debug(ctx, fmt.Sprintf("apiCap.Configuration is nil for capability %s. SystemPrompt and CompletionPrompt will be unknown.", apiCap.ID))
	}

	// Populate OutputType and SchemaDef from apiCap.Output
	if apiCap.Output != nil {
		if outputTypeVal, ok := apiCap.Output["type"].(string); ok {
			model.OutputType = types.StringValue(outputTypeVal)
		} else {
			// Per schema, output_type is required.
			model.OutputType = types.StringUnknown()
		}

		// schema_def is sourced from apiCap.Output["result"]
		// It's optional overall, but required if output_type is "schema".
		// schemaDefAPIToMap handles nil input map by returning types.DynamicNull().
		if schemaDefVal, ok := apiCap.Output["result"].(map[string]interface{}); ok {
			model.SchemaDef = schemaDefAPIToMap(schemaDefVal, diags)
		} else {
			// If "result" is not found, or not a map[string]interface{}, treat SchemaDef as null.
			// This is correct if output_type is "text" (schema_def would be absent/null)
			// or if "result" is present but malformed.
			if _, found := apiCap.Output["result"]; found && !ok {
				diags.AddAttributeWarning(
					path.Root("schema_def"), // Or a more specific path
					"Invalid Type for Schema Definition",
					fmt.Sprintf("Expected 'result' in API output to be a map, but got %T. Treating schema_def as null.", apiCap.Output["result"]),
				)
			}
			model.SchemaDef = types.DynamicNull()
		}
	} else {
		// apiCap.Output map itself is nil
		model.OutputType = types.StringUnknown()
		model.SchemaDef = types.DynamicNull()
		tflog.Debug(ctx, fmt.Sprintf("apiCap.Output is nil for capability %s. OutputType will be unknown and SchemaDef null.", apiCap.ID))
	}

	// Populate Variables from apiCap.Input
	if apiCap.Input != nil {
		if varsData, found := apiCap.Input["variables"]; found && varsData != nil {
			if vars, ok := varsData.([]interface{}); ok {
				strVars := make([]string, len(vars))
				allStrings := true
				for i, v := range vars {
					if strV, isString := v.(string); isString {
						strVars[i] = strV
					} else {
						allStrings = false
						diags.AddAttributeWarning(
							path.Root("variables"), // Or a more specific path
							"Invalid Variable Type in API Response",
							fmt.Sprintf("Variable at index %d is not a string (actual type: %T). Treating variables as null.", i, v),
						)
						break
					}
				}
				if allStrings {
					setValue, conversionDiags := types.SetValueFrom(ctx, types.StringType, strVars)
					diags.Append(conversionDiags...)
					if !conversionDiags.HasError() {
						model.Variables = setValue // Handles empty set correctly (non-null, empty set)
					} else {
						model.Variables = types.SetNull(types.StringType) // Error in types.SetValueFrom
					}
				} else {
					model.Variables = types.SetNull(types.StringType) // Non-string element found
				}
			} else if varsMap, ok := varsData.(map[string]interface{}); ok { // Handle map from GET
				strVarKeys := make([]string, 0, len(varsMap))
				for k := range varsMap {
					strVarKeys = append(strVarKeys, k)
				}
				// No need to sort keys for a set, but doesn't harm if API returns a list of keys that were sorted
				// sort.Strings(strVarKeys) // Order doesn't matter for a Set

				setValue, conversionDiags := types.SetValueFrom(ctx, types.StringType, strVarKeys)
				diags.Append(conversionDiags...)
				if !conversionDiags.HasError() {
					model.Variables = setValue
				} else {
					model.Variables = types.SetNull(types.StringType)
					diags.AddAttributeError(
						path.Root("variables"),
						"Variable Conversion Error (Map to Set)",
						fmt.Sprintf("Failed to convert variable keys from API map to set: %v", conversionDiags),
					)
				}
			} else { // apiCap.Input["variables"] is present but not []interface{} and not map[string]interface{}
				diags.AddAttributeWarning(
					path.Root("variables"),
					"Incorrect Type for Variables in API Response",
					fmt.Sprintf("Expected 'variables' in API input to be a list or map of strings, but got %T. Treating variables as null.", varsData),
				)
				model.Variables = types.SetNull(types.StringType)
			}
		} else { // "variables" key not found in apiCap.Input or its value is JSON null
			model.Variables = types.SetNull(types.StringType)
		}
	} else { // apiCap.Input map itself is nil
		model.Variables = types.SetNull(types.StringType)
		tflog.Debug(ctx, fmt.Sprintf("apiCap.Input is nil for capability %s. Variables will be null.", apiCap.ID))
	}

	model.Config = capabilityConfigAPItoModel(ctx, apiCap.Config, diags) // Common config

	model.Owner = types.StringValue(apiCap.Owner)
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
	if resp.Diagnostics.HasError() {
		return
	}

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
	if !plan.SemanticID.IsNull() && !plan.SemanticID.IsUnknown() {
		semanticID := plan.SemanticID.ValueString()
		apiPayload.SemanticID = &semanticID
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
		if resp.Diagnostics.HasError() {
			return
		}
	}
	outputType := plan.OutputType.ValueString()
	if outputType == "schema" {
		if plan.SchemaDef.IsNull() || plan.SchemaDef.IsUnknown() {
			resp.Diagnostics.AddError("Validation Error", "schema_def is required when output_type is 'schema'")
			return
		}
		apiPayload.SchemaDef = schemaDefMapToAPI(ctx, plan.SchemaDef, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if outputType == "text" {
		if !plan.SchemaDef.IsNull() && !plan.SchemaDef.IsUnknown() {
			resp.Diagnostics.AddError("Validation Error", "schema_def must not be set when output_type is 'text'")
			return
		}
	} else {
		resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("unsupported output_type '%s', must be either 'text' or 'schema'", outputType))
		return
	}

	// Common config mapping (reuse from chat capability if moved to common, or define here)
	// For now, assuming capabilityConfigModelToAPI is available (defined in chat_capability.go or common)
	apiPayload.Config = capabilityConfigModelToAPI(ctx, plan.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	createdAPICap, err := r.client.CreateCapability(ctx, apiPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create completion capability, got error: %s", err))
		return
	}

	mapAPICompletionCapabilityToModel(createdAPICap, &plan, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Completion Capability %s created successfully with ID %s", plan.Name.ValueString(), plan.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CompletionCapabilityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CompletionCapabilityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Completion Capability %s", capabilityID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CompletionCapabilityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CompletionCapabilityResourceModel
	var state CompletionCapabilityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	capabilityID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating Completion Capability with ID: %s using full plan payload", capabilityID))

	// --- Construct full update payload from plan ---
	nameValue := plan.Name.ValueString()
	typeValue := "completion" // Type is fixed for this resource
	systemPromptValue := plan.SystemPrompt.ValueString()
	completionPromptValue := plan.CompletionPrompt.ValueString()
	outputTypeValue := plan.OutputType.ValueString()

	updatePayload := coraxclient.CompletionCapabilityUpdate{
		Name:             &nameValue,
		Type:             &typeValue,
		SystemPrompt:     &systemPromptValue,
		CompletionPrompt: &completionPromptValue,
		OutputType:       &outputTypeValue,
	}

	// IsPublic
	if !plan.IsPublic.IsNull() && !plan.IsPublic.IsUnknown() {
		isPublicVal := plan.IsPublic.ValueBool()
		updatePayload.IsPublic = &isPublicVal
	} else {
		defaultIsPublic := false // As per schema default
		updatePayload.IsPublic = &defaultIsPublic
	}

	// SemanticID
	if !plan.SemanticID.IsNull() && !plan.SemanticID.IsUnknown() {
		semanticIDVal := plan.SemanticID.ValueString()
		updatePayload.SemanticID = &semanticIDVal
	} else {
		updatePayload.SemanticID = nil
	}

	// ModelID
	if !plan.ModelID.IsNull() && !plan.ModelID.IsUnknown() {
		modelIDVal := plan.ModelID.ValueString()
		updatePayload.ModelID = &modelIDVal
	} else {
		updatePayload.ModelID = nil
	}

	// ProjectID
	if !plan.ProjectID.IsNull() && !plan.ProjectID.IsUnknown() {
		projectIDVal := plan.ProjectID.ValueString()
		updatePayload.ProjectID = &projectIDVal
	} else {
		updatePayload.ProjectID = nil
	}

	// Variables
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		var vars []string
		resp.Diagnostics.Append(plan.Variables.ElementsAs(ctx, &vars, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updatePayload.Variables = vars // Assign directly, omitempty handles if vars is nil/empty based on API needs
	} else {
		// If API expects an empty list to clear, send []string{}. If omitempty on nil is preferred, send nil.
		// Assuming omitempty on nil is fine for now.
		updatePayload.Variables = nil
	}

	// SchemaDef
	if outputTypeValue == "schema" {
		if plan.SchemaDef.IsNull() || plan.SchemaDef.IsUnknown() {
			resp.Diagnostics.AddError("Validation Error", "schema_def is required when output_type is 'schema'")
			return
		}
		updatePayload.SchemaDef = schemaDefMapToAPI(ctx, plan.SchemaDef, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	} else if outputTypeValue == "text" {
		if !plan.SchemaDef.IsNull() && !plan.SchemaDef.IsUnknown() {
			resp.Diagnostics.AddError("Validation Error", "schema_def must not be set when output_type is 'text'")
			return
		}
		updatePayload.SchemaDef = nil
	} else {
		resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("unsupported output_type '%s', must be either 'text' or 'schema'", outputTypeValue))
		return
	}

	// Config
	updatePayload.Config = capabilityConfigModelToAPI(ctx, plan.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	// --- End of payload construction ---

	updatedAPICap, err := r.client.UpdateCapability(ctx, capabilityID, updatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update completion capability %s: %s", capabilityID, err))
		return
	}

	mapAPICompletionCapabilityToModel(updatedAPICap, &plan, &resp.Diagnostics, ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Completion Capability %s updated successfully", capabilityID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CompletionCapabilityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CompletionCapabilityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
