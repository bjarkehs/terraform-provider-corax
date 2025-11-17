// Copyright (c) Trifork

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestCustomParametersToAPI(t *testing.T) {
	tests := []struct {
		name          string
		input         types.Dynamic
		expectedMap   map[string]interface{}
		expectError   bool
		errorContains string
	}{
		{
			name:        "nil dynamic value",
			input:       types.DynamicNull(),
			expectedMap: nil,
			expectError: false,
		},
		{
			name:        "unknown dynamic value",
			input:       types.DynamicUnknown(),
			expectedMap: nil,
			expectError: false,
		},
		{
			name:  "valid JSON string with mixed types",
			input: types.DynamicValue(types.StringValue(`{"key1":"value1","key2":123,"key3":true}`)),
			expectedMap: map[string]interface{}{
				"key1": "value1",
				"key2": float64(123),
				"key3": true,
			},
			expectError: false,
		},
		{
			name:          "invalid JSON string",
			input:         types.DynamicValue(types.StringValue(`{invalid json}`)),
			expectedMap:   nil,
			expectError:   true,
			errorContains: "custom_parameters was provided as a string, but it's not valid JSON",
		},
		{
			name:        "null string value",
			input:       types.DynamicValue(basetypes.NewStringNull()),
			expectedMap: nil,
			expectError: false,
		},
		{
			name:  "JSON string with nested objects",
			input: types.DynamicValue(types.StringValue(`{"nested":{"innerKey":"innerValue"},"topKey":"topValue"}`)),
			expectedMap: map[string]interface{}{
				"nested": map[string]interface{}{
					"innerKey": "innerValue",
				},
				"topKey": "topValue",
			},
			expectError: false,
		},
		{
			name:  "JSON string with array values",
			input: types.DynamicValue(types.StringValue(`{"items":["item1","item2","item3"]}`)),
			expectedMap: map[string]interface{}{
				"items": []interface{}{"item1", "item2", "item3"},
			},
			expectError: false,
		},
		{
			name:          "empty string value",
			input:         types.DynamicValue(types.StringValue("")),
			expectedMap:   nil,
			expectError:   true,
			errorContains: "custom_parameters was provided as a string, but it's not valid JSON",
		},
		{
			name:        "empty JSON string",
			input:       types.DynamicValue(types.StringValue(`{}`)),
			expectedMap: map[string]interface{}{},
			expectError: false,
		},
		{
			name: "HCL object with string values (user scenario)",
			input: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"reasoning_effort": types.StringType,
					"verbosity":        types.StringType,
				},
				map[string]attr.Value{
					"reasoning_effort": types.StringValue("minimal"),
					"verbosity":        types.StringValue("low"),
				},
			)),
			expectedMap: map[string]interface{}{
				"reasoning_effort": "minimal",
				"verbosity":        "low",
			},
			expectError: false,
		},
		{
			name: "HCL object with mixed types",
			input: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"temperature": types.Float64Type,
					"max_tokens":  types.Int64Type,
					"stream":      types.BoolType,
				},
				map[string]attr.Value{
					"temperature": types.Float64Value(0.7),
					"max_tokens":  types.Int64Value(1000),
					"stream":      types.BoolValue(true),
				},
			)),
			expectedMap: map[string]interface{}{
				"temperature": 0.7,
				"max_tokens":  int64(1000),
				"stream":      true,
			},
			expectError: false,
		},
		{
			name: "HCL map with string values",
			input: types.DynamicValue(types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("value1"),
					"key2": types.StringValue("value2"),
				},
			)),
			expectedMap: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expectError: false,
		},
		{
			name: "null HCL object",
			input: types.DynamicValue(types.ObjectNull(
				map[string]attr.Type{
					"key": types.StringType,
				},
			)),
			expectedMap: nil,
			expectError: false,
		},
		{
			name: "HCL object with nested objects",
			input: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"foo": types.ObjectType{AttrTypes: map[string]attr.Type{
						"bar": types.StringType,
					}},
				},
				map[string]attr.Value{
					"foo": types.ObjectValueMust(
						map[string]attr.Type{
							"bar": types.StringType,
						},
						map[string]attr.Value{
							"bar": types.StringValue("baz"),
						},
					),
				},
			)),
			expectedMap: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			},
			expectError: false,
		},
		{
			name: "HCL object with deeply nested structure",
			input: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"level1": types.ObjectType{AttrTypes: map[string]attr.Type{
						"level2": types.ObjectType{AttrTypes: map[string]attr.Type{
							"level3": types.StringType,
							"value":  types.Int64Type,
						}},
					}},
				},
				map[string]attr.Value{
					"level1": types.ObjectValueMust(
						map[string]attr.Type{
							"level2": types.ObjectType{AttrTypes: map[string]attr.Type{
								"level3": types.StringType,
								"value":  types.Int64Type,
							}},
						},
						map[string]attr.Value{
							"level2": types.ObjectValueMust(
								map[string]attr.Type{
									"level3": types.StringType,
									"value":  types.Int64Type,
								},
								map[string]attr.Value{
									"level3": types.StringValue("deep"),
									"value":  types.Int64Value(42),
								},
							),
						},
					),
				},
			)),
			expectedMap: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": "deep",
						"value":  int64(42),
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics
			result := customParametersToAPI(tt.input, &diags)

			if tt.expectError {
				if !diags.HasError() {
					t.Errorf("expected error but got none")
				}
				if tt.errorContains != "" {
					found := false
					for _, d := range diags.Errors() {
						if contains(d.Summary(), tt.errorContains) || contains(d.Detail(), tt.errorContains) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected error containing %q, but got: %v", tt.errorContains, diags.Errors())
					}
				}
			} else {
				if diags.HasError() {
					t.Errorf("unexpected error: %v", diags.Errors())
				}
			}

			if !mapsEqual(result, tt.expectedMap) {
				t.Errorf("result = %v, want %v", result, tt.expectedMap)
			}
		})
	}
}

