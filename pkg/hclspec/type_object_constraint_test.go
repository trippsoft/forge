// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

var (
	mockObjectType = Object(
		RequiredField("field1", String),
		OptionalField("field2", String),
		OptionalField("field3", Number),
		OptionalField("field4", Bool),
	)
)

func TestObjectConstraintsValidate_Pass(t *testing.T) {
	tests := []struct {
		name        string
		constraints ObjectConstraints
		values      map[string]cty.Value
	}{
		{
			name:        "empty constraints",
			constraints: ObjectConstraints{},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
		},
		{
			name: "all constraints pass",
			constraints: ObjectConstraints{
				MutuallyExclusive("field1", "field2"),
				RequiredOneOf("field1", "field3"),
			},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
		},
		{
			name: "all constraints pass with nil constraint",
			constraints: ObjectConstraints{
				MutuallyExclusive("field1", "field2"),
				RequiredOneOf("field1", "field3"),
				nil, // nil constraint should be ignored
			},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constraints.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestObjectConstraintsValidate_Fail(t *testing.T) {
	tests := []struct {
		name          string
		constraints   ObjectConstraints
		values        map[string]cty.Value
		expectedError string
	}{
		{
			name: "first constraint fails",
			constraints: ObjectConstraints{
				MutuallyExclusive("field1", "field2"),
				RequiredOneOf("field3", "field4"),
			},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			},
			expectedError: `validation failed: mutually exclusive fields "field1", "field2" are all present`,
		},
		{
			name: "second constraint fails",
			constraints: ObjectConstraints{
				RequiredOneOf("field1", "field2"),
				RequiredOneOf("field3", "field4"),
			},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
			expectedError: `validation failed: at least one of the fields "field3", "field4" is required`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constraints.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestObjectConstraintsValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name        string
		constraints ObjectConstraints
	}{
		{
			name:        "empty constraints",
			constraints: ObjectConstraints{},
		},
		{
			name: "all constraints pass",
			constraints: ObjectConstraints{
				MutuallyExclusive("field1", "field2"),
				RequiredOneOf("field1", "field3"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			object := mockObjectType
			err := tt.constraints.ValidateSpec(object)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %v", err)
			}
		})
	}
}

func TestObjectConstraintsValidateSpec_Fail(t *testing.T) {
	tests := []struct {
		name          string
		constraints   ObjectConstraints
		expectedError string
	}{
		{
			name:          "contains nil constraint",
			constraints:   ObjectConstraints{nil},
			expectedError: "object type has a nil constraint",
		},
		{
			name: "first constraint fails",
			constraints: ObjectConstraints{
				MutuallyExclusive("nonexistantfield1", "field2"),
				RequiredOneOf("field3", "field4"),
			},
			expectedError: `constraint validation failed: field "nonexistantfield1" is not defined in the object type`,
		},
		{
			name: "second constraint fails",
			constraints: ObjectConstraints{
				RequiredOneOf("field1", "field2"),
				RequiredOneOf("field3", "nonexistantfield4"),
			},
			expectedError: `constraint validation failed: field "nonexistantfield4" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constraints.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error from Validate() to be %q, got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error from Validate() to be %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestMutuallyExclusiveValidate_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		values map[string]cty.Value
	}{
		{
			name:   "no fields present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{},
		},
		{
			name:   "one field present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
		},
		{
			name:   "one field present with null value",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.NullVal(cty.String),
			},
		},
		{
			name:   "one field present with unknown value",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.UnknownVal(cty.String),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := MutuallyExclusive(tt.fields...)
			err := constraint.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestMutuallyExclusiveValidate_MultipleFieldsPresent(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		values        map[string]cty.Value
		expectedError string
	}{
		{
			name:   "two fields present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			},
			expectedError: `mutually exclusive fields "field1", "field2" are all present`,
		},
		{
			name:   "three fields present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			},
			expectedError: `mutually exclusive fields "field1", "field2", "field3" are all present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := MutuallyExclusive(tt.fields...)

			err := constraint.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestMutuallyExclusiveValidate_Nil(t *testing.T) {
	var constraint *mutuallyExclusiveConstraint

	err := constraint.Validate(map[string]cty.Value{})
	expectedError := "mutually exclusive constraint is nil"
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestMutuallyExclusiveValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
	}{
		{
			name:   "all fields exist",
			fields: []string{"field1", "field2"},
		},
		{
			name:   "single field exists",
			fields: []string{"field1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := MutuallyExclusive(tt.fields...)

			err := constraint.ValidateSpec(mockObjectType)
			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %q", err.Error())
			}
		})
	}
}

func TestMutuallyExclusiveValidateSpec_FieldNotDefined(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		expectedError string
	}{
		{
			name:          "nonexistent field",
			fields:        []string{"field1", "nonexistent"},
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
		{
			name:          "multiple nonexistent fields",
			fields:        []string{"nonexistent1", "field2"},
			expectedError: `field "nonexistent1" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := MutuallyExclusive(tt.fields...)

			err := constraint.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestMutuallyExclusiveValidateSpec_Nil(t *testing.T) {
	var constraint *mutuallyExclusiveConstraint

	expectedError := "mutually exclusive constraint is nil"
	err := constraint.ValidateSpec(mockObjectType)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestRequiredTogetherValidate_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		values map[string]cty.Value
	}{
		{
			name:   "no fields present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{},
		},
		{
			name:   "all fields present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			},
		},
		{
			name:   "all three fields present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredTogether(tt.fields...)

			err := constraint.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestRequiredTogetherValidate_PartialFieldsPresent(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		values        map[string]cty.Value
		expectedError string
	}{
		{
			name:   "only first field present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
			expectedError: `fields "field1", "field2" are required together, but only "field1" is present`,
		},
		{
			name:   "only last field present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field3": cty.StringVal("value3"),
			},
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field3" is present`,
		},
		{
			name:   "some fields present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field3": cty.StringVal("value3"),
			},
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field1", "field3" is present`,
		},
		{
			name:   "some fields present with null field ignored",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.NullVal(cty.String),
			},
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field1", "field2" is present`,
		},
		{
			name:   "some fields present with unknown field ignored",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.UnknownVal(cty.String),
			},
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field1", "field2" is present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredTogether(tt.fields...)

			err := constraint.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q	from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestRequiredTogetherValidate_Nil(t *testing.T) {
	var constraint *requiredTogetherConstraint

	expectedError := "required together constraint is nil"
	err := constraint.Validate(map[string]cty.Value{})
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestRequiredTogetherValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
	}{
		{
			name:   "all fields exist",
			fields: []string{"field1", "field2"},
		},
		{
			name:   "single field exists",
			fields: []string{"field1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredTogether(tt.fields...)

			err := constraint.ValidateSpec(mockObjectType)
			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %q", err.Error())
			}
		})
	}
}

func TestRequiredTogetherValidateSpec_FieldNotDefined(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		expectedError string
	}{
		{
			name:          "nonexistent field",
			fields:        []string{"field1", "nonexistent"},
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
		{
			name:          "multiple nonexistent fields",
			fields:        []string{"nonexistent1", "field2"},
			expectedError: `field "nonexistent1" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredTogether(tt.fields...)

			err := constraint.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestRequiredTogetherValidateSpec_Nil(t *testing.T) {
	var constraint *requiredTogetherConstraint

	expectedError := "required together constraint is nil"
	err := constraint.ValidateSpec(mockObjectType)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestRequiredOneOfValidate_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		values map[string]cty.Value
	}{
		{
			name:   "first field present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
		},
		{
			name:   "last field present",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			},
		},
		{
			name:   "multiple fields present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			},
		},
		{
			name:   "field present with null values ignored",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.NullVal(cty.String),
			},
		},
		{
			name:   "field present with unknown values ignored",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.UnknownVal(cty.String),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredOneOf(tt.fields...)

			err := constraint.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestRequiredOneOfValidate_NoFieldsPresent(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		values        map[string]cty.Value
		expectedError string
	}{
		{
			name:          "no fields in values",
			fields:        []string{"field1", "field2"},
			values:        map[string]cty.Value{},
			expectedError: `at least one of the fields "field1", "field2" is required`,
		},
		{
			name:   "only null values",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.NullVal(cty.String),
			},
			expectedError: `at least one of the fields "field1", "field2" is required`,
		},
		{
			name:   "only unknown values",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
				"field2": cty.UnknownVal(cty.String),
			},
			expectedError: `at least one of the fields "field1", "field2" is required`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredOneOf(tt.fields...)

			err := constraint.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestRequiredOneOfValidate_Nil(t *testing.T) {
	var constraint *requiredOneOfConstraint

	expectedError := "required one of constraint is nil"
	err := constraint.Validate(map[string]cty.Value{})
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestRequiredOneOfValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
	}{
		{
			name:   "all fields exist",
			fields: []string{"field1", "field2"},
		},
		{
			name:   "single field exists",
			fields: []string{"field1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredOneOf(tt.fields...)

			err := constraint.ValidateSpec(mockObjectType)
			if err != nil {
				t.Fatalf("expected no error, got %q", err.Error())
			}
		})
	}
}

func TestRequiredOneOfValidateSpec_FieldNotDefined(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		expectedError string
	}{
		{
			name:          "nonexistent field",
			fields:        []string{"field1", "nonexistent"},
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
		{
			name:          "multiple nonexistent fields",
			fields:        []string{"nonexistent1", "field2"},
			expectedError: `field "nonexistent1" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredOneOf(tt.fields...)

			err := constraint.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestRequiredOneOfValidateSpec_Nil(t *testing.T) {
	var constraint *requiredOneOfConstraint

	expectedError := "required one of constraint is nil"
	err := constraint.ValidateSpec(mockObjectType)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestAllowedFieldValuesValidate_Pass(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		allowedValues []cty.Value
		values        map[string]cty.Value
	}{
		{
			name:          "field is null",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			values:        map[string]cty.Value{"field1": cty.NullVal(cty.String)},
		},
		{
			name:          "field is first value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			values:        map[string]cty.Value{"field1": cty.StringVal("value1")},
		},
		{
			name:          "field is second value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			values:        map[string]cty.Value{"field1": cty.StringVal("value2")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := AllowedFieldValues(tt.field, tt.allowedValues...)

			err := constraint.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestAllowedFieldValuesValidate_FieldNotPresent(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		allowedValues []cty.Value
		values        map[string]cty.Value
		expectedError string
	}{
		{
			name:          "no field in values",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			values:        map[string]cty.Value{},
			expectedError: `field "field1" is not present`,
		},
		{
			name:          "not allowed value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value3"),
				"field2": cty.NullVal(cty.String),
			},
			expectedError: `field "field1" has an invalid value, allowed values are: "value1", "value2"`,
		},
		{
			name:          "unknown value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			values: map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
				"field2": cty.UnknownVal(cty.String),
			},
			expectedError: `cannot validate unknown value`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := AllowedFieldValues(tt.field, tt.allowedValues...)

			err := constraint.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestAllowedFieldValuesValidate_Nil(t *testing.T) {
	var constraint *allowedFieldValuesConstraint

	expectedError := "allowed field values constraint is nil"
	err := constraint.Validate(map[string]cty.Value{})
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestAllowedFieldValuesValidateSpec_Pass(t *testing.T) {
	constraint := AllowedFieldValues("field1", cty.StringVal("value1"), cty.StringVal("value2"))

	err := constraint.ValidateSpec(mockObjectType)
	if err != nil {
		t.Fatalf("expected no error, got %q", err.Error())
	}
}

func TestAllowedFieldValuesValidateSpec_FieldNotDefined(t *testing.T) {
	constraint := AllowedFieldValues("notdefined", cty.StringVal("value1"), cty.StringVal("value2"))

	err := constraint.ValidateSpec(mockObjectType)
	expectedError := `field "notdefined" is not defined in the object type`
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestAllowedFieldValuesValidateSpec_Nil(t *testing.T) {
	var constraint *allowedFieldValuesConstraint

	expectedError := "allowed field values constraint is nil"
	err := constraint.ValidateSpec(mockObjectType)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestFieldPresentIsMet(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		values   map[string]cty.Value
		expected bool
	}{
		{
			name:  "field is null",
			field: "field1",
			values: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			},
			expected: true,
		},
		{
			name:     "field not in values",
			field:    "field1",
			values:   map[string]cty.Value{},
			expected: false,
		},
		{
			name:  "field is not null",
			field: "field1",
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
			expected: false,
		},
		{
			name:  "field is unknown",
			field: "field1",
			values: map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldNotPresent(tt.field)

			actual := condition.IsMet(tt.values)
			if actual != tt.expected {
				t.Errorf("expected IsMet() to return %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestFieldPresentIsMet_Nil(t *testing.T) {
	var condition *fieldPresentCondition

	actual := condition.IsMet(map[string]cty.Value{})
	if actual != false {
		t.Errorf("expected IsMet() to return %t, got %t", false, actual)
	}
}

func TestFieldPresentDescription(t *testing.T) {
	tests := []struct {
		name      string
		condition ObjectCondition
		expected  string
	}{
		{
			name:      "field1",
			condition: FieldPresent("field1"),
			expected:  `field "field1" is present`,
		},
		{
			name:      "field2",
			condition: FieldPresent("field2"),
			expected:  `field "field2" is present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.condition.Description()
			if actual != tt.expected {
				t.Errorf("expected Description() to return %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestFieldPresentDescription_Nil(t *testing.T) {
	var condition *fieldPresentCondition

	actual := condition.Description()
	if actual != "" {
		t.Errorf("expected empty string from Description(), got %q", actual)
	}
}

func TestFieldPresentValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name  string
		field string
	}{
		{
			name:  "existing field",
			field: "field1",
		},
		{
			name:  "another existing field",
			field: "field2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldNotPresent(tt.field)

			err := condition.ValidateSpec(mockObjectType)
			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %q", err.Error())
			}
		})
	}
}

func TestFieldPresentValidateSpec_FieldNotDefined(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		expectedError string
	}{
		{
			name:          "nonexistent field",
			field:         "nonexistent",
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldNotPresent(tt.field)

			err := condition.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestFieldPresentValidateSpec_Nil(t *testing.T) {
	var condition *fieldPresentCondition

	expectedError := "field present condition is nil"
	err := condition.ValidateSpec(mockObjectType)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestFieldNotPresentIsMet(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		values   map[string]cty.Value
		expected bool
	}{
		{
			name:  "field present and not null",
			field: "field1",
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
			expected: true,
		},
		{
			name:     "field not in values",
			field:    "field1",
			values:   map[string]cty.Value{},
			expected: false,
		},
		{
			name:  "different field in values",
			field: "field1",
			values: map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			},
			expected: false,
		},
		{
			name:  "field is null",
			field: "field1",
			values: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			},
			expected: false,
		},
		{
			name:  "field is unknown",
			field: "field1",
			values: map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldPresent(tt.field)

			actual := condition.IsMet(tt.values)
			if actual != tt.expected {
				t.Errorf("expected IsMet() to return %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestFieldNotPresentIsMet_Nil(t *testing.T) {
	var condition *fieldNotPresentCondition

	actual := condition.IsMet(map[string]cty.Value{})
	if actual != false {
		t.Errorf("expected IsMet() to return %t, got %t", false, actual)
	}
}

func TestFieldNotPresentDescription(t *testing.T) {
	tests := []struct {
		name      string
		condition ObjectCondition
		expected  string
	}{
		{
			name:      "field1",
			condition: FieldNotPresent("field1"),
			expected:  `field "field1" is not present`,
		},
		{
			name:      "field2",
			condition: FieldNotPresent("field2"),
			expected:  `field "field2" is not present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.condition.Description()
			if actual != tt.expected {
				t.Errorf("expected Description() to return %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestFieldNotPresentDescription_Nil(t *testing.T) {
	var condition *fieldNotPresentCondition
	actual := condition.Description()
	if actual != "" {
		t.Errorf("expected empty string from Description(), got %q", actual)
	}
}

func TestFieldNotPresentValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name  string
		field string
	}{
		{
			name:  "existing field",
			field: "field1",
		},
		{
			name:  "another existing field",
			field: "field2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldPresent(tt.field)

			err := condition.ValidateSpec(mockObjectType)
			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %q", err.Error())
			}
		})
	}
}

func TestFieldNotPresentValidateSpec_FieldNotDefined(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		expectedError string
	}{
		{
			name:          "nonexistent field",
			field:         "nonexistent",
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldPresent(tt.field)

			err := condition.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestFieldNotPresentValidateSpec_Nil(t *testing.T) {
	var condition *fieldNotPresentCondition

	expectedError := "field not present condition is nil"
	err := condition.ValidateSpec(mockObjectType)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestFieldEqualsIsMet(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    cty.Value
		values   map[string]cty.Value
		expected bool
	}{
		{
			name:  "field equals string value",
			field: "field1",
			value: cty.StringVal("expected"),
			values: map[string]cty.Value{
				"field1": cty.StringVal("expected"),
			},
			expected: true,
		},
		{
			name:  "field equals number value",
			field: "field1",
			value: cty.NumberIntVal(42),
			values: map[string]cty.Value{
				"field1": cty.NumberIntVal(42),
			},
			expected: true,
		},
		{
			name:  "field equals bool value",
			field: "field1",
			value: cty.BoolVal(true),
			values: map[string]cty.Value{
				"field1": cty.BoolVal(true),
			},
			expected: true,
		},
		{
			name:     "field not in values",
			field:    "field1",
			value:    cty.StringVal("expected"),
			values:   map[string]cty.Value{},
			expected: false,
		},
		{
			name:  "field has different string value",
			field: "field1",
			value: cty.StringVal("expected"),
			values: map[string]cty.Value{
				"field1": cty.StringVal("actual"),
			},
			expected: false,
		},
		{
			name:  "field has different type",
			field: "field1",
			value: cty.StringVal("expected"),
			values: map[string]cty.Value{
				"field1": cty.NumberIntVal(42),
			},
			expected: false,
		},
		{
			name:  "field is null",
			field: "field1",
			value: cty.StringVal("expected"),
			values: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			},
			expected: false,
		},
		{
			name:  "field is unknown",
			field: "field1",
			value: cty.StringVal("expected"),
			values: map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldEquals(tt.field, tt.value)

			actual := condition.IsMet(tt.values)
			if actual != tt.expected {
				t.Errorf("expected IsMet() to return %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestFieldEqualsIsMet_Nil(t *testing.T) {
	var condition *fieldEqualsCondition

	actual := condition.IsMet(map[string]cty.Value{})
	if actual != false {
		t.Errorf("expected IsMet() to return %t, got %t", false, actual)
	}
}

func TestFieldEqualsDescription(t *testing.T) {
	tests := []struct {
		name      string
		condition ObjectCondition
		expected  string
	}{
		{
			name:      "field1",
			condition: FieldEquals("field1", cty.StringVal("value")),
			expected:  `field "field1" is equal to "value"`,
		},
		{
			name:      "field2",
			condition: FieldEquals("field2", cty.StringVal("value2")),
			expected:  `field "field2" is equal to "value2"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.condition.Description()
			if actual != tt.expected {
				t.Errorf("expected Description() to return %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestFieldEqualsDescription_Nil(t *testing.T) {
	var condition *fieldEqualsCondition

	actual := condition.Description()
	if actual != "" {
		t.Errorf("expected empty string from Description(), got %q", actual)
	}
}

func TestFieldEqualsValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name  string
		field string
		value cty.Value
	}{
		{
			name:  "existing field with string value",
			field: "field1",
			value: cty.StringVal("test"),
		},
		{
			name:  "existing field with number value",
			field: "field3",
			value: cty.NumberIntVal(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldEquals(tt.field, tt.value)

			err := condition.ValidateSpec(mockObjectType)
			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %q", err.Error())
			}
		})
	}
}

func TestFieldEqualsValidateSpec_FieldNotDefined(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		value         cty.Value
		expectedError string
	}{
		{
			name:          "nonexistent field",
			field:         "nonexistent",
			value:         cty.StringVal("test"),
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldEquals(tt.field, tt.value)

			err := condition.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestFieldEqualsValidateSpec_Nil(t *testing.T) {
	var condition *fieldEqualsCondition

	expectedError := "field equals condition is nil"
	err := condition.ValidateSpec(mockObjectType)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestConditionalConstraintValidate_Pass(t *testing.T) {
	tests := []struct {
		name       string
		condition  ObjectCondition
		constraint ObjectConstraint
		values     map[string]cty.Value
	}{
		{
			name:       "condition not met, constraint ignored",
			condition:  FieldPresent("field1"),
			constraint: RequiredOneOf("field2", "field3"),
			values: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			},
		},
		{
			name:       "condition met, constraint passes",
			condition:  FieldPresent("field1"),
			constraint: RequiredOneOf("field2", "field3"),
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := ConditionalConstraint(tt.condition, tt.constraint)

			err := constraint.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestConditionalConstraintValidate_ConditionMetConstraintFails(t *testing.T) {
	tests := []struct {
		name          string
		condition     ObjectCondition
		constraint    ObjectConstraint
		values        map[string]cty.Value
		expectedError string
	}{
		{
			name:       "condition met, constraint fails",
			condition:  FieldPresent("field1"),
			constraint: RequiredOneOf("field2", "field3"),
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
			expectedError: `conditional constraint failed: when field "field1" is present, at least one of the fields "field2", "field3" is required`,
		},
		{
			name:       "condition met, mutually exclusive constraint fails",
			condition:  FieldEquals("field1", cty.StringVal("trigger")),
			constraint: MutuallyExclusive("field2", "field3"),
			values: map[string]cty.Value{
				"field1": cty.StringVal("trigger"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			},
			expectedError: `conditional constraint failed: when field "field1" is equal to "trigger", mutually exclusive fields "field2", "field3" are all present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := ConditionalConstraint(tt.condition, tt.constraint)

			err := constraint.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestConditionalConstraintValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name       string
		condition  ObjectCondition
		constraint ObjectConstraint
	}{
		{
			name:       "valid condition and constraint",
			condition:  FieldPresent("field1"),
			constraint: RequiredOneOf("field2", "field3"),
		},
		{
			name:       "valid field equals condition",
			condition:  FieldEquals("field1", cty.StringVal("test")),
			constraint: MutuallyExclusive("field2", "field3"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := ConditionalConstraint(tt.condition, tt.constraint)

			err := constraint.ValidateSpec(mockObjectType)
			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %q", err.Error())
			}
		})
	}
}

func TestConditionalConstraintValidateSpec_InvalidCondition(t *testing.T) {
	tests := []struct {
		name          string
		condition     ObjectCondition
		constraint    ObjectConstraint
		expectedError string
	}{
		{
			name:          "invalid condition field",
			condition:     FieldPresent("nonexistent"),
			constraint:    RequiredOneOf("field1", "field2"),
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := ConditionalConstraint(tt.condition, tt.constraint)

			err := constraint.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestConditionalConstraintValidateSpec_InvalidConstraint(t *testing.T) {
	tests := []struct {
		name          string
		condition     ObjectCondition
		constraint    ObjectConstraint
		expectedError string
	}{
		{
			name:          "invalid constraint field",
			condition:     FieldPresent("field1"),
			constraint:    RequiredOneOf("nonexistent", "field2"),
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := ConditionalConstraint(tt.condition, tt.constraint)

			err := constraint.ValidateSpec(mockObjectType)
			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}
