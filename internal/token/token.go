package token

import "fmt"

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
	default:
		return fmt.Sprintf("%s  null", t.Type)
	}
}

type TokenType string

const (
	EOF TokenType = "EOF"

	LEFT_PAREN  TokenType = "("
	RIGHT_PAREN TokenType = ")"
)
