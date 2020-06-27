package binclude_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lu4p/binclude"
	"github.com/lu4p/binclude/example"
)

var BinFS = example.BinFS

func ExampleFileSystem_Open() {
	binclude.Include("./assets")
	f, _ := BinFS.Open("./assets/asset1.txt")
	data, _ := ioutil.ReadAll(f)
	fmt.Println(string(data))
	// Output: asset1
}

func ExampleFileSystem_ReadFile() {
	binclude.Include("file.txt")
	data, _ := BinFS.ReadFile("file.txt")
	fmt.Println(string(data))
	// Output: file.txt
}

func ExampleFileSystem_ReadDir() {
	binclude.Include("./assets")
	infos, _ := BinFS.ReadDir("./assets")
	for _, info := range infos {
		fmt.Println(info.Name())
	}
	// Output: asset1.txt
	// asset2.txt
	// logo_nocompress.png
	// subdir
}

func ExampleFileSystem_CopyFile() {
	BinFS.CopyFile("./assets/asset1.txt", "asset1.txt")

	c, _ := ioutil.ReadFile("asset1.txt")

	fmt.Println(string(c))

	// Output: asset1
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
	startContent, err := BinFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	err = BinFS.Compress(binclude.None)
	if err != nil {
		t.Fatal(err)
	}

	nocompressContent, err := BinFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(startContent) != string(nocompressContent) {
		t.Fatal("Compression with binclude.None should be a noop")
	}

	err = BinFS.Compress(binclude.Gzip)
	if err != nil {
		t.Fatal(err)
	}

	if BinFS.Files["assets/logo_nocompress.png"].Compression != binclude.None {
		t.Fatal("Unexpected compressed png")
	}

	gzipContent, err := BinFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(startContent) == string(gzipContent) {
		t.Fatal("Gzip didn't compress")
	}

	err = BinFS.Decompress()
	if err != nil {
		t.Fatal(err)
	}

	decGzipContent, err := BinFS.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(startContent) != string(decGzipContent) {
		t.Fatal("File differs after compression and decompression.")
	}

}

func TestReadFile(t *testing.T) {
	_, err := BinFS.ReadFile("nonexistent.txt")
	if err == nil {
		t.Fatal("shouldn't be able to read nonexistent file")
	}
}

func TestOpen(t *testing.T) {
	f, err := BinFS.Open("file.txt")
	if err != nil {
		t.Fatal(err)
	}

	data, _ := ioutil.ReadAll(f)
	if string(data) != "file.txt" {
		t.Fatal("content does not match")
	}

	f.Close()

	f, err = BinFS.Open("nonexistent.txt")
	if err == nil {
		f.Close()
		t.Fatal("nonexistent file can be opened")
	}

	binclude.Debug = true
	f, err = BinFS.Open("go.mod")
	if err != nil {
		t.Fatal("cannot use os filesystem")
	}
	f.Close()
	binclude.Debug = false

}

func TestReadDir(t *testing.T) {
	_, err := BinFS.ReadDir("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot read directory of file")
	}

	_, err = BinFS.ReadDir("./nonexistent")
	if err == nil {
		t.Fatal("shouldn't be able to read nonexistent dir")
	}
}

func TestStat(t *testing.T) {
	_, err := BinFS.Stat("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot stat file")
	}

	_, err = BinFS.Stat("./nonexistent")
	if err == nil {
		t.Fatal("shouldn't be able to stat nonexistent dir")
	}
}

func TestFileInfo(t *testing.T) {
	file := BinFS.Files["assets/asset1.txt"]
	info, err := BinFS.Stat("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot stat file")
	}

	if info.Name() != file.Filename {
		t.Fatal("Name does not match")
	}

	if int(info.Size()) != len(file.Content) {
		t.Fatal("Size does not match")
	}

	if info.Mode() != file.Mode {
		t.Fatal("Mode does not match")
	}

	if info.IsDir() != file.Mode.IsDir() {
		t.Fatal("IsDir does not match")
	}

	if info.ModTime() != file.ModTime {
		t.Fatal("ModTime does not match")
	}

	if info.Sys() != nil {
		t.Fatal("Sys return should be nil")
	}
}
