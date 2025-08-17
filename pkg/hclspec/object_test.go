// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"testing"

	"github.com/trippsoft/forge/pkg/errorwrap"
	"github.com/zclconf/go-cty/cty"
)

func TestObjectFieldValidateSpec_Pass(t *testing.T) {

	tests := []struct {
		name      string
		field     *ObjectField
		fieldName string
	}{
		{
			name:      "valid required field",
			field:     RequiredField(String),
			fieldName: "test",
		},
		{
			name:      "valid optional field with default",
			field:     OptionalField(Number, cty.NumberIntVal(42)),
			fieldName: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.field.validateSpec(tt.fieldName)
			if err != nil {
				t.Errorf("expected no errors from validateSpec(), got: %v", err)
			}
		})
	}
}

func TestObjectFieldValidateSpec_UnknownDefaultValue(t *testing.T) {

	field := OptionalField(String, cty.UnknownVal(cty.String))

	err := field.validateSpec("test")
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	expectedError := `field "test" has an unknown default value`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_RequiredWithDefault(t *testing.T) {

	field := RequiredField(String)

	field.defaultValue = cty.StringVal("default")

	err := field.validateSpec("test")
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	errs := errorwrap.UnwrapErrors(err)

	expectedError := `field "test" is required and has a default value`
	for _, e := range errs {
		if e.Error() == expectedError {
			return
		}

		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_InvalidDefaultValue(t *testing.T) {

	field := OptionalField(Number, cty.StringVal("not-a-number"))

	err := field.validateSpec("test")
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	errs := errorwrap.UnwrapErrors(err)

	expectedError := `field "test" default value type mismatch: expected "number", got "string"`
	for _, e := range errs {
		if e.Error() == expectedError {
			return
		}

		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeCtyType(t *testing.T) {
	obj := Object(map[string]*ObjectField{
		"name":   RequiredField(String),
		"age":    OptionalField(Number, cty.NumberIntVal(0)),
		"active": OptionalField(Bool, cty.BoolVal(false)),
	})

	ctyType := obj.CtyType()

	if !ctyType.IsObjectType() {
		t.Fatalf("expected object type, got %q", ctyType.FriendlyName())
	}

	attrTypes := ctyType.AttributeTypes()
	if len(attrTypes) != 3 {
		t.Errorf("expected 3 attribute types, got %d", len(attrTypes))
	}

	nameAttr, exists := attrTypes["name"]
	if !exists {
		t.Errorf("expected attribute %q, got none", "name")
	} else if nameAttr != cty.String {
		t.Errorf(
			"expected attribute %q to be %q, got %q",
			"name",
			cty.String.FriendlyName(),
			nameAttr.FriendlyName())
	}

	ageAttr, exists := attrTypes["age"]
	if !exists {
		t.Errorf("expected attribute %q, got none", "age")
	} else if ageAttr != cty.Number {
		t.Errorf(
			"expected attribute %q to be %q, got %q",
			"age",
			cty.Number.FriendlyName(),
			ageAttr.FriendlyName())
	}

	activeAttr, exists := attrTypes["active"]
	if !exists {
		t.Errorf("expected attribute %q, got none", "active")
	} else if activeAttr != cty.Bool {
		t.Errorf(
			"expected attribute %q to be %q, got %q",
			"active",
			cty.Bool.FriendlyName(),
			activeAttr.FriendlyName())
	}
}

func TestObjectType(t *testing.T) {

	tests := []struct {
		name     string
		object   Type
		input    cty.Value
		expected cty.Value
	}{
		{
			name: "valid object",
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String, "fullname", "title"),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String, "fullname", "title"),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"field1": OptionalField(String, cty.NullVal(cty.String)),
				"field2": OptionalField(String, cty.NullVal(cty.String)),
			}, MutuallyExclusive("field1", "field2")),
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
			object: Object(map[string]*ObjectField{
				"field1": OptionalField(String, cty.NullVal(cty.String)),
				"field2": OptionalField(String, cty.NullVal(cty.String)),
			}, MutuallyExclusive("field1", "field2")),
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
			object: Object(map[string]*ObjectField{
				"field1": OptionalField(String, cty.NullVal(cty.String)),
				"field2": OptionalField(String, cty.NullVal(cty.String)),
			}, MutuallyExclusive("field1", "field2")),
			input: cty.ObjectVal(map[string]cty.Value{}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.NullVal(cty.String),
			}),
		},
		{
			name: "valid map with string values (age will be converted)",
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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
			object: Object(map[string]*ObjectField{
				"name": OptionalField(String, cty.NullVal(cty.String)),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
			input: cty.MapValEmpty(cty.String),
			expected: cty.ObjectVal(map[string]cty.Value{
				"name": cty.NullVal(cty.String),
				"age":  cty.NumberIntVal(25),
			}),
		},
		{
			name: "null object",
			object: Object(map[string]*ObjectField{
				"name": RequiredField(String, "fullname", "title"),
				"age":  OptionalField(Number, cty.NumberIntVal(25)),
			}),
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

func TestObjectType_UnknownValue(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"name": RequiredField(String),
	})
	expectedError := "cannot convert unknown value"
	verifyFailedConversion(t, obj, cty.UnknownVal(obj.CtyType()), expectedError)
}

func TestObjectType_InvalidType(t *testing.T) {
	obj := Object(map[string]*ObjectField{
		"name": RequiredField(String),
	})

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
				obj.CtyType().FriendlyName())
			verifyFailedConversion(t, obj, tt.input, expectedError)
		})
	}
}

