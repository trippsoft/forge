// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var (
	FileBase64Func = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "path",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.String),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			path := args[0].AsString()
			content, err := os.ReadFile(path)
			if err != nil {
				return cty.NullVal(cty.String), fmt.Errorf("failed to read file %q: %w", path, err)
			}

			stringBuilder := &strings.Builder{}
			encoder := base64.NewEncoder(base64.StdEncoding, stringBuilder)
			_, err = encoder.Write(content)
			if err != nil {
				return cty.NullVal(cty.String), fmt.Errorf("failed to encode file %q to base64: %w", path, err)
			}

			err = encoder.Close()
			if err != nil {
				return cty.NullVal(cty.String), fmt.Errorf(
					"failed to finalize base64 encoding for file %q: %w",
					path,
					err,
				)
			}

			return cty.StringVal(stringBuilder.String()), nil
		},
	})
)

func FileBase64(path cty.Value) (cty.Value, error) {
	return FileBase64Func.Call([]cty.Value{path})
}
