// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"strings"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func createMockObjectType() *objectType {
	return Object(map[string]*ObjectField{
		"field1": RequiredField(String),
		"field2": OptionalField(String, cty.NullVal(cty.String)),
		"field3": OptionalField(Number, cty.NullVal(cty.Number)),
		"field4": OptionalField(Bool, cty.NullVal(cty.Bool)),
	})
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constraints.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %v", err)
			}
		})
	}
}

func TestObjectConstraintsValidate_ErrorPropagation(t *testing.T) {
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
			},
			expectedError: "mutually exclusive fields",
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
			expectedError: "at least one of the fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.constraints.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error from Validate() to end with %q, got none", tt.expectedError)
			}

			if strings.HasSuffix(err.Error(), tt.expectedError) {
				t.Errorf("expected error from Validate() to end with %q, got %q", tt.expectedError, err.Error())
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
				t.Fatalf("expected no error, got %v", err)
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
			expectedError: `mutually exclusive fields ["field1" "field2"] are all present`,
		},
		{
			name:   "three fields present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			},
			expectedError: `mutually exclusive fields ["field1" "field2" "field3"] are all present`,
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

func TestMutuallyExclusiveValidateSpec_Pass(t *testing.T) {

	mockType := createMockObjectType()

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
			err := constraint.ValidateSpec(mockType)

			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %v", err)
			}
		})
	}
}

func TestMutuallyExclusiveValidateSpec_FieldNotDefined(t *testing.T) {

	mockType := createMockObjectType()

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
			fields:        []string{"nonexistent1", "nonexistent2"},
			expectedError: `field "nonexistent1" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			constraint := MutuallyExclusive(tt.fields...)
			err := constraint.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
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
				t.Fatalf("expected no error from Validate(), got %v", err)
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
			expectedError: `fields "field1, field2" are required together, but only "field1" is/are present`,
		},
		{
			name:   "only last field present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field3": cty.StringVal("value3"),
			},
			expectedError: `fields "field1, field2, field3" are required together, but only "field3" is/are present`,
		},
		{
			name:   "some fields present",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field3": cty.StringVal("value3"),
			},
			expectedError: `fields "field1, field2, field3" are required together, but only "field1, field3" is/are present`,
		},
		{
			name:   "some fields present with null field ignored",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.NullVal(cty.String),
			},
			expectedError: `fields "field1, field2, field3" are required together, but only "field1, field2" is/are present`,
		},
		{
			name:   "some fields present with unknown field ignored",
			fields: []string{"field1", "field2", "field3"},
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.UnknownVal(cty.String),
			},
			expectedError: `fields "field1, field2, field3" are required together, but only "field1, field2" is/are present`,
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

func TestRequiredTogetherValidateSpec_Pass(t *testing.T) {
	mockType := createMockObjectType()

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
			err := constraint.ValidateSpec(mockType)

			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %v", err)
			}
		})
	}
}

func TestRequiredTogetherValidateSpec_FieldNotDefined(t *testing.T) {

	mockType := createMockObjectType()

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
			fields:        []string{"nonexistent1", "nonexistent2"},
			expectedError: `field "nonexistent1" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			constraint := RequiredTogether(tt.fields...)
			err := constraint.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
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
				t.Fatalf("expected no error from Validate(), got %v", err)
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
			expectedError: `at least one of the fields ["field1" "field2"] is required`,
		},
		{
			name:   "only null values",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.NullVal(cty.String),
			},
			expectedError: `at least one of the fields ["field1" "field2"] is required`,
		},
		{
			name:   "only unknown values",
			fields: []string{"field1", "field2"},
			values: map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
				"field2": cty.UnknownVal(cty.String),
			},
			expectedError: `at least one of the fields ["field1" "field2"] is required`,
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

