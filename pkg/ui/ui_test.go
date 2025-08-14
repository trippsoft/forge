// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/trippsoft/forge/pkg/log"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		color    bool
		secrets  []string
		expected string
	}{
		{
			name:     "simple message",
			message:  "Hello World",
			color:    false,
			secrets:  []string{},
			expected: "Hello World",
		},
		{
			name:     "message with secrets",
			message:  "Password is secret123",
			color:    false,
			secrets:  []string{"secret123"},
			expected: "Password is <redacted>",
		},
		{
			name:     "message with multiple secrets",
			message:  "User admin password secret123",
			color:    false,
			secrets:  []string{"admin", "secret123"},
			expected: "User <redacted> password <redacted>",
		},
		{
			name:     "empty message",
			message:  "",
			color:    false,
			secrets:  []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			// Setup secret filter
			log.SecretFilter.Clear()
			for _, secret := range tt.secrets {
				log.SecretFilter.AddSecret(secret)
			}

			ui.Print(tt.message)

			if outBuf.String() != tt.expected {
				t.Errorf("Expected output: %q, got: %q", tt.expected, outBuf.String())
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		color    bool
		secrets  []string
		expected string
	}{
		{
			name:     "simple error message",
			message:  "Error occurred",
			color:    false,
			secrets:  []string{},
			expected: "Error occurred",
		},
		{
			name:     "error message with secrets",
			message:  "Failed to connect with apikey123",
			color:    false,
			secrets:  []string{"apikey123"},
			expected: "Failed to connect with <redacted>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			// Setup secret filter
			log.SecretFilter.Clear()
			for _, secret := range tt.secrets {
				log.SecretFilter.AddSecret(secret)
			}

			ui.Error(tt.message)

			if errBuf.String() != tt.expected {
				t.Errorf("Expected error output: %q, got: %q", tt.expected, errBuf.String())
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		text     *uiText
		color    bool
		expected string
	}{
		{
			name:     "simple text without color",
			text:     Text("Hello"),
			color:    false,
			expected: "Hello",
		},
		{
			name:     "text with left padding",
			text:     Text("Padded").WithLeftPadding(4),
			color:    false,
			expected: "    Padded",
		},
		{
			name:     "text with right padding",
			text:     Text("Right").WithRightPadding(6),
			color:    false,
			expected: "Right      ",
		},
		{
			name:     "text with both paddings",
			text:     Text("Both").WithLeftPadding(2).WithRightPadding(3),
			color:    false,
			expected: "  Both   ",
		},
		{
			name:     "text with left margin",
			text:     Text("Margin").WithLeftMargin(3),
			color:    false,
			expected: "   Margin",
		},
		{
			name:     "text with right margin",
			text:     Text("Margin").WithRightMargin(2),
			color:    false,
			expected: "Margin  ",
		},
		{
			name:     "text with all spacing",
			text:     Text("Spaced").WithLeftMargin(1).WithLeftPadding(2).WithRightPadding(1).WithRightMargin(2),
			color:    false,
			expected: "   Spaced   ",
		},
		{
			name:     "colored text with foreground",
			text:     Text("Red Text").WithForegroundColor(ForegroundRed),
			color:    true,
			expected: "\033[31mRed Text\033[0m",
		},
		{
			name:     "colored text with background",
			text:     Text("Blue BG").WithBackgroundColor(BackgroundBlue),
			color:    true,
			expected: "\033[44mBlue BG\033[0m",
		},
		{
			name:     "colored text with foreground and background",
			text:     Text("Both Colors").WithForegroundColor(ForegroundWhite).WithBackgroundColor(BackgroundRed),
			color:    true,
			expected: "\033[41;37mBoth Colors\033[0m",
		},
		{
			name:     "text with single style",
			text:     Text("Bold").WithStyle(StyleBold),
			color:    true,
			expected: "\033[1mBold\033[0m",
		},
		{
			name:     "text with multiple styles",
			text:     Text("Bold Italic").WithStyle(StyleBold).WithStyle(StyleItalic),
			color:    true,
			expected: "\033[1;3mBold Italic\033[0m",
		},
		{
			name:     "text with all formatting",
			text:     Text("Complex").WithForegroundColor(ForegroundYellow).WithBackgroundColor(BackgroundBlue).WithStyle(StyleBold).WithStyle(StyleUnderline),
			color:    true,
			expected: "\033[44;33;1;4mComplex\033[0m",
		},
		{
			name:     "colored text with padding",
			text:     Text("Padded Color").WithForegroundColor(ForegroundGreen).WithLeftPadding(2),
			color:    true,
			expected: "\033[32m  Padded Color\033[0m",
		},
		{
			name:     "colored text with margin",
			text:     Text("Margin Color").WithForegroundColor(ForegroundCyan).WithLeftMargin(3),
			color:    true,
			expected: "   \033[36mMargin Color\033[0m",
		},
		{
			name:     "color disabled but formatting provided",
			text:     Text("No Color").WithForegroundColor(ForegroundRed).WithStyle(StyleBold),
			color:    false,
			expected: "No Color",
		},
		{
			name:     "empty text",
			text:     Text(""),
			color:    false,
			expected: "",
		},
		{
			name:     "empty text with padding",
			text:     Text("").WithLeftPadding(3).WithRightPadding(2),
			color:    false,
			expected: "     ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			result := ui.Format(tt.text)

			if result != tt.expected {
				t.Errorf("Expected: %q, got: %q", tt.expected, result)
			}
		})
	}
}

func TestFormatColumns(t *testing.T) {
	tests := []struct {
		name     string
		texts    []*uiText
		color    bool
		expected string
	}{
		{
			name:     "no texts",
			texts:    []*uiText{},
			color:    false,
			expected: "\n",
		},
		{
			name:     "single text",
			texts:    []*uiText{Text("Single")},
			color:    false,
			expected: "Single\n",
		},
		{
			name:     "two simple columns",
			texts:    []*uiText{Text("Left"), Text("Right")},
			color:    false,
			expected: "Left" + strings.Repeat(" ", 69) + "Right\n",
		},
		{
			name:     "two columns with padding",
			texts:    []*uiText{Text("Left").WithLeftPadding(2), Text("Right").WithRightPadding(3)},
			color:    false,
			expected: "  Left" + strings.Repeat(" ", 64) + "Right   \n",
		},
		{
			name:     "two columns with margins",
			texts:    []*uiText{Text("Left").WithLeftMargin(1), Text("Right").WithRightMargin(2)},
			color:    false,
			expected: " Left" + strings.Repeat(" ", 66) + "Right  \n",
		},
		{
			name:     "colored columns",
			texts:    []*uiText{Text("Red").WithForegroundColor(ForegroundRed), Text("Green").WithForegroundColor(ForegroundGreen)},
			color:    true,
			expected: "\033[31mRed\033[0m" + strings.Repeat(" ", 70) + "\033[32mGreen\033[0m\n",
		},
		{
			name:     "three columns",
			texts:    []*uiText{Text("One"), Text("Two"), Text("Three")},
			color:    false,
			expected: "One" + strings.Repeat(" ", 34) + "Two" + strings.Repeat(" ", 33) + "Three\n",
		},
		{
			name:     "columns exceeding width",
			texts:    []*uiText{Text(strings.Repeat("A", 40)), Text(strings.Repeat("B", 40))},
			color:    false,
			expected: strings.Repeat("A", 40) + strings.Repeat("B", 40) + "\n",
		},
		{
			name:     "empty columns",
			texts:    []*uiText{Text(""), Text("")},
			color:    false,
			expected: strings.Repeat(" ", 78) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			result := ui.FormatColumns(tt.texts...)

			if result != tt.expected {
				t.Errorf("Expected: %q, got: %q", tt.expected, result)
			}
		})
	}
}

func TestFormatLine(t *testing.T) {
	tests := []struct {
		name      string
		character rune
		format    *textFormat
		color     bool
		expected  string
	}{
		{
			name:      "simple line with dash",
			character: '-',
			format:    nil,
			color:     false,
			expected:  strings.Repeat("-", 78) + "\n",
		},
		{
			name:      "line with equals",
			character: '=',
			format:    TextFormat(),
			color:     false,
			expected:  strings.Repeat("=", 78) + "\n",
		},
		{
			name:      "line with left margin",
			character: '-',
			format:    TextFormat().WithLeftMargin(5),
			color:     false,
			expected:  "     " + strings.Repeat("-", 73) + "\n",
		},
		{
			name:      "line with padding",
			character: '*',
			format:    TextFormat().WithLeftPadding(3).WithRightPadding(2),
			color:     false,
			expected:  "   " + strings.Repeat("*", 73) + "  \n",
		},
		{
			name:      "line with all margins and padding",
			character: '#',
			format:    TextFormat().WithLeftMargin(2).WithLeftPadding(1).WithRightPadding(1).WithRightMargin(2),
			color:     false,
			expected:  "   " + strings.Repeat("#", 72) + "   \n",
		},
		{
			name:      "colored line",
			character: '~',
			format:    TextFormat().WithForegroundColor(ForegroundBlue),
			color:     true,
			expected:  "\033[34m" + strings.Repeat("~", 78) + "\033[0m\n",
		},
		{
			name:      "line with excessive margins (zero width)",
			character: '-',
			format:    TextFormat().WithLeftMargin(40).WithRightMargin(40),
			color:     false,
			expected:  "\n",
		},
		{
			name:      "line with color disabled",
			character: '+',
			format:    TextFormat().WithForegroundColor(ForegroundRed),
			color:     false,
			expected:  strings.Repeat("+", 78) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			result := ui.FormatLine(tt.character, tt.format)

			if result != tt.expected {
				t.Errorf("Expected: %q, got: %q", tt.expected, result)
			}
		})
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		secrets  []string
		expected string
	}{
		{
			name:     "simple text",
			message:  "Hello World",
			secrets:  []string{},
			expected: "Hello World",
		},
		{
			name:     "text with secret",
			message:  "API Key: secret123",
			secrets:  []string{"secret123"},
			expected: "API Key: <redacted>",
		},
		{
			name:     "empty text",
			message:  "",
			secrets:  []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup secret filter
			log.SecretFilter.Clear()
			for _, secret := range tt.secrets {
				log.SecretFilter.AddSecret(secret)
			}

			result := Text(tt.message)

			if result.message != tt.expected {
				t.Errorf("Expected message: %q, got: %q", tt.expected, result.message)
			}
		})
	}
}

