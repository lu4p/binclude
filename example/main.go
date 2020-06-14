package main

//go:generate go run ../cmd/binclude/

import (
	"io/ioutil"
	"log"

	"github.com/lu4p/binclude"
)

func main() {
	binclude.Include("./assets")
	binclude.Include("file.txt")

	f, err := binFS.Open("file.txt")
	if err != nil {
		log.Fatalln(err)
	}

	out, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(out))

	infos, err := binFS.ReadDir("./assets")
	if err != nil {
		log.Fatalln(err)
	}

	for _, info := range infos {
		log.Println(info.Name())
	}
}
