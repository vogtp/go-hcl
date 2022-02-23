package hcl

import (
	"bytes"
	"fmt"
	"os"
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

func TestDefault(t *testing.T) {
	var buf testWriter
	initDefaultLogger()
	actLog.SetWriter(&buf)
	actLog.SetLevel(hclog.Trace)
	Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  go-hcl: text to output: string 42\n", buf.Line())
	Print("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	Println("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())

	log(hclog.Debug, "text to output")
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

type outFunc func(msg string, args ...interface{})
type isFunc func() bool

var buf testWriter

func parseLine(l string) (prefix, line string) {
	idx := strings.Index(l, "]")
	prefix = l[1:idx]
	line = l[idx+1:]
	line = strings.TrimLeft(line, " ")
	return prefix, line
}

func checkLevel(t *testing.T, tc logTestCase, logName string, level hclog.Level, out outFunc, is isFunc) {
	active := tc.level <= level
	assert.Equal(t, active, is(), "level is not active")

	prefix := strings.ToUpper(level.String())
	out("text to output")
	if !active {
		assert.Equal(t, 0, buf.Len(), "no output expected")
		return
	}
	p, l := parseLine(buf.Line())
	exp := fmt.Sprintf("%s: text to output\n", logName)
	assert.Equal(t, prefix, p, "wong prefix")
	assert.Equal(t, exp, l, "logline is wrong")

	out("text to output", "strParam", "someParam", "intParam", 42)
	p, l = parseLine(buf.Line())
	exp = fmt.Sprintf("%s: text to output: strParam=someParam intParam=42\n", logName)
	assert.Equal(t, prefix, p, "wong prefix (with params)")
	assert.Equal(t, exp, l, "logline is wrong (with params")
}

type logTestCase struct {
	level hclog.Level
}

func TestLogger(t *testing.T) {
	tt := []logTestCase{
		{hclog.Trace},
		{hclog.Debug},
		{hclog.Info},
		{hclog.Warn},
		{hclog.Error},
	}
	logName := "test-logger"
	l := New(logName, WithWriter(&buf))
	for _, tc := range tt {
		t.Run(tc.level.String(), func(t *testing.T) {
			l.SetLevel(tc.level)
			// test package functions
			checkLevel(t, tc, logName, hclog.Error, Error, IsError)
			checkLevel(t, tc, logName, hclog.Warn, Warn, IsWarn)
			checkLevel(t, tc, logName, hclog.Info, Info, IsInfo)
			checkLevel(t, tc, logName, hclog.Debug, Debug, IsDebug)
			checkLevel(t, tc, logName, hclog.Trace, Trace, IsTrace)
			// test logger funtions
			checkLevel(t, tc, logName, hclog.Error, l.Error, l.IsError)
			checkLevel(t, tc, logName, hclog.Warn, l.Warn, l.IsWarn)
			checkLevel(t, tc, logName, hclog.Info, l.Info, l.IsInfo)
			checkLevel(t, tc, logName, hclog.Debug, l.Debug, l.IsDebug)
			checkLevel(t, tc, logName, hclog.Trace, l.Trace, l.IsTrace)
		})
	}
}

func TestStdLibCompat(t *testing.T) {
	logName := "test-logger"
	l := New(logName, WithWriter(&buf))
	l.Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  test-logger: text to output: string 42\n", buf.Line())
	l.Print("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())
	l.Println("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())

	Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  test-logger: text to output: string 42\n", buf.Line())
	Print("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())
	Println("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())

}

func TestExecCheck(t *testing.T) {
	arg := os.Args[0]
	defer func() { os.Args[0] = arg }()
	os.Args[0] = "/test/tmp/go-build2932332730/b001/go-hcl.test"
	assert.True(t, IsGoRun())

	os.Args[0] = "/test/tmp/b001/go-hcl.test"
	assert.False(t, IsGoRun())

	l := New("")
	os.Args[0] = "/test/tmp/go-build2932332730/b001/go-hcl.test"
	assert.True(t, l.IsGoRun())

	os.Args[0] = "/bin/exe"
	assert.False(t, l.IsGoRun())
}
