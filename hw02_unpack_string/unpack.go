package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var res strings.Builder
	skipNext := false
	skipCheck := false
	for ind, char := range input {
		if skipNext {
			skipNext = false
			skipCheck = false
			continue
		}

		if _, err := strconv.Atoi(string(char)); err == nil && !skipCheck {
			return "", ErrInvalidString
		}

		if ind+1 == len(input) {
			res.WriteRune(char)
			continue
		}

		if string(char) == `\` && !skipCheck {
			skipCheck = true
			continue
		}

		if count, err := strconv.Atoi(string(input[ind+1])); err == nil {
			res.WriteString(strings.Repeat(string(char), count))
			skipNext = true
		} else {
			res.WriteRune(char)
		}

		skipCheck = false
	}

	return res.String(), nil
}
