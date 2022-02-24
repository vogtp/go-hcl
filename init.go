package hcl

import (
	"io"
	gologger "log"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
)

// constructs a new logger
// loglevel is Error if build and info if `go run`
// std lib logging is redirected
func New(name string, opts ...LoggerOpt) *Logger {
	actLog = &Logger{name: name}
	for _, opt := range opts {
		opt(actLog)
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
	return actLog
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

	gologger.SetOutput(l.GetWriter())
	gologger.SetPrefix("")
	gologger.SetFlags(0)
}

// return a writer to used for frameworks to output to log
func (l *Logger) GetWriter() io.Writer {
	return l.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
}

type LoggerOpt func(*Logger)

// Used to create a logger with a custom writer
func WithWriter(w io.Writer) LoggerOpt {
	return func(l *Logger) {
		l.w = w
	}
}
