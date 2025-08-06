package shell

import (
	"fmt"
	"testing"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestPluginInputSpec(t *testing.T) {

	plugin := &Plugin{}

	spec := plugin.InputSpec()

	if spec == nil {
		t.Fatal("Expected non-nil input spec")
	}
}

func TestPluginValidate(t *testing.T) {

	plugin := &Plugin{}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	err := plugin.Validate(input)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestPluginRun(t *testing.T) {

	command := "echo 'Hello, World!'"

	expectedStdout := "Hello, World!"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[command] = &transport.MockCmd{
		Stdout: fmt.Sprintf("%s\n", expectedStdout),
	}

	host := inventory.NewHost("linux", mockTransport, map[string]cty.Value{})

	p := &Plugin{}

	commonConfig := &plugin.CommonConfig{
		EscalateConfig: &transport.NoEscalate{},
		Timeout:        10,
	}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	result := p.Run(host, commonConfig, input)

	if result.Err != nil {
		t.Fatalf("Expected no error, got: %v", result.Err)
	}

	if !result.Changed {
		t.Fatal("Expected plugin to indicate changes were made")
	}

	if result.Output == nil {
		t.Fatal("Expected non-nil output")
	}

	if len(result.Output) != 2 {
		t.Fatalf("Expected output to have 2 keys, got: %d", len(result.Output))
	}

	if _, ok := result.Output["stdout"]; !ok {
		t.Fatal("Expected output to contain 'stdout' key")
	}

	if _, ok := result.Output["stderr"]; !ok {
		t.Fatal("Expected output to contain 'stderr' key")
	}

	actualStdout := result.Output["stdout"].AsString()
	if actualStdout != expectedStdout {
		t.Fatalf("Expected stdout %q, got %q", expectedStdout, actualStdout)
	}

	actualStderr := result.Output["stderr"].AsString()
	if actualStderr != "" {
		t.Fatalf("Expected stderr to be empty, got %q", actualStderr)
	}
}
