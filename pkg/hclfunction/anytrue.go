// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	AnyTrueFunc = function.New(&function.Spec{
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
			result := cty.False
			containsUnknown := false

			for iterator := args[0].ElementIterator(); iterator.Next(); {
				_, value := iterator.Element()
				if !value.IsKnown() {
					containsUnknown = true
					continue
				}

				if value.IsNull() {
					continue
				}

				result = result.Or(value)
				if result.True() {
					return cty.True, nil
				}
			}

			if containsUnknown {
				return cty.UnknownVal(cty.Bool), nil
			}

			return result, nil
		},
	})
)

func AnyTrue(list cty.Value) (cty.Value, error) {
	return AnyTrueFunc.Call([]cty.Value{list})
}
