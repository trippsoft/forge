// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package transport

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/trippsoft/forge/pkg/plugin"
)

func (l *localTransport) startEscalatedPluginSession(
	ctx context.Context,
	path string,
	escalation *Escalation,
) (plugin.Session, error) {
	user := escalation.User()
	if user == "" || user == "SYSTEM" || user == `NT AUTHORITY\SYSTEM` {
		return l.startPluginSessionAsSystem(ctx, path)
	}

	args := []string{"/c", fmt.Sprintf("gsudo.exe -u %s %q", user, path)}

	cmd := exec.CommandContext(ctx, "cmd.exe", args...)

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
		promptsAnswered := 0

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

				if strings.Contains(text, forgeGSudoPrompt) {
					if promptsAnswered >= 3 {
						errChan <- fmt.Errorf("too many gsudo password attempts for plugin at '%s'", path)
						stdin.Close()
						cmd.Process.Kill()
						cmd.Wait()
						return
					}

					promptsAnswered++
					_, err = stdin.Write([]byte(escalation.Pass() + "\n"))
					if err != nil {
						errChan <- fmt.Errorf("failed to write password to stdin for plugin at '%s': %w", path, err)
						stdin.Close()
						cmd.Process.Kill()
						cmd.Wait()
						return
					}

					accumulatedStderr.Reset()
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
			stdin:   stdin,
			stdout:  stdout,
			stderr:  stderrPipeReader,
		}, nil
	}
}

func (l *localTransport) startPluginSessionAsSystem(ctx context.Context, path string) (plugin.Session, error) {
	args := []string{"/c", fmt.Sprintf("gsudo.exe -s %q", path)}

	cmd := exec.CommandContext(ctx, "cmd.exe", args...)

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

	err = cmd.Start()
	if err != nil {
		stderrBytes, _ := io.ReadAll(stderr)
		return nil, fmt.Errorf("failed to start plugin at '%s': %w - %s", path, err, strings.TrimSpace(string(stderrBytes)))
	}

	return &localPluginSession{
		command: cmd,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
	}, nil
}
