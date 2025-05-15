package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
var _ resource.Resource = &DocumentResource{}
var _ resource.ResourceWithImportState = &DocumentResource{}

func NewDocumentResource() resource.Resource {
	return &DocumentResource{}
}

// DocumentResource defines the resource implementation.
type DocumentResource struct {
	client *coraxclient.Client
}

// DocumentResourceModel describes the resource data model.
// Based on openapi.json components.schemas.Document and DocumentIngest/DocumentUpdate
type DocumentResourceModel struct {
	ID             types.String `tfsdk:"id"`               // User-provided or API-generated Document ID
	CollectionID   types.String `tfsdk:"collection_id"`    // Required: parent collection
	TextContent    types.String `tfsdk:"text_content"`     // Optional
	JsonContent    types.String `tfsdk:"json_content"`     // Optional, user provides JSON string
	Metadata       types.Map    `tfsdk:"metadata"`         // Optional, map[string]dynamic
	Content        types.String `tfsdk:"content"`          // Computed: string or JSON string representation from API
	TokenCount     types.Int64  `tfsdk:"token_count"`      // Computed
	ChunkCount     types.Int64  `tfsdk:"chunk_count"`      // Computed
	EmbeddingsStatus types.String `tfsdk:"embeddings_status"`// Computed
	CreatedBy      types.String `tfsdk:"created_by"`       // Computed
	UpdatedBy      types.String `tfsdk:"updated_by"`       // Computed (Nullable)
	CreatedAt      types.String `tfsdk:"created_at"`       // Computed
	UpdatedAt      types.String `tfsdk:"updated_at"`       // Computed (Nullable)
}

// Helper function to map API Document to Terraform model
func mapDocumentToModel(ctx context.Context, doc *coraxclient.Document, model *DocumentResourceModel, diags *diag.Diagnostics) {
	model.ID = types.StringValue(doc.ID)
	model.CollectionID = types.StringValue(doc.CollectionID) // Assuming client populates this

	if doc.TextContent != nil {
		model.TextContent = types.StringValue(*doc.TextContent)
	} else {
		model.TextContent = types.StringNull()
	}

	if doc.JsonContent != nil {
		jsonBytes, err := json.Marshal(doc.JsonContent)
		if err != nil {
			diags.AddError("Failed to Marshal API JsonContent", fmt.Sprintf("Could not marshal json_content from API: %s", err.Error()))
			return
		}
		model.JsonContent = types.StringValue(string(jsonBytes))
	} else {
		model.JsonContent = types.StringNull()
	}

	if doc.Content != nil {
		contentBytes, err := json.Marshal(doc.Content) // Marshal whatever it is (string or object)
		if err != nil {
			diags.AddError("Failed to Marshal API Content", fmt.Sprintf("Could not marshal content from API: %s", err.Error()))
			return
		}
		model.Content = types.StringValue(string(contentBytes))
	} else {
		model.Content = types.StringNull()
	}

	if doc.Metadata != nil {
		metadataMapValue, d := types.MapValueFrom(ctx, types.DynamicType, doc.Metadata)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
		model.Metadata = metadataMapValue
	} else {
		model.Metadata = types.MapNull(types.DynamicType)
	}

	model.TokenCount = types.Int64Value(int64(doc.TokenCount))
	model.ChunkCount = types.Int64Value(int64(doc.ChunkCount))
	model.EmbeddingsStatus = types.StringValue(doc.EmbeddingsStatus)
	model.CreatedBy = types.StringValue(doc.CreatedBy)
	model.CreatedAt = types.StringValue(doc.CreatedAt)

	if doc.UpdatedBy != nil {
		model.UpdatedBy = types.StringValue(*doc.UpdatedBy)
	} else {
		model.UpdatedBy = types.StringNull()
	}
	if doc.UpdatedAt != nil {
		model.UpdatedAt = types.StringValue(*doc.UpdatedAt)
	} else {
		model.UpdatedAt = types.StringNull()
	}
}

func (r *DocumentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document"
}

func (r *DocumentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Document within a specified Collection. Documents can have text or JSON content and associated metadata.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the document (UUID). Can be user-provided or auto-generated by the API. If user-provided and changed, it forces a new resource.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(), // If user sets it and changes it, replace.
				},
				// TODO: Add validator for UUID format if desired
			},
			"collection_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the collection this document belongs to. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				// TODO: Add validator for UUID format
			},
			"text_content": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The text content of the document. One of `text_content` or `json_content` must be provided.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("json_content"),
					}...),
				},
			},
			"json_content": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The JSON content of the document, as a string. One of `text_content` or `json_content` must be provided.",
				// User must provide a valid JSON string.
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.DynamicType, // Allows for mixed types, will be marshalled to JSON for API
				Optional:            true,
				MarkdownDescription: "A map of metadata to associate with the document. Keys are strings, values can be strings, numbers, booleans, or nested maps/lists that are JSON-serializable.",
			},
			"content": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The raw content of the document as returned by the API (can be string or JSON string).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of tokens in the document's content.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"chunk_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of chunks the document was divided into for embedding.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"embeddings_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of embeddings generation for this document (e.g., 'completed', 'pending', 'failed').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who created the document.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who last updated the document.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation date and time of the document (RFC3339 format).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last update date and time of the document (RFC3339 format).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DocumentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Helper to convert types.Map (DynamicType elements) to map[string]interface{}
