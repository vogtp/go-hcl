package hcl

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

type testWriter struct {
	bytes.Buffer
}

func (tw *testWriter) Line() string {
	str := tw.String()
	dataIdx := strings.IndexByte(str, ' ')
	tw.Reset()
	return str[dataIdx+1:]
}

func TestInit(t *testing.T) {
	var buf testWriter
	l := New("test")
	l.SetWriter(&buf)
	l.Error("log line")
	assert.Equal(t, "[ERROR] test: log line\n", buf.Line())
	log.Print("std log")
	assert.Equal(t, "[INFO]  test: std log\n", buf.Line())
	Print("std log")
	assert.Equal(t, "[INFO]  test: std log\n", buf.Line())
}

func TestDefault(t *testing.T) {
	var buf testWriter
	initDefaultLogger()
	log.SetWriter(&buf)
	log.SetLevel(hclog.Trace)
	Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  go-hcl: text to output: string 42\n", buf.Line())
	Print("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	Println("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())

	Log(hclog.Debug, "text to output")
	assert.Equal(t, "[DEBUG] go-hcl: text to output\n", buf.Line())
	Trace("text to output")
	assert.Equal(t, "[TRACE] go-hcl: text to output\n", buf.Line())
	Debug("text to output")
	assert.Equal(t, "[DEBUG] go-hcl: text to output\n", buf.Line())
	Info("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	Warn("text to output")
	assert.Equal(t, "[WARN]  go-hcl: text to output\n", buf.Line())
	Error("text to output")
	assert.Equal(t, "[ERROR] go-hcl: text to output\n", buf.Line())
}

func TestLogger(t *testing.T) {
	var buf testWriter
	l := New("test-logger", WithWriter(&buf))
	l.SetLevel(hclog.Trace)
	l.Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  test-logger: text to output: string 42\n", buf.Line())
	l.Print("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())
	l.Println("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())

	l.Log(hclog.Debug, "text to output")
	assert.Equal(t, "[DEBUG] test-logger: text to output\n", buf.Line())
	l.Trace("text to output")
	assert.Equal(t, "[TRACE] test-logger: text to output\n", buf.Line())
	l.Debug("text to output")
	assert.Equal(t, "[DEBUG] test-logger: text to output\n", buf.Line())
	l.Info("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())
	l.Warn("text to output")
	assert.Equal(t, "[WARN]  test-logger: text to output\n", buf.Line())
	l.Error("text to output")
	assert.Equal(t, "[ERROR] test-logger: text to output\n", buf.Line())

	Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  test-logger: text to output: string 42\n", buf.Line())
	Print("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())
	Println("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())

	Log(hclog.Debug, "text to output")
	assert.Equal(t, "[DEBUG] test-logger: text to output\n", buf.Line())
	Trace("text to output")
	assert.Equal(t, "[TRACE] test-logger: text to output\n", buf.Line())
	Debug("text to output")
	assert.Equal(t, "[DEBUG] test-logger: text to output\n", buf.Line())
	Info("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())
	Warn("text to output")
	assert.Equal(t, "[WARN]  test-logger: text to output\n", buf.Line())
	Error("text to output")
	assert.Equal(t, "[ERROR] test-logger: text to output\n", buf.Line())
}

func TestLevel(t *testing.T) {
	tt := []hclog.Level{
		hclog.Trace,
		hclog.Debug,
		hclog.Info,
		hclog.Warn,
		hclog.Error,
	}
	l := New("test-logger")
	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			l.SetLevel(tc)
			assert.Equal(t, tc <= hclog.Trace, l.IsTrace())
			assert.Equal(t, tc <= hclog.Debug, l.IsDebug())
			assert.Equal(t, tc <= hclog.Info, l.IsInfo())
			assert.Equal(t, tc <= hclog.Warn, l.IsWarn())
			assert.Equal(t, tc <= hclog.Error, l.IsError())
			assert.Equal(t, tc <= hclog.Trace, IsTrace())
			assert.Equal(t, tc <= hclog.Debug, IsDebug())
			assert.Equal(t, tc <= hclog.Info, IsInfo())
			assert.Equal(t, tc <= hclog.Warn, IsWarn())
			assert.Equal(t, tc <= hclog.Error, IsError())
		})
	}
}
