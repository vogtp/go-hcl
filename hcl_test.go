package hcl

import (
	"bytes"
	"errors"
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
	tw.Reset()
	dataIdx := strings.IndexByte(str, ' ')
	str = str[dataIdx+1:]
	dataIdx = strings.IndexByte(str, ' ')
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
	l := Named("named")
	l.Errorf("text to output: %s", "test string")
	assert.Equal(t, "[ERROR] go-hcl.named: text to output: test string\n", buf.Line())
	l2 := ResetNamed("named")
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
	l := New(WithName(logName), WithWriter(&buf), WithLevel(hclog.Trace))
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
	l := New(WithName(logName), WithWriter(&buf), WithLevel(hclog.Trace))
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
			l := New(WithName(logName), WithWriter(&buf))
			// test package functions
			assert.Equal(t, tc.test, IsGoTest())
			assert.Equal(t, tc.run, IsGoRun())
			checkLevel(t, tc.expLevel, logName, hclog.Error, Error, IsError)
			checkLevel(t, tc.expLevel, logName, hclog.Warn, Warn, IsWarn)
			checkLevel(t, tc.expLevel, logName, hclog.Info, Info, IsInfo)
			checkLevel(t, tc.expLevel, logName, hclog.Debug, Debug, IsDebug)
			checkLevel(t, tc.expLevel, logName, hclog.Trace, Trace, IsTrace)
			// test logger funtions
			assert.Equal(t, tc.test, l.IsGoTest())
			assert.Equal(t, tc.run, l.IsGoRun())
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
	l := New(WithName(logName), WithWriter(&buf))
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

func TestWithStdlib(t *testing.T) {
	var buf testWriter
	New(WithWriter(&buf))
	Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  go-hcl: text to output: string 42\n", buf.Line())
	Print("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	Println("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	gologger.Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  go-hcl: text to output: string 42\n", buf.Line())
	gologger.Print("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	gologger.Println("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())

	New(WithWriter(&buf), WithStdlib(false))
	Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  go-hcl: text to output: string 42\n", buf.Line())
	Print("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	Println("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	gologger.Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  go-hcl: text to output: string 42\n", buf.Line())
	gologger.Print("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
	gologger.Println("text to output")
	assert.Equal(t, "[INFO]  go-hcl: text to output\n", buf.Line())
}

func TestWithLoggerOptions(t *testing.T) {
	var buf bytes.Buffer
	opts := hclog.LoggerOptions{
		DisableTime: true,
	}
	New(WithName(""), WithLevel(hclog.Info), WithLoggerOptions(&opts), WithWriter(&buf))
	Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  text to output: string 42\n", buf.String())
	buf.Reset()
	Print("text to output")
	assert.Equal(t, "[INFO]  text to output\n", buf.String())
	buf.Reset()
	Println("text to output")
	assert.Equal(t, "[INFO]  text to output\n", buf.String())
	buf.Reset()
}

func TestLibraryLogger(t *testing.T) {
	New(WithName("base"), WithLevel(hclog.Info), WithWriter(&buf))
	libLog := LibraryLogger("libName")
	libLog.Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  base.libName: text to output: string 42\n", buf.Line())

	libLog.Print("text to output")
	assert.Equal(t, "[INFO]  base.libName: text to output\n", buf.Line())

	libLog.Println("text to output")
	assert.Equal(t, "[INFO]  base.libName: text to output\n", buf.Line())

	gologger.Println("text to output")
	assert.Equal(t, "[INFO]  base: text to output\n", buf.Line())

	actLog = nil
	gologger.SetOutput(os.Stderr)
	libLog = LibraryLogger("libName")
	libLog.SetWriter(&buf)
	libLog.Printf("text to output: %s %d", "string", 42)
	assert.Equal(t, "[INFO]  libName: text to output: string 42\n", buf.Line())
	libLog.Print("text to output")
	assert.Equal(t, "[INFO]  libName: text to output\n", buf.Line())
	libLog.Println("text to output")
	assert.Equal(t, "[INFO]  libName: text to output\n", buf.Line())
	gologger.Println("text to output")
	assert.Equal(t, "", buf.Line())

}

func TestVlogCompat(t *testing.T) {
	hcl := New(WithName("base"), WithLevel(hclog.Trace), WithWriter(&buf))
	vl := hcl.Vlog()
	vl.ErrorString("Test", "something")
	assert.Equal(t, "[ERROR] base: Test something\n", buf.Line())

	vl.Error(errors.New("Test err"))
	assert.Equal(t, "[ERROR] base: Test err\n", buf.Line())

	vl.Warn("Test", "something")
	assert.Equal(t, "[WARN]  base: Test something\n", buf.Line())

	vl.Info("Test", "something")
	assert.Equal(t, "[INFO]  base: Test something\n", buf.Line())

	vl.Debug("Test", "something")
	assert.Equal(t, "[DEBUG] base: Test something\n", buf.Line())

	
	c:=vl.Trace("func")
	assert.Equal(t, "[TRACE] base: START func\n", buf.Line())
	c()
	assert.Equal(t, "[TRACE] base: END func\n", buf.Line())
}

