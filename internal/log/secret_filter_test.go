package log

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/trippsoft/forge/internal/util"
)

func TestLogSecretFilter_AddSecret(t *testing.T) {
	filter := &secretFilter{
		secrets: util.NewSet[string](),
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

func TestLogSecretFilter_SetOutput(t *testing.T) {
	filter := &secretFilter{
		secrets: util.NewSet[string](),
		writer:  io.Discard,
	}

	buf := &bytes.Buffer{}
	filter.SetOutput(buf)

	if filter.writer != buf {
		t.Error("Expected writer to be set")
	}
}

func TestLogSecretFilter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	filter := &secretFilter{
		secrets: util.NewSet[string](),
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

func TestLogSecretFilter_Filter(t *testing.T) {
	filter := &secretFilter{
		secrets: util.NewSet[string](),
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

func TestLogSecretFilter_ConcurrentAccess(t *testing.T) {
	filter := &secretFilter{
		secrets: util.NewSet[string](),
		writer:  io.Discard,
	}

	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			filter.AddSecret("secret" + string(rune(i+'0')))
		}(i)
	}

	// Concurrent filters
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filter.Filter("test message with secret5")
		}()
	}

	wg.Wait()
}

func TestLogSecretFilter_EmptySecrets(t *testing.T) {
	filter := &secretFilter{
		secrets: util.NewSet[string](),
		writer:  io.Discard,
	}

	input := "No secrets here"
	result := filter.Filter(input)

	if result != input {
		t.Errorf("Expected %q, got %q", input, result)
	}
}

func TestLogSecretFilter_WriteError(t *testing.T) {
	errorWriter := &errorWriter{}
	filter := &secretFilter{
		secrets: util.NewSet[string](),
		writer:  errorWriter,
	}

	_, err := filter.Write([]byte("test"))
	if err == nil {
		t.Error("Expected error from writer")
	}
}

type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrShortWrite
}
