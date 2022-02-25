package hcl

import (
	"bytes"
	"fmt"
	gologger "log"
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

	Tracef("text to output: %s", "test string")
	assert.Equal(t, "[TRACE] go-hcl: text to output: test string\n", buf.Line())
	Debugf("text to output: %s", "test string")
	assert.Equal(t, "[DEBUG] go-hcl: text to output: test string\n", buf.Line())
	Infof("text to output: %s", "test string")
	assert.Equal(t, "[INFO]  go-hcl: text to output: test string\n", buf.Line())
	Warnf("text to output: %s", "test string")
	assert.Equal(t, "[WARN]  go-hcl: text to output: test string\n", buf.Line())
	Errorf("text to output: %s", "test string")
	assert.Equal(t, "[ERROR] go-hcl: text to output: test string\n", buf.Line())

	GetWriter().Write([]byte("text to output"))
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	l:=Named("named")
	l.Errorf("text to output: %s", "test string")
	assert.Equal(t, "[ERROR] go-hcl.named: text to output: test string\n", buf.Line())
	l2:=ResetNamed("named")
	l2.Errorf("text to output: %s", "test string")
	assert.Equal(t, "[ERROR] named: text to output: test string\n", buf.Line())
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

func checkLevel(t *testing.T, tcLevel hclog.Level, logName string, checkLevel hclog.Level, out outFunc, is isFunc) {
	active := tcLevel <= checkLevel
	assert.Equal(t, active, is(), "level is not active")

	prefix := strings.ToUpper(checkLevel.String())
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

func TestLogLevel(t *testing.T) {
	tt := []hclog.Level{
		hclog.Trace,
		hclog.Debug,
		hclog.Info,
		hclog.Warn,
		hclog.Error,
	}
	logName := "test-logger"
	l := New(logName, WithWriter(&buf), WithLevel(hclog.Trace))
	assert.True(t, l.IsTrace())
	SetLevel(hclog.Error)
	assert.True(t, l.IsError())
	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			l.SetLevel(tc)
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

func TestSubLogger(t *testing.T) {
	tt := []hclog.Level{
		hclog.Trace,
		hclog.Debug,
		hclog.Info,
		hclog.Warn,
		hclog.Error,
	}
	logName := "test-logger"
	l := New(logName, WithWriter(&buf), WithLevel(hclog.Trace))
	assert.True(t, l.IsTrace())
	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			l.SetLevel(tc)
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
	n := "sublog"
	sublogName := fmt.Sprintf("%s.%s", logName, n)
	l = l.Named(n)
	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			l.SetLevel(tc)
			// test package functions
			checkLevel(t, tc, logName, hclog.Error, Error, IsError)
			checkLevel(t, tc, logName, hclog.Warn, Warn, IsWarn)
			checkLevel(t, tc, logName, hclog.Info, Info, IsInfo)
			checkLevel(t, tc, logName, hclog.Debug, Debug, IsDebug)
			checkLevel(t, tc, logName, hclog.Trace, Trace, IsTrace)
			// test logger funtions
			checkLevel(t, tc, sublogName, hclog.Error, l.Error, l.IsError)
			checkLevel(t, tc, sublogName, hclog.Warn, l.Warn, l.IsWarn)
			checkLevel(t, tc, sublogName, hclog.Info, l.Info, l.IsInfo)
			checkLevel(t, tc, sublogName, hclog.Debug, l.Debug, l.IsDebug)
			checkLevel(t, tc, sublogName, hclog.Trace, l.Trace, l.IsTrace)
		})
	}
	sublogName = "new-name"
	l = l.ResetNamed(sublogName)
	for _, tc := range tt {
		t.Run(tc.String(), func(t *testing.T) {
			l.SetLevel(tc)
			// test package functions
			checkLevel(t, tc, logName, hclog.Error, Error, IsError)
			checkLevel(t, tc, logName, hclog.Warn, Warn, IsWarn)
			checkLevel(t, tc, logName, hclog.Info, Info, IsInfo)
			checkLevel(t, tc, logName, hclog.Debug, Debug, IsDebug)
			checkLevel(t, tc, logName, hclog.Trace, Trace, IsTrace)
			// test logger funtions
			checkLevel(t, tc, sublogName, hclog.Error, l.Error, l.IsError)
			checkLevel(t, tc, sublogName, hclog.Warn, l.Warn, l.IsWarn)
			checkLevel(t, tc, sublogName, hclog.Info, l.Info, l.IsInfo)
			checkLevel(t, tc, sublogName, hclog.Debug, l.Debug, l.IsDebug)
			checkLevel(t, tc, sublogName, hclog.Trace, l.Trace, l.IsTrace)
		})
	}
	l = l.With("arg", "some information")
	l.Error("an other logline")
	assert.Equal(t, "[ERROR] new-name: an other logline: arg=\"some information\"\n", buf.Line())
}

