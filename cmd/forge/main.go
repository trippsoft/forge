// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package main is the entry point for the Forge CLI application.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trippsoft/forge/internal/cli"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/trippsoft/forge/pkg/workflow"
)

var (
	inventoryPaths []string
	workflowPath   string
	debug          bool
)

func main() {
	inventoryCmd := &cobra.Command{
		Use:   "inventory",
		Short: "Parse and display HCL inventory",
		Long:  "Parses HCL inventory files and displays the inventory of managed hosts.",
		Run: func(cmd *cobra.Command, args []string) {
			cli.InitUI(debug)
			i, err := parseInventory()
			if err != nil {
				os.Exit(1)
			}

			cli.UI.PrintInventoryTargets(i)
			cli.UI.PrintInventoryVars(i)
		},
	}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a workflow against an inventory",
		Long:  "Parses and runs a workflow file against the parsed inventory.",
		Run: func(cmd *cobra.Command, args []string) {
			cli.InitUI(debug)
			i, err := parseInventory()
			if err != nil {
				os.Exit(1)
			}

			moduleRegistry := module.NewRegistry()

			moduleRegistry.RegisterCoreModules()
			moduleRegistry.RegisterPluginModules()

			w, err := parseWorkflow(i, moduleRegistry)
			if err != nil {
				os.Exit(1)
			}

			workflowContext := workflow.NewWorkflowContext(cli.UI, i, debug)

			_, err = w.Run(workflowContext)

			if err != nil {
				os.Exit(1)
			}
		},
	}
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Prints the version number of Forge to the terminal.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Forge v%s%s\n", util.Version, util.VersionSuffix)
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
	inventoryCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	runCmd.Flags().StringSliceVarP(&inventoryPaths, "inventory", "i", []string{}, "Path to the HCL inventory file(s)")
	runCmd.Flags().StringVarP(&workflowPath, "workflow", "w", "", "Path to the HCL workflow file")
	runCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func parseInventory() (*inventory.Inventory, error) {
	cli.UI.Print("\nDiscovering inventory files...\n\n")

	inventoryFiles, err := inventory.DiscoverInventoryFiles(inventoryPaths...)
	if err != nil {
		cli.UI.PrintError(err.Error())
		return nil, err
	}

	if len(inventoryFiles) == 0 {
		cli.UI.PrintError("No inventory files found. Closing...\n")
		return nil, errors.New("No inventory files found.")
	}

	cli.UI.Print("Found inventory files:\n")
	for _, file := range inventoryFiles {
		cli.UI.Print(fmt.Sprintf(" - %s\n", file.Path))
	}

	cli.UI.Print("\nParsing inventory files...\n\n")
	i, diags := inventory.ParseInventoryFiles(inventoryFiles...)

	cli.UI.PrintHCLDiagnostics(diags)
	if diags.HasErrors() {
		cli.UI.PrintError("Error parsing inventory files. Closing...\n")
		return nil, diags
	}

	cli.UI.Print("Successfully parsed inventory files.\n\n")
	return i, nil
}

func parseWorkflow(inventory *inventory.Inventory, moduleRegistry *module.Registry) (*workflow.Workflow, error) {
	cli.UI.Print("Parsing workflow file...\n\n")

	if workflowPath == "" {
		cli.UI.PrintError("No workflow file specified. Closing...\n")
		return nil, errors.New("No workflow file specified.")
	}

	content, err := os.ReadFile(workflowPath)
	if err != nil {
		cli.UI.PrintError(fmt.Sprintf("Error reading workflow file: %s\n", err.Error()))
		return nil, err
	}

	parser := workflow.NewParser(inventory, moduleRegistry)
	w, diags := parser.ParseWorkflowFile(workflowPath, content)

	cli.UI.PrintHCLDiagnostics(diags)
	if diags.HasErrors() {
		cli.UI.Print("\nError parsing workflow file. Closing...\n")
		return nil, diags
	}

	cli.UI.Print("Successfully parsed workflow file.\n")
	return w, nil
}
