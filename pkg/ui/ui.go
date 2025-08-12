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
	Print(message string)                                          // Print outputs a message to the standard output.
	PrintWithFormat(message string, formatting TextFormatting)     // PrintWithFormat outputs a formatted message to the standard output.
	Error(message string)                                          // Error outputs an error message to the standard error.
	ErrorWithFormat(message string, formatting TextFormatting)     // ErrorWithFormat outputs a formatted error message to the standard error.
	PrintLine(character rune)                                      // PrintLine outputs a line of a specific character repeated to the standard output.
	PrintLineWithFormat(character rune, formatting TextFormatting) // PrintLineWithFormat outputs a line of a specific character repeated to the standard output.

	// PrintColumns outputs two columns of text with the specified formatting.
	PrintColumns(leftMessage string, leftFormatting TextFormatting, rightMessage string, rightFormatting TextFormatting)
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
	u.writeText(u.out, message, TextFormatting{})
}

// PrintWithFormat implements UI.
func (u *ui) PrintWithFormat(message string, formatting TextFormatting) {
	u.writeText(u.out, message, formatting)
}

// Error implements UI.
func (u *ui) Error(message string) {
	u.writeText(u.err, message, TextFormatting{})
}

// ErrorWithFormat implements UI.
func (u *ui) ErrorWithFormat(message string, formatting TextFormatting) {
	u.writeText(u.err, message, formatting)
}

// PrintLine implements UI.
func (u *ui) PrintLine(character rune) {
	length := u.consoleWidth()
	line := strings.Repeat(string(character), length)
	line = fmt.Sprintf("%s\n", line)
	u.writeText(u.out, line, TextFormatting{})
}

// PrintLineWithFormat implements UI.
func (u *ui) PrintLineWithFormat(character rune, formatting TextFormatting) {
	length := u.consoleWidth() - formatting.LeftPadding - formatting.RightPadding
	line := strings.Repeat(string(character), length)
	line = u.formatText(line, formatting)
	line = fmt.Sprintf("%s\n", line)
	u.writeText(u.out, line, TextFormatting{})
}

// PrintColumns implements UI.
func (u *ui) PrintColumns(leftMessage string, leftFormatting TextFormatting, rightMessage string, rightFormatting TextFormatting) {

	leftMessage = log.SecretFilter.Filter(leftMessage)
	rightMessage = log.SecretFilter.Filter(rightMessage)

	leftRuneCount := utf8.RuneCountInString(leftMessage)
	rightRuneCount := utf8.RuneCountInString(rightMessage)

	spaceCount := u.consoleWidth() - leftRuneCount - rightRuneCount - leftFormatting.LeftPadding - leftFormatting.RightPadding - rightFormatting.LeftPadding - rightFormatting.RightPadding
	left := u.formatText(leftMessage, leftFormatting)

	spaces := ""
	if spaceCount > 0 {
		spaces = strings.Repeat(" ", spaceCount)
	}

	right := u.formatText(rightMessage, rightFormatting)

	message := fmt.Sprintf("%s%s%s\n", left, spaces, right)

	u.writeText(u.out, message, TextFormatting{})
}

func (u *ui) writeText(writer io.Writer, message string, formatting TextFormatting) {

	if writer == nil {
		return
	}

	u.m.Lock()
	defer u.m.Unlock()

	message = log.SecretFilter.Filter(message)

	message = u.formatText(message, formatting)

	_, err := fmt.Fprint(writer, message)
	if err != nil {
		log.Errorf("Failed to write to UI: %v", err)
	}
}

func (u *ui) formatText(message string, formatting TextFormatting) string {

	result := ""
	if formatting.LeftPadding > 0 {
		result += strings.Repeat(" ", formatting.LeftPadding)
	}

	result += message

	if formatting.RightPadding > 0 {
		result += strings.Repeat(" ", formatting.RightPadding)
	}

	if u == nil || !u.color || len(formatting.Args) == 0 {
		return result
	}

	args := ""
	for i, arg := range formatting.Args {
		if i > 0 {
			args += ";"
		}
		args += fmt.Sprintf("%d", arg)
	}

	return fmt.Sprintf("\033[%sm%s\033[0m", args, result)
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
