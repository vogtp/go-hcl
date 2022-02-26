package hcl

import (
	"fmt"
	"io"
	gologger "log"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
)

// constructs a new logger
// loglevel is Error if build and info if `go run`
// std lib logging is redirected
func New(name string, opts ...LoggerOpt) Logger {
	actLog = Logger{name: name}
	for _, opt := range opts {
		opt(&actLog)
	}
	if actLog.w == nil {
		actLog.w = os.Stderr
	}
	if actLog.level == hclog.NoLevel {
		actLog.level = hclog.Warn
		if IsGoRun() {
			actLog.level = hclog.Info
		}
		if IsGoTest() {
			actLog.level = hclog.Debug
		}
	}
	// this creates the backend logger
	actLog.SetWriter(actLog.w)
	// sets the std lib logger to write to us
	gologger.SetOutput(actLog.GetWriter())
	gologger.SetPrefix("")
	gologger.SetFlags(0)
	return actLog
}

// Creates a sublogger that will always have the given key/value pairs
func (l Logger) With(args ...interface{}) Logger {
	sl := l.copy()
	sl.Logger = l.Logger.With(args...)
	return sl
}

// Create a sublogger with the name appended to the old name
func (l Logger) Named(name string) Logger {
	return l.ResetNamed(fmt.Sprintf("%s.%s", l.name, name))
}

// Create a logger with a new name
func (l Logger) ResetNamed(name string) Logger {
	sl := l.copy()
	sl.name = name
	sl.Logger = l.Logger.ResetNamed(name)
	return sl
}

//Sets the write to this logger and redirects the std lib log
func (l *Logger) SetWriter(w io.Writer) {
	actLog.Logger = hclog.New(&hclog.LoggerOptions{
		Name:       l.name,
		TimeFormat: time.RFC3339,
		Output:     w,
		Level:      l.level,
	})
	actLog.w = w
}

// return a writer to used for frameworks to output to log
func (l Logger) GetWriter() io.Writer {
	return l.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
}

type LoggerOpt func(*Logger)

// Used to create a logger with a custom writer
func WithWriter(w io.Writer) LoggerOpt {
	return func(l *Logger) {
		l.w = w
	}
}

// Used to create a logger with log level
func WithLevel(lvl hclog.Level) LoggerOpt {
	return func(l *Logger) {
		l.level = lvl
	}
}
