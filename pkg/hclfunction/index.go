// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"errors"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var (
	IndexFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:             "list",
				Type:             cty.DynamicPseudoType,
				AllowDynamicType: true,
			},
			{
				Name: "value",
				Type: cty.DynamicPseudoType,
			},
		},
		Type: function.StaticReturnType(cty.Number),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			list := args[0]
			value := args[1]

			if !list.Type().IsListType() && !list.Type().IsTupleType() {
				return cty.NilVal, errors.New("index function requires a list or tuple type for the first argument")
			}

			if !list.IsKnown() {
				return cty.UnknownVal(cty.Number), nil
			}

			if list.LengthInt() == 0 {
				return cty.NilVal, errors.New("index function requires a non-empty list for the first argument")
			}

			for iterator := list.ElementIterator(); iterator.Next(); {
				index, element := iterator.Element()

				isEqual, err := stdlib.Equal(element, value)
				if err != nil {
					return cty.NilVal, err
				}

				if !isEqual.IsKnown() {
					return cty.UnknownVal(cty.Number), nil
				}

				if isEqual.True() {
					return index, nil
				}
			}

			return cty.NilVal, errors.New("value not found in the list or tuple")
		},
	})
)

func Index(list, value cty.Value) (cty.Value, error) {
	return IndexFunc.Call([]cty.Value{list, value})
}
