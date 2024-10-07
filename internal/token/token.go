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
	case SEMICOLON:
		return fmt.Sprintf("SEMICOLON %s null", t.Type)
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
	SEMICOLON TokenType = ";"
)
