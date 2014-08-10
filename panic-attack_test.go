package main

import "testing"

func TestSTDLib(t *testing.T) {
	input := `package main

import "fmt"

func a() {
        i, _ := strconv.Atoi("bad")
	fmt.Println(i)
}`
	expected := `package main

import "fmt"

func a() {
        i, err := strconv.Atoi("bad")
if err != nil {
panic(err)
}
	fmt.Println(i)
}`
	output, err := ParseSource(input)
	if err != nil {
		t.Fatal(err)
	}
	if output != expected {
		t.Fatalf("got %s expected %s", output, expected)
	}
}

func TestMultiRHS(t *testing.T) {
	input := `package main

import "fmt"

func a() (int, error) {
	return 1, nil
}

func b() {
	t, _, _ := 1, a()
	fmt.Println(t)
}
`

	expected := `package main

import "fmt"

func a() (int, error) {
	return 1, nil
}

func b() {
	t, _, err := 1, a()
if err != nil {
panic(err)
}
	fmt.Println(t)
}
`

	output, err := ParseSource(input)
	if err != nil {
		t.Fatal(err)
	}
	if output != expected {
		t.Fatalf("got %s expected %s", output, expected)
	}
}
