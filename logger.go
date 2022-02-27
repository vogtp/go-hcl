package hcl

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-hclog"
)

// Logger implements hclog.Logger
type Logger struct {
	hclog.Logger

	w      io.Writer
	hcOpts *hclog.LoggerOptions

	level         hclog.Level
	name          string
	captureStdlib bool
}

//creates a copy of itslef
func (l Logger) copy() Logger {
	n := Logger{
		Logger: l.Logger,
		w:      l.w,
		hcOpts: l.hcOpts,
		level:  l.level,
		name:   l.name,
	}
	return n
}

// Errorf provides printf like logging to Error
func (l Logger) Errorf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

// Warnf provides printf like logging to Warn
func (l Logger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, v...))
}

// Infof provides printf like logging to Info
func (l Logger) Infof(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Debugf provides printf like logging to Debug
func (l Logger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, v...))
}

// Tracef provides printf like logging to Trace
func (l Logger) Tracef(format string, v ...interface{}) {
	l.Logger.Trace(fmt.Sprintf(format, v...))
}

// Printf works like Printf from stdlib
// logs to Info
func (l Logger) Printf(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Print works like Print from stdlib
// logs to Info
func (l Logger) Print(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

// Println works like hcl.Print
// logs to Info
func (l Logger) Println(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level hclog.Level) {
	l.level = level
	l.Logger.SetLevel(level)
}
