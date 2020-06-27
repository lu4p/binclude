# binclude

[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/lu4p/binclude)
[![Test](https://github.com/lu4p/binclude/workflows/Test/badge.svg)](https://github.com/lu4p/binclude/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/lu4p/binclude)](https://goreportcard.com/report/github.com/lu4p/binclude)
[![codecov](https://codecov.io/gh/lu4p/binclude/branch/master/graph/badge.svg)](https://codecov.io/gh/lu4p/binclude)

binclude is a tool for including static files into Go binaries.
- focuses on ease of use
- the bincluded files add no more than the filesize to the binary
- uses ast for typesafe parsing and code generation
- each package can have its own `binclude.FileSystem`
- `binclude.FileSystem` implements the `http.FileSystem` interface
- `ioutil` like functions `FileSystem.ReadFile`, `FileSystem.ReadDir`
- include all files/ directories under a given path by calling `binclude.Include("./path")`
- high test coverage
- supports execution of executables directly from a `binclude.FileSystem` (os/exec wrapper)
- optional compression of files with gzip `binclude -gzip`

## Install
```
GO111MODULE=on go get -u github.com/lu4p/binclude/cmd/binclude
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

	f, err := BinFS.Open("file.txt")
	if err != nil {
		log.Fatalln(err)
	}

	out, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(out))

	infos, err := BinFS.ReadDir("./assets")
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
go generate
go build
```
## Binary size
The resulting binary can get quite large, with the included files. 

You can reduce the final binary size by building without debug info (`go build -ldflags "-s -w"`) and compressing the resulting binary with [upx](https://upx.github.io/) (`upx binname`).

**Note:** If you don't need to access the compressed form of the files I would advise to just use [upx](https://upx.github.io/) and don't add seperate compression to the files. 

You can add compression to the included files with `-gzip`

**Note:** decompression is optional to allow for the scenario where you want to serve compressed files for a webapp directly.


## OS / Arch Specific Includes

Binclude supports including files/binaries only on specific architectures and OS. Binclude follows the same pattern as [Go's Build Constraints](https://golang.org/pkg/go/build/#hdr-Build_Constraints) and will generate files like `binclude_windows.go` which contains all windows specific files:
```
*_GOOS
*_GOARCH
*_GOOS_GOARCH
```

For example, if you want to include a binary only on Windows you could have a file `static_windows.go` and reference the static file:
```go
package main

import "github.com/lu4p/binclude"

func bar() {
  binclude.Include("./windows-file.dll")
}
```

## Advanced Usage
The generator can also be included into your package to allow for the code generator to run after all module dependencies are installed.
Without installing the binclude generator to the PATH seperatly.

main.go:
```go
// +build !gen

//go:generate go run -tags=gen .
package main

import (
	"github.com/lu4p/binclude"
)

func main() {
	binclude.Include("./assets")
}

```

main_gen.go:
```go
// +build gen

package main

import (
	"github.com/lu4p/binclude"
	"github.com/lu4p/binclude/bincludegen"
)

func main() {
	bincludegen.Generate(binclude.None)
	// binclude.None means no compression there are also binclude.Gzip
}
```

To build use:
```
go generate
go build
```
