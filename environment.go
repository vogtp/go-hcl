package hcl

import (
	"os"
	"strconv"
	"strings"
)

func IsGoRun() bool {
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
