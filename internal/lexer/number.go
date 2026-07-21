package lexer

import "strconv"

func MapToNumber(str []rune) (float64, bool) {
	float, err := strconv.ParseFloat(string(str), 64)
	return float, err != nil
}
