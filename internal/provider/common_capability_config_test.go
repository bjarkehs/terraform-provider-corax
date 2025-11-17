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
				objVal, ok := underlyingVal.(types.Object)
				if !ok {
					t.Errorf("expected underlying value to be types.Object, got %T", underlyingVal)
					return
				}
				attrs := objVal.Attributes()
				if len(attrs) != 0 {
					t.Errorf("expected empty object, got %d attributes", len(attrs))
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
				objVal, ok := underlyingVal.(types.Object)
				if !ok {
					t.Errorf("expected underlying value to be types.Object, got %T", underlyingVal)
					return
				}
				attrs := objVal.Attributes()
				if len(attrs) != 2 {
					t.Errorf("expected 2 attributes, got %d", len(attrs))
				}

				key1Val, ok := attrs["key1"].(types.String)
				if !ok || key1Val.ValueString() != "value1" {
					t.Errorf("expected key1=value1")
				}

				key2Val, ok := attrs["key2"].(types.String)
				if !ok || key2Val.ValueString() != "value2" {
					t.Errorf("expected key2=value2")
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
				objVal, ok := underlyingVal.(types.Object)
				if !ok {
					t.Errorf("expected underlying value to be types.Object, got %T", underlyingVal)
					return
				}
				attrs := objVal.Attributes()
				if len(attrs) != 4 {
					t.Errorf("expected 4 attributes, got %d", len(attrs))
				}

				stringVal, ok := attrs["stringKey"].(types.String)
				if !ok || stringVal.ValueString() != "test" {
					t.Errorf("expected stringKey=test")
				}

				boolVal, ok := attrs["boolKey"].(types.Bool)
				if !ok || !boolVal.ValueBool() {
					t.Errorf("expected boolKey=true")
				}

				intVal, ok := attrs["intKey"].(types.Int64)
				if !ok || intVal.ValueInt64() != 42 {
					t.Errorf("expected intKey=42")
				}

				floatVal, ok := attrs["floatKey"].(types.Float64)
				if !ok || floatVal.ValueFloat64() != 3.14 {
					t.Errorf("expected floatKey=3.14")
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
				objVal, ok := underlyingVal.(types.Object)
				if !ok {
					t.Errorf("expected underlying value to be types.Object, got %T", underlyingVal)
					return
				}
				attrs := objVal.Attributes()
				if len(attrs) != 2 {
					t.Errorf("expected 2 attributes, got %d", len(attrs))
				}

				nestedVal, ok := attrs["nested"].(types.Object)
				if !ok {
					t.Errorf("expected nested to be types.Object, got %T", attrs["nested"])
					return
				}
				nestedAttrs := nestedVal.Attributes()
				innerKeyVal, ok := nestedAttrs["innerKey"].(types.String)
				if !ok || innerKeyVal.ValueString() != "innerValue" {
					t.Errorf("expected nested.innerKey=innerValue")
				}

				arrayVal, ok := attrs["array"].(types.List)
				if !ok {
					t.Errorf("expected array to be types.List, got %T", attrs["array"])
					return
				}
				if len(arrayVal.Elements()) != 2 {
					t.Errorf("expected array to have 2 elements, got %d", len(arrayVal.Elements()))
				}
			},
		},
		{
			name: "should return Object not String (reproduces Terraform error)",
			input: map[string]interface{}{
				"reasoning": "minimal",
				"verbosity": "low",
			},
			expectNull: false,
			validateValue: func(t *testing.T, result types.Dynamic) {
				if result.IsNull() {
					t.Error("expected non-null result")
				}
				underlyingVal := result.UnderlyingValue()

				// The underlying value should be types.Object, NOT types.String
				// If it's types.String, Terraform will see it as a string type and error with:
				// "attribute custom_parameters: object required, but have string"
				objVal, ok := underlyingVal.(types.Object)
				if !ok {
					t.Errorf("expected underlying value to be types.Object, got %T (this is the bug!)", underlyingVal)
					return
				}

				// Verify the object has the correct attributes
				attrs := objVal.Attributes()
				if len(attrs) != 2 {
					t.Errorf("expected 2 attributes, got %d", len(attrs))
				}

				// Check that the values are correct
				reasoningVal, ok := attrs["reasoning"]
				if !ok {
					t.Error("expected 'reasoning' attribute")
				} else {
					reasoningStr, ok := reasoningVal.(types.String)
					if !ok {
						t.Errorf("expected reasoning to be types.String, got %T", reasoningVal)
					} else if reasoningStr.ValueString() != "minimal" {
						t.Errorf("expected reasoning='minimal', got %q", reasoningStr.ValueString())
					}
				}

				verbosityVal, ok := attrs["verbosity"]
				if !ok {
					t.Error("expected 'verbosity' attribute")
				} else {
					verbosityStr, ok := verbosityVal.(types.String)
					if !ok {
						t.Errorf("expected verbosity to be types.String, got %T", verbosityVal)
					} else if verbosityStr.ValueString() != "low" {
						t.Errorf("expected verbosity='low', got %q", verbosityStr.ValueString())
					}
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
