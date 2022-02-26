package hcl_test

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/vogtp/go-hcl"
)

// go test -bench=. -benchtime=20s

func BenchmarkHcl(b *testing.B) {
	f, _ := os.Create("temp")
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	hcl := hcl.New(hcl.WithWriter(f))
	for n := 0; n < b.N; n++ {
		hcl.Warn("hcl")
	}
}

func BenchmarkHclF(b *testing.B) {
	f, _ := os.Create("temp")
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	hcl := hcl.New(hcl.WithWriter(f))
	for n := 0; n < b.N; n++ {
		hcl.Warnf("hcl %d", n)
	}
}

func BenchmarkHclVar(b *testing.B) {
	f, _ := os.Create("temp")
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	hcl := hcl.New(hcl.WithWriter(f))
	for n := 0; n < b.N; n++ {
		hcl.Warn("hcl", "count", n)
	}
}

func BenchmarkDepHclog(b *testing.B) {

	f, _ := os.Create("temp")
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	hclog := hclog.New(&hclog.LoggerOptions{
		Output: f,
	})
	for n := 0; n < b.N; n++ {
		hclog.Error("hclog ")
	}
}

func BenchmarkDepHclogVar(b *testing.B) {

	f, _ := os.Create("temp")
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	hclog := hclog.New(&hclog.LoggerOptions{
		Output: f,
	})
	for n := 0; n < b.N; n++ {
		hclog.Error("hclog ", "count", n)
	}
}

func BenchmarkLog(b *testing.B) {
	f, _ := os.Create("temp")
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	log.SetOutput(f)
	// 7146 ns/op
	for n := 0; n < b.N; n++ {
		log.Printf("log")
	}
}

func BenchmarkLogF(b *testing.B) {
	f, _ := os.Create("temp")
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	log.SetOutput(f)
	// 7146 ns/op
	for n := 0; n < b.N; n++ {
		log.Printf("log %d", n)
	}
}

/*
cpu: AMD Ryzen 5 2600 Six-Core Processor
BenchmarkHcl-12                 10705644              2592 ns/op liblogger: check for nil
BenchmarkHcl-12                 10458598              2345 ns/op liblogger: keep track with bool
BenchmarkHclF-12                 8467414              2736 ns/op liblogger: check for nil
BenchmarkHclF-12                 8461070              2781 ns/op liblogger: check for nil
BenchmarkHclF-12                 8680538              2737 ns/op liblogger: keep track with bool
BenchmarkHclF-12                 8631999              2754 ns/op liblogger: keep track with bool
BenchmarkHclVar-12               6994975              3388 ns/op liblogger: check for nil
BenchmarkHclVar-12               6917926              3347 ns/op liblogger: keep track with bool
BenchmarkDepHclogVar-12          6730189              3397 ns/op
BenchmarkDepHclogVar-12          6982083              3413 ns/op
BenchmarkLogF-12                11285836              2205 ns/op
BenchmarkLogF-12                10238210              2210 ns/op
*/
