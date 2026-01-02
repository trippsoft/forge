// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	SensitiveFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "value",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.String),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			value := args[0].AsString()
			if value == "" {
				return args[0], nil // Skip registering empty strings
			}

			util.SecretFilter.AddSecret(value)
			return args[0], nil
		},
	})
)

func Sensitive(value cty.Value) (cty.Value, error) {
	return SensitiveFunc.Call([]cty.Value{value})
}
