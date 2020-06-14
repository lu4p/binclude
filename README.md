# binclude

[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/lu4p/binclude)
[![License](https://img.shields.io/github/license/lu4p/binclude.svg)](https://unlicense.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/lu4p/binclude)](https://goreportcard.com/report/github.com/lu4p/binclude)
[![codecov](https://codecov.io/gh/lu4p/binclude/branch/master/graph/badge.svg)](https://codecov.io/gh/lu4p/binclude)

binclude is a tool for including static files into Go binaries.
- focuses on ease of use
- the bincluded files add no more than the filesize to the binary (straight binary literals)
- uses ast for typesafe parsing and code generation
- each package can have its own `binclude.FileSystem`
- `binclude.FileSystem` implements the `http.FileSystem` interface
- `ioutil` like functions `FileSystem.ReadFile`, `FileSystem.ReadDir`


## Install
```
$ go get github.com/lu4p/binclude
```
## Usage
```go
package main

//go:generate binclude

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

```
To build use:
```
$ go generate
$ go build
```
