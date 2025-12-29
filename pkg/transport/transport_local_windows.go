// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package transport

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/trippsoft/forge/pkg/plugin"
)

func (l *localTransport) startEscalatedPlugin(
	ctx context.Context,
	path string,
	escalation *Escalation,
) (*exec.Cmd, uint16, error) {

	user := escalation.User()
	if user == "" || user == "SYSTEM" || user == `NT AUTHORITY\SYSTEM` {
		return l.startPluginAsSystem(ctx, path)
	}

	args := []string{"/c", fmt.Sprintf("gsudo.exe -u %s %q", user, path)}

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cmd.exe", args...)

	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MIN_PORT=%d", plugin.LocalPluginMinPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MAX_PORT=%d", plugin.LocalPluginMaxPort))

	var errBuf bytes.Buffer

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get stdout pipe for plugin at '%s': %w", path, err)
	}
	defer stdoutReader.Close()

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get stderr pipe for plugin at '%s': %w", path, err)
	}
	defer stderrReader.Close()

	stdinWriter, err := cmd.StdinPipe()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get stdin pipe for plugin at '%s': %w", path, err)
	}
	defer stdinWriter.Close()

	teeReader := io.TeeReader(stderrReader, &errBuf)

	outputChannel := make(chan string)

	go func() {
		scanner := bufio.NewScanner(teeReader)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, forgeGSudoPrompt) {
				_, err = stdinWriter.Write([]byte(escalation.Pass() + "\n"))
				return
			}
		}
	}()

	go func() {
		defer close(outputChannel)
		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			line := scanner.Text()
			outputChannel <- line
			return
		}
	}()

	err = cmd.Start()
	if err != nil {
		stderr := strings.TrimSpace(errBuf.String())
		return nil, 0, fmt.Errorf("failed to start plugin at '%s': %w - %s", path, err, stderr)
	}

	select {
	case <-ctx.Done():
		return nil, 0, fmt.Errorf("context cancelled while starting plugin at '%s': %w", path, ctx.Err())
	case portOutput := <-outputChannel:
		port, err := strconv.ParseUint(portOutput, 10, 16)
		if err != nil {
			stderr := strings.TrimSpace(errBuf.String())
			return nil, 0, fmt.Errorf("invalid port output from plugin at '%s': %w - %s", path, err, stderr)
		}

		return cmd, uint16(port), nil
	}
}

func (l *localTransport) startPluginAsSystem(ctx context.Context, path string) (*exec.Cmd, uint16, error) {
	args := []string{"/c", fmt.Sprintf("gsudo.exe -s %q", path)}

	var errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, "cmd.exe", args...)
	cmd.Stderr = &errBuf

	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MIN_PORT=%d", plugin.LocalPluginMinPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("FORGE_PLUGIN_MAX_PORT=%d", plugin.LocalPluginMaxPort))

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get stdout pipe for plugin at '%s': %w", path, err)
	}
	defer stdoutReader.Close()

	err = cmd.Start()
	if err != nil {
		stderr := strings.TrimSpace(errBuf.String())
		return nil, 0, fmt.Errorf("failed to start plugin at '%s': %w - %s", path, err, stderr)
	}

	scanner := bufio.NewScanner(stdoutReader)
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
		stderr := strings.TrimSpace(errBuf.String())
		return nil, 0, fmt.Errorf("invalid port output from plugin at '%s': %w - %s", path, err, stderr)
	}

	return cmd, uint16(port), nil
}
