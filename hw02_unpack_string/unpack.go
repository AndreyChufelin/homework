package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	result := strings.Builder{}
	runes := []rune(str)
	esc := false

	for i := 0; i < len(runes); i++ {
		var current rune
		isNextDigit := i < len(runes)-1 && unicode.IsDigit(runes[i+1])

		switch r := runes[i]; {
		case r == '\\':
			if esc {
				current = r
				esc = false
			} else {
				esc = true
			}
		case unicode.IsDigit(r):
			if esc {
				current = r

				esc = false
			} else {
				if i == 0 || isNextDigit {
					return "", ErrInvalidString
				}

				n, _ := strconv.Atoi(string(r))

				repeated := strings.Repeat(string(runes[i-1]), n)
				result.WriteString(repeated)
			}
		default:
			if esc {
				return "", ErrInvalidString
			}
			current = r
		}

		if current != 0 && !isNextDigit {
			result.WriteRune(current)
		}
	}

	return result.String(), nil
}
