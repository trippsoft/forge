// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"bytes"
	"io"
	"sync"
	"testing"
)

func TestSecretFilter_Secrets(t *testing.T) {
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  io.Discard,
	}

	filter.AddSecret("secret1")
	filter.AddSecret("secret2")

	secrets := filter.Secrets()

	for _, secret := range []string{"secret1", "secret2"} {
		found := false
		for _, s := range secrets {
			if s == secret {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected secret %q to be in the secrets list", secret)
		}
	}

	if len(secrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(secrets))
	}
}

func TestSecretFilter_AddSecret(t *testing.T) {
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  io.Discard,
	}

	// Test adding valid secret
	filter.AddSecret("password123")
	if !filter.secrets.Contains("password123") {
		t.Error("Expected secret to be added")
	}

	// Test adding empty secret
	filter.AddSecret("")
	if filter.secrets.Contains("") {
		t.Error("Expected empty secret to not be added")
	}
}

func TestSecretFilter_SetOutput(t *testing.T) {
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  io.Discard,
	}

	buf := &bytes.Buffer{}
	filter.SetOutput(buf)
	if filter.writer != buf {
		t.Error("Expected writer to be set")
	}
}

func TestSecretFilter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  buf,
	}

	filter.AddSecret("secret123")
	filter.AddSecret("password")

	input := []byte("User logged in with password secret123")
	n, err := filter.Write(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "User logged in with <redacted> <redacted>"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}

	if n != len(expected) {
		t.Errorf("Expected n=%d, got %d", len(expected), n)
	}
}

func TestSecretFilter_Filter(t *testing.T) {
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  io.Discard,
	}

	filter.AddSecret("mysecret")
	filter.AddSecret("token123")

	input := "API call with token123 and mysecret values"
	result := filter.Filter(input)
	expected := "API call with <redacted> and <redacted> values"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSecretFilter_ConcurrentAccess(t *testing.T) {
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  io.Discard,
	}

	var wg sync.WaitGroup

	// Concurrent writes
	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			filter.AddSecret("secret" + string(rune(i+'0')))
		}(i)
	}

	// Concurrent filters
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filter.Filter("test message with secret5")
		}()
	}

	wg.Wait()
}

func TestSecretFilter_EmptySecrets(t *testing.T) {
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  io.Discard,
	}

	input := "No secrets here"
	result := filter.Filter(input)
	if result != input {
		t.Errorf("Expected %q, got %q", input, result)
	}
}

func TestSecretFilter_WriteError(t *testing.T) {
	errorWriter := &errorWriter{}
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  errorWriter,
	}

	_, err := filter.Write([]byte("test"))
	if err == nil {
		t.Error("Expected error from writer")
	}
}

func TestSecretFilter_Clear(t *testing.T) {
	filter := &secretFilter{
		secrets: NewSet[string](),
		writer:  io.Discard,
	}

	filter.AddSecret("secret")
	filter.Clear()

	if filter.secrets.Size() != 0 {
		t.Error("Expected secrets to be cleared")
	}
}

type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrShortWrite
}
