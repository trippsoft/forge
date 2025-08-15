// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/cobra"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/module/assert"
	"github.com/trippsoft/forge/pkg/module/message"
	"github.com/trippsoft/forge/pkg/module/shell"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/trippsoft/forge/pkg/workflow"
)

const (
	versionNumber = "0.1.0"
	versionSuffix = "-dev"
)

var (
	inventoryPaths []string
	workflowPath   string

	UI          = ui.StdUI()
	ErrorFormat = ui.TextFormat().WithStyle(ui.StyleBold)
)

func main() {
	inventoryCmd := &cobra.Command{
		Use:   "inventory",
		Short: "Parse and display HCL inventory",
		Long:  "Parses HCL inventory files and displays the inventory of managed hosts.",
		Run: func(cmd *cobra.Command, args []string) {

			i, err := parseInventory()
			if err != nil {
				return
			}

			printInventoryTargets(i)
			printInventoryVars(i)
		},
	}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a workflow against an inventory",
		Long:  "Parses and runs a workflow file against the parsed inventory.",
		Run: func(cmd *cobra.Command, args []string) {

			i, err := parseInventory()
			if err != nil {
				UI.Print("Error parsing inventory files. Closing...\n")
				return
			}

			UI.Print("Successfully parsed inventory files.\n")
			UI.Print("\n")

			moduleRegistry := module.NewRegistry()

			registerLocalModules(moduleRegistry)

			w, err := parseWorkflow(i, moduleRegistry)
			if err != nil {
				UI.Print("Error parsing workflow file. Closing...\n")
				return
			}

			UI.Print("Successfully parsed workflow file.\n")
			workflowContext := workflow.WorkflowContext(UI, i)

			w.Run(workflowContext)
		},
	}
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Prints the version number of Forge to the terminal.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Forge v%s%s\n", versionNumber, versionSuffix)
		},
	}

	rootCmd := &cobra.Command{
		Use:   "forge",
		Short: "Forge CLI",
		Long:  "Forge is a configuration management tool for managing remote systems.",
	}

	rootCmd.AddCommand(inventoryCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)

	inventoryCmd.Flags().StringSliceVarP(&inventoryPaths, "inventory", "i", []string{}, "Path to the HCL inventory file(s)")

	runCmd.Flags().StringSliceVarP(&inventoryPaths, "inventory", "i", []string{}, "Path to the HCL inventory file(s)")
	runCmd.Flags().StringVarP(&workflowPath, "workflow", "w", "", "Path to the HCL workflow file")

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func parseInventory() (*inventory.Inventory, error) {

	UI.Print("\nDiscovering inventory files...\n\n")

	inventoryFiles, err := inventory.DiscoverInventoryFiles(inventoryPaths...)
	if err != nil {
		text := ui.Text(err.Error()).WithFormat(ErrorFormat)
		message := fmt.Sprintf("Error discovering inventory files: %s", UI.Format(text))
		UI.Error(message)
		return nil, err
	}

	UI.Print("Found inventory files:\n")
	for _, file := range inventoryFiles {
		UI.Print(fmt.Sprintf(" - %s\n", file.Path()))
	}

	UI.Print("\nParsing inventory files...\n\n")

	i, diags := inventory.ParseInventoryFiles(inventoryFiles)

	printHCLDiags(diags)

	if diags.HasErrors() {
		return nil, diags
	}

	return i, nil
}

func registerLocalModules(moduleRegistry *module.Registry) {

	moduleRegistry.Register("assert", module.NewLocal(&assert.Module{}))
	moduleRegistry.Register("message", module.NewLocal(&message.Module{}))
	moduleRegistry.Register("shell", module.NewLocal(&shell.Module{}))
}

func parseWorkflow(inventory *inventory.Inventory, moduleRegistry *module.Registry) (*workflow.Workflow, error) {

	UI.Print("Parsing workflow file...\n\n")

	content, err := os.ReadFile(workflowPath)
	if err != nil {
		return nil, err
	}

	parser := workflow.NewParser(inventory, moduleRegistry)
	w, diags := parser.ParseWorkflowFile(workflowPath, content)

	printHCLDiags(diags)

	if diags.HasErrors() {
		return nil, diags
	}

	return w, nil
}

