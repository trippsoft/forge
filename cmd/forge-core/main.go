// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package main is the entry point for the Forge core plugin GRPC server.
package main

import (
	"fmt"
	"os"

	"github.com/trippsoft/forge/internal/module"
	"github.com/trippsoft/forge/pkg/plugin"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/util"
	"google.golang.org/grpc"
)

var (
	plugins []pluginv1.PluginModule = []pluginv1.PluginModule{
		module.Command,
		module.Dnf,
		module.DnfInfo,
		module.Package,
		module.PackageInfo,
	}
)

func main() {
	err := realMain()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	listener, port, err := util.GetListenerAndPortInRange(plugin.GetMinimumPort(), plugin.GetMaximumPort())
	if err != nil {
		return err
	}

	fmt.Printf("%d\n", port)
	defer listener.Close()

	s := grpc.NewServer()
	pluginServer := pluginv1.NewPluginServer(plugins...)
	pluginv1.RegisterPluginV1ServiceServer(s, pluginServer)
	plugin.RegisterPluginServiceServer(s, pluginServer)

	go func() {
		err = s.Serve(listener)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	pluginServer.WaitForShutdown()
	s.GracefulStop()

	return nil
}
