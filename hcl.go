// Package hcl is a replacement for log which wraps hc-log
//
// hcl is supposed to provide advanced but painless logging
//
// - it offers simple package level functionality
//
// - exports most (all?) of the hclog features
//
// - it redirects stdlib log to itself.
//
// - it does not support a Panic function
package hcl

import (
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
)

var actLog Logger

// inits a logger with the binary name
func initDefaultLogger() {
	New(GetExecutableName())
}

func init() {
	initDefaultLogger()
}

// Fatalf provides printf like logging to Error
// it stops execution with exit code 1
func Fatalf(format string, v ...interface{}) {
	actLog.Fatalf(format, v...)
}

// Errorf provides printf like logging to Error
func Errorf(format string, v ...interface{}) {
	actLog.Errorf(format, v...)
}

// Warnf provides printf like logging to Warn
func Warnf(format string, v ...interface{}) {
	actLog.Warnf(format, v...)
}

// Infof provides printf like logging to Info
func Infof(format string, v ...interface{}) {
	actLog.Infof(format, v...)
}

// Debugf provides printf like logging to Debug
func Debugf(format string, v ...interface{}) {
	actLog.Debugf(format, v...)
}

// Tracef provides printf like logging to Trace
func Tracef(format string, v ...interface{}) {
	actLog.Tracef(format, v...)
}

// Printf works like Printf from stdlib
// logs to Info
func Printf(format string, v ...interface{}) {
	actLog.Printf(format, v...)
}

// Print works like Print from stdlib
// logs to Info
func Print(v ...interface{}) {
	actLog.Print(v...)
}

// Println works like hcl.Print
// logs to Info
func Println(v ...interface{}) {
	actLog.Print(v...)
}

// Args are alternating key, val pairs
// keys must be strings
// vals can be any type, but display is implementation specific
// Emit a message and key/value pairs at a provided log level
func log(level hclog.Level, msg string, args ...interface{}) {
	if len(args) < 1 {
		actLog.Log(level, msg)
		return
	}
	actLog.Log(level, msg, args...)
}

// Trace logs a message and key/value pairs at the TRACE level
func Trace(msg string, args ...interface{}) {
	log(hclog.Trace, msg, args...)
}

// Debug logs a message and key/value pairs at the DEBUG level
func Debug(msg string, args ...interface{}) {
	log(hclog.Debug, msg, args...)
}

// Info logs a message and key/value pairs at the INFO level
func Info(msg string, args ...interface{}) {
	log(hclog.Info, msg, args...)
}

// Warn logs a message and key/value pairs at the WARN level
func Warn(msg string, args ...interface{}) {
	log(hclog.Warn, msg, args...)
}

// Error log a message and key/value pairs at the ERROR level
func Error(msg string, args ...interface{}) {
	log(hclog.Error, msg, args...)
}

// Fatal log a message and key/value pairs at the ERROR level
// it stops execution with exit code 1
func Fatal(format string, args ...interface{}) {
	log(hclog.Error, format, args...)
	os.Exit(1)
}

// IsTrace indicates if Trace logs would be written
func IsTrace() bool {
	return actLog.IsTrace()
}

// IsDebug indicates if Debug logs would be written
func IsDebug() bool {
	return actLog.IsDebug()
}

// IsInfo indicates if Info logs would be written
func IsInfo() bool {
	return actLog.IsInfo()
}

// IsWarn indicates if Warn logs would be written
func IsWarn() bool {
	return actLog.IsWarn()
}

// IsError indicates if Error logs would be written
func IsError() bool {
	return actLog.IsError()
}

// GetWriter returns a writer
// to be used for frameworks to output to log
func GetWriter() io.Writer {
	return actLog.GetWriter()
}

// SetLevel sets the log level
func SetLevel(level hclog.Level) {
	actLog.SetLevel(level)
}

// Named creates a sublogger with the name appended to the old name
func Named(name string) Logger {
	return actLog.Named(name)
}

// ResetNamed creates a logger with a new name
func ResetNamed(name string) Logger {
	return actLog.ResetNamed(name)
}
