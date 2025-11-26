// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/gocty"
)

// ConvertHCLAttributeToString converts an HCL attribute to a string value.
func ConvertHCLAttributeToString(attribute *hcl.Attribute, evalCtx *hcl.EvalContext) (string, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}

	value, moreDiags := attribute.Expr.Value(evalCtx)
	diags = append(diags, moreDiags...)
	if diags.HasErrors() {
		return "", diags
	}

	var str string
	err := gocty.FromCtyValue(value, &str)
	if err != nil {
		return "", append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid value type",
			Detail:   fmt.Sprintf("The value for '%s' could not be converted to a string: %v", attribute.Name, err),
		})
	}

	return str, diags
}

// ConvertHCLAttributeToUint16 converts an HCL attribute to a uint16 value.
func ConvertHCLAttributeToUint16(attribute *hcl.Attribute, evalCtx *hcl.EvalContext) (uint16, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}

	value, moreDiags := attribute.Expr.Value(evalCtx)
	diags = append(diags, moreDiags...)
	if diags.HasErrors() {
		return 0, diags
	}

	var num uint16
	err := gocty.FromCtyValue(value, &num)
	if err != nil {
		return 0, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid value type",
			Detail:   fmt.Sprintf("The value for '%s' could not be converted to a uint16: %v", attribute.Name, err),
		})
	}

	return num, diags
}

// ConvertHCLAttributeToBool converts an HCL attribute to a bool value.
func ConvertHCLAttributeToBool(attribute *hcl.Attribute, evalCtx *hcl.EvalContext) (bool, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}

	value, moreDiags := attribute.Expr.Value(evalCtx)
	diags = append(diags, moreDiags...)
	if diags.HasErrors() {
		return false, diags
	}

	var b bool
	err := gocty.FromCtyValue(value, &b)
	if err != nil {
		return false, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid value type",
			Detail:   fmt.Sprintf("The value for '%s' could not be converted to a bool: %v", attribute.Name, err),
		})
	}

	return b, diags
}

// ConvertHCLAttributeToDuration converts an HCL attribute to a duration value.
func ConvertHCLAttributeToDuration(
	attribute *hcl.Attribute,
	evalCtx *hcl.EvalContext,
) (time.Duration, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}

	value, moreDiags := attribute.Expr.Value(evalCtx)
	diags = append(diags, moreDiags...)
	if diags.HasErrors() {
		return 0, diags
	}

	var durationString string
	err := gocty.FromCtyValue(value, &durationString)
	if err != nil {
		return 0, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid value type",
			Detail: fmt.Sprintf(
				"The value for '%s' could not be converted to a duration string: %v",
				attribute.Name,
				err,
			),
		})
	}

	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return 0, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid duration",
			Detail:   fmt.Sprintf("The value for '%s' could not be converted to a duration: %v", attribute.Name, err),
		})
	}

	return duration, diags
}

// GetAllCtyStrings returns all string values found within a cty.Value.
func GetAllCtyStrings(value cty.Value) []string {
	results := []string{}

	switch {
	case value.IsNull() || !value.IsWhollyKnown():
		return results
	case value.Type().IsPrimitiveType() && value.Type() == cty.String:
		valueStr := value.AsString()
		if valueStr != "" {
			results = append(results, valueStr)
		}
	case value.Type().IsListType() || value.Type().IsSetType() || value.Type().IsTupleType():

		it := value.ElementIterator()
		for it.Next() {
			_, elementValue := it.Element()
			results = append(results, GetAllCtyStrings(elementValue)...)
		}
	case value.Type().IsMapType():
		for _, elementValue := range value.AsValueMap() {
			results = append(results, GetAllCtyStrings(elementValue)...)
		}
	case value.Type().IsObjectType():
		for name := range value.Type().AttributeTypes() {
			elementValue := value.GetAttr(name)
			results = append(results, GetAllCtyStrings(elementValue)...)
		}
	}

	return results
}

// ModifyUnexpectedElementDiags modifies diagnostics for unexpected elements to be more specific and be of warning
// severity.
func ModifyUnexpectedElementDiags(diags hcl.Diagnostics, location string) hcl.Diagnostics {
	if diags == nil {
		return nil
	}

	for _, diag := range diags {
		if diag.Summary != "Unsupported argument" &&
			diag.Summary != "Unsupported block type" &&
			!(strings.HasPrefix(diag.Summary, "Unexpected ") && strings.HasSuffix(diag.Summary, " block")) {
			continue
		}
		diag.Severity = hcl.DiagWarning
		if location != "" && location != "here" {
			diag.Detail = strings.Replace(diag.Detail, "here", location, 1)
		}
	}

	return diags
}

