package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type argument struct {
	pos       int
	posInFile int
	ident     *ast.Ident
}

type arguments []argument

func (a arguments) Len() int {
	return len(a)
}

func (a arguments) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//actually reversed because I want farthest down first
func (a arguments) Less(i, j int) bool {
	return a[i].posInFile > a[j].posInFile
}

type Gatherer map[string]map[string]argument

func (g Gatherer) Visit(node ast.Node) (w ast.Visitor) {
	switch ty := node.(type) {
	case *ast.AssignStmt:
		for i, n := range ty.Lhs {
			switch t := n.(type) {
			case *ast.Ident:
				if t.Name == "_" {
					fun, ok := ty.Rhs[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr)
					if ok {
						packageName := fun.X.(*ast.Ident).Name
						funcName := fun.Sel.Name
						funcMap, ok := g[packageName]
						if !ok {
							g[packageName] = make(map[string]argument)
							funcMap, _ = g[packageName]
						}
						funcMap[funcName] = argument{i, int(t.Pos()), t}
					}
				}
			}
		}
	}
	return g
}

func intMapToBoolMap(in map[string]argument) map[string]bool {
	out := make(map[string]bool)
	for f, _ := range in {
		out[f] = true
	}
	return out
}

func main() {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, os.Args[1], nil, 0)
	var g Gatherer
	g = make(map[string]map[string]argument)
	ast.Walk(g, file)
	for pack, f := range g {
		path, _, _ := findImport(pack, intMapToBoolMap(f))
		p, _ := build.Import(path, "", build.FindOnly)
		files, _ := filepath.Glob(p.Dir + "/*.go")
		s := Searcher(f)
		for _, filee := range files {
			fset = token.NewFileSet()
			astFile, _ := parser.ParseFile(fset, filee, nil, 0)
			ast.Walk(s, astFile)
		}
	}

	args := make(arguments, 0)
	//name them err and collect them
	for _, m := range g {
		for _, arg := range m {
			arg.ident.Name = "err"
			args = append(args, arg)
		}
	}
	sort.Sort(args)
	fset = token.NewFileSet()
	buff := bytes.NewBuffer(make([]byte, 0))
	printer.Fprint(buff, fset, file)
	s := buff.String()
	//Now we insert the panic!
	for _, arg := range args {
		nextLine := strings.Index(s[arg.posInFile:], "\n")
		s = s[:nextLine+arg.posInFile] + "\nif err != nil {\npanic(err)\n}" + s[nextLine+arg.posInFile:]
	}
	fmt.Println(s)
}

type Searcher map[string]argument

func (s Searcher) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.FuncDecl:
		temp, ok := s[t.Name.Name]
		pos := temp.pos
		if ok {
			typePos := t.Type.Results.List[pos].Type.(*ast.Ident)
			if typePos.Name != "error" {
				delete(s, t.Name.Name)
			}
		}
	}
	return s
}
