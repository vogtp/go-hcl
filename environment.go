package hcl

import (
	"math"
	"os"
	"strconv"
	"strings"
)

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

func (Logger) IsGoRun() bool {
	return IsGoRun()
}

func IsGoTest() bool {
	a := os.Args[0]
	if strings.HasSuffix(a, ".test") {
		return true
	}
	if strings.HasSuffix(a, "__debug_bin") {
		return true
	}
	return false
}

func (Logger) IsGoTest() bool {
	return IsGoTest()
}

func GetExecutableName() string {
	a := os.Args[0]
	sL := strings.LastIndex(a, "/")
	sW := strings.LastIndex(a, "\\")
	s := int(math.Max(float64(sL), float64(sW)))
	e := strings.LastIndex(a, ".")
	if s < 0 {
		s = 0
	}
	if e < s {
		e = len(a)
	}
	if s+e < 1 {
		return ""
	}
	return a[s+1 : e]
}
