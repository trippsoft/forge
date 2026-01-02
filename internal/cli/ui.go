// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"golang.org/x/term"
)

const (
	defaultConsoleWidth = 78 // We will size the console output to this width by default.
)

var UI ui.UI

type CLI struct {
	m sync.Mutex

	color bool

	stdout io.Writer
	stderr io.Writer

	debug bool
}

// Print implements ui.UI.
func (c *CLI) Print(text string) {
	c.printText(c.stdout, text)
}

// PrintError implements ui.UI.
func (c *CLI) PrintError(text string) {
	if c.color {
		text = fmt.Sprintf("\033[31;1;m%s\033[0m", text) // Red, Bold
	}

	c.printText(c.stderr, text)
}

// PrintInventoryTargets implements ui.UI.
func (c *CLI) PrintInventoryTargets(i *inventory.Inventory) {
	sb := &strings.Builder{}
	sb.WriteString("Inventory Targets:\n\n")
	if c.color {
		sb.WriteString("\033[36;1m") // Cyan, Bold
	}

	sb.WriteString("all:")
	if c.color {
		sb.WriteString("\033[0m") // Reset
	}

	sb.WriteRune('\n')

	for name := range i.Hosts() {
		sb.WriteString("  - ")
		if c.color {
			sb.WriteString("\033[3m") // Italic
		}

		sb.WriteString(name)
		if c.color {
			sb.WriteString("\033[0m") // Reset
		}

		sb.WriteRune('\n')
	}

	sb.WriteRune('\n')

	for name, hosts := range i.Groups() {
		if c.color {
			sb.WriteString("\033[36;1m") // Cyan, Bold
		}

		sb.WriteString(name)
		sb.WriteString(":\n")

		for _, host := range hosts {
			sb.WriteString("  - ")
			if c.color {
				sb.WriteString("\033[3m") // Italic
			}

			sb.WriteString(host.Name())
			if c.color {
				sb.WriteString("\033[0m") // Reset
			}

			sb.WriteRune('\n')
		}

		sb.WriteRune('\n')
	}

	c.printText(c.stdout, sb.String())
}

// PrintInventoryVars implements ui.UI.
func (c *CLI) PrintInventoryVars(i *inventory.Inventory) {
	sb := &strings.Builder{}

	sb.WriteString("Inventory Variables:\n\n")

	for name, host := range i.Hosts() {
		if c.color {
			sb.WriteString("\033[36;1m") // Cyan, Bold
		}

		sb.WriteString(name)
		sb.WriteString(":\n")
		if c.color {
			sb.WriteString("\033[0m") // Reset
		}

		for key, value := range host.Vars() {
			sb.WriteString(strings.Repeat(" ", 4))
			if c.color {
				sb.WriteString("\033[1m") // Bold
			}

			sb.WriteString(key)
			if c.color {
				sb.WriteString("\033[0m") // Reset
			}

			sb.WriteString(": ")
			if c.color {
				sb.WriteString("\033[3m") // Italic
			}

			sb.WriteString(util.FormatCtyValueToIndentedString(value, 4, 4))
			if c.color {
				sb.WriteString("\033[0m") // Reset
			}

			sb.WriteRune('\n')
		}

		sb.WriteRune('\n')
	}

	c.printText(c.stdout, sb.String())
}

