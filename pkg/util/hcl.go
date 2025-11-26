// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
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

	return duration, diags
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
