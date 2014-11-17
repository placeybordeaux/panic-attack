package main

import (
	"fmt"
	"go/ast"
	"strconv"
)

type example2 struct {
	pos       int
	posInFile int
	ident     *ast.Ident
}

func a() {}

func b() {}

func c() {}

//comments
func example1() {
	b, _ := strconv.ParseBool("bad input")
	fmt.Println(b)
	_, _ = strconv.ParseInt("the number one", 64, 10)
	_ = d()
}

//comments
func d() error {
	return nil
}

func e() {}

type example struct {
	pos       int
	posInFile int
	ident     *ast.Ident
}
