// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package main is the entry point for the Forge core plugin.
package main

import (
	"fmt"
	"os"

	"github.com/trippsoft/forge/internal/module"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
)

func main() {
	plugins := []pluginv1.PluginModule{
		module.Command,
		module.Dnf,
		module.DnfInfo,
		module.LocalCopy,
		module.Package,
		module.PackageInfo,
		module.Slurp,
	}

	p := pluginv1.NewPluginV1("forge", "core", plugins...)

	err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
