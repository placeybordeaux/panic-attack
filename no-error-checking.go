package main

import (
	"fmt"
	"strconv"
)

func example1() {
	b, _ := strconv.ParseBool("bad input")
	fmt.Println(b)
	_, _ = strconv.ParseInt("the number one", 64, 10)
}
