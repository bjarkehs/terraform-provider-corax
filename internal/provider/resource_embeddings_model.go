package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"corax/internal/coraxclient"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EmbeddingsModelResource{}
var _ resource.ResourceWithImportState = &EmbeddingsModelResource{}

func NewEmbeddingsModelResource() resource.Resource {
	return &EmbeddingsModelResource{}
}

// EmbeddingsModelResource defines the resource implementation.
type EmbeddingsModelResource struct {
	client *coraxclient.Client
}

// EmbeddingsModelResourceModel describes the resource data model.
// Based on openapi.json components.schemas.EmbeddingsModel
type EmbeddingsModelResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`           // Nullable
	ModelProvider        types.String `tfsdk:"model_provider"`        // e.g., "openai", "cohere", "sentence-transformers"
	ModelNameOnProvider  types.String `tfsdk:"model_name_on_provider"`// e.g., "text-embedding-ada-002"
	Dimensions           types.Int64  `tfsdk:"dimensions"`
	MaxTokens            types.Int64  `tfsdk:"max_tokens"`            // Nullable/Optional in API, but required in TF if not computed with default
	ApiKey               types.String `tfsdk:"api_key"`               // Optional, Sensitive
	ApiBaseUrl           types.String `tfsdk:"api_base_url"`          // Optional, Nullable
	Status               types.String `tfsdk:"status"`                // Computed
	IsDefault            types.Bool   `tfsdk:"is_default"`            // Computed
	CreatedBy            types.String `tfsdk:"created_by"`            // Computed
	UpdatedBy            types.String `tfsdk:"updated_by"`            // Computed, Nullable
	CreatedAt            types.String `tfsdk:"created_at"`            // Computed
	UpdatedAt            types.String `tfsdk:"updated_at"`            // Computed, Nullable
}

// Helper function to map API EmbeddingsModel to Terraform model
func mapEmbeddingsModelToModel(modelAPI *coraxclient.EmbeddingsModel, modelTF *EmbeddingsModelResourceModel, diags *diag.Diagnostics) {
	modelTF.ID = types.StringValue(modelAPI.ID)
	modelTF.Name = types.StringValue(modelAPI.Name)

	if modelAPI.Description != nil {
		modelTF.Description = types.StringValue(*modelAPI.Description)
	} else {
		modelTF.Description = types.StringNull()
	}

	modelTF.ModelProvider = types.StringValue(modelAPI.ModelProvider)
	modelTF.ModelNameOnProvider = types.StringValue(modelAPI.ModelNameOnProvider)
	modelTF.Dimensions = types.Int64Value(int64(modelAPI.Dimensions)) // API uses int, TF uses int64

	if modelAPI.MaxTokens != nil {
		modelTF.MaxTokens = types.Int64Value(int64(*modelAPI.MaxTokens)) // API uses *int
	} else {
		modelTF.MaxTokens = types.Int64Null()
	}

	// ApiKey is sensitive and typically not returned by GET.
	// If it was just set (e.g. during Create/Update) and the API returns it,
	// we might get it here. Otherwise, we rely on the value from config/plan.
	// For now, if API returns it, we set it. Otherwise, it remains as per plan.
	// This field in TF model should be `Optional: true, Computed: true, Sensitive: true`
	// and potentially use stringplanmodifier.UseStateForUnknown() if we don't want to clear it if API doesn't return it.
	// However, since it's sensitive, it's better if the API doesn't return it on GET.
	// We will preserve the existing ApiKey from the plan/state if not explicitly returned by the API on read.
	// For Create/Update, the API response *might* include it if it was just set.
	// Let's assume GET does not return it. The TF model's ApiKey will be set from the plan during Create/Update.
	// On Read, we don't try to set modelTF.ApiKey from modelAPI.ApiKey unless explicitly present.
	// If modelAPI.ApiKey is nil, modelTF.ApiKey remains untouched (preserves state).
	// If modelAPI.ApiKey is not nil (e.g. after create), then we can set it.
	// This logic is tricky for sensitive fields not always returned.
	// For now, we'll only set it if the API explicitly returns it.
	// A common pattern is to not map it back from read if it's write-only or not returned.
	// The current schema has ApiKey as Optional, Sensitive. It does not have Computed: true.
	// This means if the user provides it, it's in the plan. If they don't, it's null.
	// The API response for EmbeddingsModel has `api_key *string json:"api_key,omitempty"`
	// So, if the API returns it, we can map it.
	if modelAPI.ApiKey != nil {
		modelTF.ApiKey = types.StringValue(*modelAPI.ApiKey)
	} else if modelTF.ApiKey.IsUnknown() { // If it was unknown in plan, and API returns nil, make it null.
		modelTF.ApiKey = types.StringNull()
	}
	// If modelTF.ApiKey was known (set by user) and API returns nil, modelTF.ApiKey retains its known value.

	if modelAPI.ApiBaseUrl != nil {
		modelTF.ApiBaseUrl = types.StringValue(*modelAPI.ApiBaseUrl)
	} else {
		modelTF.ApiBaseUrl = types.StringNull()
	}

	modelTF.Status = types.StringValue(modelAPI.Status)
	modelTF.IsDefault = types.BoolValue(modelAPI.IsDefault)
	modelTF.CreatedBy = types.StringValue(modelAPI.CreatedBy)
	modelTF.CreatedAt = types.StringValue(modelAPI.CreatedAt)

	if modelAPI.UpdatedBy != nil {
		modelTF.UpdatedBy = types.StringValue(*modelAPI.UpdatedBy)
	} else {
		modelTF.UpdatedBy = types.StringNull()
	}
	if modelAPI.UpdatedAt != nil {
		modelTF.UpdatedAt = types.StringValue(*modelAPI.UpdatedAt)
	} else {
		modelTF.UpdatedAt = types.StringNull()
	}
}


func (r *EmbeddingsModelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_embeddings_model"
}

func (r *EmbeddingsModelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Embeddings Model. Embeddings models are used to generate vector embeddings for documents.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the embeddings model (UUID).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A user-defined name for the embeddings model.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional description for the embeddings model.",
			},
			"model_provider": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The provider of the embeddings model (e.g., 'openai', 'cohere', 'sentence-transformers', 'custom').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				// Consider adding stringvalidator.OneOf if the list of providers is fixed and known.
			},
			"model_name_on_provider": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The specific model name as recognized by the provider (e.g., 'text-embedding-ada-002', 'embed-english-v2.0', 'all-MiniLM-L6-v2').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dimensions": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The number of dimensions the embeddings model outputs.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"max_tokens": schema.Int64Attribute{
				Optional:            true,
				Computed:            true, // API might have a default or derive it
				MarkdownDescription: "The maximum number of input tokens the model can handle. If not provided, a default may be used or derived by the API.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(), // If API provides a default
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The API key required for accessing the model, if it's a third-party proprietary model. Store sensitive values in a secure way (e.g., environment variables or a secrets manager).",
			},
			"api_base_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The base URL for the API of a custom or self-hosted embeddings model.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The operational status of the embeddings model (e.g., 'active', 'error', 'pending_validation').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_default": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this is the default embeddings model for the Corax instance.",
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who created the embeddings model.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who last updated the embeddings model.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation date and time of the embeddings model (RFC3339 format).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last update date and time of the embeddings model (RFC3339 format).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *EmbeddingsModelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmbeddingsModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EmbeddingsModelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating EmbeddingsModel with name: %s", data.Name.ValueString()))

	createPayload := coraxclient.EmbeddingsModelCreate{
		Name:                data.Name.ValueString(),
		ModelProvider:       data.ModelProvider.ValueString(),
		ModelNameOnProvider: data.ModelNameOnProvider.ValueString(),
		Dimensions:          int(data.Dimensions.ValueInt64()), // TF int64 to API int
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		createPayload.Description = &desc
	}
	if !data.MaxTokens.IsNull() && !data.MaxTokens.IsUnknown() {
		mt := int(data.MaxTokens.ValueInt64())
		createPayload.MaxTokens = &mt
	}
	if !data.ApiKey.IsNull() && !data.ApiKey.IsUnknown() {
		apiKey := data.ApiKey.ValueString()
		createPayload.ApiKey = &apiKey
	}
	if !data.ApiBaseUrl.IsNull() && !data.ApiBaseUrl.IsUnknown() {
		baseUrl := data.ApiBaseUrl.ValueString()
		createPayload.ApiBaseUrl = &baseUrl
	}

	createdModel, err := r.client.CreateEmbeddingsModel(ctx, createPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create embeddings model, got error: %s", err))
		return
	}

	// Preserve ApiKey from plan as API might not return it, or might return a redacted version.
	// The mapEmbeddingsModelToModel helper needs to be careful with this.
	// For create, the response *should* be complete.
	currentApiKey := data.ApiKey // Preserve from plan

	mapEmbeddingsModelToModel(createdModel, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	
	// If API did not return ApiKey (e.g. it's null in createdModel), restore from plan if it was set.
	// This ensures the state reflects what was configured if the API doesn't echo back sensitive fields.
	if createdModel.ApiKey == nil && !currentApiKey.IsNull() && !currentApiKey.IsUnknown() {
		data.ApiKey = currentApiKey
	}


	tflog.Info(ctx, fmt.Sprintf("EmbeddingsModel created successfully with ID: %s", createdModel.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmbeddingsModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EmbeddingsModelResourceModel // Current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	modelID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading EmbeddingsModel with ID: %s", modelID))

	modelAPI, err := r.client.GetEmbeddingsModel(ctx, modelID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("EmbeddingsModel with ID %s not found, removing from state", modelID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read embeddings model %s, got error: %s", modelID, err))
		return
	}

	// Preserve ApiKey from state as API GET likely won't return it.
	currentApiKeyInState := state.ApiKey

	mapEmbeddingsModelToModel(modelAPI, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If API GET doesn't return api_key (modelAPI.ApiKey is nil),
	// ensure the state's api_key is preserved from what was previously set.
	if modelAPI.ApiKey == nil && !currentApiKeyInState.IsNull() && !currentApiKeyInState.IsUnknown() {
		state.ApiKey = currentApiKeyInState
	} else if modelAPI.ApiKey != nil {
		// If API *does* return it (e.g. it's a non-sensitive placeholder or was just set), use it.
		state.ApiKey = types.StringValue(*modelAPI.ApiKey)
	} else {
		// If API returns nil and state was also nil/unknown, ensure it's null.
		state.ApiKey = types.StringNull()
	}


	tflog.Debug(ctx, fmt.Sprintf("Successfully read EmbeddingsModel with ID: %s", modelID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EmbeddingsModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EmbeddingsModelResourceModel  // Data from the plan
	var state EmbeddingsModelResourceModel // Data from the current state

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	modelID := state.ID.ValueString() // ID cannot be changed
	tflog.Debug(ctx, fmt.Sprintf("Updating EmbeddingsModel with ID: %s", modelID))

	updatePayload := coraxclient.EmbeddingsModelUpdate{}
	updateNeeded := false

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updatePayload.Name = &name
		updateNeeded = true
	}
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			var emptyDesc string // Assuming API interprets empty string as clearing, or send explicit null if supported
			updatePayload.Description = &emptyDesc
		} else {
			desc := plan.Description.ValueString()
			updatePayload.Description = &desc
		}
		updateNeeded = true
	}
	if !plan.ApiKey.Equal(state.ApiKey) { // Sensitive field, update if changed
		if !plan.ApiKey.IsNull() && !plan.ApiKey.IsUnknown() {
			apiKey := plan.ApiKey.ValueString()
			updatePayload.ApiKey = &apiKey
		} else {
			// If plan is null, and state was not, this means user wants to clear it.
			// The API schema EmbeddingsModelUpdate has api_key as *string, so sending null is possible.
			var nilStr *string = nil
			updatePayload.ApiKey = nilStr
		}
		updateNeeded = true
	}
	if !plan.ApiBaseUrl.Equal(state.ApiBaseUrl) {
		if !plan.ApiBaseUrl.IsNull() && !plan.ApiBaseUrl.IsUnknown() {
			baseUrl := plan.ApiBaseUrl.ValueString()
			updatePayload.ApiBaseUrl = &baseUrl
		} else {
			var nilStr *string = nil
			updatePayload.ApiBaseUrl = nilStr
		}
		updateNeeded = true
	}
	
	// Fields like model_provider, model_name_on_provider, dimensions, max_tokens are not updatable
	// according to EmbeddingsModelUpdate schema. If they change in plan, it should force a new resource.
	// This should be handled by `RequiresReplace` plan modifiers on those attributes if they are immutable.
	// Let's check if they are immutable. The OpenAPI spec for PUT /v1/embeddings-models/{model_id}
	// takes EmbeddingsModelUpdate, which only includes name, description, api_key, api_base_url.
	// So, other fields (provider, model_name, dimensions, max_tokens) if changed, should force replacement.
	// Add RequiresReplace to them in schema.

	if !updateNeeded {
		tflog.Debug(ctx, "No attribute changes detected for EmbeddingsModel update.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Ensure state matches plan
		return
	}

	updatedModel, err := r.client.UpdateEmbeddingsModel(ctx, modelID, updatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update embeddings model %s, got error: %s", modelID, err))
		return
	}

	currentApiKeyInPlan := plan.ApiKey // Preserve from plan before mapping

	mapEmbeddingsModelToModel(updatedModel, &plan, &resp.Diagnostics) // Update plan with response
	if resp.Diagnostics.HasError() {
		return
	}

	// If API did not return ApiKey, restore from plan if it was set.
	if updatedModel.ApiKey == nil && !currentApiKeyInPlan.IsNull() && !currentApiKeyInPlan.IsUnknown() {
		plan.ApiKey = currentApiKeyInPlan
	}


	tflog.Info(ctx, fmt.Sprintf("EmbeddingsModel updated successfully with ID: %s", modelID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EmbeddingsModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EmbeddingsModelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	modelID := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting EmbeddingsModel with ID: %s", modelID))

	err := r.client.DeleteEmbeddingsModel(ctx, modelID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("EmbeddingsModel with ID %s already deleted, removing from state", modelID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete embeddings model %s, got error: %s", modelID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("EmbeddingsModel with ID %s deleted successfully", modelID))
}

func (r *EmbeddingsModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
