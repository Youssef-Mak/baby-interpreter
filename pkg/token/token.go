package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENTIF TokenType = "IDENTIF" // add, foobar, x, y, ...
	INT     TokenType = "INT"     // 1343456
	// Operators
	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywordMap = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func IdentLookUp(id string) TokenType {
	if tok, ok := keywordMap[id]; ok {
		return tok
	}
	return IDENTIF
}
