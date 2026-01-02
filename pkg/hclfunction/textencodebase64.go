// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"encoding/base64"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"golang.org/x/text/encoding/ianaindex"
)

var (
	TextEncodeBase64Func = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "input",
				Type: cty.String,
			},
			{
				Name: "encoding",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.String),
		RefineResult: func(rb *cty.RefinementBuilder) *cty.RefinementBuilder {
			return rb.NotNull()
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			encoding, err := ianaindex.IANA.Encoding(args[1].AsString())
			if err != nil {
				return cty.UnknownVal(cty.String), function.NewArgErrorf(
					1,
					"textencodebase64 failed: invalid encoding %q",
					args[1].AsString(),
				)
			}

			encodingName, err := ianaindex.IANA.Name(encoding)
			if err != nil {
				encodingName = args[1].AsString() // Fallback to the original string if not found
			}

			input := args[0].AsString()
			encoder := encoding.NewEncoder()
			encoded, err := encoder.Bytes([]byte(input))
			if err != nil {
				return cty.UnknownVal(cty.String), function.NewArgErrorf(
					0,
					"textencodebase64 failed: failed to encode input as %q",
					encodingName,
				)
			}

			base64Encoded := base64.StdEncoding.EncodeToString(encoded)

			return cty.StringVal(base64Encoded), nil
		},
	})
)

func TextEncodeBase64(input, encoding cty.Value) (cty.Value, error) {
	return TextEncodeBase64Func.Call([]cty.Value{input, encoding})
}
