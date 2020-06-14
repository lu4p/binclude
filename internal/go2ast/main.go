package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	path := os.Args[1]

	log.Println("path", path)

	err := parseFile(path)
	if err != nil {
		log.Println(err)
	}
}

func parseFile(path string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	err = Print(fset, f.Decls)
	if err != nil {
		return err
	}

	return nil

}
