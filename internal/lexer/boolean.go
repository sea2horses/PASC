package lexer

import (
	"strings"
)

func MapToBool(str []rune) (bool, bool) {
	strl := strings.ToLower(string(str))

	switch strl {
	case "verdadero":
		return true, true
	case "falso":
		return false, true
	}

	return false, false
}
