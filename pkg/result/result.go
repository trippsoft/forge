// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package result provides the types and utilities for representing the results of operations in Forge.
package result

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

var (
	// Skipped is a singleton instance of SkippedResult, representing a skipped operation.
	Skipped Result = &SkippedResult{}
)

// Result represents the outcome of an operation, indicating whether it was a failure, an ignored failure, or a
// successful operation.
type Result interface {
	// IsFailed returns true if the result represents a failure and not ignored.
	IsFailed() bool

	// IsIgnoredFailure returns true if the result represents a failure that is ignored.
	IsIgnoredFailure() bool

	// IsSkipped returns true if the result represents a skipped operation.
	IsSkipped() bool

	// IsChanged returns true if the result represents a successful operating resulting in a change.
	IsChanged() bool

	// Warnings returns a list of warnings associated with the result, if any.
	Warnings() []string

	// Messages returns a list of messages associated with the result, if any.
	Messages() []string

	// Output returns the output value associated with the result, if any.
	//
	// For successful results, this will be the output of the operation.
	// For failed or skipped results, this will be a null value.
	Output() cty.Value

	// ToProtobuf converts the Result to its protobuf representation.
	ToProtobuf() *ResultPB
}

type result struct {
	messages []string
	warnings []string
}

func (r *result) Warnings() []string {
	return r.warnings
}

func (r *result) Messages() []string {
	return r.messages
}

// FailedResult represents a failure result, which may be either an ignored failure or a non-ignored failure.
type FailedResult struct {
	*result
	ignored      bool
	errorMessage string
	errorDetails string
}

// IsFailed implements [Result].
func (f *FailedResult) IsFailed() bool {
	return !f.ignored
}

// IsIgnoredFailure implements [Result].
func (f *FailedResult) IsIgnoredFailure() bool {
	return f.ignored
}

// IsChanged implements [Result].
func (f *FailedResult) IsChanged() bool {
	return false
}

// IsSkipped implements [Result].
func (f *FailedResult) IsSkipped() bool {
	return false
}

// Output implements [Result].
func (f *FailedResult) Output() cty.Value {
	return cty.NullVal(cty.DynamicPseudoType)
}

// ToProtobuf implements [Result].
func (f *FailedResult) ToProtobuf() *ResultPB {
	return &ResultPB{
		Warnings: f.warnings,
		Messages: f.messages,
		Result: &ResultPB_Failed{
			Failed: &FailedPB{
				Ignored: f.ignored,
				Error:   f.errorMessage,
				Detail:  f.errorDetails,
			},
		},
	}
}

// NewFailedResult creates a new FailedResult with the provided error message, error details, messages, and warnings.
func NewFailedResult(
	errorMessage string,
	errorDetails string,
	messages []string,
	warnings []string,
) Result {
	r := &result{
		warnings: warnings,
		messages: messages,
	}

	return &FailedResult{
		result:       r,
		ignored:      false,
		errorMessage: errorMessage,
		errorDetails: errorDetails,
	}
}

// NewIgnoredFailedResult creates a new FailedResult that is marked as ignored, with the provided error message, error
// details, messages, and warnings.
func NewIgnoredFailedResult(
	errorMessage string,
	errorDetails string,
	messages []string,
	warnings []string,
) Result {
	r := &result{
		warnings: warnings,
		messages: messages,
	}

	return &FailedResult{
		result:       r,
		ignored:      true,
		errorMessage: errorMessage,
		errorDetails: errorDetails,
	}
}

// SkippedResult represents a result indicating that an operation was skipped.
type SkippedResult struct{}

// IsChanged implements [Result].
func (s *SkippedResult) IsChanged() bool {
	return false
}

// IsFailed implements [Result].
func (s *SkippedResult) IsFailed() bool {
	return false
}

// IsIgnoredFailure implements [Result].
func (s *SkippedResult) IsIgnoredFailure() bool {
	return false
}

