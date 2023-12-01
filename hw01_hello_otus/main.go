package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	greeting := "Hello, OTUS!"
	reversed := reverse.String(greeting)
	fmt.Println(reversed)
}
