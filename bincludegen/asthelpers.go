package bincludegen

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func intLiteral(value int) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(value),
	}
}

func stringLiteral(value string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: fmt.Sprintf("%q", value),
	}
}

// x[index]
func indexExpr(x, index ast.Expr) *ast.IndexExpr {
	return &ast.IndexExpr{
		X:     x,
		Index: index,
	}
}

// fun(arg)
func callExpr(fun ast.Expr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  fun,
		Args: args,
	}
}

// x.sel
func selExpr(x, sel string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ident(x),
		Sel: ident(sel),
	}
}

func exprStmt(expr ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{X: expr}
}

func blockStmt(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: stmts}
}

func genDecl(tok token.Token, specs []ast.Spec) *ast.GenDecl {
	return &ast.GenDecl{
		Tok:   tok,
		Specs: specs,
	}
}

func importSpec(path string) *ast.ImportSpec {
	return &ast.ImportSpec{
		Path: stringLiteral(path),
	}
}
