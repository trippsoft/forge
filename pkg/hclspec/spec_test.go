// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestSpecConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		object   *objectType
		input    cty.Value
		expected cty.Value
	}{
		{
			name: "valid object",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			}),
		},
		{
			name: "valid map",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "object with missing optional field uses default",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Bob"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Bob"),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "object with type conversion",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Alice"),
				"age":  cty.StringVal("35"), // String to Number conversion
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Alice"),
				"age":  cty.NumberIntVal(35),
			}),
		},
		{
			name: "using primary name",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "using first alias",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"fullname": cty.StringVal("John Doe"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John Doe"),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "using second alias",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"title": cty.StringVal("Mr. John"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Mr. John"),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "constraint passes - only field1",
			object: Object(
				OptionalField("field1", String()),
				OptionalField("field2", String()),
			).WithConstraints(MutuallyExclusive("field1", "field2")).(*objectType),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.NullVal(cty.String),
			}),
		},
		{
			name: "constraint passes - only field2",
			object: Object(
				OptionalField("field1", String()),
				OptionalField("field2", String()),
			).WithConstraints(MutuallyExclusive("field1", "field2")).(*objectType),
			input: cty.ObjectVal(map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.StringVal("value2"),
			}),
		},
		{
			name: "constraint passes - neither field",
			object: Object(
				OptionalField("field1", String()),
				OptionalField("field2", String()),
			).WithConstraints(MutuallyExclusive("field1", "field2")).(*objectType),
			input: cty.ObjectVal(map[string]cty.Value{}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.NullVal(cty.String),
			}),
		},
		{
			name: "valid map with string values (age will be converted)",
			object: Object(
				RequiredField("name", String()),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.StringVal("30"), // String that can be converted to number
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			}),
		},
		{
			name: "map with missing optional field",
			object: Object(
				RequiredField("name", String()),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "empty map uses defaults",
			object: Object(
				OptionalField("name", String()),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.NullVal(cty.String),
				"age":  cty.NumberIntVal(25),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := NewSpec(tt.object)
			verifySuccessfulSpecConversion(t, spec, tt.input, tt.expected)
		})
	}
}

func TestSpecConvert_NilObject(t *testing.T) {
	spec := NewSpec(nil)
	expectedError := "object type is nil"
	verifyFailedSpecConversion(t, spec, cty.EmptyObjectVal, expectedError)
}

func TestSpecConvert_InvalidAttributes(t *testing.T) {
	spec := NewSpec(Object(RequiredField("name", String())))

	input := cty.ObjectVal(map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	})

	expectedError := `invalid indexes found: "invalid"`
	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecConvert_FieldConversionError(t *testing.T) {
	spec := NewSpec(Object(RequiredField("age", Float32())))

	input := cty.ObjectVal(map[string]cty.Value{
		"age": cty.StringVal("not-a-number"),
	})

	expectedError := fmt.Sprintf(
		`cannot convert field "age": cannot convert %q to %q: a number is required`,
		cty.String.FriendlyName(),
		cty.Number.FriendlyName(),
	)

	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecConvert_MultipleAliasesDefined(t *testing.T) {
	spec := NewSpec(Object(RequiredField("name", String()).WithAliases("fullname")))

	input := cty.ObjectVal(map[string]cty.Value{
		"name":     cty.StringVal("John"),
		"fullname": cty.StringVal("John Doe"),
	})

	expectedError := `field "name" is defined multiple times as "name", "fullname"`
	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecConvert_MissingRequiredField(t *testing.T) {
	spec := NewSpec(Object(
		RequiredField("name", String()),
		OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(0)),
	))

	input := cty.ObjectVal(map[string]cty.Value{
		"age": cty.NumberIntVal(30),
	})

	expectedError := `missing required field "name"`
	verifyFailedSpecValidation(t, spec, input, expectedError)
}

func TestSpecConvert_InvalidIndexesInMap(t *testing.T) {
	spec := NewSpec(Object(RequiredField("name", String())))

	input := cty.ObjectVal(map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	})

	expectedError := `invalid indexes found: "invalid"`
	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecConvert_Nil(t *testing.T) {
	var spec *Spec
	verifyFailedSpecConversion(t, spec, cty.EmptyObjectVal, "spec is nil")
}

func TestSpecValidate_Pass(t *testing.T) {
	tests := []struct {
		name   string
		object *objectType
		input  cty.Value
	}{
		{
			name: "valid object",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			}),
		},
		{
			name: "valid map",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			}),
		},
		{
			name: "object with missing optional field uses default",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Bob"),
			}),
		},
		{
			name: "object with type conversion",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Alice"),
				"age":  cty.StringVal("35"), // String to Number conversion
			}),
		},
		{
			name: "using primary name",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
			}),
		},
		{
			name: "using first alias",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"fullname": cty.StringVal("John Doe"),
			}),
		},
		{
			name: "using second alias",
			object: Object(
				RequiredField("name", String()).WithAliases("fullname", "title"),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"title": cty.StringVal("Mr. John"),
			}),
		},
		{
			name: "constraint passes - only field1",
			object: Object(
				OptionalField("field1", String()),
				OptionalField("field2", String()),
			).WithConstraints(MutuallyExclusive("field1", "field2")).(*objectType),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
		},
		{
			name: "constraint passes - only field2",
			object: Object(
				OptionalField("field1", String()),
				OptionalField("field2", String()),
			).WithConstraints(MutuallyExclusive("field1", "field2")).(*objectType),
			input: cty.ObjectVal(map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			}),
		},
		{
			name: "constraint passes - neither field",
			object: Object(
				OptionalField("field1", String()),
				OptionalField("field2", String()),
			).WithConstraints(MutuallyExclusive("field1", "field2")).(*objectType),
			input: cty.EmptyObjectVal,
		},
		{
			name: "valid map with string values (age will be converted)",
			object: Object(
				RequiredField("name", String()),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.StringVal("30"), // String that can be converted to number
			}),
		},
		{
			name: "map with missing optional field",
			object: Object(
				RequiredField("name", String()),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			}),
		},
		{
			name: "empty map uses defaults",
			object: Object(
				OptionalField("name", String()),
				OptionalField("age", Float32()).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.EmptyObjectVal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := NewSpec(tt.object)
			verifySuccessfulSpecValidation(t, spec, tt.input)
		})
	}
}

