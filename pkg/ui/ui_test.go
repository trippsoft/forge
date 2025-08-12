package ui

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/trippsoft/forge/pkg/log"
)

func TestPrint(t *testing.T) {

	tests := []struct {
		name     string
		message  string
		secrets  []string
		expected string
	}{
		{
			name:     "simple message without color",
			message:  "Hello World",
			secrets:  []string{},
			expected: "Hello World",
		},
		{
			name:     "message with secret filtering",
			message:  "Password is secret123",
			secrets:  []string{"secret123"},
			expected: "Password is <redacted>",
		},
		{
			name:     "message with multiple secrets",
			message:  "User admin password secret123",
			secrets:  []string{"admin", "secret123"},
			expected: "User <redacted> password <redacted>",
		},
		{
			name:     "empty message",
			message:  "",
			secrets:  []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, true)

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

func TestPrintWithFormat(t *testing.T) {

	tests := []struct {
		name         string
		message      string
		color        bool
		args         []TextArgument
		leftPadding  int
		rightPadding int
		secrets      []string
		expected     string
	}{
		{
			name:         "simple message without color",
			message:      "Hello World",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "Hello World",
		},
		{
			name:         "message with left padding",
			message:      "Padded",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  4,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "    Padded",
		},
		{
			name:         "message with right padding",
			message:      "Right",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 6,
			secrets:      []string{},
			expected:     "Right      ",
		},
		{
			name:         "message with both paddings",
			message:      "Both",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  2,
			rightPadding: 3,
			secrets:      []string{},
			expected:     "  Both   ",
		},
		{
			name:         "colored message with single argument",
			message:      "Colored Text",
			color:        true,
			args:         []TextArgument{31}, // Red
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "\033[31mColored Text\033[0m",
		},
		{
			name:         "colored message with multiple arguments",
			message:      "Bold Red",
			color:        true,
			args:         []TextArgument{1, 31}, // Bold + Red
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "\033[1;31mBold Red\033[0m",
		},
		{
			name:         "colored message with padding",
			message:      "Padded Color",
			color:        true,
			args:         []TextArgument{32}, // Green
			leftPadding:  2,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "\033[32m  Padded Color\033[0m",
		},
		{
			name:         "message with color disabled but args provided",
			message:      "No Color",
			color:        false,
			args:         []TextArgument{31}, // Should be ignored
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "No Color",
		},
		{
			name:         "message with secret filtering",
			message:      "Password is secret123",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{"secret123"},
			expected:     "Password is <redacted>",
		},
		{
			name:         "message with multiple secrets",
			message:      "User admin password secret123",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{"admin", "secret123"},
			expected:     "User <redacted> password <redacted>",
		},
		{
			name:         "colored message with secret filtering",
			message:      "API key: apikey123",
			color:        true,
			args:         []TextArgument{33}, // Yellow
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{"apikey123"},
			expected:     "\033[33mAPI key: <redacted>\033[0m",
		},
		{
			name:         "empty message",
			message:      "",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "",
		},
		{
			name:         "empty message with padding",
			message:      "",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  3,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			log.SecretFilter.Clear()
			for _, secret := range tt.secrets {
				log.SecretFilter.AddSecret(secret)
			}

			ui.PrintWithFormat(tt.message, TextFormatting{
				Args:         tt.args,
				LeftPadding:  tt.leftPadding,
				RightPadding: tt.rightPadding,
			})

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
		secrets  []string
		expected string
	}{
		{
			name:     "simple message without color",
			message:  "Hello World",
			secrets:  []string{},
			expected: "Hello World",
		},
		{
			name:     "message with secret filtering",
			message:  "Password is secret123",
			secrets:  []string{"secret123"},
			expected: "Password is <redacted>",
		},
		{
			name:     "message with multiple secrets",
			message:  "User admin password secret123",
			secrets:  []string{"admin", "secret123"},
			expected: "User <redacted> password <redacted>",
		},
		{
			name:     "empty message",
			message:  "",
			secrets:  []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, true)

			log.SecretFilter.Clear()
			for _, secret := range tt.secrets {
				log.SecretFilter.AddSecret(secret)
			}

			ui.Error(tt.message)

			if errBuf.String() != tt.expected {
				t.Errorf("Expected output: %q, got: %q", tt.expected, errBuf.String())
			}
		})
	}
}

func TestErrorWithFormat(t *testing.T) {

	tests := []struct {
		name         string
		message      string
		color        bool
		args         []TextArgument
		leftPadding  int
		rightPadding int
		secrets      []string
		expected     string
	}{
		{
			name:         "simple message without color",
			message:      "Hello World",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "Hello World",
		},
		{
			name:         "message with left padding",
			message:      "Padded",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  4,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "    Padded",
		},
		{
			name:         "message with right padding",
			message:      "Right",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 6,
			secrets:      []string{},
			expected:     "Right      ",
		},
		{
			name:         "message with both paddings",
			message:      "Both",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  2,
			rightPadding: 3,
			secrets:      []string{},
			expected:     "  Both   ",
		},
		{
			name:         "colored message with single argument",
			message:      "Colored Text",
			color:        true,
			args:         []TextArgument{31}, // Red
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "\033[31mColored Text\033[0m",
		},
		{
			name:         "colored message with multiple arguments",
			message:      "Bold Red",
			color:        true,
			args:         []TextArgument{1, 31}, // Bold + Red
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "\033[1;31mBold Red\033[0m",
		},
		{
			name:         "colored message with padding",
			message:      "Padded Color",
			color:        true,
			args:         []TextArgument{32}, // Green
			leftPadding:  2,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "\033[32m  Padded Color\033[0m",
		},
		{
			name:         "message with color disabled but args provided",
			message:      "No Color",
			color:        false,
			args:         []TextArgument{31}, // Should be ignored
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "No Color",
		},
		{
			name:         "message with secret filtering",
			message:      "Password is secret123",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{"secret123"},
			expected:     "Password is <redacted>",
		},
		{
			name:         "message with multiple secrets",
			message:      "User admin password secret123",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{"admin", "secret123"},
			expected:     "User <redacted> password <redacted>",
		},
		{
			name:         "colored message with secret filtering",
			message:      "API key: apikey123",
			color:        true,
			args:         []TextArgument{33}, // Yellow
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{"apikey123"},
			expected:     "\033[33mAPI key: <redacted>\033[0m",
		},
		{
			name:         "empty message",
			message:      "",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  0,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "",
		},
		{
			name:         "empty message with padding",
			message:      "",
			color:        false,
			args:         []TextArgument{},
			leftPadding:  3,
			rightPadding: 0,
			secrets:      []string{},
			expected:     "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			log.SecretFilter.Clear()
			for _, secret := range tt.secrets {
				log.SecretFilter.AddSecret(secret)
			}

			ui.ErrorWithFormat(tt.message, TextFormatting{
				Args:         tt.args,
				LeftPadding:  tt.leftPadding,
				RightPadding: tt.rightPadding,
			})

			if errBuf.String() != tt.expected {
				t.Errorf("Expected output: %q, got: %q", tt.expected, errBuf.String())
			}
		})
	}
}

