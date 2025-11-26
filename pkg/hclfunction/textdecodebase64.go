// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"bytes"
	"encoding/base64"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"golang.org/x/text/encoding/ianaindex"
)

var (
	TextDecodeBase64Func = function.New(&function.Spec{
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
				return cty.UnknownVal(cty.String), function.NewArgErrorf(1, "invalid encoding %q", args[1].AsString())
			}

			encodingName, err := ianaindex.IANA.Name(encoding)
			if err != nil {
				encodingName = args[1].AsString() // Fallback to the original string if not found, this should not happen
			}

			input := args[0].AsString()

			base64Decoded, err := base64.StdEncoding.DecodeString(input)
			if err != nil {
				switch err := err.(type) {
				case base64.CorruptInputError:
					return cty.UnknownVal(cty.String), function.NewArgErrorf(0, "the input has invalid base64 character at offset: %d", int(err))
				default:
					return cty.UnknownVal(cty.String), function.NewArgErrorf(0, "failed to decode input: %w", err)
				}
			}

			decoder := encoding.NewDecoder()

			decoded, err := decoder.Bytes([]byte(base64Decoded))
			if err != nil || bytes.ContainsRune(decoded, '�') {
				return cty.UnknownVal(cty.String), function.NewArgErrorf(0, "failed to decode input as %q", encodingName)
			}

			return cty.StringVal(string(decoded)), nil
		},
	})
)

func TextDecodeBase64(input, encoding cty.Value) (cty.Value, error) {
	return TextDecodeBase64Func.Call([]cty.Value{input, encoding})
}
