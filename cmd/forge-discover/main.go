// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package main is the entry point for the Forge discovery plugin GRPC server.
package main

import (
	"fmt"
	"os"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/trippsoft/forge/pkg/util"
	"google.golang.org/grpc"
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
	discoveryServer := info.NewDiscoveryServer()
	info.RegisterDiscoveryPluginServer(s, discoveryServer)

	go func() {
		err = s.Serve(listener)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	discoveryServer.WaitForShutdown()
	s.GracefulStop()

	return nil
}
