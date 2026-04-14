// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	SumFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "list",
				Type: cty.DynamicPseudoType,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			if !args[0].CanIterateElements() {
				return cty.UnknownVal(cty.Number), function.NewArgErrorf(0, "sum failed: requires an iterable type")
			}

			if args[0].LengthInt() == 0 {
				return cty.UnknownVal(cty.Number), function.NewArgErrorf(0, "sum failed: requires a non-empty iterable")
			}

			listType := args[0].Type()
			if !listType.IsListType() && !listType.IsSetType() && !listType.IsTupleType() {
				return cty.UnknownVal(cty.Number), function.NewArgErrorf(
					0,
					"sum failed: requires a list, set, or tuple type, got %q",
					listType.FriendlyName(),
				)
			}

			if !args[0].IsWhollyKnown() {
				return cty.UnknownVal(cty.Number), nil
			}

			list := args[0].AsValueSlice()

			sum := list[0]
			if sum.IsNull() {
				return cty.UnknownVal(cty.Number), function.NewArgErrorf(
					0,
					"sum failed: requires a list, set, or tuple of numbers, got null",
				)
			}

			sum, err := convert.Convert(sum, cty.Number)
			if err != nil {
				return cty.UnknownVal(cty.Number), function.NewArgErrorf(
					0,
					"sum failed: requires a list, set, or tuple of numbers",
				)
			}

			for _, value := range list[1:] {
				if value.IsNull() {
					return cty.UnknownVal(cty.Number), function.NewArgErrorf(
						0,
						"sum failed: requires a list, set, or tuple of numbers, got null",
					)
				}

				value, err := convert.Convert(value, cty.Number)
				if err != nil {
					return cty.UnknownVal(cty.Number), function.NewArgErrorf(
						0,
						"sum failed: requires a list, set, or tuple of numbers",
					)
				}

				sum = sum.Add(value)
			}

			return sum, nil
		},
	})
)

func Sum(input cty.Value) (cty.Value, error) {
	return SumFunc.Call([]cty.Value{input})
}
