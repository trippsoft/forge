// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"testing"

	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

func TestObjectTypeCtyType(t *testing.T) {
	tests := []struct {
		name     string
		object   *objectType
		expected cty.Type
	}{
		{
			name:     "empty object",
			object:   Object(),
			expected: cty.EmptyObject,
		},
		{
			name: "simple object",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(0)),
			),
			expected: cty.Object(map[string]cty.Type{
				"name": cty.String,
				"age":  cty.Number,
			}),
		},
		{
			name: "complex object",
			object: Object(
				RequiredField("name", String),
				RequiredField("address", Object(
					RequiredField("city", String),
					OptionalField("zip", Number).WithDefaultValue(cty.NumberIntVal(0)),
				)),
			),
			expected: cty.Object(map[string]cty.Type{
				"name": cty.String,
				"address": cty.Object(map[string]cty.Type{
					"city": cty.String,
					"zip":  cty.Number,
				}),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.object.CtyType()
			if !actual.IsObjectType() {
				t.Fatalf(`expected "object" type from CtyType(), got %q`, actual.FriendlyName())
			}

			if !actual.Equals(tt.expected) {
				t.Errorf("expected %q type from CtyType(), got %q", tt.expected.FriendlyName(), actual.FriendlyName())
			}
		})
	}
}

func TestObjectTypeCtyType_Nil(t *testing.T) {
	var object *objectType

	actual := object.CtyType()
	expected := cty.NilType
	if !actual.Equals(expected) {
		t.Errorf("expected nil type from CtyType(), got %q", actual.FriendlyName())
	}
}

func TestObjectTypeConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		object   Type
		input    cty.Value
		expected cty.Value
	}{
		{
			name: "valid object",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
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
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapVal(map[string]cty.Value{
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
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
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
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
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
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
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
				RequiredField("name", String).WithAliases("fullname", "title"),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
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
				RequiredField("name", String).WithAliases("fullname", "title"),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
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
				OptionalField("field1", String),
				OptionalField("field2", String),
			).WithConstraints(MutuallyExclusive("field1", "field2")),
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
				OptionalField("field1", String),
				OptionalField("field2", String),
			).WithConstraints(MutuallyExclusive("field1", "field2")),
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
				OptionalField("field1", String),
				OptionalField("field2", String),
			).WithConstraints(MutuallyExclusive("field1", "field2")),
			input: cty.ObjectVal(map[string]cty.Value{}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.NullVal(cty.String),
			}),
		},
		{
			name: "valid map with string values (age will be converted)",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapVal(map[string]cty.Value{
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
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapVal(map[string]cty.Value{
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
				OptionalField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapValEmpty(cty.String),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.NullVal(cty.String),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "null object",
			object: Object(
				RequiredField("name", String).WithAliases("fullname", "title"),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.NullVal(cty.Object(map[string]cty.Type{
				"name": cty.String,
				"age":  cty.Number,
			})),
			expected: cty.NullVal(cty.Object(map[string]cty.Type{
				"name": cty.String,
				"age":  cty.Number,
			})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, tt.object, tt.input, tt.expected)
		})
	}
}

func TestObjectTypeConvert_UnknownValue(t *testing.T) {
	object := Object(RequiredField("name", String))
	expectedError := "cannot convert unknown value"
	verifyFailedConversion(t, object, cty.UnknownVal(object.CtyType()), expectedError)
}

func TestObjectTypeConvert_InvalidType(t *testing.T) {
	object := Object(RequiredField("name", String))

	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "string",
			input: cty.StringVal("not an object"),
		},
		{
			name:  "number",
			input: cty.NumberIntVal(42),
		},
		{
			name:  "list",
			input: cty.ListVal([]cty.Value{cty.StringVal("item")}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedError := fmt.Sprintf(
				"cannot convert %q to %q",
				tt.input.Type().FriendlyName(),
				object.CtyType().FriendlyName(),
			)

			verifyFailedConversion(t, object, tt.input, expectedError)
		})
	}
}

func TestObjectTypeConvert_InvalidAttributes(t *testing.T) {
	object := Object(RequiredField("name", String))

	input := cty.ObjectVal(map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	})

	expectedError := `invalid indexes found: "invalid"`
	verifyFailedConversion(t, object, input, expectedError)
}

