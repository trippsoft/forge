// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package main is the entry point for the Forge discovery plugin GRPC server.
package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"

	"github.com/trippsoft/forge/pkg/discover"
	"google.golang.org/grpc"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	listener, port, err := getListenerAndPort()
	if err != nil {
		return err
	}

	fmt.Printf("%d\n", port)
	defer listener.Close()

	s := grpc.NewServer()
	discoveryServer := discover.NewDiscoveryServer()
	discover.RegisterDiscoveryPluginServer(s, discoveryServer)

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

func getListenerAndPort() (net.Listener, int, error) {
	minPort := getMinimumPort()
	maxPort := getMaximumPort()

	if minPort >= maxPort {
		panic("FORGE_DISCOVERY_MIN_PORT must be less than FORGE_DISCOVERY_MAX_PORT")
	}

	triedPorts := make(map[int]any)
	for {
		if len(triedPorts) >= (maxPort - minPort + 1) {
			return nil, 0, fmt.Errorf("no available ports in range %d-%d", minPort, maxPort)
		}

		randomPort := minPort + rand.Intn(maxPort+1-minPort)
		if _, tried := triedPorts[randomPort]; tried {
			continue
		}

		triedPorts[randomPort] = nil
		address := fmt.Sprintf("127.0.0.1:%d", randomPort)

		// Attempt to listen on the randomPort here.
		// If successful, return the listener and port.
		// If not, continue the loop to try another port.
		listener, err := net.Listen("tcp", address)
		if err == nil {
			return listener, randomPort, nil
		}
	}
}

func getMinimumPort() int {
	env := os.Getenv("FORGE_DISCOVERY_MIN_PORT")
	if env != "" {
		minPort, err := strconv.Atoi(env)
		if err == nil {
			return minPort
		}
	}

	return 25000
}

func getMaximumPort() int {
	env := os.Getenv("FORGE_DISCOVERY_MAX_PORT")
	if env != "" {
		maxPort, err := strconv.Atoi(env)
		if err == nil {
			return maxPort
		}
	}

	return 40000
}
