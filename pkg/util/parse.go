// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty/gocty"
)

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
	return str, nil
}

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
	return num, nil
}

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
	return b, nil
}

func ConvertHCLAttributeToDuration(attribute *hcl.Attribute, evalCtx *hcl.EvalContext) (time.Duration, hcl.Diagnostics) {

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
			Detail:   fmt.Sprintf("The value for '%s' could not be converted to a duration string: %v", attribute.Name, err),
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

	return duration, nil
}
