package main

import (
	"fmt"

	"github.com/golang/example/stringutil"
)

func main() {
	greeting := "Hello, OTUS!"
	reversed := stringutil.Reverse(greeting)
	fmt.Println(reversed)
}
