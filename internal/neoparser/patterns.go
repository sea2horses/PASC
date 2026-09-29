package neoparser

import "pseint-compiled/internal/lexer"

type TokenPattern struct {
	Type lexer.TokenType
}

func Tok(t lexer.TokenType) *TokenPattern {
	return &TokenPattern{
		Type: t,
	}
}

type KeywordPattern struct {
	Keyword lexer.Keyword
}

func Kw(k lexer.Keyword) *KeywordPattern {
	return &KeywordPattern{
		Keyword: k,
	}
}
