package hcl

import (
	"os"
	"strings"
)

func IsGoRun() bool {
	return strings.Contains(os.Args[0], "go-build")
}

func (Logger) IsGoRun() bool {
	return IsGoRun()
}
