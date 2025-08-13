package workflow

import (
	"fmt"
	"slices"

	"github.com/trippsoft/forge/pkg/ui"
)

// Process represents a series of steps in a workflow.
type Process struct {
	name  string // name represents the name of the process.
	steps []Step // steps represents the list of steps in the process.
}

// NewProcess creates a new Process.
func NewProcess(name string, steps ...Step) *Process {
	return &Process{
		name:  name,
		steps: steps,
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
func (p *Process) Run(ctx *workflowContext) error {

	nameText := ui.Text(p.name).WithForegroundColor(ui.ForegroundGreen).WithStyle(ui.StyleBold)
	name := ctx.ui.Format(nameText)
	line := ctx.ui.FormatLine('*', nil)

	message := fmt.Sprintf("\nPROCESS - %s\n%s", name, line)
	ctx.ui.Print(message)

	for _, step := range p.steps {
		err := step.Run(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
