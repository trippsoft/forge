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
	LengthFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:             "value",
				Type:             cty.DynamicPseudoType,
				AllowDynamicType: true,
				AllowUnknown:     true,
			},
		},
		Type: func(args []cty.Value) (cty.Type, error) {
			argType := args[0].Type()
			if argType == cty.String ||
				argType == cty.DynamicPseudoType ||
				argType.IsTupleType() ||
				argType.IsObjectType() ||
				argType.IsListType() ||
				argType.IsMapType() ||
				argType.IsSetType() {
				return cty.Number, nil
			}

			return cty.Number, errors.New("length function requires a string, tuple, object, list, map, or set type")
		},
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			arg := args[0]
			argType := arg.Type()

			switch {
			case argType == cty.DynamicPseudoType:
				return cty.UnknownVal(cty.Number), nil
			case argType.IsTupleType():
				length := len(arg.Type().TupleElementTypes())
				return cty.NumberIntVal(int64(length)), nil
			case argType.IsObjectType():
				length := len(argType.AttributeTypes())
				return cty.NumberIntVal(int64(length)), nil
			case argType == cty.String:
				return stdlib.Strlen(arg)
			case argType.IsListType() || argType.IsSetType() || argType.IsMapType():
				return arg.Length(), nil
			default:
				return cty.UnknownVal(cty.Number), errors.New("length function requires a string, tuple, object, list, map, or set type")
			}
		},
	})
)

func Length(value cty.Value) (cty.Value, error) {
	return LengthFunc.Call([]cty.Value{value})
}
