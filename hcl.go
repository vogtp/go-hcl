// hcl is a replacement for log which wraps hc-log
//
// it offers simple package level functionality
// it redirects stdlib log to itself.
// it does not support Fatal or Panic function
package hcl

import (
	"io"

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

// Printf like logging to Error
func Errorf(format string, v ...interface{}) {
	actLog.Errorf(format, v...)
}

// Printf like logging to Warn
func Warnf(format string, v ...interface{}) {
	actLog.Warnf(format, v...)
}

// Printf like logging to Info
func Infof(format string, v ...interface{}) {
	actLog.Infof(format, v...)
}

// Printf like logging to Debug
func Debugf(format string, v ...interface{}) {
	actLog.Debugf(format, v...)
}

// Printf like logging to Trace
func Tracef(format string, v ...interface{}) {
	actLog.Tracef(format, v...)
}

// Printf like in std lib
func Printf(format string, v ...interface{}) {
	actLog.Printf(format, v...)
}

// Print like in std lib
func Print(v ...interface{}) {
	actLog.Print(v...)
}

// Println like in std lib
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

// Emit a message and key/value pairs at the TRACE level
func Trace(msg string, args ...interface{}) {
	log(hclog.Trace, msg, args...)
}

// Emit a message and key/value pairs at the DEBUG level
func Debug(msg string, args ...interface{}) {
	log(hclog.Debug, msg, args...)
}

// Emit a message and key/value pairs at the INFO level
func Info(msg string, args ...interface{}) {
	log(hclog.Info, msg, args...)
}

// Emit a message and key/value pairs at the WARN level
func Warn(msg string, args ...interface{}) {
	log(hclog.Warn, msg, args...)
}

// Emit a message and key/value pairs at the ERROR level
func Error(msg string, args ...interface{}) {
	log(hclog.Error, msg, args...)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func IsTrace() bool {
	return actLog.IsTrace()
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func IsDebug() bool {
	return actLog.IsDebug()
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func IsInfo() bool {
	return actLog.IsInfo()
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func IsWarn() bool {
	return actLog.IsWarn()
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func IsError() bool {
	return actLog.IsError()
}

// return a writer to used for frameworks to output to log
func GetWriter() io.Writer {
	return actLog.GetWriter()
}

//set the log level
func SetLevel(level hclog.Level) {
	actLog.SetLevel(level)
}

// Create a sublogger with the name appended to the old name
func Named(name string) Logger {
	return actLog.Named(name)
}

// Create a logger with a new name
func ResetNamed(name string) Logger {
	return actLog.ResetNamed(name)
}
