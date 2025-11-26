// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"bytes"
	"io"
	"strings"
	"sync"
)

type SecretFilter struct {
	mutex   sync.Mutex
	secrets *Set[string]
	writer  io.Writer
}

func (f *SecretFilter) Secrets() []string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.secrets.Items()
}

func (f *SecretFilter) AddSecret(secret string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if secret == "" {
		return
	}

	f.secrets.Add(secret)
}

func (f *SecretFilter) SetOutput(writer io.Writer) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.writer = writer
}

func (f *SecretFilter) Write(p []byte) (n int, err error) {
	for _, secret := range f.secrets.Items() {
		if secret != "" {
			p = bytes.ReplaceAll(p, []byte(secret), []byte("<redacted>"))
		}
	}

	return f.writer.Write(p)
}

func (f *SecretFilter) Filter(message string) string {
	for _, secret := range f.secrets.Items() {
		if secret != "" {
			message = strings.ReplaceAll(message, secret, "<redacted>")
		}
	}

	return message
}

func (f *SecretFilter) Clear() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.secrets.Clear()
}
