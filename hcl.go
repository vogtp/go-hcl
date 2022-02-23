// hcl is a replacement for log which wraps hc-log
// it does not support Fatal or Panic function
package hcl

import (
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

var log *Logger

func initDefaultLogger() {
	a := os.Args[0]
	s := strings.LastIndex(a, "/")
	e := strings.LastIndex(a, ".")
	if e < s {
		e = len(a)
	}
	New(a[s+1 : e])
}

func init() {
	initDefaultLogger()
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func Println(v ...interface{}) {
	log.Print(v...)
}

// Args are alternating key, val pairs
// keys must be strings
// vals can be any type, but display is implementation specific
// Emit a message and key/value pairs at a provided log level
func Log(level hclog.Level, msg string, args ...interface{}) {
	if len(args) < 1 {
		log.Log(level, msg)
		return
	}
	log.Log(level, msg, args)
}

// Emit a message and key/value pairs at the TRACE level
func Trace(msg string, args ...interface{}) {
	if len(args) < 1 {
		log.Trace(msg)
		return
	}
	log.Trace(msg, args)
}

// Emit a message and key/value pairs at the DEBUG level
func Debug(msg string, args ...interface{}) {
	if len(args) < 1 {
		log.Debug(msg)
		return
	}
	log.Debug(msg, args)
}

// Emit a message and key/value pairs at the INFO level
func Info(msg string, args ...interface{}) {
	if len(args) < 1 {
		log.Info(msg)
		return
	}
	log.Info(msg, args)
}

// Emit a message and key/value pairs at the WARN level
func Warn(msg string, args ...interface{}) {
	if len(args) < 1 {
		log.Warn(msg)
		return
	}
	log.Warn(msg, args)
}

// Emit a message and key/value pairs at the ERROR level
func Error(msg string, args ...interface{}) {
	if len(args) < 1 {
		log.Error(msg)
		return
	}
	log.Error(msg, args)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func IsTrace() bool {
	return log.IsTrace()
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func IsDebug() bool {
	return log.IsDebug()
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func IsInfo() bool {
	return log.IsInfo()
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func IsWarn() bool {
	return log.IsWarn()
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func IsError() bool {
	return log.IsError()
}
