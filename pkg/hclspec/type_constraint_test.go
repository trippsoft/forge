// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestFieldConstraintsValidate_Pass(t *testing.T) {
	tests := []struct {
		name        string
		constraints TypeConstraints
		value       cty.Value
	}{
		{
			name:        "nil",
			constraints: nil,
			value:       cty.NilVal,
		},
		{
			name: "valid value",
			constraints: TypeConstraints{
				AllowedValues(cty.StringVal("value1"), cty.StringVal("value2")),
			},
			value: cty.StringVal("value1"),
		},
		{
			name:        "empty constraints",
			constraints: TypeConstraints{},
			value:       cty.NilVal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constraints.Validate(tt.value)
			if err != nil {
				t.Errorf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestFieldConstraintsValidate_Empty(t *testing.T) {
	constraints := TypeConstraints{}

	err := constraints.Validate(cty.NilVal)
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestFieldConstraintsValidate_Fail(t *testing.T) {
	constraints := TypeConstraints{
		AllowedValues(cty.StringVal("value1"), cty.StringVal("value2")),
	}

	expectedError := `value "value3" is not in allowed values: "value1", "value2"`
	err := constraints.Validate(cty.StringVal("value3"))
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestFieldConstraintsValidate_Nil(t *testing.T) {
	var constraints TypeConstraints

	err := constraints.Validate(cty.NilVal)
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}
