package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type DeclVisitor struct{}

func (v *DeclVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch ty := node.(type) {
	case *ast.AssignStmt:
		for i, n := range ty.Lhs {
			switch t := n.(type) {
			case *ast.Ident:
				if t.Name == "_" {
					fun := ty.Rhs[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr)
					fmt.Printf("found _ as argument number %d of func %s.%s\n", i, fun.X, fun.Sel)
				}
			}
		}
	}
	return v
}

func main() {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, os.Args[1], nil, 0)
	printer.Fprint(os.Stdout, fset, file)
	ast.Walk(new(DeclVisitor), file)
}
