package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
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

	"corax/internal/coraxclient" // TODO: Adjust if your module name is different
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CollectionResource{}
var _ resource.ResourceWithImportState = &CollectionResource{}

func NewCollectionResource() resource.Resource {
	return &CollectionResource{}
}

// CollectionResource defines the resource implementation.
type CollectionResource struct {
	client *coraxclient.Client
}

// CollectionResourceModel describes the resource data model.
// Based on openapi.json components.schemas.Collection
type CollectionResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"` // Nullable
	ProjectID           types.String `tfsdk:"project_id"`
	EmbeddingsModelID   types.String `tfsdk:"embeddings_model_id"` // Nullable
	MetadataSchema      types.Map    `tfsdk:"metadata_schema"`     // Nullable, map[string]string
	CreatedBy           types.String `tfsdk:"created_by"`
	UpdatedBy           types.String `tfsdk:"updated_by"` // Nullable
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"` // Nullable
	DocumentCount       types.Int64  `tfsdk:"document_count"`
	SizeBytes           types.Int64  `tfsdk:"size_bytes"`
	Status              types.String `tfsdk:"status"`
}

// Helper function to map API Collection to Terraform model
func mapCollectionToModel(collection *coraxclient.Collection, model *CollectionResourceModel, diags *diag.Diagnostics) {
	model.ID = types.StringValue(collection.ID)
	model.Name = types.StringValue(collection.Name)

	if collection.Description != nil {
		model.Description = types.StringValue(*collection.Description)
	} else {
		model.Description = types.StringNull()
	}

	model.ProjectID = types.StringValue(collection.ProjectID)

	if collection.EmbeddingsModelID != nil {
		model.EmbeddingsModelID = types.StringValue(*collection.EmbeddingsModelID)
	} else {
		model.EmbeddingsModelID = types.StringNull()
	}

	if collection.MetadataSchema != nil {
		elements := make(map[string]attr.Value)
		for k, v := range collection.MetadataSchema {
			elements[k] = types.StringValue(v)
		}
		mapValue, d := types.MapValue(types.StringType, elements)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
		model.MetadataSchema = mapValue
	} else {
		model.MetadataSchema = types.MapNull(types.StringType)
	}

	model.CreatedBy = types.StringValue(collection.CreatedBy)
	if collection.UpdatedBy != nil {
		model.UpdatedBy = types.StringValue(*collection.UpdatedBy)
	} else {
		model.UpdatedBy = types.StringNull()
	}
	model.CreatedAt = types.StringValue(collection.CreatedAt)
	if collection.UpdatedAt != nil {
		model.UpdatedAt = types.StringValue(*collection.UpdatedAt)
	} else {
		model.UpdatedAt = types.StringNull()
	}
	model.DocumentCount = types.Int64Value(collection.DocumentCount)
	model.SizeBytes = types.Int64Value(collection.SizeBytes)
	model.Status = types.StringValue(collection.Status)
}

