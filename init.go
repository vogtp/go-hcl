package hcl

import (
	"fmt"
	"io"
	gologger "log"
	"os"

	"github.com/hashicorp/go-hclog"
)

const (
	// TimeFormat is the default time formating of hcl
	TimeFormat = "2006/01/02 15:04:05"
)

var actLog *Logger

// inits a logger with the binary name
func initDefaultLogger() {
	if actLog == nil {
		l := New()
		actLog = &l
	}
}

// New constructs a new logger
// loglevel is Error if build and info if `go run`
// std lib logging is redirected
func New(opts ...LoggerOpt) Logger {
	l := &Logger{
		name:          GetExecutableName(),
		captureStdlib: true,
		hcOpts: &hclog.LoggerOptions{
			TimeFormat: TimeFormat,
		},
	}
	for _, opt := range opts {
		opt(l)
	}
	if l.w == nil {
		l.w = os.Stderr
	}
	if l.level == hclog.NoLevel {
		l.level = hclog.Warn
		if IsGoRun() {
			l.level = hclog.Info
		}
		if IsGoTest() {
			l.level = hclog.Debug
		}
	}
	// this creates the backend logger
	l.SetWriter(l.w)
	if l.captureStdlib {
		// sets the std lib logger to write to us
		gologger.SetOutput(l.GetWriter())
		gologger.SetPrefix("")
		gologger.SetFlags(0)
	}
	actLog = l
	return *l
}

// With sreates a sublogger
// that will always have the given key/value pairs
func (l Logger) With(args ...interface{}) Logger {
	sl := l.copy()
	sl.Logger = l.Logger.With(args...)
	return sl
}

// Named creates a sublogger with the name appended to the old name
func (l Logger) Named(name string) Logger {
	return l.ResetNamed(fmt.Sprintf("%s.%s", l.name, name))
}

// ResetNamed creates a logger with a new name
func (l Logger) ResetNamed(name string) Logger {
	sl := l.copy()
	sl.name = name
	sl.Logger = l.Logger.ResetNamed(name)
	return sl
}

// LibraryLogger creates a logger for libraries
// if the library used hcl it creates a sublogger
// otherwise it mimics stdlib
func LibraryLogger(name string) Logger {
	if actLog != nil {
		l := actLog.Named(name)
		return l
	}
	opts := hclog.LoggerOptions{TimeFormat: TimeFormat}
	l := New(
		WithName(name),
		WithLevel(hclog.Info),
		WithLoggerOptions(&opts),
		WithStdlib(false),
	)
	// keep actLog clean (we are called from a lib)
	actLog = nil
	return l
}

// SetWriter sets the write of this logger
// redirects the std lib log
func (l *Logger) SetWriter(w io.Writer) {
	l.w = w
	l.hcOpts.Name = l.name
	l.hcOpts.Output = w
	l.hcOpts.Level = l.level
	l.Logger = hclog.New(l.hcOpts)
}

// GetWriter returns a writer
// to be used for frameworks to output to log
func (l Logger) GetWriter() io.Writer {
	return l.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
}

// LoggerOpt is a func to set opts at logger creation
type LoggerOpt func(*Logger)

// WithName sets the name of the logger
func WithName(name string) LoggerOpt {
	return func(l *Logger) {
		l.name = name
	}
}

// WithWriter is used to create a logger with a custom writer
func WithWriter(w io.Writer) LoggerOpt {
	return func(l *Logger) {
		l.w = w
	}
}

// WithLevel is used to create a logger with log level
func WithLevel(lvl hclog.Level) LoggerOpt {
	return func(l *Logger) {
		l.level = lvl
	}
}

// WithLoggerOptions sets the logger options of hclog
// Name, Level, Output get overwritter by hcl options
func WithLoggerOptions(opts *hclog.LoggerOptions) LoggerOpt {
	return func(l *Logger) {
		l.hcOpts = opts
	}
}

// WithStdlib controlls if stdlib logger should be changed
func WithStdlib(b bool) LoggerOpt {
	return func(l *Logger) {
		l.captureStdlib = b
	}
}