func modelMetadataToApi(ctx context.Context, modelMetadata types.Map) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	if modelMetadata.IsNull() || modelMetadata.IsUnknown() {
		return nil, diags
	}

	apiMetadata := make(map[string]interface{})
	elements := modelMetadata.Elements()
	for k, v := range elements {
		var val interface{}
		// Attempt to convert attr.Value to a Go native type for JSON marshalling
		// This is a simplified conversion. For complex nested structures within DynamicType,
		// more sophisticated conversion might be needed, or rely on `v.As(ctx, &val, basetypes.DynamicAsOptions{})`
		// For now, we assume common primitive types or that JSON marshalling of attr.Value works.
		// A safer approach for arbitrary structures is to marshal attr.Value to JSON then unmarshal to interface{}.
		// Or, if tf-plugin-framework handles DynamicType marshalling directly when passed to the client, this might be simpler.
		// Let's try a direct approach and see if `json.Marshal` in the client handles `attr.Value` well enough.
		// If not, we'll need to convert `v` (attr.Value) to standard Go types here.
		// For simplicity, let's assume the client's JSON marshaller can handle attr.Value.
		// If not, this is where we'd convert v (type attr.Value) to a Go native type.
		// For example, if v is types.String, use v.(types.String).ValueString().
		// If v is types.Number, use v.(types.Number).ValueBigFloat() etc.
		// This is a common challenge with types.DynamicType.
		// A robust way:
		err := v.As(ctx, &val, types.DynamicAsOptions{ UnmarshalUnknown: true, UnmarshalNull: true})
		if err != nil {
			diags.AddAttributeError(
				path.Root("metadata").AtMapKey(k),
				"Failed to convert metadata value",
				fmt.Sprintf("Error converting metadata value for key %s: %s", k, err.Error()),
			)
			continue
		}
		apiMetadata[k] = val
	}
	return apiMetadata, diags
}


