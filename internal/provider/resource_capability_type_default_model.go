package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
var _ resource.Resource = &CapabilityTypeDefaultModelResource{}
var _ resource.ResourceWithImportState = &CapabilityTypeDefaultModelResource{}

func NewCapabilityTypeDefaultModelResource() resource.Resource {
	return &CapabilityTypeDefaultModelResource{}
}

// CapabilityTypeDefaultModelResource defines the resource implementation.
type CapabilityTypeDefaultModelResource struct {
	client *coraxclient.Client
}

// CapabilityTypeDefaultModelResourceModel describes the resource data model.
type CapabilityTypeDefaultModelResourceModel struct {
	CapabilityType           types.String `tfsdk:"capability_type"`             // This will also serve as the ID
	DefaultModelDeploymentID types.String `tfsdk:"default_model_deployment_id"` // UUID
	// Read-only attributes from CapabilityTypeRepresentation
	Name types.String `tfsdk:"name"`
}

func (r *CapabilityTypeDefaultModelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_capability_type_default_model"
}

func (r *CapabilityTypeDefaultModelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the default Model Deployment for a specific Capability Type (e.g., 'chat', 'completion').",
		Attributes: map[string]schema.Attribute{
			"capability_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of the capability (e.g., 'chat', 'completion', 'embedding'). This also serves as the resource ID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},                  // Changing this means managing a different capability type's default
				Validators:          []validator.String{stringvalidator.OneOf("chat", "completion", "embedding")}, // Based on CapabilityType enum
			},
			"default_model_deployment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the Model Deployment to set as the default for this capability type.",
				// TODO: Add UUID validator
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the capability type.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			// Note: No separate 'id' attribute; 'capability_type' is the identifier.
		},
	}
}

func (r *CapabilityTypeDefaultModelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create implements resource.Resource.
func (r *CapabilityTypeDefaultModelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CapabilityTypeDefaultModelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Create function for corax_capability_type_default_model (actually PUT).")
	// Create is effectively an Update (PUT) operation

	updatePayload := coraxclient.DefaultModelDeploymentUpdate{
		DefaultModelDeploymentID: plan.DefaultModelDeploymentID.ValueString(),
	}
	capabilityType := plan.CapabilityType.ValueString()

	apiResp, err := r.client.SetCapabilityTypeDefaultModel(ctx, capabilityType, updatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set default model for capability type %s: %s", capabilityType, err))
		return
	}

	plan.Name = types.StringValue(apiResp.Name)
	// The ID of this resource is the capability_type itself.

	tflog.Info(ctx, fmt.Sprintf("Default model for capability type %s set successfully.", capabilityType))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read implements resource.Resource.
func (r *CapabilityTypeDefaultModelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CapabilityTypeDefaultModelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	capabilityType := state.CapabilityType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading default model for capability type: %s", capabilityType))

	apiResp, err := r.client.GetCapabilityType(ctx, capabilityType)
	if err != nil {
		// Consider how to handle 404 for a capability type - it shouldn't happen if valid type is used.
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read capability type %s: %s", capabilityType, err))
		return
	}

	state.Name = types.StringValue(apiResp.Name)
	if apiResp.DefaultModelDeploymentID != nil {
		state.DefaultModelDeploymentID = types.StringValue(*apiResp.DefaultModelDeploymentID)
	} else {
		// If API returns null, it means no default is set.
		// Terraform might see this as a diff if the config expects one.
		// This resource *manages* the default, so if it's null, TF should try to set it on apply.
		// For Read, if it's null, we reflect that.
		state.DefaultModelDeploymentID = types.StringNull()
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully read default model for capability type %s", capabilityType))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (r *CapabilityTypeDefaultModelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CapabilityTypeDefaultModelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Since capability_type RequiresReplace, Update is only for default_model_deployment_id.
	tflog.Debug(ctx, "Update function for corax_capability_type_default_model (actually PUT).")

	updatePayload := coraxclient.DefaultModelDeploymentUpdate{
		DefaultModelDeploymentID: plan.DefaultModelDeploymentID.ValueString(),
	}
	capabilityType := plan.CapabilityType.ValueString()

	apiResp, err := r.client.SetCapabilityTypeDefaultModel(ctx, capabilityType, updatePayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update default model for capability type %s: %s", capabilityType, err))
		return
	}

	plan.Name = types.StringValue(apiResp.Name)

	tflog.Info(ctx, fmt.Sprintf("Default model for capability type %s updated successfully.", capabilityType))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (r *CapabilityTypeDefaultModelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CapabilityTypeDefaultModelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	capabilityType := state.CapabilityType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting default model for capability type (attempting to set to null): %s", capabilityType))

	// To "delete" or "unset" the default, we PUT with a null or empty DefaultModelDeploymentID.
	// The API must support this. If DefaultModelDeploymentUpdate requires the field,
	// we might need a specific client method or the API needs to allow null.
	// Assuming API allows DefaultModelDeploymentID to be explicitly set to null via the PUT endpoint.
	// The DefaultModelDeploymentUpdate struct in client uses string, not *string.
	// This implies it cannot be omitted to mean "unset". It must be sent.
	// If the API interprets an empty string as "unset", that's one way.
	// If it needs a specific "null" value or a different mechanism, this needs adjustment.
	// For now, let's try PUTting with an empty string, assuming this unsets it.
	// This is a guess and needs API behavior confirmation.

	// updatePayload := coraxclient.DefaultModelDeploymentUpdate{
	// 	DefaultModelDeploymentID: "", // Attempt to unset by sending empty string
	// }
	// _, err := r.client.SetCapabilityTypeDefaultModel(ctx, capabilityType, updatePayload)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unset default model for capability type %s: %s. Manual intervention might be needed if API does not support unsetting via empty string.", capabilityType, err))
	// 	return
	// }

	// More robust: The API should ideally support PATCH with `default_model_deployment_id: null`
	// or PUT where omitting the field or sending null clears it.
	// If PUT requires the field and an empty string is not "null", then true deletion isn't possible
	// without specific API support.
	// For now, we'll make Delete a no-op with a warning, as standard PUT might not clear it.
	tflog.Warn(ctx, fmt.Sprintf("Deletion of a capability_type_default_model resource does not actively clear the default model in Corax API due to lack of a dedicated 'unset' operation. The resource will be removed from Terraform state. If you need to clear the default, do so via the Corax API/UI if possible, or set it to a different valid model deployment ID."))
}

// ImportState implements resource.ResourceWithImportState.
func (r *CapabilityTypeDefaultModelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The ID for this resource is the capability_type itself.
	resource.ImportStatePassthroughID(ctx, path.Root("capability_type"), req, resp)
}
