/* wrapper.go
 *
 * Copyright (c) 2010, Kyle Lemons <kyle@kylelemons.net> (creator).
 * All rights reserved.
 *
 * This software may be modified and distributed under the terms
 * of the New BSD license.  See the LICENSE file for details.
 */

// Package log4go - a robust, configurable, powerful logging package
package log4go

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	Global Logger
)

func init() {
	Global = NewDefaultLogger(DEBUG)
}

func format(count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(" %v", count)[1:]
}

// LoadConfiguration Wrapper for (*Logger).LoadConfiguration
func LoadConfiguration(filename string) {
	Global.LoadConfiguration(filename)
}

// AddFilter Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

// Close Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

// Crash Logs the given message and crashes the program
func Crash(args ...interface{}) {
	format := format(len(args))
	err := Global.Critical(format, args)
	Global.Close()
	panic(err)
}

// Crashf Logs the given message and crashes the program
func Crashf(format string, args ...interface{}) {
	Global.intLog(CRITICAL, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Fatalf Compatibility with `log`
func Fatalf(format string, args ...interface{}) {
	Fatal(format, args)
}

// Fatalln Compatibility with `log`
func Fatalln(args ...interface{}) {
	Fatal(args)
}

// Fatal Compatibility with `log`
func Fatal(args ...interface{}) {
	if len(args) > 0 {
		Global.Critical(format(len(args)), args...)
	}
	Global.Close()
	os.Exit(1)
}

// Exit Compatibility with `log`
func Exit(args ...interface{}) {
	if len(args) > 0 {
		Global.intLog(ERROR, format(len(args)), args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Exitf Compatibility with `log`
func Exitf(format string, args ...interface{}) {
	Global.intLog(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Stderr Compatibility with `log`
func Stderr(args ...interface{}) {
	if len(args) > 0 {
		Global.intLog(ERROR, format(len(args)), args...)
	}
}

// Stderrf Compatibility with `log`
func Stderrf(format string, args ...interface{}) {
	Global.intLog(ERROR, format, args...)
}

// Stdout Compatibility with `log`
func Stdout(args ...interface{}) {
	if len(args) > 0 {
		Global.intLog(INFO, format(len(args)), args...)
	}
}

// Stdoutf Compatibility with `log`
func Stdoutf(format string, args ...interface{}) {
	Global.intLog(INFO, format, args...)
}

// Log Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl level, source, message string) {
	Global.Log(lvl, source, message)
}

// Logf Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl level, format string, args ...interface{}) {
	Global.intLog(lvl, format, args...)
}

// Logc Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl level, closure func() string) {
	Global.intLog(lvl, closure)
}

// Finest Utility for finest log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Finest
func Finest(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLog(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLog(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLog(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Fine Utility for fine log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Fine
func Fine(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLog(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLog(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLog(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Debug Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debug(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLog(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLog(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLog(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Trace Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Trace(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLog(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLog(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLog(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Info Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLog(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLog(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLog(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Warn Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warn(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLog(lvl, first, args...)
		return fmt.Errorf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLog(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLog(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Error Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Error(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLog(lvl, first, args...)
		return fmt.Errorf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLog(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLog(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Critical Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical
func Critical(arg0 interface{}, args ...interface{}) error {
	return Global.Critical(arg0, args)
}
