package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-corax/internal/coraxclient"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// --- Reusable Model Structs for Capability Config ---

// CapabilityConfigModel maps to components.schemas.CapabilityConfig
type CapabilityConfigModel struct {
	Temperature    types.Float64 `tfsdk:"temperature"`     // Nullable
	BlobConfig     types.Object  `tfsdk:"blob_config"`     // Nullable
	DataRetention  types.Object  `tfsdk:"data_retention"`  // Polymorphic: TimedDataRetention or InfiniteDataRetention
	ContentTracing types.Bool    `tfsdk:"content_tracing"` // Default true
}

// BlobConfigModel maps to components.schemas.BlobConfig
type BlobConfigModel struct {
	MaxFileSizeMB    types.Int64 `tfsdk:"max_file_size_mb"`   // Default 20
	MaxBlobs         types.Int64 `tfsdk:"max_blobs"`          // Default 10
	AllowedMimeTypes types.List  `tfsdk:"allowed_mime_types"` // Default ["image/png", "image/jpeg"]
}

// DataRetentionModel for the data_retention block (polymorphic)
type DataRetentionModel struct {
	Timed    types.Object `tfsdk:"timed"`    // TimedDataRetention
	Infinite types.Object `tfsdk:"infinite"` // InfiniteDataRetention (empty object)
}

// TimedDataRetentionModel maps to components.schemas.TimedDataRetention
type TimedDataRetentionModel struct {
	Hours types.Int64 `tfsdk:"hours"` // Min 1
}

// InfiniteDataRetentionModel maps to components.schemas.InfiniteDataRetention (effectively a marker)
type InfiniteDataRetentionModel struct {
	Enabled types.Bool `tfsdk:"enabled"` // User sets this to true to indicate infinite
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
		"timed":    types.ObjectType{AttrTypes: timedDataRetentionAttributeTypes()},
		"infinite": types.ObjectType{AttrTypes: infiniteDataRetentionAttributeTypes()},
	}
}

func timedDataRetentionAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{"hours": types.Int64Type}
}

func infiniteDataRetentionAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{"enabled": types.BoolType}
}

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
			MarkdownDescription: "Defines how long execution input and output data should be kept.",
			Attributes: map[string]schema.Attribute{
				"timed": schema.SingleNestedAttribute{
					Optional:            true,
					MarkdownDescription: "Retain data for a specific duration.",
					Attributes: map[string]schema.Attribute{
						"hours": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "Duration in hours to retain data (minimum 1).",
							Validators:          []validator.Int64{int64validator.AtLeast(1)},
						},
					},
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("infinite")),
					},
				},
				"infinite": schema.SingleNestedAttribute{
					Optional:            true,
					MarkdownDescription: "Retain data indefinitely.",
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Set to true to enable infinite data retention.",
						},
					},
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("timed")),
					},
				},
			},
			Validators: []validator.Object{
				objectvalidator.ExactlyOneOf(
					path.MatchRelative().AtName("timed"),
					path.MatchRelative().AtName("infinite"),
				),
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
		if !drModel.Timed.IsNull() && !drModel.Timed.IsUnknown() {
			var timedModel TimedDataRetentionModel
			respDiags := drModel.Timed.As(ctx, &timedModel, basetypes.ObjectAsOptions{})
			diags.Append(respDiags...)
			if diags.HasError() {
				return nil
			}

			apiDR.Type = "timed"
			val := int(timedModel.Hours.ValueInt64())
			apiDR.Hours = &val
			drChanges = true
		} else if !drModel.Infinite.IsNull() && !drModel.Infinite.IsUnknown() {
			var infModel InfiniteDataRetentionModel
			respDiags := drModel.Infinite.As(ctx, &infModel, basetypes.ObjectAsOptions{})
			diags.Append(respDiags...)
			if diags.HasError() {
				return nil
			}

			if infModel.Enabled.ValueBool() { // Only set if enabled is true
				apiDR.Type = "infinite"
				drChanges = true
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
		drAttrs["timed"] = types.ObjectNull(timedDataRetentionAttributeTypes())
		drAttrs["infinite"] = types.ObjectNull(infiniteDataRetentionAttributeTypes())

		if apiConfig.DataRetention.Type == "timed" && apiConfig.DataRetention.Hours != nil {
			timedAttrs := map[string]attr.Value{
				"hours": types.Int64Value(int64(*apiConfig.DataRetention.Hours)),
			}
			timedObj, timedObjDiags := types.ObjectValue(timedDataRetentionAttributeTypes(), timedAttrs)
			diags.Append(timedObjDiags...)
			drAttrs["timed"] = timedObj
		} else if apiConfig.DataRetention.Type == "infinite" {
			infiniteAttrs := map[string]attr.Value{
				"enabled": types.BoolValue(true), // Mark as enabled if type is infinite
			}
			infObj, infObjDiags := types.ObjectValue(infiniteDataRetentionAttributeTypes(), infiniteAttrs)
			diags.Append(infObjDiags...)
			drAttrs["infinite"] = infObj
		}
		drObj, drObjDiags := types.ObjectValue(dataRetentionAttributeTypes(), drAttrs)
		diags.Append(drObjDiags...)
		attrs["data_retention"] = drObj
	} else {
		attrs["data_retention"] = types.ObjectNull(dataRetentionAttributeTypes())
	}

	objVal, objDiags := types.ObjectValue(capabilityConfigAttributeTypes(), attrs)
	diags.Append(objDiags...)
	return objVal
}
