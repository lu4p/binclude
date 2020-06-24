package main

import (
	"errors"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lu4p/binclude"
)

var fset *token.FileSet
var compress = binclude.None

var brotli bool
var gzip bool

func init() {
	flag.BoolVar(&gzip, "gzip", false, "compress files with gzip")
	flag.BoolVar(&brotli, "brotli", false, "compress files with brotli")
}

func main() {
	os.Exit(main1())
}

func main1() int {
	flag.Parse()
	if brotli {
		compress = binclude.Brotli
	} else if gzip {
		compress = binclude.Gzip
	}
	log.SetPrefix("[binclude] ")

	err := mainErr()
	if err != nil {
		log.Println("failed:", err)
		return 1
	}

	return 0
}

func mainErr() error {
	paths, _ := filepath.Glob("*.go")

	if len(paths) == 0 {
		return errors.New("No .go files found in the current directory")
	}

	fset = token.NewFileSet()

	var files []*ast.File
	for _, path := range paths {
		if strings.HasSuffix(path, "binclude.go") {
			continue
		}
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		files = append(files, file)
	}

	pkgName := files[0].Name

	paths, err := detectIncluded(files)
	if err != nil {
		return err
	}

	fs, err := buildFS(paths)
	if err != nil {
		return err
	}

	if err := fs.Compress(compress); err != nil {
		return err
	}

	return generateFile(pkgName, fs)
}

func buildFS(paths []string) (binclude.FileSystem, error) {
	fs := make(binclude.FileSystem)

	var walkFn filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var content []byte

		if !info.IsDir() {
			content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		}

		path = filepath.ToSlash(path)

		fs[path] = &binclude.File{
			Filename: info.Name(),
			Mode:     info.Mode(),
			ModTime:  info.ModTime(),
			Content:  content,
		}

		return nil
	}

	for _, path := range paths {
		err := filepath.Walk(path, walkFn)
		if err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func detectIncluded(files []*ast.File) ([]string, error) {
	var includedPaths []string

	visit := func(node ast.Node) bool {
		if node == nil {
			return true
		}

		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		v, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if sel.Sel.Name != "Include" || v.Name != "binclude" {
			return true
		}

		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			log.Fatalln("argument is not string literal")
		}

		value, err := strconv.Unquote(lit.Value)
		if err != nil {
			log.Fatalln("cannot unquote string:", err)
		}

		includedPaths = append(includedPaths, value)

		return true
	}

	for _, file := range files {
		ast.Inspect(file, visit)
	}

	for i, path := range includedPaths {
		var err error

		if filepath.IsAbs(path) {
			return nil, errors.New("only supports relative include paths")
		}

		_, err = os.Stat(path)
		if err != nil {
			return nil, err
		}

		includedPaths[i] = strings.TrimPrefix(path, "./")
	}

	return includedPaths, nil
}