// PrintHCLDiagnostics implements ui.UI.
func (c *CLI) PrintHCLDiagnostics(diagnostics hcl.Diagnostics) {
	if len(diagnostics) == 0 {
		return
	}

	c.m.Lock()
	defer c.m.Unlock()

	for _, diag := range diagnostics {
		if diag.Severity == hcl.DiagError {
			fmt.Fprint(c.stderr, strings.Repeat(" ", 2))
			if c.color {
				fmt.Fprint(c.stderr, "\033[31;1;m") // Red, Bold
			}

			fmt.Fprint(c.stderr, "ERROR")
			if c.color {
				fmt.Fprint(c.stderr, "\033[0m") // Reset
			}

			fmt.Fprint(c.stderr, ":   ")
			fmt.Fprint(c.stderr, diag.Summary)
			fmt.Fprint(c.stderr, "\n")
			fmt.Fprint(c.stderr, strings.Repeat(" ", 4))
			if c.color {
				fmt.Fprint(c.stderr, "\033[3m") // Italic
			}

			fmt.Fprint(c.stderr, diag.Detail)
			if c.color {
				fmt.Fprint(c.stderr, "\033[0m") // Reset
			}

			fmt.Fprint(c.stderr, "\n")
		} else {
			fmt.Fprint(c.stdout, strings.Repeat(" ", 2))
			if c.color {
				fmt.Fprint(c.stdout, "\033[33;1;m") // Yellow, Bold
			}

			fmt.Fprint(c.stdout, "WARNING")

			if c.color {
				fmt.Fprint(c.stdout, "\033[0m") // Reset
			}

			fmt.Fprint(c.stdout, ": ")
			fmt.Fprint(c.stdout, diag.Summary)
			fmt.Fprint(c.stdout, "\n")
			fmt.Fprint(c.stdout, strings.Repeat(" ", 4))
			if c.color {
				fmt.Fprint(c.stdout, "\033[3m") // Italic
			}

			fmt.Fprint(c.stdout, diag.Detail)
			if c.color {
				fmt.Fprint(c.stdout, "\033[0m") // Reset
			}

			fmt.Fprint(c.stdout, "\n")
		}
	}
}

// PrintHeader implements ui.UI.
func (c *CLI) PrintHeader(level ui.HeaderLevel, prefix, text string) {

	var lineChar string
	var indentation int
	switch level {
	case ui.HeaderLevel1:
		lineChar = "*"
		indentation = 0
	case ui.HeaderLevel2:
		lineChar = "="
		indentation = 1
	default:
		lineChar = "-"
		indentation = 2
	}
	sb := &strings.Builder{}
	sb.WriteRune('\n')
	if indentation > 0 {
		sb.WriteString(strings.Repeat(" ", indentation))
	}

	sb.WriteString(prefix)

	if c.color {
		sb.WriteString("\033[1m") // Bold
	}

	sb.WriteString(text)

	if c.color {
		sb.WriteString("\033[0m") // Reset
	}

	sb.WriteRune('\n')
	if indentation > 0 {
		sb.WriteString(strings.Repeat(" ", indentation))
	}

	sb.WriteString(strings.Repeat(lineChar, c.getConsoleWidth()-indentation))
	sb.WriteRune('\n')

	c.printText(c.stdout, sb.String())
}

// PrintHostResult implements ui.UI.
func (c *CLI) PrintHostResult(hostname string, r *result.Result) {
	c.printResult(hostname, r)
}

// PrintIterationResult implements ui.UI.
func (c *CLI) PrintIterationResult(hostname string, iterationLabel string, r *result.Result) {
	var label string
	if iterationLabel == "" {
		label = hostname
	} else {
		label = fmt.Sprintf("%s -> %s", hostname, iterationLabel)
	}

	c.printResult(label, r)
}

func (c *CLI) printText(writer io.Writer, text string) {
	if writer == nil {
		return
	}

	c.m.Lock()
	defer c.m.Unlock()

	text = util.SecretFilter.Filter(text)

	_, _ = fmt.Fprint(writer, text)
}

func (c *CLI) getConsoleWidth() int {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return defaultConsoleWidth
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return defaultConsoleWidth
	}

	return width - 2
}

