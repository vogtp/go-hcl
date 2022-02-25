package hcl

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-hclog"
)

// Logger implements hclog.Logger
type Logger struct {
	hclog.Logger

	w io.Writer

	level hclog.Level
	name  string
}

//creates a copy of itslef
func (l *Logger) copy() Logger {
	n := Logger{
		Logger: l.Logger,
		w:      l.w,
		level:  l.level,
		name:   l.name,
	}
	return n
}

// Printf like logging to Error
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

// Printf like logging to Warn
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, v...))
}

// Printf like logging to Info
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Printf like logging to Debug
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, v...))
}

// Printf like logging to Trace
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Logger.Trace(fmt.Sprintf(format, v...))
}

// Printf logging to Info
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

// Print logging to Info
func (l *Logger) Print(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

// Println logging to Info
func (l *Logger) Println(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

// set the log level
func (l *Logger) SetLevel(level hclog.Level) {
	l.level = level
	l.Logger.SetLevel(level)
}
