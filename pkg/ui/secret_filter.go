// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package ui

import (
	"bytes"
	"io"
	"strings"
	"sync"

	"github.com/trippsoft/forge/pkg/util"
)

var (
	SecretFilter = &secretFilter{
		secrets: util.NewSet[string](),
		writer:  io.Discard,
	}
)

// secretFilter wraps an io.Writer to filter out sensitive information from logs and output.
type secretFilter struct {
	mutex   sync.RWMutex
	secrets *util.Set[string]
	writer  io.Writer
}

// Secrets returns a slice of all secrets currently being filtered.
func (f *secretFilter) Secrets() []string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.secrets.Items()
}

// AddSecret adds a secret to be filtered from output.
func (f *secretFilter) AddSecret(secret string) {
	if secret == "" {
		return
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.secrets.Add(secret)
}

// SetOutput sets the io.Writer where filtered output will be written.
func (f *secretFilter) SetOutput(writer io.Writer) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.writer = writer
}

// Write implements io.Writer, filtering out any added secrets before writing to the underlying writer.
func (f *secretFilter) Write(p []byte) (n int, err error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	for _, secret := range f.secrets.Items() {
		if secret != "" {
			p = bytes.ReplaceAll(p, []byte(secret), []byte("<redacted>"))
		}
	}

	return f.writer.Write(p)
}

// Filter filters out any added secrets from the provided message string.
func (f *secretFilter) Filter(message string) string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	for _, secret := range f.secrets.Items() {
		if secret != "" {
			message = strings.ReplaceAll(message, secret, "<redacted>")
		}
	}

	return message
}

// Clear removes all secrets from the filter.
func (f *secretFilter) Clear() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.secrets.Clear()
}
