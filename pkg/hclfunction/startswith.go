// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	StartsWithFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:         "value",
				Type:         cty.String,
				AllowUnknown: false,
			},
			{
				Name: "prefix",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.Bool),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			value := args[0].AsString()
			prefix := args[1].AsString()
			return cty.BoolVal(strings.HasPrefix(value, prefix)), nil
		},
	})
)

func StartsWith(value, prefix cty.Value) (cty.Value, error) {
	return StartsWithFunc.Call([]cty.Value{value, prefix})
}