func TestCustomParametersAPIToTerraform(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]interface{}
		expectNull    bool
		expectError   bool
		validateValue func(t *testing.T, result types.Dynamic)
	}{
		{
			name:       "nil input",
			input:      nil,
			expectNull: true,
		},
		{
			name:       "empty map",
			input:      map[string]interface{}{},
			expectNull: false,
			validateValue: func(t *testing.T, result types.Dynamic) {
				if result.IsNull() {
					t.Error("expected non-null result for empty map")
				}
				underlyingVal := result.UnderlyingValue()
				strVal, ok := underlyingVal.(types.String)
				if !ok {
					t.Errorf("expected underlying value to be types.String, got %T", underlyingVal)
					return
				}
				if strVal.ValueString() != "{}" {
					t.Errorf("expected '{}', got %q", strVal.ValueString())
				}
			},
		},
		{
			name: "map with string values",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expectNull: false,
			validateValue: func(t *testing.T, result types.Dynamic) {
				if result.IsNull() {
					t.Error("expected non-null result")
				}
				underlyingVal := result.UnderlyingValue()
				strVal, ok := underlyingVal.(types.String)
				if !ok {
					t.Errorf("expected underlying value to be types.String, got %T", underlyingVal)
					return
				}
				jsonStr := strVal.ValueString()
				if !contains(jsonStr, `"key1":"value1"`) || !contains(jsonStr, `"key2":"value2"`) {
					t.Errorf("expected JSON to contain key-value pairs, got %q", jsonStr)
				}
			},
		},
		{
			name: "map with mixed value types",
			input: map[string]interface{}{
				"stringKey": "test",
				"boolKey":   true,
				"intKey":    42,
				"floatKey":  3.14,
			},
			expectNull: false,
			validateValue: func(t *testing.T, result types.Dynamic) {
				if result.IsNull() {
					t.Error("expected non-null result")
				}
				underlyingVal := result.UnderlyingValue()
				strVal, ok := underlyingVal.(types.String)
				if !ok {
					t.Errorf("expected underlying value to be types.String, got %T", underlyingVal)
					return
				}
				jsonStr := strVal.ValueString()
				if !contains(jsonStr, `"stringKey":"test"`) {
					t.Errorf("expected stringKey in JSON, got %q", jsonStr)
				}
				if !contains(jsonStr, `"boolKey":true`) {
					t.Errorf("expected boolKey in JSON, got %q", jsonStr)
				}
				if !contains(jsonStr, `"intKey":42`) {
					t.Errorf("expected intKey in JSON, got %q", jsonStr)
				}
			},
		},
		{
			name: "map with nested structure",
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"innerKey": "innerValue",
				},
				"array": []interface{}{"item1", "item2"},
			},
			expectNull: false,
			validateValue: func(t *testing.T, result types.Dynamic) {
				if result.IsNull() {
					t.Error("expected non-null result")
				}
				underlyingVal := result.UnderlyingValue()
				strVal, ok := underlyingVal.(types.String)
				if !ok {
					t.Errorf("expected underlying value to be types.String, got %T", underlyingVal)
					return
				}
				jsonStr := strVal.ValueString()
				if !contains(jsonStr, `"nested"`) || !contains(jsonStr, `"innerKey"`) {
					t.Errorf("expected nested structure in JSON, got %q", jsonStr)
				}
				if !contains(jsonStr, `"array"`) {
					t.Errorf("expected array in JSON, got %q", jsonStr)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics
			result := customParametersAPIToTerraform(tt.input, &diags)

			if tt.expectError {
				if !diags.HasError() {
					t.Error("expected error but got none")
				}
			} else {
				if diags.HasError() {
					t.Errorf("unexpected error: %v", diags.Errors())
				}
			}

			if tt.expectNull {
				if !result.IsNull() {
					t.Error("expected null result")
				}
			} else {
				if tt.validateValue != nil {
					tt.validateValue(t, result)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func mapsEqual(a, b map[string]interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !valuesEqual(v, bv) {
			return false
		}
	}
	return true
}

func valuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok {
			return false
		}
		return mapsEqual(aVal, bVal)
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !valuesEqual(aVal[i], bVal[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
