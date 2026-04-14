// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package powershell

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// Shell defines the interface for executing PowerShell commands and managing the PowerShell session.
type Shell interface {
	// Execute runs a PowerShell command and returns the standard output, standard error, and any execution error.
	Execute(cmd string) (string, string, error)
	// Exit terminates the PowerShell session and releases any associated resources.
	Exit()
}

type localShell struct {
	cmd    *exec.Cmd
	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader
}

// Execute implements [Shell].
func (l *localShell) Execute(cmd string) (string, string, error) {
	if l.cmd == nil {
		return "", "", errors.New("PowerShell session is not active")
	}

	outMarker, err := createMarker()
	if err != nil {
		return "", "", fmt.Errorf("failed to create output marker: %w", err)
	}

	errMarker, err := createMarker()
	if err != nil {
		return "", "", fmt.Errorf("failed to create error marker: %w", err)
	}

	fullCmd := fmt.Sprintf("%s; echo '%s'; [Console]::Error.WriteLine('%s')\r\n", cmd, outMarker, errMarker)

	_, err = l.stdin.Write([]byte(fullCmd))
	if err != nil {
		return "", "", fmt.Errorf("failed to write command to PowerShell stdin: %w", err)
	}

	stdout := ""
	stderr := ""

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go readUntilMarker(l.stdout, outMarker, &stdout, wg)
	go readUntilMarker(l.stderr, errMarker, &stderr, wg)

	wg.Wait()

	if len(stderr) > 0 {
		return stdout, stderr, fmt.Errorf("PowerShell command execution error: %s", stderr)
	}

	return stdout, stderr, nil
}

// Exit implements [Shell].
func (l *localShell) Exit() {
	l.stdin.Write([]byte("exit\r\n"))

	closer, ok := l.stdout.(io.Closer)
	if ok {
		closer.Close()
	}

	closer, ok = l.stderr.(io.Closer)
	if ok {
		closer.Close()
	}

	l.cmd.Wait()

	l.cmd = nil
	l.stdin = nil
	l.stdout = nil
	l.stderr = nil
}

// NewLocal creates a new local PowerShell session and returns a Shell interface for executing commands.
func NewLocal() (Shell, error) {
	cmd := exec.Command("powershell.exe", "-NoExit", "-Command", "-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("Could not get standard input stream from PowerShell command: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Could not get standard output stream from PowerShell command: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Could not get standard error stream from PowerShell command: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Failed to start PowerShell command: %w", err)
	}

	return &localShell{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func createMarker() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes for marker string: %w", err)
	}

	return fmt.Sprintf("$MARKER%x$", b), nil
}

func readUntilMarker(r io.Reader, marker string, buffer *string, wg *sync.WaitGroup) error {
	defer wg.Done()

	output := ""
	bufferSize := 64
	marker = fmt.Sprintf("%s\r\n", marker)

	for {
		buf := make([]byte, bufferSize)
		read, err := r.Read(buf)
		if err != nil {
			return err
		}

		output += string(buf[:read])

		if strings.HasSuffix(output, marker) {
			break
		}
	}

	*buffer = strings.TrimSuffix(output, marker)

	return nil
}