func (r *CollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (r *CollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Knowledge Collection. Collections store documents and their embeddings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the collection (UUID).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the collection. Must be at least 1 character long.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional description for the collection.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the project this collection belongs to.",
				// TODO: Add validator for UUID format if available and desired
			},
			"embeddings_model_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The UUID of the embeddings model to use for this collection. If not provided, a default model may be used by the API.",
				// TODO: Add validator for UUID format if available and desired
			},
			"metadata_schema": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map defining the schema for metadata attached to documents in this collection. Keys are property names, and values are their types (e.g., 'string', 'number', 'boolean', 'array', 'object').",
				Validators: []validator.Map{
					mapvalidator.ValueStringsAre(stringvalidator.OneOf("string", "number", "boolean", "array", "object")),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who created the collection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who last updated the collection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation date and time of the collection (RFC3339 format).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last update date and time of the collection (RFC3339 format).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"document_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of documents in the collection.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"size_bytes": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total size of the documents in the collection in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the collection (e.g., 'ready', 'indexing', 'failed').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CollectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating Collection with name: %s for project %s", data.Name.ValueString(), data.ProjectID.ValueString()))

	collectionCreatePayload := coraxclient.CollectionCreate{
		Name:      data.Name.ValueString(),
		ProjectID: data.ProjectID.ValueString(),
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		collectionCreatePayload.Description = &desc
	}
	if !data.EmbeddingsModelID.IsNull() && !data.EmbeddingsModelID.IsUnknown() {
		embID := data.EmbeddingsModelID.ValueString()
		collectionCreatePayload.EmbeddingsModelID = &embID
	}
	if !data.MetadataSchema.IsNull() && !data.MetadataSchema.IsUnknown() {
		metadataSchema := make(map[string]string)
		diags := data.MetadataSchema.ElementsAs(ctx, &metadataSchema, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		collectionCreatePayload.MetadataSchema = metadataSchema
	}

	createdCollection, err := r.client.CreateCollection(ctx, collectionCreatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create collection, got error: %s", err))
		return
	}

	mapCollectionToModel(createdCollection, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Collection created successfully with ID: %s", createdCollection.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CollectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionID := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading Collection with ID: %s", collectionID))

	collection, err := r.client.GetCollection(ctx, collectionID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Collection with ID %s not found, removing from state", collectionID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read collection %s, got error: %s", collectionID, err))
		return
	}

	mapCollectionToModel(collection, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully read Collection with ID: %s", collectionID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CollectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CollectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionID := state.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating Collection with ID: %s", collectionID))

	collectionUpdatePayload := coraxclient.CollectionUpdate{}
	updateNeeded := false

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		collectionUpdatePayload.Name = &name
		updateNeeded = true
	}
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			// To clear description, send an empty string or null, depending on API.
			// Assuming API interprets an empty string as clearing. If it needs explicit null, this needs adjustment.
			// For CollectionUpdate, Description is *string, so setting it to a pointer to an empty string.
			// If the API expects null to clear, then we'd need to handle that.
			// Let's assume for now that sending an empty string clears it, or sending null if the plan is null.
			var emptyDesc string // Default to empty string to clear
			collectionUpdatePayload.Description = &emptyDesc
		} else {
			desc := plan.Description.ValueString()
			collectionUpdatePayload.Description = &desc
		}
		updateNeeded = true
	}
	if !plan.EmbeddingsModelID.Equal(state.EmbeddingsModelID) {
		if plan.EmbeddingsModelID.IsNull() {
			// API might not allow clearing embeddings_model_id or changing it post-creation.
			// For now, assume it can be set to null.
			// This field is *string in CollectionUpdate.
			collectionUpdatePayload.EmbeddingsModelID = nil // Explicitly set to null
		} else {
			embID := plan.EmbeddingsModelID.ValueString()
			collectionUpdatePayload.EmbeddingsModelID = &embID
		}
		updateNeeded = true
	}

	if !plan.MetadataSchema.Equal(state.MetadataSchema) {
		if plan.MetadataSchema.IsNull() {
			// Send null to clear the schema if API supports it.
			// CollectionUpdate.MetadataSchema is *map[string]string
			var nilMap map[string]string = nil
			collectionUpdatePayload.MetadataSchema = &nilMap // Pointer to a nil map
		} else if !plan.MetadataSchema.IsUnknown() {
			metadataSchema := make(map[string]string)
			diags := plan.MetadataSchema.ElementsAs(ctx, &metadataSchema, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			collectionUpdatePayload.MetadataSchema = &metadataSchema
		}
		updateNeeded = true
	}


	if !updateNeeded {
		tflog.Debug(ctx, "No attribute changes detected for Collection update.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Ensure state matches plan
		return
	}

	updatedCollection, err := r.client.UpdateCollection(ctx, collectionID, collectionUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update collection %s, got error: %s", collectionID, err))
		return
	}

	mapCollectionToModel(updatedCollection, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Collection updated successfully with ID: %s", collectionID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CollectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionID := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting Collection with ID: %s", collectionID))

	err := r.client.DeleteCollection(ctx, collectionID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Collection with ID %s already deleted, removing from state", collectionID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete collection %s, got error: %s", collectionID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Collection with ID %s deleted successfully", collectionID))
}

func (r *CollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