func printInventoryTargets(i *inventory.Inventory) {

	UI.Print("Inventory targets:\n\n")

	allHosts := i.Hosts()

	group := ui.Text("all").WithStyle(ui.StyleBold).WithForegroundColor(ui.ForegroundCyan)
	message := fmt.Sprintf("%s:\n", UI.Format(group))

	UI.Print(message)

	for name := range allHosts {
		hostText := ui.Text(name).WithStyle(ui.StyleItalic)
		hostMessage := fmt.Sprintf("  - %s\n", UI.Format(hostText))
		UI.Print(hostMessage)
	}

	UI.Print("\n")

	groups := i.Groups()

	for name, hosts := range groups {

		group := ui.Text(name).WithStyle(ui.StyleBold).WithForegroundColor(ui.ForegroundCyan)

		message := fmt.Sprintf("%s:\n", UI.Format(group))

		UI.Print(message)

		for _, host := range hosts {
			hostText := ui.Text(host.Name()).WithStyle(ui.StyleItalic)
			hostMessage := fmt.Sprintf("  - %s\n", UI.Format(hostText))
			UI.Print(hostMessage)
		}
		UI.Print("\n")
	}
}

func printInventoryVars(i *inventory.Inventory) {

	UI.Print("Inventory variables:\n\n")

	hosts := i.Hosts()

	for name, host := range hosts {
		hostText := ui.Text(name).WithStyle(ui.StyleBold).WithForegroundColor(ui.ForegroundCyan)
		hostMessage := fmt.Sprintf("%s:\n", UI.Format(hostText))
		UI.Print(hostMessage)

		for key, value := range host.Vars() {
			varText := ui.Text(key).WithStyle(ui.StyleBold)
			valueText := ui.Text(util.FormatCtyValueToString(value, 4, 4)).WithStyle(ui.StyleItalic)

			message := fmt.Sprintf("    %s: %s\n", UI.Format(varText), UI.Format(valueText))
			UI.Print(message)
		}

		UI.Print("\n")
	}
}

func printDiags(diags util.Diags) {

	if len(diags) == 0 {
		return
	}

	for _, diag := range diags {

		severityMessage := ""
		if diag.Severity == util.DiagError {
			severityText := ui.Text("ERROR").WithForegroundColor(ui.ForegroundRed).WithStyle(ui.StyleBold)
			severityMessage = fmt.Sprintf("%s:  ", UI.Format(severityText))
		} else {
			severityText := ui.Text("WARNING").WithForegroundColor(ui.ForegroundYellow).WithStyle(ui.StyleBold)
			severityMessage = fmt.Sprintf("%s:", UI.Format(severityText))
		}

		detailText := ui.Text(diag.Detail).WithStyle(ui.StyleItalic)
		detailMessage := UI.Format(detailText)

		message := fmt.Sprintf("  %s %s\n    %s\n", severityMessage, diag.Summary, detailMessage)
		UI.Print(message)
	}
}

func printHCLDiags(diags hcl.Diagnostics) {

	if len(diags) == 0 {
		return
	}

	for _, diag := range diags {

		severityMessage := ""
		if diag.Severity == hcl.DiagError {
			severityText := ui.Text("ERROR").WithForegroundColor(ui.ForegroundRed).WithStyle(ui.StyleBold)
			severityMessage = fmt.Sprintf("%s:  ", UI.Format(severityText))
		} else {
			severityText := ui.Text("WARNING").WithForegroundColor(ui.ForegroundYellow).WithStyle(ui.StyleBold)
			severityMessage = fmt.Sprintf("%s:", UI.Format(severityText))
		}

		detailText := ui.Text(diag.Detail).WithStyle(ui.StyleItalic)
		detailMessage := UI.Format(detailText)

		message := fmt.Sprintf("  %s %s\n    %s\n", severityMessage, diag.Summary, detailMessage)
		UI.Print(message)
	}
}
