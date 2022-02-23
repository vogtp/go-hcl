package hcl

import (
	"io"
	gologger "log"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
)

func New(name string, opts ...LoggerOpt) *Logger {
	log = &Logger{name: name}
	for _, opt := range opts {
		opt(log)
	}
	if log.w == nil {
		log.SetWriter(os.Stderr)
	}
	return log
}

func (l *Logger) SetWriter(w io.Writer) {
	log.Logger = hclog.New(&hclog.LoggerOptions{
		Name:       l.name,
		TimeFormat: time.RFC3339,
		Output:     w,
	})
	log.w = log.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})

	gologger.SetOutput(log.w)
	gologger.SetPrefix("")
	gologger.SetFlags(0)
}

type LoggerOpt func(*Logger)

func WithWriter(w io.Writer) LoggerOpt {
	return func(l *Logger) {
		l.SetWriter(w)
	}
}