func TestSpecValidate_Nil(t *testing.T) {
	var spec *Spec

	expectedError := "spec is nil"
	err := spec.Validate(cty.NilVal)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err)
	}
}

func TestSpecValidate_NilObject(t *testing.T) {
	spec := NewSpec(nil)

	expectedError := "object type is nil"
	err := spec.Validate(cty.NilVal)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err)
	}
}

func TestSpecConvert_ConstraintFailure(t *testing.T) {
	spec := NewSpec(Object(
		OptionalField("field1", String()),
		OptionalField("field2", String()),
	).WithConstraints(MutuallyExclusive("field1", "field2")).(*objectType))

	input := cty.ObjectVal(map[string]cty.Value{
		"field1": cty.StringVal("value1"),
		"field2": cty.StringVal("value2"),
	})

	expectedError := `mutually exclusive fields "field1", "field2" are all present`
	verifyFailedSpecValidation(t, spec, input, expectedError)
}

func verifySuccessfulSpecConversion(t *testing.T, spec *Spec, input cty.Value, expected cty.Value) {
	t.Helper()

	actual, err := spec.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got error %q", err.Error())
	}

	if !actual.Type().IsObjectType() {
		t.Fatalf("expected output to be an object, got %q", actual.Type().FriendlyName())
	}

	actualMap := actual.AsValueMap()
	expectedMap := expected.AsValueMap()

	if len(expectedMap) != len(actualMap) {
		t.Fatalf("expected %d fields, got %d", len(expectedMap), len(actualMap))
	}

	for key, expectedValue := range expectedMap {
		actualValue, exists := actualMap[key]
		if !exists {
			t.Errorf("expected field %q to be present in output, but it was not found", key)
			continue
		}

		if expectedValue.Equals(actualValue) != cty.True {
			t.Errorf(
				"expected field %q to have value %q, got %q",
				key,
				expectedValue.GoString(),
				actualValue.GoString(),
			)
		}
	}
}

func verifyFailedSpecConversion(t *testing.T, spec *Spec, input cty.Value, expectedError string) {
	t.Helper()

	_, err := spec.Convert(input)
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}
}

func verifySuccessfulSpecValidation(t *testing.T, spec *Spec, input cty.Value) {
	t.Helper()

	converted, err := spec.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got error %q", err.Error())
	}

	err = spec.Validate(converted)
	if err != nil {
		t.Fatalf("expected no error from Validate(), got error %q", err.Error())
	}
}

func verifyFailedSpecValidation(t *testing.T, spec *Spec, input cty.Value, expectedError string) {
	t.Helper()

	converted, err := spec.Convert(input)
	if err != nil {
		t.Errorf("expected no error from Convert(), got error %q", err.Error())
	}

	err = spec.Validate(converted)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
