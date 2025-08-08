// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	EndsWithFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "value",
				Type: cty.String,
			},
			{
				Name: "suffix",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.Bool),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			value := args[0].AsString()
			suffix := args[1].AsString()

			return cty.BoolVal(strings.HasSuffix(value, suffix)), nil
		},
	})
)

func EndsWith(value, suffix cty.Value) (cty.Value, error) {
	return EndsWithFunc.Call([]cty.Value{value, suffix})
}
