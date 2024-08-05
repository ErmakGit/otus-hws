package main

import (
	"errors"
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	err := Copy(from, to, offset, limit)
	if err != nil {
		if errors.Is(err, ErrOffsetOrLimitBellowZero) ||
			errors.Is(err, ErrPathIsEmpty) ||
			errors.Is(err, ErrUnsupportedFile) ||
			errors.Is(err, ErrOffsetExceedsFileSize) {
			fmt.Printf("Cannot copy file: %v", err)
			return
		}

		fmt.Printf("Unrecognized error: %v", err)
		return
	}
}