func TestRequiredOneOfValidateSpec_Pass(t *testing.T) {

	mockType := createMockObjectType()

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
			err := constraint.ValidateSpec(mockType)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestRequiredOneOfValidateSpec_FieldNotDefined(t *testing.T) {

	mockType := createMockObjectType()

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
			fields:        []string{"nonexistent1", "nonexistent2"},
			expectedError: `field "nonexistent1" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			constraint := RequiredOneOf(tt.fields...)
			err := constraint.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestFieldIsNotNullIsMet(t *testing.T) {

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

			condition := FieldIsNotNull(tt.field)
			result := condition.IsMet(tt.values)

			if result != tt.expected {
				t.Errorf("expected IsMet() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFieldIsNotNullValidateSpec_Pass(t *testing.T) {

	mockType := createMockObjectType()

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

			condition := FieldIsNotNull(tt.field)
			err := condition.ValidateSpec(mockType)

			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %v", err)
			}
		})
	}
}

func TestFieldIsNotNullValidateSpec_FieldNotDefined(t *testing.T) {

	mockType := createMockObjectType()

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

			condition := FieldIsNotNull(tt.field)
			err := condition.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestFieldIsNullIsMet(t *testing.T) {

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

			condition := FieldIsNull(tt.field)
			result := condition.IsMet(tt.values)

			if result != tt.expected {
				t.Errorf("expected IsMet() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFieldIsNullValidateSpec_Pass(t *testing.T) {

	mockType := createMockObjectType()

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

			condition := FieldIsNull(tt.field)
			err := condition.ValidateSpec(mockType)

			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %v", err)
			}
		})
	}
}

func TestFieldIsNullValidateSpec_FieldNotDefined(t *testing.T) {

	mockType := createMockObjectType()

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

			condition := FieldIsNull(tt.field)
			err := condition.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
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
			result := condition.IsMet(tt.values)

			if result != tt.expected {
				t.Errorf("expected IsMet() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFieldEqualsValidateSpec_Pass(t *testing.T) {

	mockType := createMockObjectType()

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
			err := condition.ValidateSpec(mockType)

			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %v", err)
			}
		})
	}
}

func TestFieldEqualsValidateSpec_FieldNotDefined(t *testing.T) {

	mockType := createMockObjectType()

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
			err := condition.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error %q from ValidateSpec(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from ValidateSpec(), got %q", tt.expectedError, err.Error())
			}
		})
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
			condition:  FieldIsNotNull("field1"),
			constraint: RequiredOneOf("field2", "field3"),
			values: map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			},
		},
		{
			name:       "condition met, constraint passes",
			condition:  FieldIsNotNull("field1"),
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
				t.Fatalf("expected no error from Validate(), got %v", err)
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
			condition:  FieldIsNotNull("field1"),
			constraint: RequiredOneOf("field2", "field3"),
			values: map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			},
			expectedError: "at least one of the fields",
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
			expectedError: "mutually exclusive fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			constraint := ConditionalConstraint(tt.condition, tt.constraint)
			err := constraint.Validate(tt.values)

			if err == nil {
				t.Fatalf("expected error from Validate() to end with %q, got none", tt.expectedError)
			}

			if strings.HasSuffix(err.Error(), tt.expectedError) {
				t.Errorf("expected error from Validate() to end with %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestConditionalConstraintValidateSpec_Pass(t *testing.T) {

	mockType := createMockObjectType()

	tests := []struct {
		name       string
		condition  ObjectCondition
		constraint ObjectConstraint
	}{
		{
			name:       "valid condition and constraint",
			condition:  FieldIsNotNull("field1"),
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
			err := constraint.ValidateSpec(mockType)

			if err != nil {
				t.Fatalf("expected no error from ValidateSpec(), got %v", err)
			}
		})
	}
}

func TestConditionalConstraintValidateSpec_InvalidCondition(t *testing.T) {

	mockType := createMockObjectType()

	tests := []struct {
		name          string
		condition     ObjectCondition
		constraint    ObjectConstraint
		expectedError string
	}{
		{
			name:          "invalid condition field",
			condition:     FieldIsNotNull("nonexistent"),
			constraint:    RequiredOneOf("field1", "field2"),
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			constraint := ConditionalConstraint(tt.condition, tt.constraint)
			err := constraint.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error from ValidateSpec() to end with %q, got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error from ValidateSpec() to end with %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestConditionalConstraintValidateSpec_InvalidConstraint(t *testing.T) {

	mockType := createMockObjectType()

	tests := []struct {
		name          string
		condition     ObjectCondition
		constraint    ObjectConstraint
		expectedError string
	}{
		{
			name:          "invalid constraint field",
			condition:     FieldIsNotNull("field1"),
			constraint:    RequiredOneOf("nonexistent", "field2"),
			expectedError: `field "nonexistent" is not defined in the object type`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			constraint := ConditionalConstraint(tt.condition, tt.constraint)
			err := constraint.ValidateSpec(mockType)

			if err == nil {
				t.Fatalf("expected error from ValidateSpec() to end with %q, got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error from ValidateSpec() to end with %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}
