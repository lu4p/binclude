package main

import (
	"github.com/lu4p/binclude"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"strconv"
)

func dataAsByteSlice(data []byte) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: "`" + binclude.ByteToStr(data) + "`",
	}
}

func fileToAst(path string, file *binclude.File, num int) (c *ast.ValueSpec, m *ast.KeyValueExpr) {
	constName := "_binclude" + strconv.Itoa(num)
	var x ast.Expr = &ast.Ident{Name: "nil"}

	if !file.Mode.IsDir() {
		c = &ast.ValueSpec{
			Names: []*ast.Ident{
				{
					Name: constName,
				},
			},
			Values: []ast.Expr{
				dataAsByteSlice(file.Content),
			},
		}

		x = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "b",
				},
				Sel: &ast.Ident{
					Name: "StrToByte",
				},
			},
			Args: []ast.Expr{
				&ast.Ident{
					Name: constName,
				},
			},
		}
	}

	m = &ast.KeyValueExpr{
		Key: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + path + `"`,
		},
		Value: &ast.CompositeLit{
			Elts: []ast.Expr{
				&ast.KeyValueExpr{
					Key: &ast.Ident{
						Name: "Filename",
					},
					Value: &ast.BasicLit{
						Kind:  token.STRING,
						Value: `"` + file.Filename + `"`,
					},
				},
				&ast.KeyValueExpr{
					Key: &ast.Ident{
						Name: "Mode",
					},
					Value: &ast.BasicLit{
						Kind:  token.INT,
						Value: strconv.Itoa(int(file.Mode)),
					},
				},
				&ast.KeyValueExpr{
					Key: &ast.Ident{
						Name: "ModTime",
					},
					Value: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "time",
							},
							Sel: &ast.Ident{
								Name: "Unix",
							},
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: strconv.Itoa(int(file.ModTime.Unix())),
							},
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "0",
							},
						},
						Ellipsis: 0,
					},
				},
				&ast.KeyValueExpr{
					Key: &ast.Ident{
						Name: "Content",
					},
					Value: x,
				},
			},
		},
	}

	return c, m
}

func generateFile(pkgName *ast.Ident, fs binclude.FileSystem) error {
	var (
		astConsts []ast.Spec
		astFiles  []ast.Expr
	)

	num := 0
	for path, file := range fs {
		astConst, astFile := fileToAst(path, file, num)
		if !file.Mode.IsDir() {
			astConsts = append(astConsts, astConst)
			num++
		}

		astFiles = append(astFiles, astFile)
	}

	bincludeFile := &ast.File{
		Name: pkgName,
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: "b",
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/lu4p/binclude\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"time\"",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok:   token.CONST,
				Specs: astConsts,
			},
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "binFS",
							},
						},
						Values: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "b",
									},
									Sel: &ast.Ident{
										Name: "FileSystem",
									},
								},
								Elts: astFiles,
							},
						},
					},
				},
			},
		},
	}

	f, err := os.OpenFile("binclude.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	err = printer.Fprint(f, fset, bincludeFile)
	if err != nil {
		return err
	}

	return nil

}
