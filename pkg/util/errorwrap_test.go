// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"errors"
	"slices"
	"testing"
)

func TestUnwrapErrors_Nil(t *testing.T) {
	var err error
	result := UnwrapErrors(err)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestUnwrapErrors_SingleError(t *testing.T) {
	err := errors.New("single error")
	result := UnwrapErrors(err)
	if len(result) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result))
	}

	if result[0] != err {
		t.Errorf("expected %v, got %v", err, result[0])
	}
}

func TestUnwrapErrors_MultipleErrors(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	err := errors.Join(err1, err2, err3)
	result := UnwrapErrors(err)
	if len(result) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(result))
	}

	if !slices.Contains(result, err1) {
		t.Errorf("expected %v to be in result", err1)
	}

	if !slices.Contains(result, err2) {
		t.Errorf("expected %v to be in result", err2)
	}

	if !slices.Contains(result, err3) {
		t.Errorf("expected %v to be in result", err3)
	}
}

func TestUnwrapErrors_NestedErrorJoins(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")
	err4 := errors.New("error 4")

	err := errors.Join(err1, errors.Join(err2, errors.Join(err3, err4)))
	result := UnwrapErrors(err)
	if len(result) != 4 {
		t.Fatalf("expected 4 errors, got %d", len(result))
	}

	if !slices.Contains(result, err1) {
		t.Errorf("expected %v to be in result", err1)
	}

	if !slices.Contains(result, err2) {
		t.Errorf("expected %v to be in result", err2)
	}

	if !slices.Contains(result, err3) {
		t.Errorf("expected %v to be in result", err3)
	}

	if !slices.Contains(result, err4) {
		t.Errorf("expected %v to be in result", err4)
	}
}
