package main

import (
	"binclude"
	"io/ioutil"
	"log"
)

func main() {
	binclude.Debug = true
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

	f, err = binFS.Open("./assets")
	f.Stat()

	log.Println(string(out))

	infos, err := binFS.ReadDir("./assets")
	if err != nil {
		log.Fatalln(err)
	}

	for _, info := range infos {
		log.Println(info.Name())
	}
}
