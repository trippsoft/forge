// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"fmt"

	"github.com/hashicorp/go-cty-funcs/filesystem"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	TemplateFileFunc = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "path",
				Type: customdecode.ExpressionClosureType, // Use the ExpressionClosureType to pass in evaluation context
			},
		},
		Type: function.StaticReturnType(cty.String),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			closure := customdecode.ExpressionClosureFromVal(args[0])
			ctx := closure.EvalContext

			path, diags := closure.Value()
			if diags.HasErrors() {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"templatefile error: failed to evaluate template file path: %w",
					diags,
				)
			}

			if !path.IsWhollyKnown() || path.IsNull() || path.Type() != cty.String {
				return cty.UnknownVal(cty.String), function.NewArgErrorf(
					0,
					"templatefile error: path must be a known non-null string",
				)
			}

			templateValue, err := filesystem.File(".", path)
			if err != nil {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"templatefile error: failed to read template file %q: %w",
					path.AsString(),
					err,
				)
			}

			expression, diags := hclsyntax.ParseTemplate(
				[]byte(templateValue.AsString()),
				path.AsString(),
				hcl.Pos{Line: 1, Column: 1},
			)

			if diags.HasErrors() {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"templatefile error: failed to parse template file %q: %w",
					path.AsString(),
					diags,
				)
			}

			result, diags := expression.Value(ctx)
			if diags.HasErrors() {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"templatefile error: failed to evaluate template file %q: %w",
					path.AsString(),
					diags,
				)
			}

			if !result.IsWhollyKnown() || result.IsNull() || result.Type() != cty.String {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"templatefile error: template file %q did not evaluate to a known non-null string",
					path.AsString(),
				)
			}

			return result, nil
		},
	})
)

func TemplateFile(path cty.Value) (cty.Value, error) {
	return TemplateFileFunc.Call([]cty.Value{path})
}
