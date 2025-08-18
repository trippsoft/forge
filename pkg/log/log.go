// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"fmt"
	"io"
	"log"
)

type LogLevel uint8

const (
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var (
	AddTimestamp bool     = true
	Verbosity    LogLevel = LevelInfo

	traceLogger *log.Logger
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	Init(io.Discard)
}

func Init(w io.Writer) {
	flags := log.Lmsgprefix
	if Verbosity >= LevelDebug {
		flags |= log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	} else if AddTimestamp {
		flags |= log.Ldate | log.Ltime
	}

	traceLogger = log.New(w, "TRACE:\t", flags)
	debugLogger = log.New(w, "DEBUG:\t", flags)
	infoLogger = log.New(w, "INFO :\t", flags)
	warnLogger = log.New(w, "WARN :\t", flags)
	errorLogger = log.New(w, "ERROR:\t", flags)
}

func Trace(v ...any) {
	if Verbosity < LevelTrace {
		return
	}

	traceLogger.Output(2, fmt.Sprintln(v...))
}

func Tracef(format string, v ...any) {
	if Verbosity < LevelTrace {
		return
	}

	traceLogger.Output(2, fmt.Sprintf(format, v...))
}

func Debug(v ...any) {
	if Verbosity < LevelDebug {
		return
	}

	debugLogger.Output(2, fmt.Sprintln(v...))
}

func Debugf(format string, v ...any) {
	if Verbosity < LevelDebug {
		return
	}

	debugLogger.Output(2, fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	if Verbosity < LevelInfo {
		return
	}

	infoLogger.Output(2, fmt.Sprintln(v...))
}

func Infof(format string, v ...any) {
	if Verbosity < LevelInfo {
		return
	}

	infoLogger.Output(2, fmt.Sprintf(format, v...))
}

func Warn(v ...any) {
	if Verbosity < LevelWarn {
		return
	}

	warnLogger.Output(2, fmt.Sprintln(v...))
}

func Warnf(format string, v ...any) {
	if Verbosity < LevelWarn {
		return
	}

	warnLogger.Output(2, fmt.Sprintf(format, v...))
}

func Error(v ...any) {
	errorLogger.Output(2, fmt.Sprintln(v...))
}

func Errorf(format string, v ...any) {
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}
