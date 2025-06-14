// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator" // Added
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corax/internal/coraxclient"
)

// --- Reusable Model Structs for Capability Config ---

// CapabilityConfigModel maps to components.schemas.CapabilityConfig.
type CapabilityConfigModel struct {
	Temperature    types.Float64 `tfsdk:"temperature"`     // Nullable
	BlobConfig     types.Object  `tfsdk:"blob_config"`     // Nullable
	DataRetention  types.Object  `tfsdk:"data_retention"`  // Polymorphic: TimedDataRetention or InfiniteDataRetention
	ContentTracing types.Bool    `tfsdk:"content_tracing"` // Default true
}

// BlobConfigModel maps to components.schemas.BlobConfig.
type BlobConfigModel struct {
	MaxFileSizeMB    types.Int64 `tfsdk:"max_file_size_mb"`   // Default 20
	MaxBlobs         types.Int64 `tfsdk:"max_blobs"`          // Default 10
	AllowedMimeTypes types.List  `tfsdk:"allowed_mime_types"` // Default ["image/png", "image/jpeg"]
}

// DataRetentionModel for the data_retention block.
type DataRetentionModel struct {
	Type  types.String `tfsdk:"type"`  // Will store "timed" or "infinite"
	Hours types.Int64  `tfsdk:"hours"` // Nullable, only used if type is "timed"
}

// TimedDataRetentionModel (Removed)
// InfiniteDataRetentionModel (Removed)

// --- Custom Validator for DataRetention ---

// dataRetentionValidator validates the DataRetentionModel object.
// It ensures that 'hours' is set if 'type' is 'timed',
// and 'hours' is not set if 'type' is 'infinite'.
type dataRetentionValidator struct{}

func (v dataRetentionValidator) Description(ctx context.Context) string {
	return "Validates that 'hours' is configured correctly based on the 'type' of data retention. " +
		"If 'type' is 'timed', 'hours' must be set. If 'type' is 'infinite', 'hours' must not be set."
}

func (v dataRetentionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v dataRetentionValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// If the object is null or unknown, we cannot validate its attributes.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var dataRetention DataRetentionModel
	diags := req.ConfigValue.As(ctx, &dataRetention, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return // Error converting to model, can't proceed with this validation.
	}

	// The 'type' attribute itself has a stringvalidator.OneOf("timed", "infinite")
	// and is 'Required', so we can assume it's one of these if not null/unknown.
	if dataRetention.Type.IsNull() || dataRetention.Type.IsUnknown() {
		// This case should ideally be caught by the 'type' attribute's own Required validator.
		// If it still reaches here as null/unknown, we can't perform conditional validation.
		return
	}
	retentionType := dataRetention.Type.ValueString()

	hoursIsSet := !dataRetention.Hours.IsNull() && !dataRetention.Hours.IsUnknown()

	switch retentionType {
	case "timed":
		if !hoursIsSet {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("hours"), // Path to the 'hours' attribute within the data_retention object
				"Missing 'hours' for timed data retention",
				"The 'hours' attribute must be configured when data retention 'type' is 'timed'.",
			)
		}
	case "infinite":
		if hoursIsSet {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("hours"), // Path to the 'hours' attribute
				"Unexpected 'hours' for infinite data retention",
				"The 'hours' attribute must not be configured when data retention 'type' is 'infinite'.",
			)
		}
	}
}

// --- Reusable Attribute Type Definitions ---

func capabilityConfigAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"temperature":     types.Float64Type,
		"blob_config":     types.ObjectType{AttrTypes: blobConfigAttributeTypes()},
		"data_retention":  types.ObjectType{AttrTypes: dataRetentionAttributeTypes()},
		"content_tracing": types.BoolType,
	}
}

func blobConfigAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"max_file_size_mb":   types.Int64Type,
		"max_blobs":          types.Int64Type,
		"allowed_mime_types": types.ListType{ElemType: types.StringType},
	}
}

func dataRetentionAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":  types.StringType,
		"hours": types.Int64Type,
	}
}

// func timedDataRetentionAttributeTypes() (Removed)
// func infiniteDataRetentionAttributeTypes() (Removed)

// --- Reusable Schema Definition for Config Block ---

func capabilityConfigSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"temperature": schema.Float64Attribute{
			Optional:            true,
			MarkdownDescription: "Controls randomness in response generation (0.0 to 1.0). Higher values make output more random.",
			// TODO: Add float validator for range 0.0-1.0
		},
		"blob_config": schema.SingleNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Configuration for handling file uploads (blobs) if the capability supports it.",
			Attributes: map[string]schema.Attribute{
				"max_file_size_mb": schema.Int64Attribute{
					Optional:            true,
					Computed:            true, // API might have its own defaults
					MarkdownDescription: "Maximum file size in megabytes for uploaded blobs.",
				},
				"max_blobs": schema.Int64Attribute{
					Optional:            true,
					Computed:            true, // API might have its own defaults
					MarkdownDescription: "Maximum number of blobs that can be uploaded.",
				},
				"allowed_mime_types": schema.ListAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					Computed:            true, // API might have its own defaults
					MarkdownDescription: "List of allowed MIME types for uploaded blobs.",
				},
			},
		},
		"data_retention": schema.SingleNestedAttribute{
			Optional:            true,
			MarkdownDescription: "Defines how long execution input and output data should be kept. Configure with 'type' and optionally 'hours'.",
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of data retention. Must be 'timed' or 'infinite'.",
					Validators:          []validator.String{stringvalidator.OneOf("timed", "infinite")},
				},
				"hours": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "Duration in hours to retain data. Required if type is 'timed'. Must not be set if type is 'infinite'. Minimum 1.",
					Validators:          []validator.Int64{int64validator.AtLeast(1)},
				},
			},
			Validators: []validator.Object{
				dataRetentionValidator{}, // Use the custom validator
			},
		},
		"content_tracing": schema.BoolAttribute{
			Optional:            true,
			Computed:            true, // API default is true
			MarkdownDescription: "Whether content (prompts, completion data, variables) should be recorded in observability systems. Automatically set to false by the API for timed data retention.",
		},
	}
}

// --- Reusable Mapping Functions ---

func capabilityConfigModelToAPI(ctx context.Context, modelConfig types.Object, diags *diag.Diagnostics) *coraxclient.CapabilityConfig {
	if modelConfig.IsNull() || modelConfig.IsUnknown() {
		return nil
	}

	var cfgModel CapabilityConfigModel
	respDiags := modelConfig.As(ctx, &cfgModel, basetypes.ObjectAsOptions{})
	diags.Append(respDiags...)
	if diags.HasError() {
		return nil
	}

	apiConfig := &coraxclient.CapabilityConfig{}
	hasChanges := false // Track if any field in config is actually set to avoid sending empty config object

	if !cfgModel.Temperature.IsNull() && !cfgModel.Temperature.IsUnknown() {
		val := cfgModel.Temperature.ValueFloat64()
		apiConfig.Temperature = &val
		hasChanges = true
	}
	if !cfgModel.ContentTracing.IsNull() && !cfgModel.ContentTracing.IsUnknown() {
		val := cfgModel.ContentTracing.ValueBool()
		apiConfig.ContentTracing = &val
		hasChanges = true
	}

	if !cfgModel.BlobConfig.IsNull() && !cfgModel.BlobConfig.IsUnknown() {
		var blobCfgModel BlobConfigModel
		respDiags := cfgModel.BlobConfig.As(ctx, &blobCfgModel, basetypes.ObjectAsOptions{})
		diags.Append(respDiags...)
		if diags.HasError() {
			return nil
		}

		apiBlobCfg := &coraxclient.BlobConfig{}
		blobChanges := false
		if !blobCfgModel.MaxFileSizeMB.IsNull() && !blobCfgModel.MaxFileSizeMB.IsUnknown() {
			val := int(blobCfgModel.MaxFileSizeMB.ValueInt64())
			apiBlobCfg.MaxFileSizeMB = &val
			blobChanges = true
		}
		if !blobCfgModel.MaxBlobs.IsNull() && !blobCfgModel.MaxBlobs.IsUnknown() {
			val := int(blobCfgModel.MaxBlobs.ValueInt64())
			apiBlobCfg.MaxBlobs = &val
			blobChanges = true
		}
		if !blobCfgModel.AllowedMimeTypes.IsNull() && !blobCfgModel.AllowedMimeTypes.IsUnknown() {
			diags.Append(blobCfgModel.AllowedMimeTypes.ElementsAs(ctx, &apiBlobCfg.AllowedMimeTypes, false)...)
			if diags.HasError() {
				return nil
			}
			blobChanges = true
		}
		if blobChanges {
			apiConfig.BlobConfig = apiBlobCfg
			hasChanges = true
		}
	}

	if !cfgModel.DataRetention.IsNull() && !cfgModel.DataRetention.IsUnknown() {
		var drModel DataRetentionModel
		respDiags := cfgModel.DataRetention.As(ctx, &drModel, basetypes.ObjectAsOptions{})
		diags.Append(respDiags...)
		if diags.HasError() {
			return nil
		}

		apiDR := &coraxclient.DataRetention{}
		drChanges := false

		if !drModel.Type.IsNull() && !drModel.Type.IsUnknown() {
			retentionType := drModel.Type.ValueString()
			apiDR.Type = retentionType
			drChanges = true // Setting the type is a change

			switch retentionType {
			case "timed":
				// Schema ensures Hours is non-null and valid if Type is "timed"
				if !drModel.Hours.IsNull() && !drModel.Hours.IsUnknown() {
					val := int(drModel.Hours.ValueInt64())
					apiDR.Hours = &val
				}
				// If Hours were null/unknown here despite schema, it's an issue.
				// The API requires 'hours' for 'timed' type.
			case "infinite":
				apiDR.Hours = nil // Explicitly ensure Hours is not sent for infinite type
			}
		}

		if drChanges {
			apiConfig.DataRetention = apiDR
			hasChanges = true
		}
	}

	if !hasChanges {
		return nil
	} // If no actual values were set in config, return nil to omit it from API payload
	return apiConfig
}

