// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	StrContainsFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "value",
				Type: cty.String,
			},
			{
				Name: "substring",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.Bool),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			value := args[0].AsString()
			substring := args[1].AsString()
			if strings.Contains(value, substring) {
				return cty.True, nil
			}

			return cty.False, nil
		},
	})
)

func StrContains(value, substring cty.Value) (cty.Value, error) {
	return StrContainsFunc.Call([]cty.Value{value, substring})
}
