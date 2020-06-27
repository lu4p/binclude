package bincludegen_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/lu4p/binclude/bincludegen"
	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"binclude": bincludegen.Main1,
	}))
}

var update = flag.Bool("u", false, "update testscript output files")

func TestScripts(t *testing.T) {
	t.Parallel()

	p := testscript.Params{
		Dir: filepath.Join("testdata", "scripts"),
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"bincmp": bincmp,
		},
		UpdateScripts: *update,
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}

func bincmp(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: bincmp file1 file2")
	}
	data1 := ts.ReadFile(args[0])
	data2 := ts.ReadFile(args[1])
	if neg {
		if data1 == data2 {
			ts.Fatalf("%s and %s don't differ",
				args[0], args[1])
		}
		return
	}
	if data1 != data2 {
		sizeDiff := len(data2) - len(data1)
		ts.Fatalf("%s and %s differ; size diff: %+d",
			args[0], args[1], sizeDiff)
	}
}