func capabilityConfigAPItoModel(ctx context.Context, apiConfig *coraxclient.CapabilityConfig, diags *diag.Diagnostics) types.Object {
	if apiConfig == nil {
		return types.ObjectNull(capabilityConfigAttributeTypes())
	}

	attrs := make(map[string]attr.Value)

	if apiConfig.Temperature != nil {
		attrs["temperature"] = types.Float64Value(*apiConfig.Temperature)
	} else {
		attrs["temperature"] = types.Float64Null()
	}

	if apiConfig.ContentTracing != nil {
		attrs["content_tracing"] = types.BoolValue(*apiConfig.ContentTracing)
	} else {
		// Default to true as per schema, if API omits it (meaning default)
		attrs["content_tracing"] = types.BoolValue(true)
	}

	if apiConfig.BlobConfig != nil {
		blobAttrs := make(map[string]attr.Value)
		if apiConfig.BlobConfig.MaxFileSizeMB != nil {
			blobAttrs["max_file_size_mb"] = types.Int64Value(int64(*apiConfig.BlobConfig.MaxFileSizeMB))
		} else {
			blobAttrs["max_file_size_mb"] = types.Int64Null()
		}
		if apiConfig.BlobConfig.MaxBlobs != nil {
			blobAttrs["max_blobs"] = types.Int64Value(int64(*apiConfig.BlobConfig.MaxBlobs))
		} else {
			blobAttrs["max_blobs"] = types.Int64Null()
		}
		if apiConfig.BlobConfig.AllowedMimeTypes != nil {
			listVal, listDiags := types.ListValueFrom(ctx, types.StringType, apiConfig.BlobConfig.AllowedMimeTypes)
			diags.Append(listDiags...)
			blobAttrs["allowed_mime_types"] = listVal
		} else {
			blobAttrs["allowed_mime_types"] = types.ListNull(types.StringType)
		}
		blobObj, objDiags := types.ObjectValue(blobConfigAttributeTypes(), blobAttrs)
		diags.Append(objDiags...)
		attrs["blob_config"] = blobObj
	} else {
		attrs["blob_config"] = types.ObjectNull(blobConfigAttributeTypes())
	}

	if apiConfig.DataRetention != nil {
		drAttrs := make(map[string]attr.Value)
		retentionType := apiConfig.DataRetention.Type

		drAttrs["type"] = types.StringValue(retentionType)

		if retentionType == "timed" && apiConfig.DataRetention.Hours != nil {
			drAttrs["hours"] = types.Int64Value(int64(*apiConfig.DataRetention.Hours))
		} else {
			// For "infinite", or if "timed" but hours is missing from API (which would be an API inconsistency for "timed")
			// or if type is unknown from API.
			drAttrs["hours"] = types.Int64Null()
		}

		// Use the new dataRetentionAttributeTypes() which expects "type" and "hours"
		drObj, drObjDiags := types.ObjectValue(dataRetentionAttributeTypes(), drAttrs)
		diags.Append(drObjDiags...)
		attrs["data_retention"] = drObj
	} else {
		// Use the new dataRetentionAttributeTypes()
		attrs["data_retention"] = types.ObjectNull(dataRetentionAttributeTypes())
	}

	objVal, objDiags := types.ObjectValue(capabilityConfigAttributeTypes(), attrs)
	diags.Append(objDiags...)
	return objVal
}