func TestObjectTypeConvert_FieldConversionError(t *testing.T) {
	object := Object(RequiredField("age", Number))

	input := cty.ObjectVal(map[string]cty.Value{
		"age": cty.StringVal("not-a-number"),
	})

	expectedError := fmt.Sprintf(
		`cannot convert field "age": cannot convert %q to %q: a number is required`,
		cty.String.FriendlyName(),
		cty.Number.FriendlyName(),
	)

	verifyFailedConversion(t, object, input, expectedError)
}

func TestObjectTypeConvert_MultipleAliasesDefined(t *testing.T) {
	object := Object(RequiredField("name", String).WithAliases("fullname"))

	input := cty.ObjectVal(map[string]cty.Value{
		"name":     cty.StringVal("John"),
		"fullname": cty.StringVal("John Doe"),
	})

	expectedError := `field "name" is defined multiple times as "name", "fullname"`
	verifyFailedConversion(t, object, input, expectedError)
}

func TestObjectTypeConvert_ConversionError(t *testing.T) {
	object := Object(RequiredField("name", String))

	input := cty.StringVal("not-an-object")

	expectedError := fmt.Sprintf(
		"cannot convert %q to %q",
		input.Type().FriendlyName(),
		object.CtyType().FriendlyName(),
	)

	verifyFailedConversion(t, object, input, expectedError)
}

func TestObjectTypeConvert_InvalidIndexesInMap(t *testing.T) {
	object := Object(RequiredField("name", String))

	input := cty.MapVal(map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	})

	expectedError := `invalid indexes found: "invalid"`
	verifyFailedConversion(t, object, input, expectedError)
}

func TestObjectTypeConvert_AliasConversionError(t *testing.T) {
	object := Object(RequiredField("age", Number).WithAliases("alias"))

	input := cty.ObjectVal(map[string]cty.Value{
		"alias": cty.StringVal("not-a-number"),
	})

	expectedError := fmt.Sprintf(
		`cannot convert field "age" (alias "alias"): cannot convert %q to %q: a number is required`,
		cty.String.FriendlyName(),
		cty.Number.FriendlyName(),
	)

	verifyFailedConversion(t, object, input, expectedError)
}

func TestObjectTypeConvert_Nil(t *testing.T) {
	var object *objectType

	expectedError := "object type is nil"
	converted, err := object.Convert(cty.StringVal("test"))
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	if converted.Equals(cty.NilVal) != cty.True {
		t.Fatalf("expected nil value from Convert(), got %s", util.FormatCtyValueToString(converted))
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidate_Pass(t *testing.T) {
	tests := []struct {
		name   string
		object Type
		input  cty.Value
	}{
		{
			name: "valid object",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			}),
		},
		{
			name: "valid map",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			}),
		},
		{
			name: "object with missing optional field uses default",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Bob"),
			}),
		},
		{
			name: "object with type conversion",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Alice"),
				"age":  cty.StringVal("35"), // String to Number conversion
			}),
		},
		{
			name: "using primary name",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
			}),
		},
		{
			name: "using first alias",
			object: Object(
				RequiredField("name", String).WithAliases("fullname", "title"),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"fullname": cty.StringVal("John Doe"),
			}),
		},
		{
			name: "using second alias",
			object: Object(
				RequiredField("name", String).WithAliases("fullname", "title"),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"title": cty.StringVal("Mr. John"),
			}),
		},
		{
			name: "constraint passes - only field1",
			object: Object(
				OptionalField("field1", String),
				OptionalField("field2", String),
			).WithConstraints(MutuallyExclusive("field1", "field2")),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
		},
		{
			name: "constraint passes - only field2",
			object: Object(
				OptionalField("field1", String),
				OptionalField("field2", String),
			).WithConstraints(MutuallyExclusive("field1", "field2")),
			input: cty.ObjectVal(map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			}),
		},
		{
			name: "constraint passes - neither field",
			object: Object(
				OptionalField("field1", String),
				OptionalField("field2", String),
			).WithConstraints(MutuallyExclusive("field1", "field2")),
			input: cty.ObjectVal(map[string]cty.Value{}),
		},
		{
			name: "valid map with string values (age will be converted)",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.StringVal("30"), // String that can be converted to number
			}),
		},
		{
			name: "map with missing optional field",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapVal(map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			}),
		},
		{
			name: "empty map uses defaults",
			object: Object(
				OptionalField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.MapValEmpty(cty.String),
		},
		{
			name: "null object",
			object: Object(
				RequiredField("name", String).WithAliases("fullname", "title"),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(25)),
			),
			input: cty.NullVal(cty.Object(map[string]cty.Type{
				"name": cty.String,
				"age":  cty.Number,
			})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, tt.object, tt.input)
		})
	}
}

