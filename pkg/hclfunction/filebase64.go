// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func MakeFileBase64Func(baseDir string) function.Function {
	return function.New(&function.Spec{
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
			if !filepath.IsAbs(path) {
				path = filepath.Join(baseDir, path)
			}

			path = filepath.Clean(path)
			content, err := os.ReadFile(path)
			if err != nil {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"filebase64 failed: failed to read file %q: %w",
					path,
					err,
				)
			}

			stringBuilder := &strings.Builder{}
			encoder := base64.NewEncoder(base64.StdEncoding, stringBuilder)
			_, err = encoder.Write(content)
			if err != nil {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"filebase64 failed: failed to encode file %q to base64: %w",
					path,
					err,
				)
			}

			err = encoder.Close()
			if err != nil {
				return cty.UnknownVal(cty.String), fmt.Errorf(
					"filebase64 failed: failed to finalize base64 encoding for file %q: %w",
					path,
					err,
				)
			}

			return cty.StringVal(stringBuilder.String()), nil
		},
	})
}

func FileBase64(path cty.Value) (cty.Value, error) {
	workingDir, _ := os.Getwd()
	return MakeFileBase64Func(workingDir).Call([]cty.Value{path})
}
