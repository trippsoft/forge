// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"runtime"

	"github.com/trippsoft/forge/pkg/info"
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

// Connect implements Transport.
func (w *MockTransport) Connect() error {
	return nil
}

// Close implements Transport.
func (w *MockTransport) Close() error {
	return nil
}

// GetRuntimeInfo implements Transport.
func (w *MockTransport) GetRuntimeInfo() (*info.RuntimeInfo, error) {
	return info.NewRuntimeInfo(runtime.GOOS, runtime.GOARCH), nil
}

// NewMockTransport creates a new instance of MockTransport with default settings.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		TransportType: TransportTypeSSH,
	}
}
