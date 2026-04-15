// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package result provides the types and utilities for representing the results of operations in Forge.
package result

import (
	"fmt"

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

	// AddWarnings adds the provided warnings to the result.
	AddWarnings(warnings ...string)

	// Messages returns a list of messages associated with the result, if any.
	Messages() []string

	// AddMessages adds the provided messages to the result.
	AddMessages(messages ...string)

	// ErrorMessage returns the error message associated with the result, if it is a failure result.
	ErrorMessage() string

	// ErrorDetails returns the error details associated with the result, if it is a failure result.
	ErrorDetails() string

	// IgnoreFailure ensures that if the result represents a failure, it is treated as an ignored failure instead.
	IgnoreFailure()

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

func (r *result) AddWarnings(warnings ...string) {
	if r.warnings == nil {
		r.warnings = warnings
		return
	}

	r.warnings = append(r.warnings, warnings...)
}

func (r *result) Messages() []string {
	return r.messages
}

func (r *result) AddMessages(messages ...string) {
	if r.messages == nil {
		r.messages = messages
		return
	}

	r.messages = append(r.messages, messages...)
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

// ErrorMessage implements [Result].
func (f *FailedResult) ErrorMessage() string {
	return f.errorMessage
}

// ErrorDetails implements [Result].
func (f *FailedResult) ErrorDetails() string {
	return f.errorDetails
}

// IgnoreFailure implements [Result].
func (f *FailedResult) IgnoreFailure() {
	f.ignored = true
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
func NewFailedResult(errorMessage string, errorDetails string) Result {
	return &FailedResult{
		result:       &result{},
		ignored:      false,
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

// Warnings implements [Result].
func (s *SkippedResult) Warnings() []string {
	return nil
}

// AddWarnings implements [Result].
func (s *SkippedResult) AddWarnings(warnings ...string) {
	// No-op since skipped results do not have warnings.
}

// Messages implements [Result].
func (s *SkippedResult) Messages() []string {
	return nil
}

// AddMessages implements [Result].
func (s *SkippedResult) AddMessages(messages ...string) {
	// No-op since skipped results do not have messages.
}

// ErrorMessage implements [Result].
func (s *SkippedResult) ErrorMessage() string {
	return ""
}

// ErrorDetails implements [Result].
func (s *SkippedResult) ErrorDetails() string {
	return ""
}

// IgnoreFailure implements [Result].
func (s *SkippedResult) IgnoreFailure() {
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

// ErrorMessage implements [Result].
func (s *SuccessResult) ErrorMessage() string {
	return ""
}

// ErrorDetails implements [Result].
func (s *SuccessResult) ErrorDetails() string {
	return ""
}

// IgnoreFailure implements [Result].
func (s *SuccessResult) IgnoreFailure() {
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
func NewNotChangedResult(output cty.Value) Result {
	return &SuccessResult{
		result:  &result{},
		changed: false,
		output:  output,
	}
}

// NewChangedResult creates a new SuccessResult indicating that the operation was successful and resulted in changes,
// with the provided output, messages, and warnings.
func NewChangedResult(output cty.Value) Result {
	return &SuccessResult{
		result:  &result{},
		changed: true,
		output:  output,
	}
}

// ToResult converts the FailedPB to a FailedResult, using the provided warnings and messages.
func (f *FailedPB) ToResult(warnings []string, messages []string) Result {
	result := NewFailedResult(f.Error, f.Detail)
	result.AddWarnings(warnings...)
	result.AddMessages(messages...)
	return result
}

// ToResult converts the SkippedPB to a SkippedResult.
func (s *SkippedPB) ToResult() Result {
	return Skipped
}

// ToResult converts the SuccessPB to a SuccessResult, using the provided warnings and messages.
func (s *SuccessPB) ToResult(warnings []string, messages []string) Result {
	output, err := json.Unmarshal(s.Output, cty.DynamicPseudoType)
	if err != nil {
		result := NewFailedResult(fmt.Sprintf("failed to unmarshal output: %v", err), "")
		result.AddWarnings(warnings...)
		result.AddMessages(messages...)
		return result
	}

	if s.Changed {
		result := NewChangedResult(output)
		result.AddWarnings(warnings...)
		result.AddMessages(messages...)
		return result
	}

	result := NewNotChangedResult(output)
	result.AddWarnings(warnings...)
	result.AddMessages(messages...)
	return result
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
		result := NewFailedResult("invalid protobuf result: unknown result type", "")
		result.AddWarnings(r.Warnings...)
		result.AddMessages(r.Messages...)
		return result
	}
}
