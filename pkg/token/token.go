package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENTIF TokenType = "IDENTIF" // add, foobar, x, y, ...
	INT     TokenType = "INT"     // 1343456
	STRING  TokenType = "STRING"  // "Hello World"
	// Operators
	ASSIGN     TokenType = "="
	REF_ASSIGN TokenType = "=&"
	VAL_ASSIGN TokenType = "=*"
	PLUS       TokenType = "+"
	MINUS      TokenType = "-"
	SLASH      TokenType = "/"
	ASTERIX    TokenType = "*"
	NOT        TokenType = "!"
	DOT        TokenType = "."
	// Logic
	LESSTHAN      TokenType = "<"
	GREATERTHAN   TokenType = ">"
	LTEQUAL       TokenType = "<="
	GTEQUAL       TokenType = ">="
	AND           TokenType = "&"
	OR            TokenType = "|"
	REF_EQUALS    TokenType = "=&="
	REF_NOTEQUALS TokenType = "!&="
	VAL_EQUALS    TokenType = "=*="
	VAL_NOTEQUALS TokenType = "!*="
	TRUE          TokenType = "TRUE"
	FALSE         TokenType = "FALSE"
	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	WHILE    TokenType = "WHILE"
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
	"fun":    FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
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
