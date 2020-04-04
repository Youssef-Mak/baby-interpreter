package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENTIF TokenType = "IDENTIF" // add, foobar, x, y, ...
	INT     TokenType = "INT"     // 1343456
	// Operators
	ASSIGN  TokenType = "="
	PLUS    TokenType = "+"
	MINUS   TokenType = "-"
	SLASH   TokenType = "/"
	ASTERIX TokenType = "*"
	NOT     TokenType = "!"
	// Logic
	LESSTHAN    TokenType = "<"
	GREATERTHAN TokenType = ">"
	AND         TokenType = "&"
	OR          TokenType = "|"
	TRUE        TokenType = "TRUE"
	FALSE       TokenType = "FALSE"
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
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	RETURN   TokenType = "return"
)

type Token struct {
	Type    TokenType
	Literal string
}

var valMap = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
}

var keywordMap = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func IdentLookUp(id string) TokenType {
	if tok, ok := keywordMap[id]; ok {
		return tok
	} else if tok, ok := valMap[id]; ok {
		return tok
	}
	return IDENTIF
}