func (c *CLI) printResult(label string, r *result.Result) {
	if r == nil {
		r = result.NewFailure(errors.New("no result returned from module"), "")
	}
	errorStringBuilder := &strings.Builder{}
	if r.Err != nil {
		errorStringBuilder.WriteString(strings.Repeat(" ", 6))

		if c.color {
			errorStringBuilder.WriteString("\033[31;1;3m") // Red, Bold, Italic
		}

		errorStringBuilder.WriteString("ERROR:   ")
		errorStringBuilder.WriteString(strings.ReplaceAll(r.Err.Error(), "\n", "\n"+strings.Repeat(" ", 8)))

		if c.color {
			errorStringBuilder.WriteString("\033[0m") // Reset
		}

		errorStringBuilder.WriteRune('\n')

		if r.ErrDetail != "" && c.debug {
			errorStringBuilder.WriteString(strings.Repeat(" ", 8))
			if c.color {
				errorStringBuilder.WriteString("\033[40;3m") // Black, Italic
			}

			errorStringBuilder.WriteString("DETAIL:  ")
			errorStringBuilder.WriteString(strings.ReplaceAll(r.ErrDetail, "\n", "\n"+strings.Repeat(" ", 10)))

			if c.color {
				errorStringBuilder.WriteString("\033[0m") // Reset
			}

			errorStringBuilder.WriteRune('\n')
		}
	}

	outStringBuilder := &strings.Builder{}
	outStringBuilder.WriteString(strings.Repeat(" ", 4))
	outStringBuilder.WriteString(label)
	outStringBuilder.WriteString(": ")

	statusText, statusFormat := c.getStatusTextAndFormat(r)
	labelRunes := utf8.RuneCountInString(label)
	statusRunes := utf8.RuneCountInString(statusText)
	dotsCount := max(c.getConsoleWidth()-11-labelRunes-statusRunes, 1)
	outStringBuilder.WriteString(strings.Repeat(".", dotsCount))
	outStringBuilder.WriteString(" ")
	outStringBuilder.WriteString(statusFormat)
	outStringBuilder.WriteString(statusText)

	if c.color {
		outStringBuilder.WriteString("\033[0m") // Reset
	}

	outStringBuilder.WriteRune('\n')

	if r.Failed && r.IgnoredFailure {
		outStringBuilder.WriteString(strings.Repeat(" ", 6))
		if c.color {
			outStringBuilder.WriteString("\033[44;3m") // Blue, Italic
		}

		outStringBuilder.WriteString("Failure ignored")
		if c.color {
			outStringBuilder.WriteString("\033[0m") // Reset
		}

		outStringBuilder.WriteRune('\n')
	}

	for _, warning := range r.Warnings {
		outStringBuilder.WriteString(strings.Repeat(" ", 6))
		if c.color {
			outStringBuilder.WriteString("\033[33;1;3m") // Yellow, Bold, Italic
		}

		outStringBuilder.WriteString("WARNING: ")
		outStringBuilder.WriteString(strings.ReplaceAll(warning, "\n", "\n"+strings.Repeat(" ", 8)))
		if c.color {
			outStringBuilder.WriteString("\033[0m") // Reset
		}

		outStringBuilder.WriteRune('\n')
	}

	for _, message := range r.Messages {
		outStringBuilder.WriteString(strings.Repeat(" ", 6))
		if c.color {
			outStringBuilder.WriteString("\033[32;1;3m") // Green, Bold, Italic
		}

		outStringBuilder.WriteString("MESSAGE: ")
		outStringBuilder.WriteString(strings.ReplaceAll(message, "\n", "\n"+strings.Repeat(" ", 8)))
		if c.color {
			outStringBuilder.WriteString("\033[0m") // Reset
		}

		outStringBuilder.WriteRune('\n')
	}

	c.m.Lock()
	defer c.m.Unlock()

	stdout := util.SecretFilter.Filter(outStringBuilder.String())
	stderr := util.SecretFilter.Filter(errorStringBuilder.String())

	_, _ = fmt.Fprint(c.stdout, stdout)
	_, _ = fmt.Fprint(c.stderr, stderr)
}

func (c *CLI) getStatusTextAndFormat(r *result.Result) (string, string) {
	if r.Skipped {
		if c.color {
			return "SKIPPED", "\033[36;1m" // Cyan, Italic
		}

		return "SKIPPED", ""
	}

	if r.Failed {
		if c.color {
			return "FAILED", "\033[31;1m" // Red, Bold
		}

		return "FAILED", ""
	}

	if r.Changed {
		if c.color {
			return "CHANGED", "\033[33;1m" // Yellow, Bold
		}

		return "CHANGED", ""
	}

	if c.color {
		return "NOT CHANGED", "\033[32;1m" // Green, Bold
	}

	return "NOT CHANGED", ""
}
