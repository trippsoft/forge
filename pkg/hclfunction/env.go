// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"os"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	EnvFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:         "name",
				Type:         cty.String,
				AllowNull:    false,
				AllowUnknown: false,
			},
		},
		Type: function.StaticReturnType(cty.String),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			name := args[0].AsString()

			return cty.StringVal(os.Getenv(name)), nil
		},
	})
)

func Env(name cty.Value) (cty.Value, error) {
	return EnvFunc.Call([]cty.Value{name})
}
