package scanner

import (
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Scanner struct {
	input string
	index int
}

func NewScanner(input string) Scanner {
	return Scanner{
		input: input,
		index: 0,
	}
}

func (s *Scanner) Scan() (token.Token, error) {
	if s.index >= len(s.input) {
		return token.Token{Type: token.EOF, Lexeme: "EOF", Literal: nil}, nil
	}

	switch s.input[s.index] {
	case '(':
		s.index++
		return token.Token{Type: token.LEFT_PAREN, Lexeme: "(", Literal: nil}, nil
	case ')':
		s.index++
		return token.Token{Type: token.RIGHT_PAREN, Lexeme: ")", Literal: nil}, nil
	default:
		return token.Token{Type: token.EOF, Lexeme: "EOF", Literal: nil}, nil
	}
}