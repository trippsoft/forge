// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/customdecode"
	"github.com/zclconf/go-cty/cty"
)

func getTemplateFileTestCases() []struct {
	name         string
	templateText string
	contextVars  map[string]cty.Value
	expectedText string
} {
	return []struct {
		name         string
		templateText string
		contextVars  map[string]cty.Value
		expectedText string
	}{
		{
			name:         "simple text template",
			templateText: "Hello, World!",
			contextVars:  map[string]cty.Value{},
			expectedText: "Hello, World!",
		},
		{
			name:         "template with variable interpolation",
			templateText: "Hello, ${name}!",
			contextVars: map[string]cty.Value{
				"name": cty.StringVal("Alice"),
			},
			expectedText: "Hello, Alice!",
		},
		{
			name:         "empty template",
			templateText: "",
			contextVars:  map[string]cty.Value{},
			expectedText: "",
		},
		{
			name:         "template with multiple variables",
			templateText: "${greeting} ${name}, you are ${age} years old.",
			contextVars: map[string]cty.Value{
				"greeting": cty.StringVal("Hello"),
				"name":     cty.StringVal("Bob"),
				"age":      cty.NumberIntVal(30),
			},
			expectedText: "Hello Bob, you are 30 years old.",
		},
		{
			name:         "template with conditional",
			templateText: "%{ if enabled }Feature is ON%{ else }Feature is OFF%{ endif }",
			contextVars: map[string]cty.Value{
				"enabled": cty.True,
			},
			expectedText: "Feature is ON",
		},
		{
			name:         "template with for loop",
			templateText: "%{ for item in items }${item},%{ endfor }",
			contextVars: map[string]cty.Value{
				"items": cty.ListVal([]cty.Value{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				}),
			},
			expectedText: "a,b,c,",
		},
	}
}

func createTempTemplateFile(t *testing.T, dir string, content string) string {
	t.Helper()

	file, err := os.CreateTemp(dir, "template_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	err = file.Sync()
	if err != nil {
		t.Fatalf("Failed to sync temp file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return file.Name()
}

func createEvalContext(vars map[string]cty.Value) *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: vars,
	}
}

type mockExpr struct {
	value cty.Value
	err   hcl.Diagnostics
}

func (m *mockExpr) Value(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return m.value, m.err
}

func (m *mockExpr) Variables() []hcl.Traversal {
	return nil
}

func (m *mockExpr) Range() hcl.Range {
	return hcl.Range{}
}

func (m *mockExpr) StartRange() hcl.Range {
	return hcl.Range{}
}

func createExpressionClosure(t *testing.T, path string, evalCtx *hcl.EvalContext) cty.Value {
	t.Helper()

	return customdecode.ExpressionClosureVal(
		&customdecode.ExpressionClosure{
			Expression: &mockExpr{
				value: cty.StringVal(path),
			},
			EvalContext: evalCtx,
		},
	)
}

func TestTemplateFile(t *testing.T) {
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	tests := getTemplateFileTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTempTemplateFile(t, tempDir, tt.templateText)
			evalCtx := createEvalContext(tt.contextVars)
			closure := createExpressionClosure(t, path, evalCtx)

			result, err := TemplateFile(closure)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !result.IsKnown() || result.IsNull() {
				t.Fatalf("expected known non-null string, got %v", result)
			}

			actual := result.AsString()
			if actual != tt.expectedText {
				t.Errorf("expected %q, got %q", tt.expectedText, actual)
			}
		})
	}
}

func TestTemplateFile_NonExistentFile(t *testing.T) {
	evalCtx := createEvalContext(map[string]cty.Value{})
	closure := createExpressionClosure(t, "non_existent_template.tpl", evalCtx)

	_, err := TemplateFile(closure)
	if err == nil {
		t.Fatal("expected an error for non-existent file, got nil")
	}

	expectedErr := "failed to read template file"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error containing %q, got %q", expectedErr, err.Error())
	}
}

func TestTemplateFile_InvalidTemplate(t *testing.T) {
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	// Template with invalid syntax (unclosed interpolation)
	invalidTemplate := "Hello, ${name"
	path := createTempTemplateFile(t, tempDir, invalidTemplate)

	evalCtx := createEvalContext(map[string]cty.Value{
		"name": cty.StringVal("Alice"),
	})
	closure := createExpressionClosure(t, path, evalCtx)

	_, err := TemplateFile(closure)
	if err == nil {
		t.Fatal("expected an error for invalid template, got nil")
	}

	expectedErr := "failed to parse template file"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error containing %q, got %q", expectedErr, err.Error())
	}
}

func TestTemplateFileFunc(t *testing.T) {
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	tests := getTemplateFileTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTempTemplateFile(t, tempDir, tt.templateText)
			evalCtx := createEvalContext(tt.contextVars)
			closure := createExpressionClosure(t, path, evalCtx)

			result, err := TemplateFileFunc.Call([]cty.Value{closure})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !result.IsKnown() || result.IsNull() {
				t.Fatalf("expected known non-null string, got %v", result)
			}

			actual := result.AsString()
			if actual != tt.expectedText {
				t.Errorf("expected %q, got %q", tt.expectedText, actual)
			}
		})
	}
}

func TestTemplateFileFunc_NonExistentFile(t *testing.T) {
	evalCtx := createEvalContext(map[string]cty.Value{})
	closure := createExpressionClosure(t, "non_existent_template.tpl", evalCtx)

	_, err := TemplateFileFunc.Call([]cty.Value{closure})
	if err == nil {
		t.Fatal("expected an error for non-existent file, got nil")
	}

	expectedErr := "failed to read template file"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error containing %q, got %q", expectedErr, err.Error())
	}
}

func TestTemplateFileFunc_InvalidTemplate(t *testing.T) {
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	// Template with invalid syntax (unclosed interpolation)
	invalidTemplate := "Hello, ${name"
	path := createTempTemplateFile(t, tempDir, invalidTemplate)

	evalCtx := createEvalContext(map[string]cty.Value{
		"name": cty.StringVal("Alice"),
	})
	closure := createExpressionClosure(t, path, evalCtx)

	_, err := TemplateFileFunc.Call([]cty.Value{closure})
	if err == nil {
		t.Fatal("expected an error for invalid template, got nil")
	}

	expectedErr := "failed to parse template file"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Fatalf("expected error containing %q, got %q", expectedErr, err.Error())
	}
}