func TestPrintLine(t *testing.T) {

	tests := []struct {
		name         string
		character    rune
		color        bool
		leftPadding  int
		rightPadding int
		expected     string
	}{
		{
			name:      "hyphen",
			character: '-',
			expected:  fmt.Sprintf("%s\n", strings.Repeat("-", 78)),
		},
		{
			name:      "equal sign",
			character: '=',
			expected:  fmt.Sprintf("%s\n", strings.Repeat("=", 78)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, false)

			ui.PrintLine(tt.character)

			if outBuf.String() != tt.expected {
				t.Errorf("Expected output: %q, got: %q", tt.expected, outBuf.String())
			}
		})
	}
}

func TestPrintLineWithFormat(t *testing.T) {

	tests := []struct {
		name         string
		character    rune
		color        bool
		args         []TextArgument
		leftPadding  int
		rightPadding int
		expected     string
	}{
		{
			name:         "hyphen",
			character:    '-',
			color:        false,
			leftPadding:  0,
			rightPadding: 0,
			expected:     fmt.Sprintf("%s\n", strings.Repeat("-", 78)),
		},
		{
			name:         "equal sign",
			character:    '=',
			color:        false,
			leftPadding:  0,
			rightPadding: 0,
			expected:     fmt.Sprintf("%s\n", strings.Repeat("=", 78)),
		},
		{
			name:         "hyphen with padding",
			character:    '-',
			color:        false,
			leftPadding:  2,
			rightPadding: 2,
			expected:     fmt.Sprintf("  %s  \n", strings.Repeat("-", 74)),
		},
		{
			name:         "colored asterisk",
			character:    '*',
			color:        true,
			args:         []TextArgument{1, 31}, // Bold + Red
			leftPadding:  0,
			rightPadding: 0,
			expected:     fmt.Sprintf("\033[1;31m%s\033[0m\n", strings.Repeat("*", 78)),
		},
		{
			name:         "colored equal sign with padding",
			character:    '=',
			color:        true,
			args:         []TextArgument{32}, // Green
			leftPadding:  2,
			rightPadding: 0,
			expected:     fmt.Sprintf("\033[32m  %s\033[0m\n", strings.Repeat("=", 76)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			ui.PrintLineWithFormat(tt.character, TextFormatting{
				Args:         tt.args,
				LeftPadding:  tt.leftPadding,
				RightPadding: tt.rightPadding,
			})

			if outBuf.String() != tt.expected {
				t.Errorf("Expected output: %q, got: %q", tt.expected, outBuf.String())
			}
		})
	}
}

