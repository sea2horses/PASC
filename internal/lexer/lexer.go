package lexer

import (
	"fmt"
	"unicode"
)

type ErrOutOfBounds struct {
	Max uint32
	Value uint32
}

func (e ErrOutOfBounds) Error() string {
	return fmt.Sprintf("value %d is out of bounds (%d)", e.Value, e.Max)
}

type ErrUnexpectedChar struct {
	Char rune
}

func (e ErrUnexpectedChar) Error() string {
	return fmt.Sprintf("unexpected char: %c", e.Char)
}

type LexerErr string

func (e LexerErr) Error() string {
	return string(e)
}

const (
	ErrExpectedIdentifier = LexerErr("expected identifier")
	ErrExpectedString = LexerErr("expected string")
	ErrUnterminatedString = LexerErr("unterminated string")
)

/* Lexer struct */
type Lexer struct {
	Src []rune /* Loaded source code */
	positionInSrc uint32
	positionInFile Position
}

/* Get character at current position */
func (l *Lexer) getch() (rune, error) {
	return l.at(l.positionInSrc)
}

/* Get character at any position */
func (l *Lexer) at(position uint32) (rune, error) {
	if int(position) >= len(l.Src) {
		return 0, ErrOutOfBounds{Max: uint32(len(l.Src)), Value: position}
	}

	return l.Src[position], nil
}

/* Advance function */
func (l *Lexer) advance() error {
	ch, err := l.getch()

	if err != nil {
		return err
	}

	/* If there is a line break, reset the column and add to the line in file */
	if ch == '\n' {
		l.positionInFile.Line++;
		l.positionInFile.Column = 1;
	} else {
		l.positionInFile.Column++;
	}

	l.positionInSrc++;
	return nil
}

/* Skip Whitespace */
func (l *Lexer) skip_whitespace() {
	for {
		ch, err := l.getch()

		if err != nil || !unicode.IsSpace(ch) { break }

		l.advance()
	}
}

/* Parse identifiers */
func (l *Lexer) parse_identifier() (*Token, error) {
	init_position := l.positionInSrc;

	/* Run until reaching the end of the file or putting a stop */
	for {
		ch, err := l.getch()

		if err != nil || (!unicode.IsLetter(ch) && !unicode.IsNumber(ch) && ch != '_') { break }

		l.advance()
	}

	end_position := l.positionInSrc;

	init_ch, err := l.at(init_position);

	if err != nil {
		return nil, err
	}

	/* If first char is NOT an identifier friendly character, throw an error */
	if !unicode.IsLetter(init_ch) && init_ch != '_' {
		return nil, ErrExpectedIdentifier
	}
	
	/* Else, emmit the token */
	value := l.Src[init_position:end_position]

	_, ok := MapToKeyword(value);

	/* TODO: Parse boolean literals */
	if ok {
		return &Token{ Type: KEYWORD, Value: string(value) }, nil
	} else {
		return &Token{ Type: IDENTIFIER, Value: string(value) }, nil
	}
}

/* Parse string literal */
func (l *Lexer) parse_string_literal() (*Token, error) {
	init_position := l.positionInSrc;

	/* Skip first " */
	l.advance();

	/* Run until reaching the end of the file or putting a stop */
	for {
		ch, err := l.getch()

		if err != nil || ch == '"' { break }

		l.advance()
	}

	end_position := l.positionInSrc;

	/* Check that the init position and the end poisiton exist */
	init_ch, err := l.at(init_position);

	if err != nil {
		return nil, err
	}

	end_ch, err := l.at(end_position);

	if err != nil {
		return nil, err
	}

	/* Check that both ends of our string are good */
	if init_ch != '"' {
		return nil, ErrExpectedString
	}

	if end_ch != '"' {
		return nil, ErrUnterminatedString
	}

	/* Emit string token */
	value := l.Src[init_position + 1:end_position]

	/* Skip last '"' */
	l.advance()

	return &Token{Type: STRING_LITERAL, Value: string(value)}, nil
}

/* TODO: Parse char */

/* Tokenize function */
func (l *Lexer) Tokenize() ([]Token, error) {
	/* Make empty token list */
	tokens := []Token{}

	/* While there is source code to parse */
	for l.positionInSrc < uint32(len(l.Src)) {
		ch, _ := l.getch()

		if unicode.IsSpace(ch) {
			l.skip_whitespace()
		} else if ch == '"' {
			token, err := l.parse_string_literal()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *token)
		} else if unicode.IsLetter(ch) || ch == '_' {
			token, err := l.parse_identifier()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *token)
		} else {
			return nil, ErrUnexpectedChar{Char: ch}
		}
	}

	return tokens, nil
}
