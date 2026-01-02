// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package result

import (
	"github.com/zclconf/go-cty/cty"
)

// Result holds the result of a module execution.
//
// It includes whether the module made any changes, any error encountered, and the output data.
type Result struct {
	Failed         bool      // Indicates if the module execution failed.
	IgnoredFailure bool      // Indicates if the failure was ignored.
	Skipped        bool      // Indicates if the module was skipped.
	Changed        bool      // Indicates if the module made any changes.
	Err            error     // Error encountered during module execution, if any.
	ErrDetail      string    // Detailed error message, if any.
	Output         cty.Value // Output data from the module execution.
	Warnings       []string  // Warning message, if any.
	Messages       []string  // Informational message, if any.
}

// NewSuccess creates a new success result.
func NewSuccess(changed bool, output cty.Value) *Result {
	return &Result{
		Changed: changed,
		Output:  output,
	}
}

// NewSkipped creates a new skipped result.
func NewSkipped() *Result {
	return &Result{
		Skipped: true,
	}
}

// NewFailure creates a new failure result.
func NewFailure(err error, errDetail string) *Result {
	return &Result{
		Err:       err,
		ErrDetail: errDetail,
		Failed:    true,
	}
}
