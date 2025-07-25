package log

import (
	"bytes"
	"io"
	"strings"
	"sync"

	"github.com/trippsoft/forge/internal/util"
)

var (
	LogSecretFilter *secretFilter = &secretFilter{
		secrets: util.NewSet[string](),
		writer:  io.Discard,
	}
)

type secretFilter struct {
	mutex   sync.Mutex
	secrets *util.Set[string]
	writer  io.Writer
}

func (f *secretFilter) Secrets() []string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.secrets.Items()
}

func (f *secretFilter) AddSecret(secret string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if secret == "" {
		return
	}

	f.secrets.Add(secret)
}

func (f *secretFilter) SetOutput(writer io.Writer) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.writer = writer
}

func (f *secretFilter) Write(p []byte) (n int, err error) {

	for _, secret := range f.secrets.Items() {
		if secret != "" {
			p = bytes.ReplaceAll(p, []byte(secret), []byte("<redacted>"))
		}
	}

	return f.writer.Write(p)
}

func (f *secretFilter) Filter(message string) string {
	for _, secret := range f.secrets.Items() {
		if secret != "" {
			message = strings.ReplaceAll(message, secret, "<redacted>")
		}
	}

	return message
}

func (f *secretFilter) Clear() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.secrets.Clear()
}
