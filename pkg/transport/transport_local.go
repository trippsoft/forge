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

	"github.com/trippsoft/forge/pkg/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	LocalTransport = &localTransport{}
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

// StartPlugin implements Transport.
func (l *localTransport) StartPlugin(
	namespace string,
	pluginName string,
	escalation *Escalation,
) (*grpc.ClientConn, func(), error) {
	// TODO - handle escalation if needed

	pluginPath, err := plugin.FindPluginPath(namespace, pluginName, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, nil, err
	}

	cmd, port, err := l.startPlugin(pluginPath)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		cmd.Process.Kill()
		cmd.Wait()
	}

	address := fmt.Sprintf("127.0.0.1:%d", port)
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return connection, cleanup, nil
}

func (l *localTransport) startPlugin(path string) (*exec.Cmd, uint16, error) {
	cmd := exec.Command(path)

	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MIN_PORT=%d", plugin.LocalPluginMinPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MAX_PORT=%d", plugin.LocalPluginMaxPort))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, 0, fmt.Errorf(
			"failed to get stdout pipe for plugin at '%s': %w",
			path,
			err,
		)
	}

	errBuf := &bytes.Buffer{}
	cmd.Stderr = errBuf

	err = cmd.Start()
	if err != nil {
		return nil, 0, fmt.Errorf(
			"failed to start plugin at '%s': %w - %s",
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
			"invalid port output from plugin at '%s': %w - %s",
			path,
			err,
			errBuf.String(),
		)
	}

	return cmd, uint16(port), nil
}
