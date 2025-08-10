// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestSpecConvert(t *testing.T) {

	tests := []struct {
		name     string
		object   *objectType
		input    map[string]cty.Value
		expected map[string]cty.Value
	}{
		{
			name: "valid object",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			},
		},
		{
			name: "valid map",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("Jane"),
				"age":  cty.NumberIntVal(25),
			},
		},
		{
			name: "object with missing optional field uses default",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"name": cty.StringVal("Bob"),
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("Bob"),
				"age":  cty.NumberIntVal(25),
			},
		},
		{
			name: "object with type conversion",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"name": cty.StringVal("Alice"),
				"age":  cty.StringVal("35"), // String to Number conversion
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("Alice"),
				"age":  cty.NumberIntVal(35),
			},
		},
		{
			name: "using primary name",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"name": cty.StringVal("John"),
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(25),
			},
		},
		{
			name: "using first alias",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"fullname": cty.StringVal("John Doe"),
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("John Doe"),
				"age":  cty.NumberIntVal(25),
			},
		},
		{
			name: "using second alias",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"title": cty.StringVal("Mr. John"),
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("Mr. John"),
				"age":  cty.NumberIntVal(25),
			},
		},
		{
			name: "constraint passes - only field1",
			object: Object(map[string]*ObjectField{
				"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
			}, MutuallyExclusive("field1", "field2")),
			input: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
			expected: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.NullVal(cty.String),
			},
		},
		{
			name: "constraint passes - only field2",
			object: Object(map[string]*ObjectField{
				"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
			}, MutuallyExclusive("field1", "field2")),
			input: map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			},
			expected: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.StringVal("value2"),
			},
		},
		{
			name: "constraint passes - neither field",
			object: Object(map[string]*ObjectField{
				"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
			}, MutuallyExclusive("field1", "field2")),
			input: map[string]cty.Value{},
			expected: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.NullVal(cty.String),
			},
		},
		{
			name: "valid map with string values (age will be converted)",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.StringVal("30"), // String that can be converted to number
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("John"),
				"age":  cty.NumberIntVal(30),
			},
		},
		{
			name: "map with missing optional field",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{
				"name": cty.StringVal("Jane"),
			},
			expected: map[string]cty.Value{
				"name": cty.StringVal("Jane"),
				"age":  cty.NumberIntVal(25),
			},
		},
		{
			name: "empty map uses defaults",
			object: Object(map[string]*ObjectField{
				"name": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
			}),
			input: map[string]cty.Value{},
			expected: map[string]cty.Value{
				"name": cty.NullVal(cty.String),
				"age":  cty.NumberIntVal(25),
			},
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
	verifyFailedSpecConversion(t, spec, map[string]cty.Value{}, expectedError)
}

func TestSpecConvert_NilValues(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
	}))

	expectedError := "cannot convert nil map"
	verifyFailedSpecConversion(t, spec, nil, expectedError)
}

func TestSpecConvert_InvalidAttributes(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
	}))

	input := map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	}

	expectedError := fmt.Sprintf(
		"invalid indexes found: %v",
		[]string{"invalid"},
	)

	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecConvert_FieldConversionError(t *testing.T) {
	spec := NewSpec(Object(map[string]*ObjectField{
		"age": {Type: Number, Required: true, DefaultValue: cty.NullVal(cty.Number)},
	}))

	input := map[string]cty.Value{
		"age": cty.StringVal("not-a-number"),
	}

	expectedError := fmt.Sprintf(
		"cannot convert field %q: cannot convert %q to %q: a number is required",
		"age",
		cty.String.FriendlyName(),
		cty.Number.FriendlyName())

	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecConvert_MultipleAliasesDefined(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname"}},
	}))

	input := map[string]cty.Value{
		"name":     cty.StringVal("John"),
		"fullname": cty.StringVal("John Doe"),
	}

	expectedError := fmt.Sprintf("field %q is defined multiple times as %v", "name", []string{"name", "fullname"})
	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecConvert_MissingRequiredField(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
		"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(0)},
	}))

	input := map[string]cty.Value{
		"age": cty.NumberIntVal(30),
	}

	expectedError := fmt.Sprintf("missing required field %q", "name")
	verifyFailedSpecValidation(t, spec, input, expectedError)
}