// FormatCtyValueToString formats a cty.Value to a string representation.
func FormatCtyValueToString(value cty.Value) string {
	if value.IsNull() || !value.IsWhollyKnown() {
		return "null"
	}

	switch {
	case value.IsNull() || !value.IsWhollyKnown():
		return "null"
	case value.Type().Equals(cty.String):
		return fmt.Sprintf("%q", value.AsString())
	case value.Type().Equals(cty.Bool):
		return fmt.Sprintf("%t", value.True())
	case value.Type().Equals(cty.Number):
		converted, _ := convert.Convert(value, cty.String)
		return converted.AsString()
	case value.Type().IsListType() || value.Type().IsSetType() || value.Type().IsTupleType():
		length := value.LengthInt()
		if length == 0 {
			return "[]"
		}

		stringBuilder := &strings.Builder{}
		it := value.ElementIterator()
		stringBuilder.WriteString("[")
		i := 0
		for it.Next() {
			_, elemValue := it.Element()
			stringBuilder.WriteString(FormatCtyValueToString(elemValue))
			if i < length-1 {
				stringBuilder.WriteString(", ")
			}

			i++
		}

		stringBuilder.WriteString("]")

		return stringBuilder.String()

	case value.Type().IsMapType() || value.Type().IsObjectType():
		length := value.LengthInt()
		if length == 0 {
			return "{}"
		}

		stringBuilder := &strings.Builder{}
		it := value.ElementIterator()
		stringBuilder.WriteString("{")
		i := 0
		for it.Next() {
			key, elemValue := it.Element()
			fmt.Fprintf(stringBuilder, "%q: %s", key.AsString(), FormatCtyValueToString(elemValue))
			if i < length-1 {
				stringBuilder.WriteString(", ")
			}

			i++
		}

		stringBuilder.WriteString("}")

		return stringBuilder.String()

	default:
		return "unsupported type"
	}
}

// FormatCtyValueToIndentedString formats a cty.Value to a string with indentation for nested structures.
func FormatCtyValueToIndentedString(value cty.Value, currentIndent int, indentSize int) string {
	if value.IsNull() || !value.IsWhollyKnown() {
		return "null"
	}

	switch {
	case value.IsNull() || !value.IsWhollyKnown() || value.Type().IsPrimitiveType():
		return FormatCtyValueToString(value)
	case value.Type().IsListType() || value.Type().IsSetType() || value.Type().IsTupleType():
		length := value.LengthInt()
		if length == 0 {
			return "[]"
		}

		stringBuilder := &strings.Builder{}
		it := value.ElementIterator()
		stringBuilder.WriteString("[\n")
		i := 0
		for it.Next() {
			stringBuilder.WriteString(strings.Repeat(" ", currentIndent+indentSize))
			_, elemValue := it.Element()
			stringBuilder.WriteString(FormatCtyValueToIndentedString(elemValue, currentIndent+indentSize, indentSize))
			if i < length-1 {
				stringBuilder.WriteString(",\n")
			}

			i++
		}

		stringBuilder.WriteString("\n")
		stringBuilder.WriteString(strings.Repeat(" ", currentIndent))
		stringBuilder.WriteString("]")

		return stringBuilder.String()

	case value.Type().IsMapType() || value.Type().IsObjectType():
		length := value.LengthInt()
		if length == 0 {
			return "{}"
		}

		stringBuilder := &strings.Builder{}
		it := value.ElementIterator()
		stringBuilder.WriteString("{\n")
		i := 0
		for it.Next() {
			stringBuilder.WriteString(strings.Repeat(" ", currentIndent+indentSize))
			key, elemValue := it.Element()
			fmt.Fprintf(
				stringBuilder,
				"%q: %s",
				key.AsString(),
				FormatCtyValueToIndentedString(elemValue, currentIndent+indentSize, indentSize),
			)

			if i < length-1 {
				stringBuilder.WriteString(",\n")
			}

			i++
		}

		stringBuilder.WriteString("\n")
		stringBuilder.WriteString(strings.Repeat(" ", currentIndent))
		stringBuilder.WriteString("}")

		return stringBuilder.String()

	default:
		return "unsupported type"
	}
}
