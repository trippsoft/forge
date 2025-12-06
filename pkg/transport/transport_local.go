// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/trippsoft/forge/pkg/discover"
)

type localTransport struct {
	discoveryPluginBasePath string
	minPluginPort           uint16
	maxPluginPort           uint16
}

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

// StartDiscovery implements Transport.
func (l *localTransport) StartDiscovery() (*discover.DiscoveryClient, error) {
	var extension string
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}

	pluginPath := fmt.Sprintf(
		"%sforge-discover_%s_%s%s",
		l.discoveryPluginBasePath,
		runtime.GOOS,
		runtime.GOARCH,
		extension,
	)

	cmd, port, err := l.startDiscoveryPlugin(pluginPath)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		cmd.Process.Kill()
		cmd.Wait()
	}

	discoveryClient := discover.NewDiscoveryClient(port, cleanup)

	return discoveryClient, nil
}

func (l *localTransport) startDiscoveryPlugin(path string) (*exec.Cmd, uint16, error) {
	cmd := exec.Command(path)

	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MIN_PORT=%d", l.minPluginPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MAX_PORT=%d", l.maxPluginPort))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, 0, fmt.Errorf(
			"failed to get stdout pipe for discovery plugin at '%s': %w",
			path,
			err,
		)
	}

	errBuf := &bytes.Buffer{}
	cmd.Stderr = errBuf

	err = cmd.Start()
	if err != nil {
		return nil, 0, fmt.Errorf(
			"failed to start discovery plugin at '%s': %w - %s",
			path,
			err,
			errBuf.String(),
		)
	}

	scanner := bufio.NewScanner(stdout)
	var portOutput string
	for scanner.Scan() {
		line := scanner.Text()
		portOutput = line
		break
	}

	port, err := strconv.ParseUint(portOutput, 10, 16)
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return nil, 0, fmt.Errorf(
			"invalid port output from discovery plugin at '%s': %w - %s",
			path,
			err,
			errBuf.String(),
		)
	}

	return cmd, uint16(port), nil
}

// LocalTransportBuilder is a builder for LocalTransport.
type LocalTransportBuilder struct {
	discoveryPluginBasePath string
	minPluginPort           uint16
	maxPluginPort           uint16
}

// WithDiscoveryPluginBasePath sets the discovery plugin base path for the local transport.
func (b *LocalTransportBuilder) WithDiscoveryPluginBasePath(path string) *LocalTransportBuilder {
	b.discoveryPluginBasePath = path
	return b
}

// WithPluginPortRange sets the plugin port range for the local transport.
func (b *LocalTransportBuilder) WithPluginPortRange(minPluginPort, maxPluginPort uint16) *LocalTransportBuilder {
	b.minPluginPort = minPluginPort
	b.maxPluginPort = maxPluginPort
	return b
}

// Build constructs the LocalTransport based on the builder's configuration.
func (b *LocalTransportBuilder) Build() Transport {
	return &localTransport{
		discoveryPluginBasePath: b.discoveryPluginBasePath,
		minPluginPort:           b.minPluginPort,
		maxPluginPort:           b.maxPluginPort,
	}
}

// NewLocalTransportBuilder creates a new LocalTransportBuilder.
func NewLocalTransportBuilder() *LocalTransportBuilder {
	return &LocalTransportBuilder{
		discoveryPluginBasePath: discover.DefaultDiscoverPluginBasePath(),
	}
}
