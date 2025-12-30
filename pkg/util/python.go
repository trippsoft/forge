// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// RemoveEmptyLinesAndComments removes empty lines and comments from a script.
func RemoveEmptyLinesAndComments(script string) string {
	lines := strings.Split(script, "\n")
	filteredLines := []string{}
	for _, line := range lines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if strings.HasPrefix(strings.TrimSpace(line), "#") || line == "" {
			continue
		}

		filteredLines = append(filteredLines, line)
	}

	return strings.Join(filteredLines, "\n")
}

// EncodePythonAsBase64 encodes a Python command as a base64 string.
func EncodePythonAsBase64(python string) (string, error) {
	if python == "" {
		return "", errors.New("input Python command cannot be empty")
	}

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(python))
	if base64Encoded == "" {
		return "", errors.New("failed to encode Python command to base64")
	}

	return base64Encoded, nil
}

// FormatInputForPython formats input to be inserted into a Python script.
func FormatInputForPython(input map[string]cty.Value, whatIf bool) string {
	builder := &strings.Builder{}
	builder.WriteString("WHAT_IF = ")

	if whatIf {
		builder.WriteString("True\n")
	} else {
		builder.WriteString("False\n")
	}

	builder.WriteString("INPUT = {")

	first := true
	for key, value := range input {
		if !first {
			builder.WriteString(", ")
		}
		first = false

		builder.WriteString(`"`)
		builder.WriteString(key)
		builder.WriteString(`": `)
		builder.WriteString(formatCtyValueForPython(value))
	}

	builder.WriteString("}\n")
	return builder.String()
}

func formatCtyValueForPython(value cty.Value) string {
	if !value.IsWhollyKnown() || value.IsNull() {
		return "None"
	}

	switch {
	case !value.IsWhollyKnown() || value.IsNull():
		return "None"
	case value.Type().Equals(cty.String):
		return fmt.Sprintf("%q", value.AsString())
	case value.Type().Equals(cty.Bool) && value.True():
		return "True"
	case value.Type().Equals(cty.Bool) && !value.True():
		return "False"
	case value.Type().Equals(cty.Number):
		converted, err := convert.Convert(value, cty.String)
		if err != nil {
			return "None"
		}

		return fmt.Sprintf("%s", converted.AsString())

	case value.Type().IsListType() || value.Type().IsTupleType():
		stringBuilder := &strings.Builder{}
		stringBuilder.WriteString("[")

		i := 0
		it := value.ElementIterator()
		for it.Next() {
			_, elem := it.Element()
			if i > 0 {
				stringBuilder.WriteString(", ")
			}

			stringBuilder.WriteString(formatCtyValueForPython(elem))
			i++
		}

		stringBuilder.WriteString("]")

		return stringBuilder.String()

	case value.Type().IsSetType():
		stringBuilder := &strings.Builder{}
		stringBuilder.WriteString("{")

		i := 0
		it := value.ElementIterator()
		for it.Next() {
			_, elem := it.Element()
			if i > 0 {
				stringBuilder.WriteString(", ")
			}

			stringBuilder.WriteString(formatCtyValueForPython(elem))
			i++
		}

		stringBuilder.WriteString("}")

		return stringBuilder.String()

	case value.Type().IsMapType() || value.Type().IsObjectType():
		stringBuilder := &strings.Builder{}
		stringBuilder.WriteString("{")

		i := 0
		it := value.ElementIterator()
		for it.Next() {
			key, elem := it.Element()
			if i > 0 {
				stringBuilder.WriteString(", ")
			}

			stringBuilder.WriteString(formatCtyValueForPython(key))
			stringBuilder.WriteString(": ")
			stringBuilder.WriteString(formatCtyValueForPython(elem))
			i++
		}

		stringBuilder.WriteString("}")

		return stringBuilder.String()
	default:
		return "None"
	}
}
