package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type Gatherer map[string]map[string]bool

func (g Gatherer) Visit(node ast.Node) (w ast.Visitor) {
	switch ty := node.(type) {
	case *ast.AssignStmt:
		for i, n := range ty.Lhs {
			switch t := n.(type) {
			case *ast.Ident:
				if t.Name == "_" {
					fun := ty.Rhs[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr)
					fmt.Printf("found _ as argument number %d of func %s.%s\n", i, fun.X, fun.Sel)
					packageName := fun.X.(*ast.Ident).Name
					funcName := fun.Sel.Name
					funcMap, ok := g[packageName]
					if !ok {
						g[packageName] = make(map[string]bool)
						funcMap, _ = g[packageName]
					}
					funcMap[funcName] = true
				}
			}
		}
	}
	return g
}

type ImportVisitor struct{}

func (v *ImportVisitor) Visit(node ast.Node) (w ast.Visitor) {

	switch ty := node.(type) {
	case *ast.ImportSpec:
		spew.Dump(ty.Path.Value)
	}
	return v
}

func main() {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, os.Args[1], nil, 0)
	var g Gatherer
	g = make(map[string]map[string]bool)
	ast.Walk(g, file)
	for pack, f := range g {
		path, b, err := findImport(pack, f)
		spew.Dump(path)
		spew.Dump(loadExports(path))
		spew.Dump(b)
		spew.Dump(err)
		p, err := build.Import(path, "", build.FindOnly)
		fmt.Printf("Needs to check package %s for functions %v", p.Dir, f)
	}

}
