// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/plugin"
)

func main() {
	err := realMain()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	_, err := fmt.Fprintln(os.Stderr, plugin.PluginReadyMessage)
	if err != nil {
		return err
	}

	var request info.DiscoverRequest
	err = plugin.Read(os.Stdin, &request)
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	response := info.DiscoverHostInfo()
	err = plugin.Write(os.Stdout, response)
	if err != nil {
		return err
	}

	return nil
}
