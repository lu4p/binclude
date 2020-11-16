package bincludegen_test

import (
	"testing"

	"github.com/lu4p/binclude"
	"github.com/lu4p/binclude/bincludegen"
)

func BenchmarkGenerate(b *testing.B) {
	err := bincludegen.Generate(binclude.None, "./testdata/bench")
	if err != nil {
		b.Fatal(err)
	}
}