func TestSpecConvert_FieldValidationFailure(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"duration": {Type: Duration, Required: true, DefaultValue: cty.NullVal(cty.String)},
	}))

	input := map[string]cty.Value{
		"duration": cty.StringVal("invalid-duration"),
	}

	expectedError := fmt.Sprintf(
		"field %q validation failed: time: invalid duration %q",
		"duration",
		"invalid-duration")

	verifyFailedSpecValidation(t, spec, input, expectedError)
}

func TestSpecConvert_ConstraintFailure(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
		"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
	}, MutuallyExclusive("field1", "field2")))

	input := map[string]cty.Value{
		"field1": cty.StringVal("value1"),
		"field2": cty.StringVal("value2"),
	}

	expectedError := fmt.Sprintf(
		"validation failed: mutually exclusive fields %q are all present",
		[]string{"field1", "field2"})
	verifyFailedSpecValidation(t, spec, input, expectedError)
}

func TestSpecConvert_InvalidIndexesInMap(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
	}))

	input := map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	}

	expectedError := fmt.Sprintf("invalid indexes found: %v", []string{"invalid"})
	verifyFailedSpecConversion(t, spec, input, expectedError)
}

func TestSpecValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name        string
		fields      map[string]*ObjectField
		constraints []objectConstraint
	}{
		{
			name: "valid spec",
			fields: map[string]*ObjectField{
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(0)},
			},
			constraints: nil,
		},
		{
			name: "valid spec with constraints",
			fields: map[string]*ObjectField{
				"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
			},
			constraints: []objectConstraint{MutuallyExclusive("field1", "field2")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			spec := NewSpec(Object(tt.fields, tt.constraints...))
			errs := spec.ValidateSpec()

			if len(errs) != 0 {
				t.Errorf("expected no errors from ValidateSpec(), got %d errors", len(errs))

				for _, err := range errs {
					t.Errorf("expected no errors from ValidateSpec(), got error: %v", err)
				}
			}
		})
	}
}

func TestSpecValidateSpec_NilObject(t *testing.T) {

	spec := NewSpec(nil)

	expectedError := "object type is nil"
	errs := spec.ValidateSpec()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from ValidateSpec(), got %d errors: %v", len(errs), errs)
	}

	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func TestSpecValidateSpec_FieldErrors(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"invalid": {
			Type:         String,
			Required:     true,
			DefaultValue: cty.StringVal("default"), // Required field with default
		},
	}))

	errs := spec.ValidateSpec()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from ValidateSpec(), got %d errors", len(errs))
	}

	expectedError := fmt.Sprintf("field %q is required and has a default value", "invalid")
	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func TestSpecValidateSpec_DuplicateFieldNames(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"name"}},
	}))

	errs := spec.ValidateSpec()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from ValidateSpec(), got %d errors", len(errs))
	}

	expectedError := fmt.Sprintf(
		"field %q is defined multiple times (aliases: %v)",
		"name",
		[]string{"name", "name"})

	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func TestSpecValidateSpec_InvalidConstraint(t *testing.T) {

	spec := NewSpec(Object(map[string]*ObjectField{
		"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
	}, MutuallyExclusive("field1", "nonexistent")))

	errs := spec.ValidateSpec()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from ValidateSpec(), got %d errors", len(errs))
	}

	expectedError := fmt.Sprintf(
		"constraint validation failed: field %q is not defined in the object type",
		"nonexistent")

	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func verifySuccessfulSpecConversion(t *testing.T, spec *Spec, input map[string]cty.Value, expected map[string]cty.Value) {
	t.Helper()

	actual, err := spec.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got error: %v", err)
	}

	if len(expected) != len(actual) {
		t.Fatalf("expected %d fields, got %d", len(expected), len(actual))
	}

	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			t.Errorf("expected field %q to be present in output, but it was not found", key)
			continue
		}

		if expectedValue.Equals(actualValue) != cty.True {
			t.Errorf("expected field %q to have value %v, got %v", key, expectedValue.GoString(), actualValue.GoString())
		}
	}
}

func verifyFailedSpecConversion(t *testing.T, spec *Spec, input map[string]cty.Value, expectedError string) {

	_, err := spec.Convert(input)
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got no error", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}

	err = spec.Validate(input)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got no error", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func verifyFailedSpecValidation(t *testing.T, spec *Spec, input map[string]cty.Value, expectedError string) {

	_, err := spec.Convert(input)
	if err != nil {
		t.Errorf("expected no error from Convert(), got error: %v", err)
	}

	err = spec.Validate(input)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got no error", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
