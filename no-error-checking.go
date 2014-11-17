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
	_, _, _ = d()
	_, _, err := f()
	fmt.Println(err)
	_ = e()
}

var f = d

//comments
func d() (int, bool, bool) {
	return 0, false, false
}

func e() bool {
	return false
}

type example struct {
	pos       int
	posInFile int
	ident     *ast.Ident
}
