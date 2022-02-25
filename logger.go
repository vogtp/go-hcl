// hcl is a replacement for log which wraps hc-log
// in does not support Fatal or Panic function
package hcl

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-hclog"
)

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

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, v...))
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Logger.Trace(fmt.Sprintf(format, v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Print(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

func (l *Logger) Println(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

func (l *Logger) SetLevel(level hclog.Level) {
	l.level = level
	l.Logger.SetLevel(level)
}
