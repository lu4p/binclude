package binexec_test

import (
	"context"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/lu4p/binclude/binexec"
	"github.com/lu4p/binclude/binexec/example"
)

var BinFS = example.BinFS

var testprg = "testprg/testprg"

func init() {
	if runtime.GOOS == "windows" {
		testprg = filepath.FromSlash(testprg) + ".exe"
	}

	log.Println("Testprg:", testprg)
}

func TestRun(t *testing.T) {
	cmd, err := binexec.Command(BinFS, testprg)
	if err != nil {
		t.Fatal("cannot initialize cmd", err)
	}

	err = cmd.Run()
	if err != nil {
		t.Fatal("cannot execute cmd", err)
	}
}

func TestNonexistent(t *testing.T) {
	_, err := binexec.Command(BinFS, "Nonexistent")
	if err == nil {
		t.Fatal("can initialize cmd on Nonexistent file")
	}

	_, err = binexec.CommandContext(context.Background(), BinFS, "Nonexistent")
	if err == nil {
		t.Fatal("can initialize cmd on Nonexistent file")
	}

}

func TestStart(t *testing.T) {
	cmd, err := binexec.Command(BinFS, testprg)
	if err != nil {
		t.Fatal("cannot initialize cmd", err)
	}
	t.Log(cmd)

	_, err = cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		t.Fatal("cannot execute cmd", err)
	}

	out, err := ioutil.ReadAll(stderr)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ioutil.ReadAll(stdout); err != nil {
		t.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		t.Fatal("wait:", err)
	}

	if string(out) != "Hello world!\n" {
		t.Fatal("unexpected output:", string(out))
	}
}

func TestRunContext(t *testing.T) {
	cmd, err := binexec.CommandContext(context.Background(), BinFS, testprg)
	if err != nil {
		t.Fatal("cannot initialize cmd", err)
	}

	err = cmd.Run()
	if err != nil {
		t.Fatal("cannot execute cmd", err)
	}
}
