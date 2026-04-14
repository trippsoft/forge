// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"context"

	"github.com/trippsoft/forge/pkg/plugin"
)

const (
	forgeSudoPrompt  = "forge_sudo_prompt"
	forgeGSudoPrompt = "Password for user "
)

// TransportType represents the type of transport used for connecting to managed systems.
type TransportType string

const (
	TransportTypeLocal TransportType = "local"
	TransportTypeSSH   TransportType = "ssh"
)

// Transport defines the transport mechanism for interacting a managed system.
type Transport interface {
	Type() TransportType   // Type returns the type of transport.  This is used for testing only.
	OS() (string, error)   // OS returns the operating system of the managed system.
	Arch() (string, error) // Arch returns the architecture of the managed system.

	Connect() error // Connect establishes the transport connection.
	Close() error   // Close terminates the transport connection.

	// StartPluginSession initializes a plugin session.
	//
	// basePath specifies the base path for plugins.
	//
	// namespace specifies the namespace of the plugin in the filename.
	//
	// pluginName specifies the name of the plugin in the filename.
	// The OS, architecture, and extension, if applicable, will be appended to this path.
	//
	// escalation specifies the escalation method to use when starting the plugin session, if any.
	// This is nil if no escalation is needed.
	//
	// It returns a plugin.Session that represents the running plugin process and provides access to its standard input,
	// output, and error streams.
	StartPluginSession(
		ctx context.Context,
		basePath string,
		namespace string,
		pluginName string,
		escalation *Escalation,
	) (plugin.Session, error)
}
