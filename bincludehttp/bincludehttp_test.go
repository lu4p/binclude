package bincludehttp_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lu4p/binclude"
	"github.com/lu4p/binclude/bincludehttp"
	"github.com/lu4p/binclude/example"
)

var BinFS = example.BinFS
var wrappedFS = bincludehttp.Wrap(BinFS)

func ExampleFileSystem_Open() {
	binclude.Include("./assets")

	fs := bincludehttp.Wrap(BinFS)
	f, _ := fs.Open("./assets/asset1.txt")
	data, _ := ioutil.ReadAll(f)
	fmt.Println(string(data))
	// Output: asset1
}

func ExampleFileSystem_ReadFile() {
	binclude.Include("file.txt")
	fs := bincludehttp.Wrap(BinFS)
	data, _ := fs.ReadFile("file.txt")
	fmt.Println(string(data))
	// Output: file.txt
}

func ExampleFileSystem_ReadDir() {
	binclude.Include("./assets")
	fs := bincludehttp.Wrap(BinFS)
	infos, _ := fs.ReadDir("./assets")
	for _, info := range infos {
		fmt.Println(info.Name())
	}
	// Output: asset1.txt
	// asset2.txt
	// logo_nocompress.png
	// subdir
}

func TestCopyFile(t *testing.T) {
	err := BinFS.CopyFile("./assets/asset1.txt", "asset1.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("asset1.txt")

	c, err := ioutil.ReadFile("asset1.txt")
	if err != nil {
		t.Fatal(err)
	}

	if string(c) != "asset1" {
		t.Fatal("content doesn't match", c)
	}

	err = BinFS.CopyFile("nonexistent.txt", "nonexistent.txt")
	if err == nil {
		t.Fatal("can copy nonexistent file")
	}
}

func TestCompression(t *testing.T) {
	testPath := "assets/asset1.txt"
	startContent, err := wrappedFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	err = wrappedFS.Compress(binclude.None)
	if err != nil {
		t.Fatal(err)
	}

	nocompressContent, err := wrappedFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(startContent) != string(nocompressContent) {
		t.Fatal("Compression with binclude.None should be a noop")
	}

	err = wrappedFS.Compress(binclude.Gzip)
	if err != nil {
		t.Fatal(err)
	}

	gzipContent, err := wrappedFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(startContent) == string(gzipContent) {
		t.Fatal("Gzip didn't compress")
	}

	err = wrappedFS.Decompress()
	if err != nil {
		t.Fatal(err)
	}

	decGzipContent, err := wrappedFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(startContent) != string(decGzipContent) {
		t.Fatal("File differs after compression and decompression.")
	}

}

func TestReadFile(t *testing.T) {
	_, err := wrappedFS.ReadFile("nonexistent.txt")
	if err == nil {
		t.Fatal("shouldn't be able to read nonexistent file")
	}
}

func TestOpen(t *testing.T) {
	f, err := wrappedFS.Open("file.txt")
	if err != nil {
		t.Fatal(err)
	}

	data, _ := ioutil.ReadAll(f)
	if string(data) != "file.txt" {
		t.Fatal("content does not match")
	}

	f.Close()

	f, err = wrappedFS.Open("nonexistent.txt")
	if err == nil {
		f.Close()
		t.Fatal("nonexistent file can be opened")
	}

	binclude.Debug = true
	f, err = wrappedFS.Open("../go.mod")
	if err != nil {
		t.Fatal("cannot use os filesystem")
	}
	f.Close()
	binclude.Debug = false

}

func TestReadDir(t *testing.T) {
	_, err := wrappedFS.ReadDir("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot read directory of file")
	}

	_, err = wrappedFS.ReadDir("./nonexistent")
	if err == nil {
		t.Fatal("shouldn't be able to read nonexistent dir")
	}
}

func TestStat(t *testing.T) {
	_, err := wrappedFS.Stat("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot stat file")
	}

	_, err = wrappedFS.Stat("./nonexistent")
	if err == nil {
		t.Fatal("shouldn't be able to stat nonexistent dir")
	}
}

func TestFileInfo(t *testing.T) {
	info, err := wrappedFS.Stat("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot stat file")
	}

	if info.Name() == "" {
		t.Fatal("No name")
	}

	if int(info.Size()) == 0 {
		t.Fatal("Size does not match")
	}

	if info.IsDir() {
		t.Fatal("IsDir does not match")
	}

	if info.Sys() != nil {
		t.Fatal("Sys return should be nil")
	}
}
