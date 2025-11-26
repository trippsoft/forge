// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import "fmt"

// DiagSeverity represents the severity of a diagnostic message.
type DiagSeverity int

const (
	// DiagError indicates that the problem reported by a diagnostic will likely prevent the host from being managed.
	DiagError DiagSeverity = iota
	// DiagWarning indicates that the problem reported by a diagnostic is not fatal, but may not be intended.
	DiagWarning
	// DiagDebug indicates normal operations for debugging purposes.
	DiagDebug
)

func (s DiagSeverity) String() string {
	switch s {
	case DiagError:
		return "error"
	case DiagWarning:
		return "warning"
	case DiagDebug:
		return "debug"
	default:
		return "unknown"
	}
}

// Diag represents a diagnostic message that can be used to report issues or information about the host.
// It includes a severity level, a summary, and detailed information.
type Diag struct {
	// Severity indicates the severity of the diagnostic.
	Severity DiagSeverity
	// Summary is a short description of the diagnostic.
	Summary string
	// Detail is a longer description of the diagnostic, which may include remediation steps.
	Detail string
}

// Error returns a formatted string representation of the diagnostic message.
// It includes the summary and detail, formatted according to the severity.
func (d *Diag) Error() string {
	return fmt.Sprintf("%s; %s", d.Summary, d.Detail)
}

type Diags []*Diag

// Error returns a formatted string representation of the diagnostics.
// It summarizes the diagnostics, indicating how many there are and providing details for each.
func (d Diags) Error() string {
	length := len(d)

	switch length {
	case 0:
		return "no diagnostics to report"
	case 1:
		return d[0].Error()
	default:
		return fmt.Sprintf("%s, and %d more diagnostics", d[0].Error(), length-1)
	}
}

// HasErrors checks if there are any diagnostics with severity DiagError.
func (d Diags) HasErrors() bool {
	for _, diag := range d {
		if diag.Severity == DiagError {
			return true
		}
	}

	return false
}

// Errors returns a new Diags containing only the diagnostics with severity DiagError.
func (d Diags) Errors() Diags {
	var errors Diags
	for _, diag := range d {
		if diag.Severity == DiagError {
			errors = append(errors, diag)
		}
	}
	return errors
}

// HasWarnings checks if there are any diagnostics with severity DiagWarning.
func (d Diags) HasWarnings() bool {
	for _, diag := range d {
		if diag.Severity == DiagWarning {
			return true
		}
	}

	return false
}

// Warnings returns a new Diags containing only the diagnostics with severity DiagWarning.
func (d Diags) Warnings() Diags {
	var warnings Diags
	for _, diag := range d {
		if diag.Severity == DiagWarning {
			warnings = append(warnings, diag)
		}
	}
	return warnings
}

// HasDebug checks if there are any diagnostics with severity DiagDebug.
func (d Diags) HasDebug() bool {
	for _, diag := range d {
		if diag.Severity == DiagDebug {
			return true
		}
	}

	return false
}

// Debugs returns a new Diags containing only the diagnostics with severity DiagDebug.
func (d Diags) Debugs() Diags {
	var debugs Diags
	for _, diag := range d {
		if diag.Severity == DiagDebug {
			debugs = append(debugs, diag)
		}
	}
	return debugs
}

// Append appends a new diagnostic to the Diags collection.
func (d Diags) Append(diag *Diag) Diags {
	if diag == nil {
		return d
	}
	return append(d, diag)
}

// AppendAll appends all diagnostics from another Diags collection to the current one.
func (d Diags) AppendAll(diags Diags) Diags {
	if diags == nil {
		return d
	}

	return append(d, diags...)
}