func TestTextFormat(t *testing.T) {
	format := TextFormat()

	if format == nil {
		t.Error("Expected non-nil textFormat")
	}

	if format.styles == nil {
		t.Error("Expected non-nil styles slice")
	}

	if len(format.styles) != 0 {
		t.Error("Expected empty styles slice initially")
	}
}

func TestTextFormatChaining(t *testing.T) {
	format := TextFormat().
		WithBackgroundColor(BackgroundRed).
		WithForegroundColor(ForegroundWhite).
		WithStyle(StyleBold).
		WithStyle(StyleUnderline).
		WithLeftPadding(2).
		WithRightPadding(3).
		WithLeftMargin(1).
		WithRightMargin(4)

	if format.backgroundColor != BackgroundRed {
		t.Errorf("Expected background color %s, got %s", BackgroundRed, format.backgroundColor)
	}

	if format.foregroundColor != ForegroundWhite {
		t.Errorf("Expected foreground color %s, got %s", ForegroundWhite, format.foregroundColor)
	}

	if len(format.styles) != 2 {
		t.Errorf("Expected 2 styles, got %d", len(format.styles))
	}

	if format.styles[0] != StyleBold || format.styles[1] != StyleUnderline {
		t.Errorf("Expected styles [Bold, Underline], got %v", format.styles)
	}

	if format.leftPadding != 2 {
		t.Errorf("Expected left padding 2, got %d", format.leftPadding)
	}

	if format.rightPadding != 3 {
		t.Errorf("Expected right padding 3, got %d", format.rightPadding)
	}

	if format.leftMargin != 1 {
		t.Errorf("Expected left margin 1, got %d", format.leftMargin)
	}

	if format.rightMargin != 4 {
		t.Errorf("Expected right margin 4, got %d", format.rightMargin)
	}
}

