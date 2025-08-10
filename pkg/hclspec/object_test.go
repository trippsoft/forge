package hclspec

import (
	"fmt"
	"slices"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestObjectField(t *testing.T) {

	tests := []struct {
		name             string
		t                Type
		required         bool
		defaultValue     cty.Value
		aliases          []string
		expectedRequired bool
		expectedType     Type
	}{
		{
			name:             "required string field with no default",
			t:                String,
			required:         true,
			defaultValue:     cty.NullVal(cty.String),
			aliases:          nil,
			expectedRequired: true,
			expectedType:     String,
		},
		{
			name:             "optional number field with default",
			t:                Number,
			required:         false,
			defaultValue:     cty.NumberIntVal(42),
			aliases:          []string{"num", "number"},
			expectedRequired: false,
			expectedType:     Number,
		},
		{
			name:             "bool field with aliases",
			t:                Bool,
			required:         false,
			defaultValue:     cty.BoolVal(false),
			aliases:          []string{"enabled", "active"},
			expectedRequired: false,
			expectedType:     Bool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := &ObjectField{
				Type:         tt.t,
				Required:     tt.required,
				DefaultValue: tt.defaultValue,
				Aliases:      tt.aliases,
			}

			if field.Required != tt.expectedRequired {
				t.Errorf("expected required to be \"%v\", got \"%v\"", tt.expectedRequired, field.Required)
			}

			if field.Type != tt.expectedType {
				t.Errorf("expected type to be %q, got %q", tt.expectedType, field.Type)
			}

			if !field.DefaultValue.Equals(tt.defaultValue).True() {
				t.Errorf(
					"expected default value to be %q, got %q",
					tt.defaultValue.GoString(),
					field.DefaultValue.GoString())
			}

			if len(field.Aliases) != len(tt.aliases) {
				t.Errorf("expected %d aliases, got %d", len(tt.aliases), len(field.Aliases))
			}

			for _, alias := range tt.aliases {
				if !slices.Contains(field.Aliases, alias) {
					t.Errorf("expected alias %q to be present, but it was not", alias)
				}
			}
		})
	}
}

func TestObject(t *testing.T) {

	fields := map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
		"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(0)},
	}

	constraint := MutuallyExclusive("name", "age")

	obj := Object(fields, constraint)

	if len(obj.fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(obj.fields))
	}

	if len(obj.constraints) != 1 {
		t.Errorf("expected 1 constraint, got %d", len(obj.constraints))
	}

	if _, exists := obj.fields["name"]; !exists {
		t.Errorf("expected field %q to be present, but it was not", "name")
	}

	if _, exists := obj.fields["age"]; !exists {
		t.Errorf("expected field %q to be present, but it was not", "age")
	}
}

func TestObjectFieldValidateSpec_Pass(t *testing.T) {

	tests := []struct {
		name      string
		field     *ObjectField
		fieldName string
	}{
		{
			name: "valid required field",
			field: &ObjectField{
				Type:         String,
				Required:     true,
				DefaultValue: cty.NullVal(cty.String),
			},
			fieldName: "test",
		},
		{
			name: "valid optional field with default",
			field: &ObjectField{
				Type:         Number,
				Required:     false,
				DefaultValue: cty.NumberIntVal(42),
			},
			fieldName: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			errs := tt.field.validateSpec(tt.fieldName)

			if len(errs) != 0 {
				t.Errorf("expected no errors from validateSpec(), got %d", len(errs))

				for _, err := range errs {
					t.Errorf("expected no errors from validateSpec(), got %v", err)
				}
			}
		})
	}
}

func TestObjectFieldValidateSpec_UnknownDefaultValue(t *testing.T) {

	field := &ObjectField{
		Type:         String,
		Required:     false,
		DefaultValue: cty.UnknownVal(cty.String),
	}

	errs := field.validateSpec("test")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from validateSpec(), got %d", len(errs))
	}

	expectedError := `field "test" has an unknown default value`
	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func TestObjectFieldValidateSpec_RequiredWithDefault(t *testing.T) {

	field := &ObjectField{
		Type:         String,
		Required:     true,
		DefaultValue: cty.StringVal("default"),
	}

	errs := field.validateSpec("test")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from validateSpec(), got %d", len(errs))
	}

	expectedError := `field "test" is required and has a default value`
	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func TestObjectFieldValidateSpec_InvalidDefaultValue(t *testing.T) {

	field := &ObjectField{
		Type:         Number,
		Required:     false,
		DefaultValue: cty.StringVal("not-a-number"),
	}

	errs := field.validateSpec("test")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from validateSpec(), got %d", len(errs))
	}

	expectedError := `field "test" default value validation failed: cannot convert "string" to "number": a number is required`
	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func TestObjectTypeCtyType(t *testing.T) {
	obj := Object(map[string]*ObjectField{
		"name":   {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
		"age":    {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(0)},
		"active": {Type: Bool, Required: false, DefaultValue: cty.BoolVal(false)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
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
				"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
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
				"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
				"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname", "title"}},
				"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(25)},
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
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
	})
	expectedError := "cannot convert unknown value"
	verifyFailedConversion(t, obj, cty.UnknownVal(obj.CtyType()), expectedError)
}

func TestObjectType_InvalidType(t *testing.T) {
	obj := Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
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
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
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
		"age": {Type: Number, Required: true, DefaultValue: cty.NullVal(cty.Number)},
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
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"fullname"}},
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
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
		"age":  {Type: Number, Required: false, DefaultValue: cty.NumberIntVal(0)},
	})

	input := cty.ObjectVal(map[string]cty.Value{
		"age": cty.NumberIntVal(30),
	})

	expectedError := fmt.Sprintf("missing required field %q", "name")
	verifyFailedValidation(t, obj, input, expectedError)
}

func TestObjectType_FieldValidationFailure(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"duration": {Type: Duration, Required: true, DefaultValue: cty.NullVal(cty.String)},
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
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
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
		"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
		"field2": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
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
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String)},
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
			obj := Object(tt.fields, tt.constraints...)
			errs := obj.ValidateSpec()
			if len(errs) != 0 {
				t.Errorf("expected no errors from ValidateSpec(), got %d errors", len(errs))
				for _, err := range errs {
					t.Errorf("expected no errors from ValidateSpec(), got error: %v", err)
				}
			}
		})
	}
}

func TestObjectTypeValidateSpec_FieldErrors(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"invalid": {
			Type:         String,
			Required:     true,
			DefaultValue: cty.StringVal("default"), // Required field with default
		},
	})

	errs := obj.ValidateSpec()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from ValidateSpec(), got %d errors", len(errs))
	}

	expectedError := fmt.Sprintf("field %q is required and has a default value", "invalid")
	if errs[0].Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, errs[0].Error())
	}
}

func TestObjectTypeValidateSpec_DuplicateFieldNames(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"name": {Type: String, Required: true, DefaultValue: cty.NullVal(cty.String), Aliases: []string{"name"}},
	})

	errs := obj.ValidateSpec()
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

func TestObjectTypeValidateSpec_InvalidConstraint(t *testing.T) {

	obj := Object(map[string]*ObjectField{
		"field1": {Type: String, Required: false, DefaultValue: cty.NullVal(cty.String)},
	}, MutuallyExclusive("field1", "nonexistent"))

	errs := obj.ValidateSpec()
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
