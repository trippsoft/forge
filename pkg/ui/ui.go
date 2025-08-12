package ui

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/trippsoft/forge/pkg/log"
	"golang.org/x/term"
)

const (
	defaultConsoleWidth = 78 // We will size the console output to this width.
)

var (
	_ UI = (*ui)(nil)
)

// UI wraps the standard output and error streams to allow them to be
// used by multiple goroutines.
type UI interface {
	Print(message string) // Print outputs a message to the standard output.
	Error(message string) // Error outputs an error message to the standard error.

	// Format formats the text with the specified formatting.
	Format(text *uiText) string
	// FormatColumns formats the text into two columns.
	FormatColumns(text ...*uiText) string
	// FormatLine formats a character as a line with the specified formatting.
	FormatLine(character rune, format *textFormat) string
}

type ui struct {
	m sync.Mutex

	color bool

	out io.Writer
	err io.Writer
}

// StdUI returns a new UI that writes to the standard output and error.
func StdUI() UI {

	color := false
	if runtime.GOOS != "windows" {
		color = true
	}

	return &ui{
		color: color,
		out:   os.Stdout,
		err:   os.Stderr,
	}
}

// MockUI returns a new UI that writes to the provided output and error writers.
func MockUI(outWriter, errWriter io.Writer, color bool) UI {
	return &ui{
		out:   outWriter,
		err:   errWriter,
		color: color,
	}
}

// Print implements UI.
func (u *ui) Print(message string) {
	u.writeText(u.out, message)
}

// Error implements UI.
func (u *ui) Error(message string) {
	u.writeText(u.err, message)
}

// Format implements UI.
func (u *ui) Format(text *uiText) string {

	leftMargin := strings.Repeat(" ", text.leftMargin)
	leftPadding := strings.Repeat(" ", text.leftPadding)
	rightPadding := strings.Repeat(" ", text.rightPadding)
	rightMargin := strings.Repeat(" ", text.rightMargin)

	message := fmt.Sprintf("%s%s%s", leftPadding, text.message, rightPadding)

	if u.color {
		terminalArgs := string(text.backgroundColor)

		if terminalArgs == "" {
			terminalArgs = string(text.foregroundColor)
		} else if text.foregroundColor != "" {
			terminalArgs = fmt.Sprintf("%s;%s", terminalArgs, text.foregroundColor)
		}

		if len(text.styles) > 0 {
			for _, style := range text.styles {
				if terminalArgs == "" {
					terminalArgs = string(style)
				} else {
					terminalArgs = fmt.Sprintf("%s;%s", terminalArgs, style)
				}
			}
		}

		if terminalArgs != "" {
			message = fmt.Sprintf("\033[%sm%s\033[0m", terminalArgs, message)
		}
	}

	return fmt.Sprintf("%s%s%s", leftMargin, message, rightMargin)
}

// FormatColumns implements UI.
func (u *ui) FormatColumns(text ...*uiText) string {

	if len(text) == 0 {
		return "\n"
	}

	if len(text) == 1 {
		return fmt.Sprintf("%s\n", u.Format(text[0]))
	}

	totalRuneCount := 0
	for _, t := range text {
		totalRuneCount += getRuneCount(t)
	}

	spacing := (u.consoleWidth() - totalRuneCount) / (len(text) - 1)
	remainder := (u.consoleWidth() - totalRuneCount) % (len(text) - 1)

	message := ""
	if spacing <= 0 {
		for _, t := range text {
			message = fmt.Sprintf("%s%s", message, u.Format(t))
		}

		return fmt.Sprintf("%s\n", message)
	}

	for i, t := range text {
		message = fmt.Sprintf("%s%s", message, u.Format(t))
		if i >= len(text)-1 {
			break
		}

		message = fmt.Sprintf("%s%s", message, strings.Repeat(" ", spacing))

		if i < remainder {
			message = fmt.Sprintf("%s ", message) // Add an extra space for the remainder.
		}
	}

	return fmt.Sprintf("%s\n", message)
}

// FormatLine implements UI.
func (u *ui) FormatLine(character rune, format *textFormat) string {
	if format == nil {
		format = TextFormat()
	}

	runeCount := u.consoleWidth() - format.leftMargin - format.leftPadding - format.rightPadding - format.rightMargin
	if runeCount <= 0 {
		return "\n"
	}

	message := strings.Repeat(string(character), runeCount)
	text := &uiText{message: message, textFormat: format}
	message = u.Format(text)
	return fmt.Sprintf("%s\n", message)
}

func (u *ui) writeText(writer io.Writer, message string) {

	if writer == nil {
		return
	}

	u.m.Lock()
	defer u.m.Unlock()

	message = log.SecretFilter.Filter(message)

	_, err := fmt.Fprint(writer, message)
	if err != nil {
		log.Errorf("Failed to write to UI: %v", err)
	}
}

func (u *ui) consoleWidth() int {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return defaultConsoleWidth // If not a terminal, return the default width.
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Errorf("Failed to get terminal size: %v", err)
		return defaultConsoleWidth
	}

	return width - 2
}

func getRuneCount(text *uiText) int {

	messageRuneCount := utf8.RuneCountInString(text.message)

	messageRuneCount += text.leftPadding
	messageRuneCount += text.rightPadding

	messageRuneCount += text.leftMargin
	messageRuneCount += text.rightMargin

	return messageRuneCount
}
