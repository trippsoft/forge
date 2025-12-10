// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package main is the entry point for the Forge discovery plugin GRPC server.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/network"
	"google.golang.org/grpc"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	listener, port, err := network.GetListenerAndPortInRange(getMinimumPort(), getMaximumPort())
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

func getMinimumPort() uint16 {
	env := os.Getenv("FORGE_PLUGIN_MIN_PORT")
	if env != "" {
		minPort, err := strconv.ParseUint(env, 10, 16)
		if err == nil {
			return uint16(minPort)
		}
	}

	return 25000
}

func getMaximumPort() uint16 {
	env := os.Getenv("FORGE_PLUGIN_MAX_PORT")
	if env != "" {
		maxPort, err := strconv.ParseUint(env, 10, 16)
		if err == nil {
			return uint16(maxPort)
		}
	}

	return 40000
}
