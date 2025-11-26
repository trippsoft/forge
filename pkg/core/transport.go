// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package core

import "github.com/trippsoft/forge/pkg/info"

// TransportType represents the type of transport used for connecting to managed systems.
type TransportType string

const (
	TransportTypeLocal TransportType = "local"
	TransportTypeSSH   TransportType = "ssh"
)

// Transport defines the transport mechanism for interacting a managed system.
type Transport interface {
	Type() TransportType // Type returns the type of transport.

	Connect() error // Connect establishes the transport connection.
	Close() error   // Close terminates the transport connection.

	GetRuntimeInfo() (*info.RuntimeInfo, error) // GetRuntimeInfo retrieves OS and architecture information.
}
