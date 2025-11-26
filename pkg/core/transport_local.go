// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"runtime"

	"github.com/trippsoft/forge/pkg/info"
)

var (
	// LocalTransport is a transport that represents local system access.
	LocalTransport Transport = &localTransport{}
)

type localTransport struct{}

// Type implements Transport.
func (l *localTransport) Type() TransportType {
	return TransportTypeLocal
}

// Connect implements Transport.
func (l *localTransport) Connect() error {
	return nil
}

// Close implements Transport.
func (l *localTransport) Close() error {
	return nil
}

// GetRuntimeInfo implements Transport.
func (l *localTransport) GetRuntimeInfo() (*info.RuntimeInfo, error) {
	return info.NewRuntimeInfo(runtime.GOOS, runtime.GOARCH), nil
}
