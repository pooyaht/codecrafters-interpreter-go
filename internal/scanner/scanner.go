package scanner

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Scanner struct {
	input string
	index int
	line  int
}

func NewScanner(input string) Scanner {
	return Scanner{
		input: input,
		index: 0,
		line:  1,
	}
}

func (s *Scanner) Scan() (*token.Token, error) {
	if s.isAtEnd() {
		return &token.Token{Type: token.EOF, Lexeme: "EOF", Literal: nil}, nil
	}

	switch s.input[s.index] {
	case '(':
		s.index++
		return &token.Token{Type: token.LEFT_PAREN, Lexeme: "(", Literal: nil}, nil
	case ')':
		s.index++
		return &token.Token{Type: token.RIGHT_PAREN, Lexeme: ")", Literal: nil}, nil
	case '{':
		s.index++
		return &token.Token{Type: token.LEFT_BRACE, Lexeme: "{", Literal: nil}, nil
	case '}':
		s.index++
		return &token.Token{Type: token.RIGHT_BRACE, Lexeme: "}", Literal: nil}, nil
	case ',':
		s.index++
		return &token.Token{Type: token.COMMA, Lexeme: ",", Literal: nil}, nil
	case '.':
		s.index++
		return &token.Token{Type: token.DOT, Lexeme: ".", Literal: nil}, nil
	case '*':
		s.index++
		return &token.Token{Type: token.STAR, Lexeme: "*", Literal: nil}, nil
	case '+':
		s.index++
		return &token.Token{Type: token.PLUS, Lexeme: "+", Literal: nil}, nil
	case '-':
		s.index++
		return &token.Token{Type: token.MINUS, Lexeme: "-", Literal: nil}, nil
	case ';':
		s.index++
		return &token.Token{Type: token.SEMICOLON, Lexeme: ";", Literal: nil}, nil
	default:
		var err = fmt.Errorf("[line %d] Error: Unexpected character: %c", s.line, s.input[s.index])
		s.index++
		return nil, err
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.index >= len(s.input)
}
