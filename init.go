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
		actLog.SetWriter(os.Stderr)
	}
	if actLog.level == hclog.NoLevel {
		actLog.level = hclog.Warn
		if IsGoRun() {
			actLog.level = hclog.Info
		}
	}
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
	actLog.w = actLog.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})

	gologger.SetOutput(actLog.w)
	gologger.SetPrefix("")
	gologger.SetFlags(0)
}

type LoggerOpt func(*Logger)

// Used to create a logger with a custom writer
func WithWriter(w io.Writer) LoggerOpt {
	return func(l *Logger) {
		l.SetWriter(w)
	}
}
