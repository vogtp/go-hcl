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
