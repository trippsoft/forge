// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"runtime"
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

// OS implements Transport.
func (l *localTransport) OS() (string, error) {
	return runtime.GOOS, nil
}

// Arch implements Transport.
func (l *localTransport) Arch() (string, error) {
	return runtime.GOARCH, nil
}

// Connect implements Transport.
func (l *localTransport) Connect() error {
	return nil
}

// Close implements Transport.
func (l *localTransport) Close() error {
	return nil
}
