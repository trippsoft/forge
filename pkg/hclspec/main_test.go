// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/trippsoft/forge/pkg/errorwrap"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

func verifySuccessfulConversion(t *testing.T, ty Type, input, expected cty.Value) {
	t.Helper()

	actual, err := ty.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got %v", err)
	}

	if !actual.Type().Equals(ty.CtyType()) {
		t.Errorf(
			"expected Convert() to produce type %q, got %q",
			ty.CtyType().FriendlyName(),
			actual.Type().FriendlyName())
	}

	if actual.Equals(expected) != cty.True {
		t.Errorf(
			"expected Convert() to produce value %s, got %s",
			util.FormatCtyValueToString(expected, 0, 0),
			util.FormatCtyValueToString(actual, 0, 0))
	}
}

func verifyFailedConversion(t *testing.T, ty Type, input cty.Value, expectedError string) {
	t.Helper()

	_, err := ty.Convert(input)
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	errs := errorwrap.UnwrapErrors(err)
	for _, e := range errs {
		if e.Error() == expectedError {
			return
		}
	}

	t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
}

func verifySuccessfulValidation(t *testing.T, ty Type, input cty.Value) {
	t.Helper()

	converted, err := ty.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got %v", err)
	}

	err = ty.Validate(converted)
	if err != nil {
		t.Fatalf("expected no error from Validate(), got %v", err)
	}
}

func verifyFailedValidation(t *testing.T, ty Type, input cty.Value, expectedError string) {
	t.Helper()

	converted, err := ty.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got %v", err)
	}

	err = ty.Validate(converted)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	errs := errorwrap.UnwrapErrors(err)
	for _, e := range errs {
		if e.Error() == expectedError {
			return
		}
	}

	t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
}
