package token

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/internal/util"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
}

func (t Token) String() string {
	switch t.Type {
	case EOF:
		return fmt.Sprintf("%s  null", t.Type)
	case LEFT_PAREN:
		return fmt.Sprintf("LEFT_PAREN %s null", t.Type)
	case RIGHT_PAREN:
		return fmt.Sprintf("RIGHT_PAREN %s null", t.Type)
	case LEFT_BRACE:
		return fmt.Sprintf("LEFT_BRACE %s null", t.Type)
	case RIGHT_BRACE:
		return fmt.Sprintf("RIGHT_BRACE %s null", t.Type)
	case COMMA:
		return fmt.Sprintf("COMMA %s null", t.Type)
	case DOT:
		return fmt.Sprintf("DOT %s null", t.Type)
	case STAR:
		return fmt.Sprintf("STAR %s null", t.Type)
	case PLUS:
		return fmt.Sprintf("PLUS %s null", t.Type)
	case MINUS:
		return fmt.Sprintf("MINUS %s null", t.Type)
	case SLASH:
		return fmt.Sprintf("SLASH %s null", t.Type)
	case SEMICOLON:
		return fmt.Sprintf("SEMICOLON %s null", t.Type)
	case EQUAL:
		return fmt.Sprintf("EQUAL %s null", t.Type)
	case EQUAL_EQUAL:
		return fmt.Sprintf("EQUAL_EQUAL %s null", t.Type)
	case BANG:
		return fmt.Sprintf("BANG %s null", t.Type)
	case BANG_EQUAL:
		return fmt.Sprintf("BANG_EQUAL %s null", t.Type)
	case LESS:
		return fmt.Sprintf("LESS %s null", t.Type)
	case LESS_EQUAL:
		return fmt.Sprintf("LESS_EQUAL %s null", t.Type)
	case GREATER:
		return fmt.Sprintf("GREATER %s null", t.Type)
	case GREATER_EQUAL:
		return fmt.Sprintf("GREATER_EQUAL %s null", t.Type)
	case STRING:
		return fmt.Sprintf("STRING %s %s", t.Lexeme, t.Literal)
	case NUMBER:
		return fmt.Sprintf("NUMBER %s %s", t.Lexeme, util.FormatFloat(t.Literal.(float64)))
	case IDENTIFIER:
		return fmt.Sprintf("IDENTIFIER %s null", t.Lexeme)
	default:
		return fmt.Sprintf("%s  null", t.Type)
	}
}

type TokenType string

const (
	EOF TokenType = "EOF"

	LEFT_PAREN  TokenType = "("
	RIGHT_PAREN TokenType = ")"

	LEFT_BRACE  TokenType = "{"
	RIGHT_BRACE TokenType = "}"

	COMMA     TokenType = ","
	DOT       TokenType = "."
	STAR      TokenType = "*"
	PLUS      TokenType = "+"
	MINUS     TokenType = "-"
	SLASH     TokenType = "/"
	SEMICOLON TokenType = ";"

	EQUAL       TokenType = "="
	EQUAL_EQUAL TokenType = "=="

	BANG       TokenType = "!"
	BANG_EQUAL TokenType = "!="

	LESS       TokenType = "<"
	LESS_EQUAL TokenType = "<="

	GREATER       TokenType = ">"
	GREATER_EQUAL TokenType = ">="

	STRING TokenType = "STRING"

	NUMBER TokenType = "NUMBER"

	IDENTIFIER TokenType = "IDENTIFIER"
)
