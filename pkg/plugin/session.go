// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"io"
)

const (
	// PluginReadyMessage is the message that a plugin process must write to its standard error stream to indicate that
	// it is ready to receive requests.
	PluginReadyMessage = "PLUGIN_READY"
)

// Session represents a running plugin process and provides access to its standard input, output, and error streams.
type Session interface {
	Close() error // Close terminates the plugin process and releases any associated resources.

	Stdin() io.WriteCloser // Stdin returns a pipe that can be used to write to the plugin's standard input.
	Stdout() io.Reader     // Stdout returns a pipe that can be used to read the plugin's standard output.
	Stderr() io.Reader     // Stderr returns a pipe that can be used to read the plugin's standard error.
}
