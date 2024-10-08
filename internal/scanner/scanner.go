package scanner

import (
	"fmt"
	"strconv"

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

	switch s.peak() {
	case '(':
		s.advance()
		return &token.Token{Type: token.LEFT_PAREN, Lexeme: "(", Literal: nil}, nil
	case ')':
		s.advance()
		return &token.Token{Type: token.RIGHT_PAREN, Lexeme: ")", Literal: nil}, nil
	case '{':
		s.advance()
		return &token.Token{Type: token.LEFT_BRACE, Lexeme: "{", Literal: nil}, nil
	case '}':
		s.advance()
		return &token.Token{Type: token.RIGHT_BRACE, Lexeme: "}", Literal: nil}, nil
	case ',':
		s.advance()
		return &token.Token{Type: token.COMMA, Lexeme: ",", Literal: nil}, nil
	case '.':
		s.advance()
		return &token.Token{Type: token.DOT, Lexeme: ".", Literal: nil}, nil
	case '*':
		s.advance()
		return &token.Token{Type: token.STAR, Lexeme: "*", Literal: nil}, nil
	case '+':
		s.advance()
		return &token.Token{Type: token.PLUS, Lexeme: "+", Literal: nil}, nil
	case '-':
		s.advance()
		return &token.Token{Type: token.MINUS, Lexeme: "-", Literal: nil}, nil
	case ';':
		s.advance()
		return &token.Token{Type: token.SEMICOLON, Lexeme: ";", Literal: nil}, nil
	case '=':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=", Literal: nil}, nil
		}
		return &token.Token{Type: token.EQUAL, Lexeme: "=", Literal: nil}, nil
	case '!':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.BANG_EQUAL, Lexeme: "!=", Literal: nil}, nil
		}
		return &token.Token{Type: token.BANG, Lexeme: "!", Literal: nil}, nil
	case '<':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.LESS_EQUAL, Lexeme: "<=", Literal: nil}, nil
		}
		return &token.Token{Type: token.LESS, Lexeme: "<", Literal: nil}, nil
	case '>':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.GREATER_EQUAL, Lexeme: ">=", Literal: nil}, nil
		}
		return &token.Token{Type: token.GREATER, Lexeme: ">", Literal: nil}, nil
	case '/':
		s.advance()
		if s.peak() == '/' {
			for ; !s.isAtEnd(); s.advance() {
				if s.peak() == '\n' {
					s.advance()
					s.line++
					break
				}
			}
			return nil, nil
		}
		return &token.Token{Type: token.SLASH, Lexeme: "/", Literal: nil}, nil
	case ' ', '\r', '\t':
		s.advance()
		return nil, nil
	case '\n':
		s.advance()
		s.line++
		return nil, nil
	case '"':
		s.advance()
		var literal string
		for s.peak() != '"' {
			if s.isAtEnd() {
				var err = fmt.Errorf("[line %d] Error: Unterminated string.", s.line)
				return nil, err
			}
			literal += string(s.peak())
			s.advance()
		}
		s.advance()
		return &token.Token{Type: token.STRING, Lexeme: fmt.Sprintf("\"%s\"", literal), Literal: literal}, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var literal string
		for s.isDigit(s.peak()) {
			literal += string(s.peak())
			s.advance()
		}

		if s.peak() == '.' && s.isDigit(s.peakNext()) {
			literal += string(s.peak())
			s.advance()
			for s.isDigit(s.peak()) {
				literal += string(s.peak())
				s.advance()
			}
		}

		num, _ := strconv.ParseFloat(literal, 64)

		return &token.Token{Type: token.NUMBER, Lexeme: literal, Literal: num}, nil
	default:
		if s.isAlpha(s.peak()) {
			var literal string
			for s.isAlphaNumeric(s.peak()) {
				literal += string(s.peak())
				s.advance()
			}
			return &token.Token{Type: token.IDENTIFIER, Lexeme: literal, Literal: nil}, nil
		} else {
			var err = fmt.Errorf("[line %d] Error: Unexpected character: %c", s.line, s.peak())
			s.advance()
			return nil, err
		}
	}
}

func (s *Scanner) advance() byte {
	s.index++
	return s.input[s.index-1]
}

func (s *Scanner) peak() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.input[s.index]
}

func (s *Scanner) peakNext() byte {
	if s.index+1 >= len(s.input) {
		return 0
	}
	return s.input[s.index+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.index >= len(s.input)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
