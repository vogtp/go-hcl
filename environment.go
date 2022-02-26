package hcl

import (
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// IsGoRun checks if run by go run
// it does this by checking arg[0]
func IsGoRun() bool {
	if IsGoTest() {
		return false
	}
	a := os.Args[0]
	i := strings.Index(a, "go-build")
	if i == -1 {
		return false
	}
	i1 := strings.Index(a, "/go-build")
	i2 := strings.Index(a, "\\go-build")
	if i1+i2 < 0 {
		return false
	}
	s := string(a[i+len("go-build")])
	_, err := strconv.Atoi(s)
	return err == nil
}

// IsGoRun checks if run by go run
// it does this by checking arg[0]
func (Logger) IsGoRun() bool {
	return IsGoRun()
}

// IsGoTest checks if run by go test
// it does this by checking arg[0]
func IsGoTest() bool {
	a := os.Args[0]
	if strings.HasSuffix(a, ".test") {
		return true
	}
	if strings.HasSuffix(a, ".test.exe") {
		return true
	}
	if strings.HasSuffix(a, "__debug_bin") {
		return true
	}
	return false
}

// IsGoTest checks if run by go test
// it does this by checking arg[0]
func (Logger) IsGoTest() bool {
	return IsGoTest()
}

// GetExecutableName extracts the name of the executable
// removes path and suffix
func GetExecutableName() string {
	a := os.Args[0]
	sL := strings.LastIndex(a, "/")
	sW := strings.LastIndex(a, "\\")
	s := int(math.Max(float64(sL), float64(sW)))
	e := strings.LastIndex(a, ".")
	if s < 0 {
		s = 0
	}
	if e < s || s+1 == e {
		e = len(a)
	}
	if s+e < 1 {
		return ""
	}
	n := a[s+1 : e]
	e = strings.LastIndex(n, ".")
	if e < 2 {
		return n
	}
	return n[:e]
}

// GetCaller reports the package and fuction called from
func GetCaller() (pkg string, fun string) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	return funcName[:lastDot], funcName[lastDot+1:]
}
