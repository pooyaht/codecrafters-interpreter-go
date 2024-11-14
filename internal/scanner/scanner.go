package scanner

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Scanner struct {
	input   string
	tokens  []*token.Token
	start   int
	current int
	line    int
}

func NewScanner(input string) Scanner {
	return Scanner{
		input:   input,
		tokens:  make([]*token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]*token.Token, error) {
	for {
		t, err := s.Scan()
		if err != nil {
			return nil, err
		}

		if t == nil {
			continue
		}

		s.tokens = append(s.tokens, t)

		if t.Type == token.EOF {
			break
		}
	}

	return s.tokens, nil
}

func (s *Scanner) Scan() (*token.Token, error) {
	if s.isAtEnd() {
		return &token.Token{Type: token.EOF, Lexeme: "EOF", Literal: nil, Line: s.line}, nil
	}

	defer func() {
		s.start = s.current
	}()

	switch s.peak() {
	case '(':
		s.advance()
		return &token.Token{Type: token.LEFT_PAREN, Lexeme: "(", Literal: nil, Line: s.line}, nil
	case ')':
		s.advance()
		return &token.Token{Type: token.RIGHT_PAREN, Lexeme: ")", Literal: nil, Line: s.line}, nil
	case '{':
		s.advance()
		return &token.Token{Type: token.LEFT_BRACE, Lexeme: "{", Literal: nil, Line: s.line}, nil
	case '}':
		s.advance()
		return &token.Token{Type: token.RIGHT_BRACE, Lexeme: "}", Literal: nil, Line: s.line}, nil
	case ',':
		s.advance()
		return &token.Token{Type: token.COMMA, Lexeme: ",", Literal: nil, Line: s.line}, nil
	case '.':
		s.advance()
		return &token.Token{Type: token.DOT, Lexeme: ".", Literal: nil, Line: s.line}, nil
	case '*':
		s.advance()
		return &token.Token{Type: token.STAR, Lexeme: "*", Literal: nil, Line: s.line}, nil
	case '+':
		s.advance()
		return &token.Token{Type: token.PLUS, Lexeme: "+", Literal: nil, Line: s.line}, nil
	case '-':
		s.advance()
		return &token.Token{Type: token.MINUS, Lexeme: "-", Literal: nil, Line: s.line}, nil
	case ';':
		s.advance()
		return &token.Token{Type: token.SEMICOLON, Lexeme: ";", Literal: nil, Line: s.line}, nil
	case '=':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.EQUAL_EQUAL, Lexeme: "==", Literal: nil, Line: s.line}, nil
		}
		return &token.Token{Type: token.EQUAL, Lexeme: "=", Literal: nil, Line: s.line}, nil
	case '!':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.BANG_EQUAL, Lexeme: "!=", Literal: nil, Line: s.line}, nil
		}
		return &token.Token{Type: token.BANG, Lexeme: "!", Literal: nil, Line: s.line}, nil
	case '<':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.LESS_EQUAL, Lexeme: "<=", Literal: nil, Line: s.line}, nil
		}
		return &token.Token{Type: token.LESS, Lexeme: "<", Literal: nil, Line: s.line}, nil
	case '>':
		s.advance()
		if s.peak() == '=' {
			s.advance()
			return &token.Token{Type: token.GREATER_EQUAL, Lexeme: ">=", Literal: nil, Line: s.line}, nil
		}
		return &token.Token{Type: token.GREATER, Lexeme: ">", Literal: nil, Line: s.line}, nil
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
		return &token.Token{Type: token.SLASH, Lexeme: "/", Literal: nil, Line: s.line}, nil
	case ' ', '\r', '\t':
		s.advance()
		return nil, nil
	case '\n':
		s.advance()
		s.line++
		return nil, nil
	case '"':
		s.advance()
		for s.peak() != '"' {
			if s.isAtEnd() {
				var err = fmt.Errorf("[line %d] Error: Unterminated string.", s.line)
				return nil, err
			}
			s.advance()
		}

		literal := s.input[s.start+1 : s.current]
		s.advance()

		return &token.Token{Type: token.STRING, Lexeme: fmt.Sprintf("\"%s\"", literal), Literal: literal, Line: s.line}, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		for s.isDigit(s.peak()) {
			s.advance()
		}

		if s.peak() == '.' && s.isDigit(s.peakNext()) {
			s.advance()
			for s.isDigit(s.peak()) {
				s.advance()
			}
		}

		literal := s.input[s.start:s.current]
		num, _ := strconv.ParseFloat(literal, 64)

		return &token.Token{Type: token.NUMBER, Lexeme: literal, Literal: num, Line: s.line}, nil
	default:
		if s.isAlpha(s.peak()) {
			for s.isAlphaNumeric(s.peak()) {
				s.advance()
			}

			literal := s.input[s.start:s.current]
			if keyword, ok := token.Keywords[literal]; ok {
				return &token.Token{Type: keyword, Lexeme: literal, Literal: nil, Line: s.line}, nil
			}

			return &token.Token{Type: token.IDENTIFIER, Lexeme: literal, Literal: nil, Line: s.line}, nil
		} else {
			var err = fmt.Errorf("[line %d] Error: Unexpected character: %c", s.line, s.peak())
			s.advance()
			return nil, err
		}
	}
}

func (s *Scanner) advance() {
	s.current++
}

func (s *Scanner) peak() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.input[s.current]
}

func (s *Scanner) peakNext() byte {
	if s.current+1 >= len(s.input) {
		return 0
	}
	return s.input[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.input)
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
