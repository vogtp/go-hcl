package hcl_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vogtp/go-hcl"
)

func TestIsGoRun(t *testing.T) {
	assert.False(t, hcl.IsGoRun(), "Selftest", os.Args[0])
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
			assert.Equal(t, tc.goRun, hcl.IsGoRun())
		})
	}
}

func TestIsGoTest(t *testing.T) {
	assert.True(t, hcl.IsGoTest(), "Selftest", os.Args[0])
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
			assert.Equal(t, tc.exp, hcl.IsGoTest())
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
		{`c:\Users\Administrator\some.test.exe`, "some"},
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
			assert.Equal(t, tc.name, hcl.GetExecutableName())
		})
	}
}
