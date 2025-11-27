// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package ui

type HeaderLevel uint8

const (
	HeaderLevel1 HeaderLevel = 1
	HeaderLevel2 HeaderLevel = 2
	HeaderLevel3 HeaderLevel = 3
)

type ResultType uint8

const (
	ResultTypeSuccess ResultType = iota
	ResultTypeChanged
	ResultTypeSkipped
	ResultTypeFailure
)

type MessageLevel uint8

const (
	MessageLevelInfo MessageLevel = iota
	MessageLevelWarning
	MessageLevelError
)

// UI represents a user interface for text output.
//
// Each implementation will be specific to the type of UI (e.g. CLI, Packer plugin, web).
// The implementation should handle secret filtering and text formatting.
type UI interface {
	// Print prints a general text message.
	//
	// This is usually used for startup messages or other non-structured output.
	Print(text string)

	// PrintHeader prints a header with a specified level.
	//
	// The level (1-3) indicates the importance of the header.
	// Level 1 is the highest level (e.g. main section title), while level 3 is the lowest (e.g. subsection title).
	PrintHeader(level HeaderLevel, text string)

	// PrintHostResult prints the step result for a host without multiple iterations.
	//
	// The hostname is the name of the managed system.
	// The result indicates the outcome of the step execution.
	PrintHostResult(hostname string, result ResultType)

	// PrintIterationResult prints the step result for a single iteration on a host.
	//
	// The hostname is the name of the managed system.
	// The iterationLabel is the label for the specific iteration.
	// The result indicates the outcome of the step execution for that iteration.
	PrintIterationResult(hostname, iterationLabel string, result ResultType)

	// PrintMessage prints a general message with a specified level.
	//
	// The message is the content to be printed.
	// The level indicates the severity or importance of the message.
	PrintMessage(message string, level MessageLevel)

	// PrintErrorDetail prints a detailed error message.
	//
	// The message is the detailed error content to be printed.
	// This will be supressed in non-debug modes.
	PrintErrorDetail(message string)
}