func TestTextFormatClone(t *testing.T) {
	original := TextFormat().
		WithBackgroundColor(BackgroundGreen).
		WithForegroundColor(ForegroundBlack).
		WithStyle(StyleItalic).
		WithLeftPadding(5)

	cloned := original.Clone()

	// Modify original
	original.WithForegroundColor(ForegroundRed).WithStyle(StyleBold)

	// Check that clone is independent
	if cloned.foregroundColor != ForegroundBlack {
		t.Errorf("Clone should not be affected by changes to original")
	}

	if len(cloned.styles) != 1 || cloned.styles[0] != StyleItalic {
		t.Errorf("Clone styles should be independent")
	}
}

func TestUITextChaining(t *testing.T) {
	text := Text("Test").
		WithBackgroundColor(BackgroundYellow).
		WithForegroundColor(ForegroundBlue).
		WithStyle(StyleDim).
		WithLeftPadding(1).
		WithRightPadding(2).
		WithLeftMargin(3).
		WithRightMargin(4)

	if text.message != "Test" {
		t.Errorf("Expected message 'Test', got %s", text.message)
	}

	if text.backgroundColor != BackgroundYellow {
		t.Errorf("Expected background color %s, got %s", BackgroundYellow, text.backgroundColor)
	}

	if text.foregroundColor != ForegroundBlue {
		t.Errorf("Expected foreground color %s, got %s", ForegroundBlue, text.foregroundColor)
	}

	if len(text.styles) != 1 || text.styles[0] != StyleDim {
		t.Errorf("Expected styles [Dim], got %v", text.styles)
	}
}

func TestUITextWithFormat(t *testing.T) {
	format := TextFormat().WithForegroundColor(ForegroundMagenta).WithStyle(StyleBlink)
	text := Text("Formatted").WithFormat(format)

	if text.foregroundColor != ForegroundMagenta {
		t.Errorf("Expected foreground color %s, got %s", ForegroundMagenta, text.foregroundColor)
	}

	if len(text.styles) != 1 || text.styles[0] != StyleBlink {
		t.Errorf("Expected styles [Blink], got %v", text.styles)
	}

	// Test with nil format
	text2 := Text("Default").WithFormat(nil)
	if text2.textFormat == nil {
		t.Error("Expected non-nil textFormat when nil format provided")
	}
}
