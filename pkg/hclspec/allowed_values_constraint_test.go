// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestAllowedValuesValidate_Pass(t *testing.T) {
	constraints := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	err := constraints.Validate(cty.StringVal("value1"))
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestAllowedValuesValidate_NullValue(t *testing.T) {
	constraints := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	err := constraints.Validate(cty.NullVal(cty.String))
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestAllowedValuesValidate_Fail(t *testing.T) {
	constraints := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	expectedError := `value "value3" is not in allowed values: "value1", "value2"`
	err := constraints.Validate(cty.StringVal("value3"))
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestAllowedValuesValidate_Nil(t *testing.T) {
	var constraint *allowedValuesConstraint

	expectedError := "allowed values constraint is nil"
	err := constraint.Validate(cty.NilVal)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
