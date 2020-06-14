package binclude_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/lu4p/binclude"
)

var (
	_binclude0 = []byte{102, 105, 108, 101, 46, 116, 120, 116}
	_binclude1 = []byte{97, 115, 115, 101, 116, 49}
	_binclude2 = []byte{97, 115, 115, 101, 116, 49}
	_binclude3 = []byte{115, 117, 98, 100, 105, 114, 97, 115, 115, 101, 116, 49}
	_binclude4 = []byte{115, 117, 98, 100, 105, 114, 97, 115, 115, 101, 116, 49}
)
var binFS = binclude.FileSystem{"file.txt": {Filename: "file.txt", Mode: 420, ModTime: time.Unix(1592104011, 0), Content: _binclude0}, "assets": {Filename: "assets", Mode: 2147484141, ModTime: time.Unix(1592147369, 0), Content: nil}, "assets/asset1.txt": {Filename: "asset1.txt", Mode: 420, ModTime: time.Unix(1592103949, 0), Content: _binclude1}, "assets/asset2.txt": {Filename: "asset2.txt", Mode: 420, ModTime: time.Unix(1592103958, 0), Content: _binclude2}, "assets/subdir": {Filename: "subdir", Mode: 2147484141, ModTime: time.Unix(1592104033, 0), Content: nil}, "assets/subdir/subdirasset1.txt": {Filename: "subdirasset1.txt", Mode: 420, ModTime: time.Unix(1592104027, 0), Content: _binclude3}, "assets/subdir/subdirasset2.txt": {Filename: "subdirasset2.txt", Mode: 420, ModTime: time.Unix(1592104033, 0), Content: _binclude4}}

func ExampleFileSystem_Open() {
	binclude.Include("./assets")
	f, _ := binFS.Open("./assets/asset1.txt")
	data, _ := ioutil.ReadAll(f)
	fmt.Println(string(data))
	// Output: asset1
}

func ExampleFileSystem_ReadFile() {
	binclude.Include("file.txt")
	data, _ := binFS.ReadFile("file.txt")
	fmt.Println(string(data))
	// Output: file.txt
}

func ExampleFileSystem_ReadDir() {
	binclude.Include("./assets")
	infos, _ := binFS.ReadDir("./assets")
	for _, info := range infos {
		fmt.Println(info.Name())
	}
	// Output: asset1.txt
	// asset2.txt
	// subdir
}

func TestReadFile(t *testing.T) {
	_, err := binFS.ReadFile("nonexistent.txt")
	if err == nil {
		t.Fatal("shouldn't be able to read nonexistent file")
	}
}

func TestOpen(t *testing.T) {
	f, err := binFS.Open("file.txt")
	if err != nil {
		t.Fatal(err)
	}

	data, _ := ioutil.ReadAll(f)
	if string(data) != "file.txt" {
		t.Fatal("content does not match")
	}

	f.Close()

	f, err = binFS.Open("nonexistent.txt")
	if err == nil {
		f.Close()
		t.Fatal("nonexistent file can be opened")
	}

	binclude.Debug = true
	f, err = binFS.Open("go.mod")
	if err != nil {
		t.Fatal("cannot use os filesystem")
	}
	f.Close()
	binclude.Debug = false

}

func TestReadDir(t *testing.T) {
	_, err := binFS.ReadDir("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot read directory of file")
	}

	_, err = binFS.ReadDir("./nonexistent")
	if err == nil {
		t.Fatal("shouldn't be able to read nonexistent dir")
	}
}

func TestStat(t *testing.T) {
	_, err := binFS.Stat("./assets/asset1.txt")
	if err != nil {
		t.Fatal("cannot stat file")
	}

	_, err = binFS.Stat("./nonexistent")
	if err == nil {
		t.Fatal("shouldn't be able to stat nonexistent dir")
	}
}

func TestFileInfo(t *testing.T) {
	file := binFS["assets/asset1.txt"]
	info, err := binFS.Stat("./assets/asset1.txt")
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
