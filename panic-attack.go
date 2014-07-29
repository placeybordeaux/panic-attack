package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
)

type Gatherer map[string]map[string]int

func (g Gatherer) Visit(node ast.Node) (w ast.Visitor) {
	switch ty := node.(type) {
	case *ast.AssignStmt:
		for i, n := range ty.Lhs {
			switch t := n.(type) {
			case *ast.Ident:
				if t.Name == "_" {
					callExpr, ok := ty.Rhs[0].(*ast.CallExpr)
					if ok {
						fun, ok := callExpr.Fun.(*ast.SelectorExpr)
						if ok {
							fmt.Printf("found _ as argument number %d of func %s.%s\n", i, fun.X, fun.Sel)
							packageName := fun.X.(*ast.Ident).Name
							funcName := fun.Sel.Name
							funcMap, ok := g[packageName]
							if !ok {
								g[packageName] = make(map[string]int)
								funcMap, _ = g[packageName]
							}
							funcMap[funcName] = i
						}
					}
				}
			}
		}
	}
	return g
}

func intMapToBoolMap(in map[string]int) map[string]bool {
	out := make(map[string]bool)
	for f, _ := range in {
		out[f] = true
	}
	return out
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
	g = make(map[string]map[string]int)
	ast.Walk(g, file)
	for pack, f := range g {
		path, _, _ := findImport(pack, intMapToBoolMap(f))
		p, _ := build.Import(path, "", build.FindOnly)
		fmt.Printf("Needs to check package %s for functions %v\n", p.Dir, f)
		files, _ := filepath.Glob(p.Dir + "/*.go")
		for _, file := range files {
			fset = token.NewFileSet()
			astFile, _ := parser.ParseFile(fset, file, nil, 0)
			s := Searcher(f)
			ast.Walk(s, astFile)
		}
	}
}

type Searcher map[string]int

func (s Searcher) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.FuncDecl:
		pos, ok := s[t.Name.Name]
		if ok {
			typePos := t.Type.Results.List[pos].Type.(*ast.Ident)
			if typePos.Name == "error" {
				fmt.Printf("must replace %+v at %v\n", t, t.Pos())
			}
		}
	}
	return s
}