func TestPrintColumns(t *testing.T) {

	tests := []struct {
		name            string
		leftMessage     string
		leftFormatting  TextFormatting
		rightMessage    string
		rightFormatting TextFormatting
		color           bool
		secrets         []string
		expected        string
	}{
		{
			name:            "simple two columns",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("Left%sRight\n", strings.Repeat(" ", 69)),
		},
		{
			name:            "columns with left padding",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{LeftPadding: 2},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("  Left%sRight\n", strings.Repeat(" ", 67)),
		},
		{
			name:            "columns with right padding",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{RightPadding: 3},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("Left%sRight   \n", strings.Repeat(" ", 66)),
		},
		{
			name:            "columns with both left and right padding",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{LeftPadding: 2, RightPadding: 1},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{LeftPadding: 1, RightPadding: 2},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("  Left %s Right  \n", strings.Repeat(" ", 63)),
		},
		{
			name:            "colored left column",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{Args: []TextArgument{ForegroundColorRed}},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{},
			color:           true,
			secrets:         []string{},
			expected:        fmt.Sprintf("\033[31mLeft\033[0m%sRight\n", strings.Repeat(" ", 69)),
		},
		{
			name:            "colored right column",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{Args: []TextArgument{ForegroundColorGreen}},
			color:           true,
			secrets:         []string{},
			expected:        fmt.Sprintf("Left%s\033[32mRight\033[0m\n", strings.Repeat(" ", 69)),
		},
		{
			name:            "both columns colored",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{Args: []TextArgument{Bold, ForegroundColorRed}},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{Args: []TextArgument{Bold, ForegroundColorGreen}},
			color:           true,
			secrets:         []string{},
			expected:        fmt.Sprintf("\033[1;31mLeft\033[0m%s\033[1;32mRight\033[0m\n", strings.Repeat(" ", 69)),
		},
		{
			name:            "colored columns with padding",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{Args: []TextArgument{31}, LeftPadding: 1},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{Args: []TextArgument{32}, RightPadding: 1},
			color:           true,
			secrets:         []string{},
			expected:        fmt.Sprintf("\033[31m Left\033[0m%s\033[32mRight \033[0m\n", strings.Repeat(" ", 67)),
		},
		{
			name:            "color disabled with args provided",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{Args: []TextArgument{31}},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{Args: []TextArgument{32}},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("Left%sRight\n", strings.Repeat(" ", 69)),
		},
		{
			name:            "columns with secret filtering",
			leftMessage:     "User: admin",
			leftFormatting:  TextFormatting{},
			rightMessage:    "Pass: secret123",
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{"admin", "secret123"},
			expected:        fmt.Sprintf("User: <redacted>%sPass: <redacted>\n", strings.Repeat(" ", 46)),
		},
		{
			name:            "colored columns with secret filtering",
			leftMessage:     "API Key",
			leftFormatting:  TextFormatting{Args: []TextArgument{ForegroundColorYellow}},
			rightMessage:    "apikey123",
			rightFormatting: TextFormatting{Args: []TextArgument{ForegroundColorRed}},
			color:           true,
			secrets:         []string{"apikey123"},
			expected:        fmt.Sprintf("\033[33mAPI Key\033[0m%s\033[31m<redacted>\033[0m\n", strings.Repeat(" ", 61)),
		},
		{
			name:            "empty left column",
			leftMessage:     "",
			leftFormatting:  TextFormatting{},
			rightMessage:    "Right",
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("%sRight\n", strings.Repeat(" ", 73)),
		},
		{
			name:            "empty right column",
			leftMessage:     "Left",
			leftFormatting:  TextFormatting{},
			rightMessage:    "",
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("Left%s\n", strings.Repeat(" ", 74)),
		},
		{
			name:            "both columns empty",
			leftMessage:     "",
			leftFormatting:  TextFormatting{},
			rightMessage:    "",
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("%s\n", strings.Repeat(" ", 78)),
		},
		{
			name:            "long messages that fill width",
			leftMessage:     strings.Repeat("A", 39),
			leftFormatting:  TextFormatting{},
			rightMessage:    strings.Repeat("B", 39),
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("%s%s\n", strings.Repeat("A", 39), strings.Repeat("B", 39)),
		},
		{
			name:            "messages that exceed width",
			leftMessage:     strings.Repeat("A", 40),
			leftFormatting:  TextFormatting{},
			rightMessage:    strings.Repeat("B", 40),
			rightFormatting: TextFormatting{},
			color:           false,
			secrets:         []string{},
			expected:        fmt.Sprintf("%s%s\n", strings.Repeat("A", 40), strings.Repeat("B", 40)), // No spaces when exceeding width
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var outBuf, errBuf bytes.Buffer
			ui := MockUI(&outBuf, &errBuf, tt.color)

			log.SecretFilter.Clear()
			for _, secret := range tt.secrets {
				log.SecretFilter.AddSecret(secret)
			}

			ui.PrintColumns(tt.leftMessage, tt.leftFormatting, tt.rightMessage, tt.rightFormatting)

			if outBuf.String() != tt.expected {
				t.Errorf("Expected output: %q, got: %q", tt.expected, outBuf.String())
			}
		})
	}
}
