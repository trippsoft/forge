// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	AllTrueFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "list",
				Type: cty.List(cty.Bool),
			},
		},
		Type: function.StaticReturnType(cty.Bool),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			result := cty.True
			for iterator := args[0].ElementIterator(); iterator.Next(); {
				_, value := iterator.Element()
				if !value.IsKnown() {
					return cty.UnknownVal(cty.Bool), nil
				}

				if value.IsNull() {
					return cty.False, nil
				}

				result = result.And(value)

				if result.False() {
					return cty.False, nil
				}
			}

			return result, nil
		},
	})
)

func AllTrue(list cty.Value) (cty.Value, error) {
	return AllTrueFunc.Call([]cty.Value{list})
}
