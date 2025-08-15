// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"fmt"
	"slices"
	"sync"
	"unicode/utf8"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
)

// Process represents a series of steps in a workflow.
type Process struct {
	name       string // name represents the name of the process.
	gatherInfo bool   // gatherInfo indicates whether to gather information during the process.
	steps      []Step // steps represents the list of steps in the process.
}

// NewProcess creates a new Process.
func NewProcess(name string, gatherInfo bool, steps ...Step) *Process {
	return &Process{
		name:       name,
		gatherInfo: gatherInfo,
		steps:      steps,
	}
}

// Name returns the name of the process.
// This is used primarily for testing purposes.
func (p *Process) Name() string {
	return p.name
}

// Steps returns a copy of the steps in the process.
// This is used primarily for testing purposes.
func (p *Process) Steps() []Step {
	steps := slices.Clone(p.steps)
	return steps
}

// Run executes the process with the given workflow context.
func (p *Process) Run(ctx *workflowContext) {

	nameText := ui.Text(p.name).WithStyle(ui.StyleBold)
	name := ctx.ui.Format(nameText)
	line := ctx.ui.FormatLine('*', nil)

	message := fmt.Sprintf("\nPROCESS - %s\n%s", name, line)
	ctx.ui.Print(message)

	if p.gatherInfo {
		nameText = ui.Text("Gathering Information").WithStyle(ui.StyleBold)
		name := ctx.ui.Format(nameText)
		line := ctx.ui.FormatLine('=', nil)

		message := fmt.Sprintf("\n%s\n%s", name, line)
		ctx.ui.Print(message)

		wg := sync.WaitGroup{}

		for _, host := range ctx.inventory.Hosts() {
			wg.Add(1)
			go func(host *inventory.Host) {
				defer wg.Done()

				t := host.Transport()

				var resultCode stepResultCode
				var diags util.Diags
				if t.Type() != transport.TransportTypeNone {
					diags = host.Info().Populate(host.Transport())

					if diags.HasErrors() {
						resultCode = stepResultFailure
					} else {
						resultCode = stepResultNotChanged
					}
				} else {
					resultCode = stepResultSkipped
				}

				hostName := host.Name()
				hostMessage := fmt.Sprintf("%s:", hostName)
				runeCount := utf8.RuneCountInString(hostName)
				hostText := ui.Text(hostMessage).WithLeftMargin(2).WithRightMargin(65 - runeCount)

				statusMessage := stepResultText[resultCode]
				runeCount = utf8.RuneCountInString(statusMessage)
				statusText := ui.Text(statusMessage).WithFormat(stepResultFormat[resultCode]).WithLeftMargin(12 - runeCount)

				ctx.ui.Print(fmt.Sprintf("%s%s\n", ctx.ui.Format(hostText), ctx.ui.Format(statusText)))
				if len(diags) > 0 {
					printDiags(ctx, diags)
				}
			}(host)
		}

		wg.Wait()
	}

	for _, step := range p.steps {
		step.Run(ctx)
	}
}

func printDiags(ctx *workflowContext, diags util.Diags) {

	if len(diags) == 0 {
		return
	}

	for _, diag := range diags {

		severityMessage := ""
		if diag.Severity == util.DiagError {
			severityText := ui.Text("ERROR").WithForegroundColor(ui.ForegroundRed).WithStyle(ui.StyleBold)
			severityMessage = fmt.Sprintf("%s:  ", ctx.ui.Format(severityText))
		} else {
			severityText := ui.Text("WARNING").WithForegroundColor(ui.ForegroundYellow).WithStyle(ui.StyleBold)
			severityMessage = fmt.Sprintf("%s:", ctx.ui.Format(severityText))
		}

		detailText := ui.Text(diag.Detail).WithStyle(ui.StyleItalic)
		detailMessage := ctx.ui.Format(detailText)

		message := fmt.Sprintf("  %s %s\n    %s\n", severityMessage, diag.Summary, detailMessage)
		ctx.ui.Print(message)
	}
}
