package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lu4p/binclude"
)

var fset *token.FileSet

func main() {
	now := time.Now()
	defer fmt.Println("binclude finished in:", time.Since(now))

	exitCode := main1()

	defer os.Exit(exitCode)
}

func main1() int {
	log.SetPrefix("[binclude] ")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("could not get working directory:", err)
	}

	err = mainErr(wd)
	if err != nil {
		log.Println("failed:", err)
		return 1
	}

	return 0
}

func mainErr(path string) error {
	paths, err := filepath.Glob(filepath.Join(path, "*.go"))
	if err != nil {
		return err
	}

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

	paths, err = detectIncluded(files)
	if err != nil {
		return err
	}

	fs, err := buildFS(paths)
	if err != nil {
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

		fs[path] = &binclude.File{
			Filename: info.Name(),
			Mode:     info.Mode(),
			ModTime:  info.ModTime(),
			Content:  content,
		}

		return nil
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		err = filepath.Walk(path, walkFn)
		if err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func detectIncluded(files []*ast.File) ([]string, error) {
	var includedPaths []string

	wd, _ := os.Getwd()

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
		path = filepath.Join(wd, path)

		path, err = filepath.Rel(wd, path)
		if err != nil {
			return nil, err
		}

		includedPaths[i] = strings.TrimPrefix(path, "./")
	}

	return includedPaths, nil
}