func TestObjectTypeValidate_UnknownValue(t *testing.T) {
	tests := []struct {
		name   string
		object Type
		input  cty.Value
	}{
		{
			name: "unknown object",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(0)),
			),
			input: cty.UnknownVal(cty.Object(map[string]cty.Type{
				"name": cty.String,
				"age":  cty.Number,
			})),
		},
		{
			name: "object with unknown element",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(0)),
			),
			input: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.UnknownVal(cty.Number),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedError := "cannot convert unknown value"
			err := tt.object.Validate(tt.input)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", expectedError)
			}

			if err.Error() != expectedError {
				t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
			}
		})
	}
}

func TestObjectTypeValidate_NotMapOrObject(t *testing.T) {
	tests := []struct {
		name          string
		object        Type
		input         cty.Value
		expectedError string
	}{
		{
			name: "string",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(0)),
			),
			input:         cty.StringVal("not-an-object"),
			expectedError: `cannot convert "string" to "object"`,
		},
		{
			name: "number",
			object: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(0)),
			),
			input:         cty.NumberIntVal(42),
			expectedError: `cannot convert "number" to "object"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.object.Validate(tt.input)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestObjectTypeValidate_MissingRequiredField(t *testing.T) {
	object := Object(
		RequiredField("name", String),
		OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(0)),
	)

	input := cty.ObjectVal(map[string]cty.Value{
		"age": cty.NumberIntVal(30),
	})

	expectedError := `missing required field "name"`
	verifyFailedValidation(t, object, input, expectedError)
}

func TestObjectTypeValidate_FieldValidationFailure(t *testing.T) {
	object := Object(RequiredField("duration", Duration))

	input := cty.ObjectVal(map[string]cty.Value{
		"duration": cty.StringVal("invalid-duration"),
	})

	expectedError := `field "duration" validation failed: time: invalid duration "invalid-duration"`
	verifyFailedValidation(t, object, input, expectedError)
}

func TestObjectTypeValidate_ConstraintFailure(t *testing.T) {
	object := Object(
		OptionalField("field1", String),
		OptionalField("field2", String),
	).WithConstraints(MutuallyExclusive("field1", "field2"))

	input := cty.ObjectVal(map[string]cty.Value{
		"field1": cty.StringVal("value1"),
		"field2": cty.StringVal("value2"),
	})

	expectedError := `validation failed: mutually exclusive fields "field1", "field2" are all present`
	verifyFailedValidation(t, object, input, expectedError)
}

func TestObjectTypeValidate_Nil(t *testing.T) {
	var object *objectType

	expectedError := "object type is nil"
	err := object.Validate(cty.NilVal)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name string
		obj  *objectType
	}{
		{
			name: "valid spec",
			obj: Object(
				RequiredField("name", String),
				OptionalField("age", Number).WithDefaultValue(cty.NumberIntVal(0)),
			),
		},
		{
			name: "valid spec with constraints",
			obj: Object(
				OptionalField("field1", String),
				OptionalField("field2", String),
			).WithConstraints(MutuallyExclusive("field1", "field2")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.obj.ValidateSpec()
			if err != nil {
				t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
			}
		})
	}
}

func TestObjectTypeValidateSpec_FieldErrors(t *testing.T) {
	obj := Object(
		RequiredField("invalid", String).WithDefaultValue(cty.StringVal("default")), // Required field with default
	)

	expectedError := `field "invalid" is required and has a default value`
	err := obj.ValidateSpec()
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_DuplicateFieldName(t *testing.T) {
	object := Object(RequiredField("name", String), RequiredField("name", String))

	expectedError := `field "name" is defined multiple times`
	err := object.ValidateSpec()
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_DuplicateAliases(t *testing.T) {
	obj := Object(RequiredField("name", String).WithAliases("name"))

	expectedError := `field "name" is defined multiple times (aliases: "name", "name")`
	err := obj.ValidateSpec()
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_NilField(t *testing.T) {
	object := Object(RequiredField("name", String), OptionalField("age", Number), nil)

	expectedError := "nil field definition found"
	err := object.ValidateSpec()
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_InvalidConstraint(t *testing.T) {
	object := Object(
		OptionalField("field1", String),
	).WithConstraints(MutuallyExclusive("field1", "nonexistent"))

	expectedError := `constraint validation failed: field "nonexistent" is not defined in the object type`
	err := object.ValidateSpec()
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_Nil(t *testing.T) {
	var object *objectType

	expectedError := "object type is nil"
	err := object.ValidateSpec()
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
