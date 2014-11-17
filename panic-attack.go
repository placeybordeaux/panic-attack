package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type argument struct {
	pos       int
	posInFile int
	ident     *ast.Ident
	verified  bool
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
					callExpr, ok := ty.Rhs[0].(*ast.CallExpr)
					if ok {
						fun, ok := callExpr.Fun.(*ast.SelectorExpr)
						if ok { //Is contained in another package
							packageName := fun.X.(*ast.Ident).Name
							funcName := fun.Sel.Name
							funcMap, ok := g[packageName]
							if !ok {
								g[packageName] = make(map[string]argument)
								funcMap, _ = g[packageName]
							}
							funcMap[funcName] = argument{i, int(t.Pos()), t, false}
						} else { //Is local
							packageName := "LOCAL"
							funcName := callExpr.Fun.(*ast.Ident).Name
							funcMap, ok := g[packageName]
							if !ok {
								g[packageName] = make(map[string]argument)
								funcMap, _ = g[packageName]
							}
							funcMap[funcName] = argument{i, int(t.Pos()), t, false}
						}
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
	var s string
	//for no file passed in
	if len(os.Args) == 1 {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		s, _ = ParseSource(string(b))
	} else {
		s, _ = ParseFile(os.Args[1])
	}
	fmt.Println(s)
}

func ParseSource(source string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, 0)
	if err != nil {
		return "", err
	}
	var g Gatherer
	g = make(map[string]map[string]argument)
	ast.Walk(g, file)
	spew.Dump(g)
	g.trimNonErrors()

	args := make(arguments, 0)
	//name them err and collect them
	for _, m := range g {
		for _, arg := range m {
			args = append(args, arg)
		}
	}
	sort.Sort(args)
	//fset = token.NewFileSet()
	s := source
	//Now we insert the panic!
	for _, arg := range args {
		if arg.verified == false {
			continue
		}
		nextLine := strings.Index(s[arg.posInFile:], "\n")
		s = s[:nextLine+arg.posInFile] + "\nif err != nil {\npanic(err)\n}" + s[nextLine+arg.posInFile:] //insert the err
		s = s[:arg.posInFile-1] + "err" + s[arg.posInFile:]
	}
	return s, nil

}

func ParseFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return ParseSource(string(b))
}

func (g *Gatherer) trimNonErrors() {
	for pack, f := range *g {
		var files []string
		var err error
		if pack == "LOCAL" {
			files, err = filepath.Glob("*.go") //all files in this import
			if err != nil {
				panic(err)
			}
		} else {
			path, _, err := findImport(pack, intMapToBoolMap(f))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to find the import path for %v\n, skipping it", pack)
				continue
			}
			p, err := build.Import(path, "", build.FindOnly) //find the import's path
			if err != nil {
				if err.Error() == `import "": invalid import path` {
					continue
				}
				panic(err)
			}
			files, err = filepath.Glob(p.Dir + "/*.go") //all files in this import
			if err != nil {
				panic(err)
			}
		}
		s := Trimmer(f)
		for _, file := range files {
			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, file, nil, 0)
			if err != nil {
				panic(err)
			}
			ast.Walk(s, astFile)
		}
	}
}

type Trimmer map[string]argument

func (s Trimmer) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.FuncDecl:
		temp, ok := s[t.Name.Name]
		pos := temp.pos
		if ok {
			typePos := t.Type.Results.List[pos].Type.(*ast.Ident)
			if typePos.Name == "error" {
				temp.verified = true
				s[t.Name.Name] = temp
			}
		}
	}
	return s
}
