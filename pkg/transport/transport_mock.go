// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"runtime"

	"github.com/trippsoft/forge/pkg/discover"
)

var (
	_ Transport = (*MockTransport)(nil)
)

// MockTransport is a mock implementation of the Transport interface for testing purposes.
type MockTransport struct {
	TransportType TransportType
}

// Type implements Transport.
func (w *MockTransport) Type() TransportType {
	return w.TransportType
}

// OS implements Transport.
func (w *MockTransport) OS() (string, error) {
	return runtime.GOOS, nil
}

// Arch implements Transport.
func (w *MockTransport) Arch() (string, error) {
	return runtime.GOARCH, nil
}

// Connect implements Transport.
func (w *MockTransport) Connect() error {
	return nil
}

// Close implements Transport.
func (w *MockTransport) Close() error {
	return nil
}

// StartDiscovery implements Transport.
func (w *MockTransport) StartDiscovery() (*discover.DiscoveryClient, error) {
	panic("unimplemented") // TODO: Implement mock discovery client if needed
}

// NewMockTransport creates a new instance of MockTransport with default settings.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		TransportType: TransportTypeSSH,
	}
}
