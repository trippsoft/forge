// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package result provides the types and utilities for representing the results of operations in Forge.
package result

import (
	"errors"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// Result represents the outcome of an operation, indicating whether it was a failure, an ignored failure, or a
// successful operation.
type Result struct {
	Failed         bool      // Indicates if the module execution failed.
	IgnoredFailure bool      // Indicates if the failure was ignored.
	Skipped        bool      // Indicates if the module was skipped.
	Changed        bool      // Indicates if the module made any changes.
	Error          error     // Error encountered during module execution, if any.
	ErrorDetail    string    // Detailed error message, if any.
	Output         cty.Value // Output data from the module execution.
	Warnings       []string  // Warning message, if any.
	Messages       []string  // Informational message, if any.
}

// NewNotChanged creates a new Result indicating that the operation was successful but did not result in any changes,
// with the provided output.
func NewNotChanged(output cty.Value) *Result {
	return &Result{
		Failed:         false,
		IgnoredFailure: false,
		Skipped:        false,
		Changed:        false,
		Error:          nil,
		ErrorDetail:    "",
		Output:         output,
		Warnings:       nil,
		Messages:       nil,
	}
}

// NewChanged creates a new Result indicating that the operation was successful and resulted in changes, with the
// provided output.
func NewChanged(output cty.Value) *Result {
	return &Result{
		Failed:         false,
		IgnoredFailure: false,
		Skipped:        false,
		Changed:        true,
		Error:          nil,
		ErrorDetail:    "",
		Output:         output,
		Warnings:       nil,
		Messages:       nil,
	}
}

// NewSkipped creates a new Result indicating that the operation was skipped.
func NewSkipped() *Result {
	return &Result{
		Failed:         false,
		IgnoredFailure: false,
		Skipped:        true,
		Changed:        false,
		Error:          nil,
		ErrorDetail:    "",
		Output:         cty.NilVal,
		Warnings:       nil,
		Messages:       nil,
	}
}

// NewFailure creates a new Result indicating that the operation failed, with the provided error.
func NewFailure(err error, errDetail string) *Result {
	return &Result{
		Failed:         true,
		IgnoredFailure: false,
		Skipped:        false,
		Changed:        false,
		Error:          err,
		ErrorDetail:    errDetail,
		Output:         cty.NilVal,
		Warnings:       nil,
		Messages:       nil,
	}
}

// ToResult converts the ResultPB to a Result, using the appropriate conversion based on the type of result contained in
// the protobuf.
func (r *ModuleResult) ToResult() *Result {
	switch res := r.Result.(type) {
	case *ModuleResult_Failure:
		result := NewFailure(errors.New(res.Failure.Error), res.Failure.Detail)
		result.Warnings = r.Warnings
		result.Messages = r.Messages
		return result
	case *ModuleResult_Success:
		output, err := json.Unmarshal(res.Success.Output, cty.DynamicPseudoType)
		if err != nil {
			result := NewFailure(errors.New("failed to unmarshal output from protobuf result: "+err.Error()), "")
			result.Warnings = r.Warnings
			result.Messages = r.Messages
			return result
		}

		var result *Result
		if res.Success.Changed {
			result = NewChanged(output)
		} else {
			result = NewNotChanged(output)
		}

		result.Warnings = r.Warnings
		result.Messages = r.Messages
		return result
	default:
		result := NewFailure(errors.New("invalid protobuf result: unknown result type"), "")
		result.Warnings = r.Warnings
		result.Messages = r.Messages
		return result
	}
}
