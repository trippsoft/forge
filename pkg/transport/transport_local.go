// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"

	"github.com/trippsoft/forge/pkg/plugin"
)

var (
	LocalTransport Transport = &localTransport{}
)

type localPluginSession struct {
	command *exec.Cmd
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	stdin   io.WriteCloser
}

// Close implements [plugin.Session].
func (l *localPluginSession) Close() error {
	l.stdout.Close()
	l.stderr.Close()
	l.stdin.Close()

	err := l.command.Process.Kill()
	if err != nil {
		return err
	}

	return l.command.Wait()
}

// Stdout implements [plugin.Session].
func (l *localPluginSession) Stdout() io.Reader {
	return l.stdout
}

// Stderr implements [plugin.Session].
func (l *localPluginSession) Stderr() io.Reader {
	return l.stderr
}

// Stdin implements [plugin.Session].
func (l *localPluginSession) Stdin() io.WriteCloser {
	return l.stdin
}

type localTransport struct{}

// Type implements [Transport].
func (l *localTransport) Type() TransportType {
	return TransportTypeLocal
}

// OS implements [Transport].
func (l *localTransport) OS() (string, error) {
	return runtime.GOOS, nil
}

// Arch implements [Transport].
func (l *localTransport) Arch() (string, error) {
	return runtime.GOARCH, nil
}

// Connect implements [Transport].
func (l *localTransport) Connect() error {
	return nil
}

// Close implements [Transport].
func (l *localTransport) Close() error {
	return nil
}

// StartPluginSession implements [Transport].
func (l *localTransport) StartPluginSession(
	ctx context.Context,
	basePath string,
	namespace string,
	pluginName string,
	escalation *Escalation,
) (plugin.Session, error) {
	path, err := plugin.FindPluginPath(basePath, namespace, pluginName, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, err
	}

	if escalation != nil {
		return l.startEscalatedPluginSession(ctx, path, escalation)
	}

	cmd := exec.CommandContext(ctx, path)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for plugin at '%s': %w", path, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe for plugin at '%s': %w", path, err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe for plugin at '%s': %w", path, err)
	}

	stderrPipeReader, stderrPipeWriter := io.Pipe()

	readyChan := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		defer stderrPipeWriter.Close()

		var accumulatedStderr strings.Builder
		buf := make([]byte, 4096)

		for {
			n, readErr := stderr.Read(buf)
			if n > 0 {
				accumulatedStderr.Write(buf[:n])
				text := accumulatedStderr.String()

				if strings.Contains(text, plugin.PluginReadyMessage) {
					close(readyChan)
					io.Copy(stderrPipeWriter, stderr)
					return
				}
			}

			if readErr != nil {
				if readErr != io.EOF {
					errChan <- fmt.Errorf("error reading stderr for plugin at '%s': %w", path, readErr)
					stdin.Close()
					cmd.Process.Kill()
					cmd.Wait()
				}
				return
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start plugin at '%s': %w", path, err)
	}

	select {
	case <-ctx.Done():
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
		return nil, fmt.Errorf("context cancelled while starting plugin at '%s': %w", path, ctx.Err())
	case err := <-errChan:
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
		return nil, err
	case <-readyChan:
		return &localPluginSession{
			command: cmd,
			stdout:  stdout,
			stderr:  stderrPipeReader,
			stdin:   stdin,
		}, nil
	}
}
