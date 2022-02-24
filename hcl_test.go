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
	line := buf.Line()
	p, l := parseLine(line)
	exp := fmt.Sprintf("%s: text to output\n", logName)
	assert.Equal(t, prefix, p, "wong prefix", line)
	assert.Equal(t, exp, l, "logline is wrong", line)

	out("text to output", "strParam", "someParam", "intParam", 42)
	line = buf.Line()
	p, l = parseLine(line)
	exp = fmt.Sprintf("%s: text to output: strParam=someParam intParam=42\n", logName)
	assert.Equal(t, prefix, p, "wong prefix (with params)", line)
	assert.Equal(t, exp, l, "logline is wrong (with params", line)
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

func TestIsGoRun(t *testing.T) {
	assert.False(t, IsGoRun(), "Selftest", os.Args[0])
	l := New("")
	assert.False(t, l.IsGoRun(), "Selftest", os.Args[0])
	arg := os.Args[0]
	defer func() { os.Args[0] = arg }()
	tests := []struct {
		arg0  string
		goRun bool
	}{
		{"/test/tmp/go-build2932332730/b001/go-hcl", true},
		{"/test/tmp/go-build2932332730/b001/go-hcl.test", false},
		{"/test/tmp/b001/go-hcl", false},
		{"/test/tmp/go-build/2932332730/b001/go-hcl", false},
		{"/test/tmp/go-build/2932332730/b001/go-hcl.test", false},
		{"/bin/exe", false},
		{"./main.go", false},
		{`c:\Users\Administrator\some.exe`, false},
		{`c:\Users\Administrator\AppData\Local\Temp\go-build607140747/b001/go-hcl.exe`, true},
		{`/test/gogo-build/someThing`, false},
		{`c:\Temp\thisgo-build\a.exe`, false},
		{`c:\Temp\this\go-build\a.exe`, false},
		{`\\go-build-server\someshare`, false},
		{``, false},
	}

	for _, tc := range tests {
		t.Run(tc.arg0, func(t *testing.T) {

			os.Args[0] = tc.arg0
			assert.Equal(t, tc.goRun, IsGoRun())
		})
	}
}

func TestIsGoTest(t *testing.T) {
	assert.True(t, IsGoTest(), "Selftest", os.Args[0])
	l := New("")
	assert.True(t, l.IsGoTest(), "Selftest", os.Args[0])
	arg := os.Args[0]
	defer func() { os.Args[0] = arg }()
	tests := []struct {
		arg0 string
		exp  bool
	}{
		{"/test/tmp/go-build2932332730/b001/go-hcl", false},
		{"/test/tmp/go-build2932332730/b001/go-hcl.test", true},
		{"/test/tmp/b001/go-hcl", false},
		{"/test/tmp/go-build/2932332730/b001/go-hcl", false},
		{"/test/tmp/go-build/2932332730/b001/go-hcl.test", true},
		{"/bin/exe", false},
		{"./main.go", false},
		{`c:\Users\Administrator\some.exe`, false},
		{`c:\Users\Administrator\AppData\Local\Temp\go-build607140747/b001/go-hcl.exe`, false},
		{`c:\Users\Administrator\AppData\Local\Temp\go-build607140747/b001/go-hcl.test`, true},
		{`/test/gogo-build/someThing`, false},
		{`c:\Temp\thisgo-build\a.exe`, false},
		{`c:\Temp\this\go-build\a.exe`, false},
		{`\\go-build-server\someshare`, false},
		{``, false},
	}

	for _, tc := range tests {
		t.Run(tc.arg0, func(t *testing.T) {

			os.Args[0] = tc.arg0
			assert.Equal(t, tc.exp, IsGoTest())
		})
	}
}

func TestGetExecutableName(t *testing.T) {
	arg := os.Args[0]
	defer func() { os.Args[0] = arg }()
	tests := []struct {
		arg0 string
		name string
	}{
		{"/test/tmp/go-build2932332730/b001/go-hcl.test", "go-hcl"},
		{"/test/tmp/b001/go-hcl.test", "go-hcl"},
		{"/test/tmp/go-build/2932332730/b001/go-hcl.test", "go-hcl"},
		{"/bin/exe", "exe"},
		{"./main.go", "main"},
		{`c:\Users\Administrator\some.exe`, "some"},
		{`c:\Users\Administrator\AppData\Local\Temp\go-build607140747/b001/go-hcl.exe`, "go-hcl"},
		{`/test/gogo-build/someThing`, "someThing"},
		{`c:\Temp\thisgo-build\a.exe`, "a"},
		{`\\go-build-server\someshare`, "someshare"},
		{``, ""},
	}

	for _, tc := range tests {
		t.Run(tc.arg0, func(t *testing.T) {

			os.Args[0] = tc.arg0
			assert.Equal(t, tc.name, GetExecutableName())
		})
	}
}
