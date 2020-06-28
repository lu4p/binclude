package example

//go:generate go run ../cmd/binclude

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/lu4p/binclude"
)

var assetPath = binclude.Include("./assets") // include ./assets with all files and subdirectories

func main() {
	binclude.Include("file.txt") // include file.txt

	// like os.Open
	f, err := BinFS.Open("file.txt")
	if err != nil {
		log.Fatalln(err)
	}

	out, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(out))

	// like ioutil.Readfile
	content, err := BinFS.ReadFile(filepath.Join(assetPath, "asset1.txt"))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(content))

	// like ioutil.ReadDir
	infos, err := BinFS.ReadDir(assetPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, info := range infos {
		log.Println(info.Name())
	}
}
