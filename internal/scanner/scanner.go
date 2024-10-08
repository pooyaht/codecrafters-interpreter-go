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
	case '=':
		s.index++
		if !s.isAtEnd() && s.input[s.index] == '=' {
			s.index++
			return &token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=", Literal: nil}, nil
		}
		return &token.Token{Type: token.EQUAL, Lexeme: "=", Literal: nil}, nil
	case '!':
		s.index++
		if !s.isAtEnd() && s.input[s.index] == '=' {
			s.index++
			return &token.Token{Type: token.BANG_EQUAL, Lexeme: "!=", Literal: nil}, nil
		}
		return &token.Token{Type: token.BANG, Lexeme: "!", Literal: nil}, nil
	case '<':
		s.index++
		if !s.isAtEnd() && s.input[s.index] == '=' {
			s.index++
			return &token.Token{Type: token.LESS_EQUAL, Lexeme: "<=", Literal: nil}, nil
		}
		return &token.Token{Type: token.LESS, Lexeme: "<", Literal: nil}, nil
	case '>':
		s.index++
		if !s.isAtEnd() && s.input[s.index] == '=' {
			s.index++
			return &token.Token{Type: token.GREATER_EQUAL, Lexeme: ">=", Literal: nil}, nil
		}
		return &token.Token{Type: token.GREATER, Lexeme: ">", Literal: nil}, nil
	case '/':
		s.index++
		if !s.isAtEnd() && s.input[s.index] == '/' {
			for ; !s.isAtEnd(); s.index++ {
				if s.input[s.index] == '\n' {
					s.index++
					s.line++
					break
				}
			}
			return nil, nil
		}
		return &token.Token{Type: token.SLASH, Lexeme: "/", Literal: nil}, nil
	case ' ', '\r', '\t':
		s.index++
		return nil, nil
	case '\n':
		s.index++
		s.line++
		return nil, nil
	case '"':
		s.index++
		var literal string
		for !s.isAtEnd() && s.input[s.index] != '"' {
			literal += string(s.input[s.index])
			s.index++
		}
		if s.isAtEnd() {
			var err = fmt.Errorf("[line %d] Error: Unterminated string.", s.line)
			return nil, err
		}
		s.index++
		return &token.Token{Type: token.STRING, Lexeme: fmt.Sprintf("\"%s\"", literal), Literal: literal}, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var literal string
		for s.isDigit(s.input[s.index]) {
			literal += string(s.input[s.index])
			s.index++
		}

		if s.input[s.index] == '.' && s.isDigit(s.peak()) {
			literal += string(s.input[s.index])
			s.index++
			for s.isDigit(s.input[s.index]) {
				literal += string(s.input[s.index])
				s.index++
			}
		}

		num, _ := strconv.ParseFloat(literal, 64)

		return &token.Token{Type: token.NUMBER, Lexeme: literal, Literal: num}, nil
	default:
		var err = fmt.Errorf("[line %d] Error: Unexpected character: %c", s.line, s.input[s.index])
		s.index++
		return nil, err
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.index >= len(s.input)
}

func (s *Scanner) peak() byte {
	if s.index+1 >= len(s.input) {
		return 0
	}
	return s.input[s.index+1]
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
