package binclude_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/lu4p/binclude"
)

const (
	_binclude0 = `H4sIAAAAAAAA/0osLk4tMQQEAAD//zJdA+EGAAAA`
	_binclude1 = `H4sIAAAAAAAA/0osLk4tMQQEAAD//zJdA+EGAAAA`
	_binclude2 = `H4sIAAAAAAAA/youTUrJLEosLk4tMQQEAAD//zzpe5cMAAAA`
	_binclude3 = `H4sIAAAAAAAA/youTUrJLEosLk4tMQQEAAD//zzpe5cMAAAA`
	_binclude4 = `H4sIAAAAAAAA/0rLzEnVK6koAQQAAP//JRb34AgAAAA=`
)

var binFS = binclude.FileSystem{"assets/asset1.txt": {Filename: "asset1.txt", Mode: 420, ModTime: time.Unix(1592103949, 0), Content: binclude.StrToByte(_binclude0)}, "assets/asset2.txt": {Filename: "asset2.txt", Mode: 420, ModTime: time.Unix(1592103958, 0), Content: binclude.StrToByte(_binclude1)}, "assets/subdir": {Filename: "subdir", Mode: 2147484141, ModTime: time.Unix(1592104033, 0), Content: nil}, "assets/subdir/subdirasset1.txt": {Filename: "subdirasset1.txt", Mode: 420, ModTime: time.Unix(1592104027, 0), Content: binclude.StrToByte(_binclude2)}, "assets/subdir/subdirasset2.txt": {Filename: "subdirasset2.txt", Mode: 420, ModTime: time.Unix(1592104033, 0), Content: binclude.StrToByte(_binclude3)}, "file.txt": {Filename: "file.txt", Mode: 420, ModTime: time.Unix(1592104011, 0), Content: binclude.StrToByte(_binclude4)}, "assets": {Filename: "assets", Mode: 2147484141, ModTime: time.Unix(1592097652, 0), Content: nil}}

func TestBinary(t *testing.T) {
	data := make([]byte, 512)
	rand.Read(data)
	dataStr := binclude.ByteToStr(data)

	dataAfter := binclude.StrToByte(dataStr)

	if !bytes.Equal(data, dataAfter) {
		t.Fatal("data is different from dataAfter")
	}
}

func ExampleFileSystem_Open() {
	f, _ := binFS.Open("file.txt")
	data, _ := ioutil.ReadAll(f)
	fmt.Println(string(data))
	// Output: file.txt
}

func ExampleFileSystem_ReadFile() {
	data, _ := binFS.ReadFile("file.txt")
	fmt.Println(string(data))
	// Output: file.txt
}

func ExampleFileSystem_ReadDir() {
	infos, _ := binFS.ReadDir("./assets")
	for _, info := range infos {
		fmt.Println(info.Name())
	}
	// Output: asset1.txt
	// asset2.txt
	// subdir
}
