package hcl_test

import (
	"log"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/vogtp/go-hcl"
)

func BenchmarkHcl(b *testing.B) {
	// 13716 ns/op
	for n := 0; n < b.N; n++ {
		hcl.Warnf("hcl %d", n)
	}
}
func BenchmarkHclog(b *testing.B) {
	// 15635 ns/op
	hclog := hclog.Default()
	for n := 0; n < b.N; n++ {
		hclog.Error("hclog ", "count", n)
	}
}
func BenchmarkLog(b *testing.B) {
	// 7146 ns/op
	for n := 0; n < b.N; n++ {
		log.Printf("log %d", n)
	}
}
