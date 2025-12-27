// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"

	"google.golang.org/grpc"
)

var (
	_ Transport = (*MockTransport)(nil)
)

type MockCmd struct {
	completed bool

	stdin io.Reader

	Stdout string
	Stderr string
	Err    error
}

// OutputWithError implements Cmd.
func (m *MockCmd) OutputWithError(ctx context.Context) (stdout string, stderr string, err error) {
	if m.completed {
		return "", "", fmt.Errorf("command already completed")
	}

	m.completed = true

	m.Stdout = strings.TrimSpace(m.Stdout)
	m.Stderr = strings.TrimSpace(m.Stderr)

	if m.Err != nil {
		return m.Stdout, m.Stderr, m.Err
	}

	return m.Stdout, m.Stderr, nil
}

// Output implements Cmd.
func (m *MockCmd) Output(ctx context.Context) (string, error) {
	if m.completed {
		return "", fmt.Errorf("command already completed")
	}

	m.completed = true

	m.Stdout = strings.TrimSpace(m.Stdout)

	if m.Err != nil {
		return m.Stdout, m.Err
	}

	return m.Stdout, nil
}

// Run implements Cmd.
func (m *MockCmd) Run(ctx context.Context) error {
	if m.completed {
		return fmt.Errorf("command already completed")
	}

	m.completed = true

	if m.Err != nil {
		return m.Err
	}

	return nil
}

// MockTransport is a mock implementation of the Transport interface for testing purposes.
type MockTransport struct {
	TransportType TransportType

	CommandResults map[string]*MockCmd

	Files map[string][]byte
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

// StartPlugin implements Transport.
func (w *MockTransport) StartPlugin(
	basePath string,
	namespace string,
	pluginName string,
	escalation *Escalation,
) (*grpc.ClientConn, func(), error) {

	panic("unimplemented") // TODO: Implement mock discovery client if needed
}

// NewMockTransport creates a new instance of MockTransport with default settings.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		TransportType: TransportTypeSSH,
	}
}