// IsSkipped implements [Result].
func (s *SkippedResult) IsSkipped() bool {
	return true
}

// Messages implements [Result].
func (s *SkippedResult) Messages() []string {
	return nil
}

// Warnings implements [Result].
func (s *SkippedResult) Warnings() []string {
	return nil
}

// Output implements [Result].
func (s *SkippedResult) Output() cty.Value {
	return cty.NullVal(cty.DynamicPseudoType)
}

// ToProtobuf implements [Result].
func (s *SkippedResult) ToProtobuf() *ResultPB {
	return &ResultPB{
		Result: &ResultPB_Skipped{
			Skipped: &SkippedPB{},
		},
	}
}

type SuccessResult struct {
	*result
	changed bool
	output  cty.Value
}

// IsFailed implements [Result].
func (s *SuccessResult) IsFailed() bool {
	return false
}

// IsIgnoredFailure implements [Result].
func (s *SuccessResult) IsIgnoredFailure() bool {
	return false
}

// IsChanged implements [Result].
func (s *SuccessResult) IsChanged() bool {
	return s.changed
}

// IsSkipped implements [Result].
func (s *SuccessResult) IsSkipped() bool {
	return false
}

// Output implements [Result].
func (s *SuccessResult) Output() cty.Value {
	return s.output
}

// ToProtobuf implements [Result].
func (s *SuccessResult) ToProtobuf() *ResultPB {
	output, _ := json.Marshal(s.output, cty.DynamicPseudoType)

	return &ResultPB{
		Warnings: s.warnings,
		Messages: s.messages,
		Result: &ResultPB_Success{
			Success: &SuccessPB{
				Changed: s.changed,
				Output:  output,
			},
		},
	}
}

// NewNotChangedResult creates a new SuccessResult indicating that the operation was successful but did not result in
// any changes, with the provided output, messages, and warnings.
func NewNotChangedResult(output cty.Value, messages []string, warnings []string) Result {
	r := &result{
		warnings: warnings,
		messages: messages,
	}

	return &SuccessResult{
		result:  r,
		changed: false,
		output:  output,
	}
}

// NewChangedResult creates a new SuccessResult indicating that the operation was successful and resulted in changes,
// with the provided output, messages, and warnings.
func NewChangedResult(output cty.Value, messages []string, warnings []string) Result {
	r := &result{
		warnings: warnings,
		messages: messages,
	}

	return &SuccessResult{
		result:  r,
		changed: true,
		output:  output,
	}
}

// ToResult converts the FailedPB to a FailedResult, using the provided warnings and messages.
func (f *FailedPB) ToResult(warnings []string, messages []string) Result {
	return NewFailedResult(f.Error, f.Detail, messages, warnings)
}

// ToResult converts the SkippedPB to a SkippedResult.
func (s *SkippedPB) ToResult() Result {
	return Skipped
}

// ToResult converts the SuccessPB to a SuccessResult, using the provided warnings and messages.
func (s *SuccessPB) ToResult(warnings []string, messages []string) Result {
	output, err := json.Unmarshal(s.Output, cty.DynamicPseudoType)
	if err != nil {
		output = cty.NullVal(cty.DynamicPseudoType)
	}

	if s.Changed {
		return NewChangedResult(output, messages, warnings)
	}

	return NewNotChangedResult(output, messages, warnings)
}

// ToResult converts the ResultPB to a Result, using the appropriate conversion based on the type of result contained in
// the protobuf.
func (r *ResultPB) ToResult() Result {
	switch res := r.Result.(type) {
	case *ResultPB_Failed:
		return res.Failed.ToResult(r.Warnings, r.Messages)
	case *ResultPB_Skipped:
		return res.Skipped.ToResult()
	case *ResultPB_Success:
		return res.Success.ToResult(r.Warnings, r.Messages)
	default:
		return NewFailedResult("invalid protobuf result: unknown result type", "", r.Messages, r.Warnings)
	}
}
