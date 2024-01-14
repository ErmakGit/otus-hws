package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

const slashStr = `\`

func Unpack(input string) (string, error) {
	var res strings.Builder
	skipNext := false
	skipCheck := false
	runes := []rune(input)
	for ind, char := range runes {
		if skipNext {
			skipNext = false
			skipCheck = false
			continue
		}

		if _, err := strconv.Atoi(string(char)); err == nil && !skipCheck {
			return "", ErrInvalidString
		}

		if ind+1 == len(runes) {
			if string(char) == slashStr && !skipCheck {
				return "", ErrInvalidString
			}
			res.WriteRune(char)
			continue
		}

		if string(char) == slashStr && !skipCheck {
			if _, err := strconv.Atoi(string(runes[ind+1])); err != nil && string(runes[ind+1]) != slashStr {
				return "", ErrInvalidString
			}
			skipCheck = true
			continue
		}

		if count, err := strconv.Atoi(string(runes[ind+1])); err == nil {
			res.WriteString(strings.Repeat(string(char), count))
			skipNext = true
		} else {
			res.WriteRune(char)
		}

		skipCheck = false
	}

	return res.String(), nil
}
