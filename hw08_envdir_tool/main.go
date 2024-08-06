package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("too few arguments")
	}
	path := os.Args[1]

	environment, err := ReadDir(path)
	if err != nil {
		log.Fatal("error reading dir")
	}
	returnCode := RunCmd(os.Args[2:], environment)

	os.Exit(returnCode)
}
