package lexer

import "strconv"

func MapToNumber(str []rune) (uint64, bool) {
	num, err := strconv.ParseUint(string(str), 10, 64)
	return num, err == nil
}