func TestObjectType_InvalidAttributes(t *testing.T) {
	obj := Object(map[string]*ObjectField{
		"name": RequiredField(String),
	})

	input := cty.ObjectVal(map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	})

	expectedError := fmt.Sprintf(
		"invalid indexes found: %v",
		[]string{"invalid"},
	)

	verifyFailedConversion(t, obj, input, expectedError)
}

func TestObjectType_FieldConversionError(t *testing.T) {
	obj := Object(map[string]*ObjectField{
		"age": RequiredField(Number),
	})

	input := cty.ObjectVal(map[string]cty.Value{
		"age": cty.StringVal("not-a-number"),
	})

	expectedError := fmt.Sprintf(
		"cannot convert field %q: cannot convert %q to %q: a number is required",
		"age",
		cty.String.FriendlyName(),
		cty.Number.FriendlyName())

	verifyFailedConversion(t, obj, input, expectedError)
}

func TestObjectType_MultipleAliasesDefined(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"name": RequiredField(String, "fullname"),
	})

	input := cty.ObjectVal(map[string]cty.Value{
		"name":     cty.StringVal("John"),
		"fullname": cty.StringVal("John Doe"),
	})

	expectedError := fmt.Sprintf("field %q is defined multiple times as %v", "name", []string{"name", "fullname"})
	verifyFailedConversion(t, obj, input, expectedError)
}

func TestObjectType_MissingRequiredField(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"name": RequiredField(String),
		"age":  OptionalField(Number, cty.NumberIntVal(0)),
	})

	input := cty.ObjectVal(map[string]cty.Value{
		"age": cty.NumberIntVal(30),
	})

	expectedError := fmt.Sprintf("missing required field %q", "name")
	verifyFailedValidation(t, obj, input, expectedError)
}

func TestObjectType_FieldValidationFailure(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"duration": RequiredField(Duration),
	})

	input := cty.ObjectVal(map[string]cty.Value{
		"duration": cty.StringVal("invalid-duration"),
	})

	expectedError := fmt.Sprintf(
		"field %q validation failed: time: invalid duration %q",
		"duration",
		"invalid-duration")

	verifyFailedValidation(t, obj, input, expectedError)
}

func TestObjectType_ConversionError(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"name": RequiredField(String),
	})

	input := cty.StringVal("not-an-object")

	expectedError := fmt.Sprintf(
		"cannot convert %q to %q",
		input.Type().FriendlyName(),
		obj.CtyType().FriendlyName())

	verifyFailedConversion(t, obj, input, expectedError)
}

func TestObjectType_ConstraintFailure(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"field1": OptionalField(String, cty.NullVal(cty.String)),
		"field2": OptionalField(String, cty.NullVal(cty.String)),
	}, MutuallyExclusive("field1", "field2"))

	input := cty.ObjectVal(map[string]cty.Value{
		"field1": cty.StringVal("value1"),
		"field2": cty.StringVal("value2"),
	})

	expectedError := fmt.Sprintf(
		"validation failed: mutually exclusive fields %q are all present",
		[]string{"field1", "field2"})
	verifyFailedValidation(t, obj, input, expectedError)
}

func TestObjectType_InvalidIndexesInMap(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"name": {t: String, required: true, defaultValue: cty.NullVal(cty.String)},
	})

	input := cty.MapVal(map[string]cty.Value{
		"name":    cty.StringVal("John"),
		"invalid": cty.StringVal("not allowed"),
	})

	expectedError := fmt.Sprintf("invalid indexes found: %v", []string{"invalid"})
	verifyFailedConversion(t, obj, input, expectedError)
}

func TestObjectTypeValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name        string
		fields      map[string]*ObjectField
		constraints []ObjectConstraint
	}{
		{
			name: "valid spec",
			fields: map[string]*ObjectField{
				"name": RequiredField(String),
				"age":  OptionalField(Number, cty.NumberIntVal(0)),
			},
			constraints: nil,
		},
		{
			name: "valid spec with constraints",
			fields: map[string]*ObjectField{
				"field1": OptionalField(String, cty.NullVal(cty.String)),
				"field2": OptionalField(String, cty.NullVal(cty.String)),
			},
			constraints: []ObjectConstraint{MutuallyExclusive("field1", "field2")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := Object(tt.fields, tt.constraints...)
			err := obj.ValidateSpec()
			if err != nil {
				t.Errorf("expected no errors from ValidateSpec(), got error: %v", err)
			}
		})
	}
}

func TestObjectTypeValidateSpec_FieldErrors(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"invalid": {
			t:            String,
			required:     true,
			defaultValue: cty.StringVal("default"), // Required field with default
		},
	})

	err := obj.ValidateSpec()
	if err == nil {
		t.Fatal("expected error from ValidateSpec(), got none")
	}

	expectedError := fmt.Sprintf("field %q is required and has a default value", "invalid")
	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_DuplicateFieldNames(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"name": RequiredField(String, "name"),
	})

	err := obj.ValidateSpec()
	if err == nil {
		t.Fatal("expected error from ValidateSpec(), got none")
	}

	expectedError := fmt.Sprintf(
		"field %q is defined multiple times (aliases: %v)",
		"name",
		[]string{"name", "name"})

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectTypeValidateSpec_InvalidConstraint(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"field1": OptionalField(String, cty.NullVal(cty.String)),
	}, MutuallyExclusive("field1", "nonexistent"))

	err := obj.ValidateSpec()
	if err == nil {
		t.Fatal("expected error from ValidateSpec(), got none")
	}

	expectedError := fmt.Sprintf(
		"constraint validation failed: field %q is not defined in the object type",
		"nonexistent")

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
