package hcl_test

import (
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/vogtp/go-hcl"
)

func ExampleLogger() {

	// log to a logger named after the executable
	hcl.Print("I am a logger named after the executable")
	hcl.Printf("But I go rid of some parts in %s", os.Args[0])
	hcl.Println("I log to Info")

	hcl.Errorf("I am getting %s bored", "really")
	hcl.Error("I got called with", "args", os.Args)

	hcl.Warn("I am visible in build code")
	hcl.Info("I am visible when started by go run")
	hcl.Warn("I am visible when started by go test")
	hcl.Trace("I am not visible")
	hcl.SetLevel(hclog.Trace)
	hcl.Trace("now you can see me")

	log.Printf("I look the same as %s", "hcl.Printf")

	//create a sublogger
	webLogger := hcl.Named("web")
	webLogger.Info("Start of web logs")

	// 2022-02-25T09:40:12+01:00 [INFO]  hcl: I am a logger named after the executable
	// 2022-02-25T09:40:12+01:00 [INFO]  hcl: But I go rid of some parts in /tmp/go-build2029833176/b001/hcl.test
	// 2022-02-25T09:40:12+01:00 [INFO]  hcl: I log to Info
	// 2022-02-25T09:40:12+01:00 [ERROR] hcl: I am getting really bored
	// 2022-02-25T09:40:12+01:00 [ERROR] hcl: I got called with: args=["/tmp/go-build2029833176/b001/hcl.test", "-test.paniconexit0", "-test.timeout=10m0s"]
	// 2022-02-25T09:40:12+01:00 [WARN]  hcl: I am visible in build code
	// 2022-02-25T09:40:12+01:00 [INFO]  hcl: I am visible when started by go run
	// 2022-02-25T09:40:12+01:00 [WARN]  hcl: I am visible when started by go test
	// 2022-02-25T09:40:12+01:00 [TRACE] hcl: now you can see me
	// 2022-02-25T09:40:12+01:00 [INFO]  hcl: I look the same as hcl.Printf
	// 2022-02-25T09:40:12+01:00 [INFO]  hcl.web: Start of web logs
}

func ExampleNew() {
	log := hcl.New()
	log.Info("I am a logger named after the application")

	log2 := hcl.ResetNamed("app-name")
	log2.Warnf("Mostly the same as above: %s", "but not as clear")
}
