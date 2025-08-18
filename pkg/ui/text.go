// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package ui

import "github.com/trippsoft/forge/pkg/log"

type BackgroundColor string

const (
	BackgroundBlack   BackgroundColor = "40"
	BackgroundRed     BackgroundColor = "41"
	BackgroundGreen   BackgroundColor = "42"
	BackgroundYellow  BackgroundColor = "43"
	BackgroundBlue    BackgroundColor = "44"
	BackgroundMagenta BackgroundColor = "45"
	BackgroundCyan    BackgroundColor = "46"
	BackgroundWhite   BackgroundColor = "47"
)

type ForegroundColor string

const (
	ForegroundBlack   ForegroundColor = "30"
	ForegroundRed     ForegroundColor = "31"
	ForegroundGreen   ForegroundColor = "32"
	ForegroundYellow  ForegroundColor = "33"
	ForegroundBlue    ForegroundColor = "34"
	ForegroundMagenta ForegroundColor = "35"
	ForegroundCyan    ForegroundColor = "36"
	ForegroundWhite   ForegroundColor = "37"
)

type TextStyle string

const (
	StyleReset     TextStyle = "0"
	StyleBold      TextStyle = "1"
	StyleDim       TextStyle = "2"
	StyleItalic    TextStyle = "3"
	StyleUnderline TextStyle = "4"
	StyleBlink     TextStyle = "5"
)

type textFormat struct {
	backgroundColor BackgroundColor
	foregroundColor ForegroundColor

	styles []TextStyle

	leftPadding  int // padding is inside the formatting
	rightPadding int

	leftMargin  int // margin is outside the formatting
	rightMargin int
}

type TextFormatMap[T comparable] map[T]*textFormat

func TextFormat() *textFormat {
	return &textFormat{styles: []TextStyle{}}
}

func (t *textFormat) WithBackgroundColor(color BackgroundColor) *textFormat {
	t.backgroundColor = color
	return t
}

func (t *textFormat) WithForegroundColor(color ForegroundColor) *textFormat {
	t.foregroundColor = color
	return t
}

func (t *textFormat) WithStyle(style TextStyle) *textFormat {
	if t.styles == nil {
		t.styles = []TextStyle{}
	}

	t.styles = append(t.styles, style)
	return t
}

func (t *textFormat) WithLeftPadding(padding int) *textFormat {
	t.leftPadding = padding
	return t
}

func (t *textFormat) WithRightPadding(padding int) *textFormat {
	t.rightPadding = padding
	return t
}

func (t *textFormat) WithLeftMargin(margin int) *textFormat {
	t.leftMargin = margin
	return t
}

func (t *textFormat) WithRightMargin(margin int) *textFormat {
	t.rightMargin = margin
	return t
}

func (t *textFormat) Clone() *textFormat {
	return &textFormat{
		backgroundColor: t.backgroundColor,
		foregroundColor: t.foregroundColor,
		styles:          append([]TextStyle{}, t.styles...),
		leftPadding:     t.leftPadding,
		rightPadding:    t.rightPadding,
		leftMargin:      t.leftMargin,
		rightMargin:     t.rightMargin,
	}
}

type uiText struct {
	*textFormat
	message string
}

func Text(message string) *uiText {
	return &uiText{
		message:    log.SecretFilter.Filter(message),
		textFormat: TextFormat(),
	}
}

func (t *uiText) WithBackgroundColor(color BackgroundColor) *uiText {
	t.textFormat.WithBackgroundColor(color)
	return t
}

func (t *uiText) WithForegroundColor(color ForegroundColor) *uiText {
	t.textFormat.WithForegroundColor(color)
	return t
}

func (t *uiText) WithStyle(style TextStyle) *uiText {
	t.textFormat.WithStyle(style)
	return t
}

func (t *uiText) WithLeftPadding(padding int) *uiText {
	t.textFormat.WithLeftPadding(padding)
	return t
}

func (t *uiText) WithRightPadding(padding int) *uiText {
	t.textFormat.WithRightPadding(padding)
	return t
}

func (t *uiText) WithLeftMargin(margin int) *uiText {
	t.textFormat.WithLeftMargin(margin)
	return t
}

func (t *uiText) WithRightMargin(margin int) *uiText {
	t.textFormat.WithRightMargin(margin)
	return t
}

func (t *uiText) WithFormat(format *textFormat) *uiText {
	if format != nil {
		t.textFormat = format.Clone()
	} else {
		t.textFormat = TextFormat()
	}

	return t
}
