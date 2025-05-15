package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-corax/internal/coraxclient" // TODO: Adjust if your module name is different
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

// TODO: Add ResourceWithConfigure if client is needed (it is)

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *coraxclient.Client
}

// ProjectResourceModel describes the resource data model.
// Based on openapi.json components.schemas.Project
type ProjectResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	IsPublic        types.Bool   `tfsdk:"is_public"`
	CreatedBy       types.String `tfsdk:"created_by"`
	UpdatedBy       types.String `tfsdk:"updated_by"` // Nullable
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"` // Nullable
	Owner           types.String `tfsdk:"owner"`
	CollectionCount types.Int64  `tfsdk:"collection_count"`
	CapabilityCount types.Int64  `tfsdk:"capability_count"`
	// _links is not typically managed by Terraform directly as an attribute
}

// Helper function to map API Project to Terraform model
func mapProjectToModel(project *coraxclient.Project, model *ProjectResourceModel) {
	model.ID = types.StringValue(project.ID)
	model.Name = types.StringValue(project.Name)
	if project.Description != nil {
		model.Description = types.StringValue(*project.Description)
	} else {
		model.Description = types.StringNull()
	}
	model.IsPublic = types.BoolValue(project.IsPublic)
	model.CreatedBy = types.StringValue(project.CreatedBy)
	if project.UpdatedBy != nil {
		model.UpdatedBy = types.StringValue(*project.UpdatedBy)
	} else {
		model.UpdatedBy = types.StringNull()
	}
	model.CreatedAt = types.StringValue(project.CreatedAt)
	if project.UpdatedAt != nil {
		model.UpdatedAt = types.StringValue(*project.UpdatedAt)
	} else {
		model.UpdatedAt = types.StringNull()
	}
	model.Owner = types.StringValue(project.Owner)
	model.CollectionCount = types.Int64Value(int64(project.CollectionCount))
	model.CapabilityCount = types.Int64Value(int64(project.CapabilityCount))
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Corax Project. Projects are used to organize collections and capabilities.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the project (UUID).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the project. Must be at least 1 character long.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional description for the project.",
			},
			"is_public": schema.BoolAttribute{
				Optional:            true,
				Computed:            true, // API defaults to false if not provided
				MarkdownDescription: "Indicates whether the project is public. Defaults to false.",
				PlanModifiers: []planmodifier.Bool{
					// Use the API's default if the user doesn't specify a value.
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who created the project.",
			},
			"updated_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user who last updated the project.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation date and time of the project (RFC3339 format).",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last update date and time of the project (RFC3339 format).",
			},
			"owner": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The owner of the project.",
			},
			"collection_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of collections associated with this project.",
			},
			"capability_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of capabilities associated with this project.",
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating Project with name: %s", data.Name.ValueString()))

	projectCreatePayload := coraxclient.ProjectCreate{
		Name: data.Name.ValueString(),
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		projectCreatePayload.Description = &desc
	}
	if !data.IsPublic.IsNull() && !data.IsPublic.IsUnknown() {
		isPublic := data.IsPublic.ValueBool()
		projectCreatePayload.IsPublic = &isPublic
	}

	createdProject, err := r.client.CreateProject(ctx, projectCreatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create project, got error: %s", err))
		return
	}

	mapProjectToModel(createdProject, &data)
	tflog.Info(ctx, fmt.Sprintf("Project created successfully with ID: %s", createdProject.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading Project with ID: %s", projectID))

	project, err := r.client.GetProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Project with ID %s not found, removing from state", projectID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project %s, got error: %s", projectID, err))
		return
	}

	mapProjectToModel(project, &data)
	tflog.Debug(ctx, fmt.Sprintf("Successfully read Project with ID: %s", projectID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ID.ValueString() // ID comes from state, not plan
	tflog.Debug(ctx, fmt.Sprintf("Updating Project with ID: %s", projectID))

	projectUpdatePayload := coraxclient.ProjectUpdate{}
	// Only include fields in the payload if they have changed.
	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		projectUpdatePayload.Name = &name
	}
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			// To clear the description, the API might expect null or an empty string.
			// Assuming API handles "" as "clear description". If it needs explicit null, this needs adjustment.
			// For now, let's send what's in the plan. If plan.Description is null, it won't be set.
			// If it's an empty string, it will be sent.
			// The ProjectUpdate struct uses *string, so if plan.Description is null, omitempty works.
			// If plan.Description is "", we send "".
			if !plan.Description.IsUnknown() { // only set if not unknown
				desc := plan.Description.ValueString()
				projectUpdatePayload.Description = &desc
			}
		} else {
			desc := plan.Description.ValueString()
			projectUpdatePayload.Description = &desc
		}
	}
	if !plan.IsPublic.Equal(state.IsPublic) {
		if !plan.IsPublic.IsUnknown() { // only set if not unknown
			isPublic := plan.IsPublic.ValueBool()
			projectUpdatePayload.IsPublic = &isPublic
		}
	}

	// Check if there's anything to update
	// This check might be too simplistic if the API interprets absence of a field vs. explicit null differently.
	// For now, if the payload is empty, we might skip the API call.
	// However, it's safer to let the API handle an empty PUT if that's its behavior.
	// The openapi spec for ProjectUpdate has all fields as optional.
	// if projectUpdatePayload.Name == nil && projectUpdatePayload.Description == nil && projectUpdatePayload.IsPublic == nil {
	// 	tflog.Debug(ctx, "No changes detected for Project update.")
	// 	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Ensure state matches plan
	// 	return
	// }

	updatedProject, err := r.client.UpdateProject(ctx, projectID, projectUpdatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project %s, got error: %s", projectID, err))
		return
	}

	mapProjectToModel(updatedProject, &plan) // Update plan with response
	tflog.Info(ctx, fmt.Sprintf("Project updated successfully with ID: %s", projectID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting Project with ID: %s", projectID))

	err := r.client.DeleteProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, coraxclient.ErrNotFound) {
			tflog.Warn(ctx, fmt.Sprintf("Project with ID %s already deleted, removing from state", projectID))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project %s, got error: %s", projectID, err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Project with ID %s deleted successfully", projectID))
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
