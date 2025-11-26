// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"time"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	TimestampFunc = function.New(&function.Spec{
		Params: []function.Parameter{},
		Type:   function.StaticReturnType(cty.String),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			// This function returns the current timestamp as a string.
			// In a real implementation, you would use time.Now() or similar.
			return cty.StringVal(time.Now().UTC().Format(time.RFC3339)), nil
		},
	})
)

func Timestamp() (cty.Value, error) {
	return TimestampFunc.Call([]cty.Value{})
}
