// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"context"
	"runtime"

	"github.com/trippsoft/forge/pkg/plugin"
)

var (
	_ Transport = (*MockTransport)(nil)
)

// MockTransport is a mock implementation of the Transport interface for testing purposes.
type MockTransport struct {
	TransportType TransportType
}

// Type implements [Transport].
func (m *MockTransport) Type() TransportType {
	return m.TransportType
}

// OS implements [Transport].
func (m *MockTransport) OS() (string, error) {
	return runtime.GOOS, nil
}

// Arch implements [Transport].
func (m *MockTransport) Arch() (string, error) {
	return runtime.GOARCH, nil
}

// Connect implements [Transport].
func (m *MockTransport) Connect() error {
	return nil
}

// Close implements [Transport].
func (m *MockTransport) Close() error {
	return nil
}

// StartPluginSession implements [Transport].
func (w *MockTransport) StartPluginSession(
	ctx context.Context,
	basePath string,
	namespace string,
	pluginName string,
	escalation *Escalation,
) (plugin.Session, error) {
	panic("unimplemented") // TODO: Implement mock plugin if needed
}

// NewMockTransport creates a new instance of MockTransport with default settings.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		TransportType: TransportTypeSSH,
	}
}