type logTestCase struct {
	name      string
	expLevel  hclog.Level
	arg0      string
	test, run bool
}

func TestLogger(t *testing.T) {
	arg := os.Args[0]
	defer func() { os.Args[0] = arg }()
	tt := []logTestCase{
		{"Build/Production", hclog.Warn, "exe", false, false},
		{"go run", hclog.Info, "/test/tmp/go-build2932332730/b001/go-hcl", false, true},
		{"go test", hclog.Debug, "exe.test", true, false},
	}
	logName := "test-logger"
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			os.Args[0] = tc.arg0
			l := New(logName, WithWriter(&buf))
			// test package functions
			assert.Equal(t, tc.test, l.IsGoTest())
			assert.Equal(t, tc.run, l.IsGoRun())
			checkLevel(t, tc.expLevel, logName, hclog.Error, Error, IsError)
			checkLevel(t, tc.expLevel, logName, hclog.Warn, Warn, IsWarn)
			checkLevel(t, tc.expLevel, logName, hclog.Info, Info, IsInfo)
			checkLevel(t, tc.expLevel, logName, hclog.Debug, Debug, IsDebug)
			checkLevel(t, tc.expLevel, logName, hclog.Trace, Trace, IsTrace)
			// test logger funtions
			assert.Equal(t, tc.test, IsGoTest())
			assert.Equal(t, tc.run, IsGoRun())
			checkLevel(t, tc.expLevel, logName, hclog.Error, l.Error, l.IsError)
			checkLevel(t, tc.expLevel, logName, hclog.Warn, l.Warn, l.IsWarn)
			checkLevel(t, tc.expLevel, logName, hclog.Info, l.Info, l.IsInfo)
			checkLevel(t, tc.expLevel, logName, hclog.Debug, l.Debug, l.IsDebug)
			checkLevel(t, tc.expLevel, logName, hclog.Trace, l.Trace, l.IsTrace)
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

	gologger.Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  test-logger: text to output: string 42\n", buf.Line())
	gologger.Print("text to output")
	assert.Equal(t, "[INFO]  test-logger: text to output\n", buf.Line())
	gologger.Println("text to output")
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
		{`c:\Temp\this\go-build\a.exe.test`, false},
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
		{"exe.text", false},
		{"exe.test", true},
		{"/test/tmp/go-build2932332730/b001/go-hcl", false},
		{"/test/tmp/go-build2932332730/b001/__debug_bin", true},
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
		{`c:\Temp\thisgo-build\a..test.exe`, true},
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
		{"/bin/.exe", ".exe"},
		{"/bin/.exe.ext", ".exe"},
		{"/bin/exe.ext", "exe"},
		{"./main.go", "main"},
		{`c:\Users\Administrator\some.exe`, "some"},
		{`c:\Users\Administrator\.some.exe`, ".some"},
		{`c:\Users\Administrator\.some`, ".some"},
		{`c:\Users\Administrator\AppData\Local\Temp\go-build607140747/b001/go-hcl.exe`, "go-hcl"},
		{`/test/gogo-build/someThing`, "someThing"},
		{`c:\Temp\thisgo-build\a.exe`, "a"},
		{`\\go-build-server\someshare`, "someshare"},
		{`./exe`, "exe"},
		{``, ""},
	}

	for _, tc := range tests {
		t.Run(tc.arg0, func(t *testing.T) {

			os.Args[0] = tc.arg0
			assert.Equal(t, tc.name, GetExecutableName())
		})
	}
}