func (r *DocumentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DocumentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionID := data.CollectionID.ValueString()
	documentID := ""

	if data.ID.IsUnknown() || data.ID.IsNull() {
		// Generate UUID if not provided by user.
		// The API's PUT endpoint for documents acts as an upsert and requires a document ID in the path.
		uid, err := uuid.NewRandom()
		if err != nil {
			resp.Diagnostics.AddError("Failed to generate Document ID", fmt.Sprintf("Error generating UUID for document: %s", err.Error()))
			return
		}
		documentID = uid.String()
		tflog.Debug(ctx, fmt.Sprintf("Generated Document ID: %s for collection %s", documentID, collectionID))
	} else {
		documentID = data.ID.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Using user-provided Document ID: %s for collection %s", documentID, collectionID))
	}

	docUpdatePayload := coraxclient.DocumentUpdate{}

	if !data.TextContent.IsNull() && !data.TextContent.IsUnknown() {
		textContent := data.TextContent.ValueString()
		docUpdatePayload.TextContent = &textContent
	}

	if !data.JsonContent.IsNull() && !data.JsonContent.IsUnknown() {
		var jsonMap map[string]interface{}
		err := json.Unmarshal([]byte(data.JsonContent.ValueString()), &jsonMap)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("json_content"),
				"Invalid JSON String",
				fmt.Sprintf("Failed to unmarshal json_content: %s. Please provide a valid JSON string.", err.Error()),
			)
			return
		}
		docUpdatePayload.JsonContent = jsonMap
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		apiMetadata, diags := modelMetadataToApi(ctx, data.Metadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		docUpdatePayload.Metadata = apiMetadata
	}

	tflog.Debug(ctx, fmt.Sprintf("Upserting Document ID: %s in Collection ID: %s", documentID, collectionID))
	upsertedDoc, err := r.client.UpsertDocument(ctx, collectionID, documentID, docUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create/update document, got error: %s", err))
		return
	}

	mapDocumentToModel(ctx, upsertedDoc, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	// Ensure the ID used (generated or user-provided) is set in the state.
	// The mapDocumentToModel should handle this from upsertedDoc.ID.
	// If we generated an ID and the API confirmed it, it's fine.
	// If user provided an ID, API confirmed it, also fine.

	tflog.Info(ctx, fmt.Sprintf("Document %s/%s upserted successfully", collectionID, upsertedDoc.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DocumentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionID := data.CollectionID.ValueString()
	documentID := data.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading Document ID: %s from Collection ID: %s", documentID, collectionID))

	doc, err := r.client.GetDocument(ctx, collectionID, documentID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Document %s/%s not found, removing from state", collectionID, documentID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read document %s/%s, got error: %s", collectionID, documentID, err))
		return
	}

	mapDocumentToModel(ctx, doc, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Document %s/%s", collectionID, documentID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DocumentResourceModel // Data from the plan
	var state DocumentResourceModel // Data from the current state

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionID := state.CollectionID.ValueString() // Collection ID cannot change
	documentID := state.ID.ValueString()         // Document ID cannot change for an update (it's ForceNew if user changes it in config)

	docUpdatePayload := coraxclient.DocumentUpdate{}
	updateNeeded := false

	// Check TextContent
	if !plan.TextContent.Equal(state.TextContent) {
		if plan.TextContent.IsNull() || plan.TextContent.IsUnknown() {
			// To clear, send null or an empty string based on API behavior.
			// Assuming sending null is the way to clear if it was previously set.
			// The client type DocumentUpdate.TextContent is *string.
			// If API needs an empty string to clear, this needs adjustment.
			// For now, if plan is null, we don't set it in payload, relying on omitempty.
			// If it was set and now plan is null, we might need to send explicit null.
			// Let's assume if it's null in plan, we want it null in API.
			// If it's empty string in plan, we send empty string.
			if plan.TextContent.IsNull() {
				var nilStr *string = nil
				docUpdatePayload.TextContent = nilStr // Explicitly set to null
			} else {
				// This case should not happen if IsNull is true.
				// If it's an empty string and different from state, set it.
				tc := plan.TextContent.ValueString()
				docUpdatePayload.TextContent = &tc
			}
		} else {
			tc := plan.TextContent.ValueString()
			docUpdatePayload.TextContent = &tc
		}
		updateNeeded = true
	}

	// Check JsonContent
	if !plan.JsonContent.Equal(state.JsonContent) {
		if plan.JsonContent.IsNull() || plan.JsonContent.IsUnknown() {
			var nilMap map[string]interface{} = nil
			docUpdatePayload.JsonContent = nilMap // Explicitly set to null
		} else {
			var jsonMap map[string]interface{}
			err := json.Unmarshal([]byte(plan.JsonContent.ValueString()), &jsonMap)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("json_content"),
					"Invalid JSON String for Update",
					fmt.Sprintf("Failed to unmarshal json_content: %s", err.Error()),
				)
				return
			}
			docUpdatePayload.JsonContent = jsonMap
		}
		updateNeeded = true
	}

	// Check Metadata
	if !plan.Metadata.Equal(state.Metadata) {
		if plan.Metadata.IsNull() || plan.Metadata.IsUnknown() {
			var nilMap map[string]interface{} = nil
			docUpdatePayload.Metadata = nilMap // Explicitly set to null
		} else {
			apiMetadata, diags := modelMetadataToApi(ctx, plan.Metadata)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			docUpdatePayload.Metadata = apiMetadata
		}
		updateNeeded = true
	}

	if !updateNeeded {
		tflog.Debug(ctx, "No attribute changes detected for Document update.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Ensure state matches plan
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating Document ID: %s in Collection ID: %s", documentID, collectionID))
	updatedDoc, err := r.client.UpsertDocument(ctx, collectionID, documentID, docUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update document %s/%s, got error: %s", collectionID, documentID, err))
		return
	}

	mapDocumentToModel(ctx, updatedDoc, &plan, &resp.Diagnostics) // Update plan with response
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Document %s/%s updated successfully", collectionID, documentID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DocumentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DocumentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionID := data.CollectionID.ValueString()
	documentID := data.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Deleting Document ID: %s from Collection ID: %s", documentID, collectionID))

	err := r.client.DeleteDocument(ctx, collectionID, documentID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Document %s/%s already deleted, removing from state", collectionID, documentID))
			resp.State.RemoveResource(ctx) // Remove from state if not found
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete document %s/%s, got error: %s", collectionID, documentID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Document %s/%s deleted successfully", collectionID, documentID))
	// No need to remove from state here, Terraform does it automatically on successful delete.
}

func (r *DocumentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// ID for import will be "collection_id/document_id"
	// Example: resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	// This needs custom logic to parse collection_id and document_id.
	importIDParts := path.SplitAttributePath(path.Root(req.ID))
	if len(importIDParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'collection_id/document_id', got: %s", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_id"), importIDParts[0].String())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), importIDParts[1].String())...)
}
