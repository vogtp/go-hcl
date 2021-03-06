# go-hcl [![Go](https://github.com/vogtp/go-hcl/actions/workflows/go.yml/badge.svg)](https://github.com/vogtp/go-hcl/actions/workflows/go.yml)[![codecov](https://codecov.io/gh/vogtp/go-hcl/branch/main/graph/badge.svg?token=DV0IDZ2FXE)](https://codecov.io/gh/vogtp/go-hcl)[![Go Report Card](https://goreportcard.com/badge/github.com/vogtp/go-hcl)](https://goreportcard.com/report/github.com/vogtp/go-hcl)[![Release](https://img.shields.io/github/release/vogtp/go-hcl.svg?style=flat-square)](https://github.com/vogtp/go-hcl/releases)[![GoDoc](https://pkg.go.dev/badge/github.com/vogtp/go-hcl?status.svg)](https://pkg.go.dev/github.com/vogtp/go-hcl?tab=doc)

hcl is a replacement for log which wraps hc-log 

hcl is supposed to provide advanced but painless logging


## Features

- it offers simple package level functionality
- exports most (all?) of the hclog features 
- it redirects stdlib log to itself
- it does not support Fatal or Panic functions

## Example

```go
package main


import (
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/vogtp/go-hcl"
)

func main() {
	// log to a logger named after the executable
	hcl.Print("I am a logger named after the executable")
	hcl.Printf("But I go rid of some parts in %s", os.Args[0])
	hcl.Println("I log to Info")

	hcl.Errorf("I am getting %s bored", "really")
	hcl.Error("I got called with", "args", os.Args)

	hcl.Info("I am visible when started by go run")
	hcl.Warn("I am visible when started by go build")
	hcl.Trace("I am not visible")
	hcl.SetLevel(hclog.Trace)
	hcl.Trace("now you can see me")

	log.Printf("I look the same as %s", "hcl.Printf")

	//create a sublogger
	webLogger := hcl.Named("web")
	webLogger.Info("Start of web logs")
}


```

Output: 

```
2022-02-25T09:40:12+01:00 [INFO]  hcl: I am a logger named after the executable
2022-02-25T09:40:12+01:00 [INFO]  hcl: But I go rid of some parts in /tmp/go-build2029833176/b001/hcl.test
2022-02-25T09:40:12+01:00 [INFO]  hcl: I log to Info
2022-02-25T09:40:12+01:00 [ERROR] hcl: I am getting really bored
2022-02-25T09:40:12+01:00 [ERROR] hcl: I got called with: args=["/tmp/go-build2029833176/b001/hcl.test", "-test.paniconexit0", "-test.timeout=10m0s"]
2022-02-25T09:40:12+01:00 [INFO]  hcl: I am visible when started by go run
2022-02-25T09:40:12+01:00 [WARN]  hcl: I am visible when started by go build
2022-02-25T09:40:12+01:00 [TRACE] hcl: now you can see me
2022-02-25T09:40:12+01:00 [INFO]  hcl: I look the same as hcl.Printf
2022-02-25T09:40:12+01:00 [INFO]  hcl.web: Start of web logs
```
