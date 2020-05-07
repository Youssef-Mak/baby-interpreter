package tokenizer

import (
	"regexp"

	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
)

type Tokenizer struct {
	input        string
	position     int  // index of current character being processed
	readPosition int  // index of next character to be processed
	ch           byte // Current character being processed (ASCII)
}

func New(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.readChar() // Initializes Indices
	return t
}

// Returns Token based on input at current index(position)
func (t *Tokenizer) NextToken() token.Token {
	var tok token.Token
	t.consumeWhitespace()

	switch t.ch {
	case '+':
		tok = newToken(token.PLUS, t.ch)
	case '-':
		tok = newToken(token.MINUS, t.ch)
	case '/':
		tok = newToken(token.SLASH, t.ch)
	case '*':
		tok = newToken(token.ASTERIX, t.ch)
	case '=':
		switch t.peekChar() {
		case '&':
			t.readChar()
			if t.peekChar() == '=' {
				t.readChar()
				tok = token.Token{Type: token.REF_EQUALS, Literal: "=&="}
			} else {
				tok = token.Token{Type: token.REF_ASSIGN, Literal: "=&"}
			}
		case '*':
			t.readChar()
			if t.peekChar() == '=' {
				t.readChar()
				tok = token.Token{Type: token.VAL_EQUALS, Literal: "=*="}
			} else {
				tok = token.Token{Type: token.VAL_ASSIGN, Literal: "=*"}
			}
		default:
			tok = newToken(token.ASSIGN, t.ch)
		}
	case '!':
		switch t.peekChar() {
		case '&':
			t.readChar()
			if t.peekChar() == '=' {
				t.readChar()
				tok = token.Token{Type: token.REF_NOTEQUALS, Literal: "!&="}
			} else {
				tok = newToken(token.ILLEGAL, t.ch)
			}
		case '*':
			t.readChar()
			if t.peekChar() == '=' {
				t.readChar()
				tok = token.Token{Type: token.VAL_NOTEQUALS, Literal: "!*="}
			} else {
				tok = newToken(token.ILLEGAL, t.ch)
			}
		default:
			tok = newToken(token.NOT, t.ch)
		}
	case '>':
		if t.peekChar() == '=' {
			t.readChar()
			tok = token.Token{Type: token.GTEQUAL, Literal: ">="}
		} else {
			tok = newToken(token.GREATERTHAN, t.ch)
		}
	case '<':
		if t.peekChar() == '=' {
			t.readChar()
			tok = token.Token{Type: token.LTEQUAL, Literal: "<="}
		} else {
			tok = newToken(token.LESSTHAN, t.ch)
		}
	case '&':
		tok = newToken(token.AND, t.ch)
	case '|':
		tok = newToken(token.OR, t.ch)
	case ',':
		tok = newToken(token.COMMA, t.ch)
	case '.':
		tok = newToken(token.DOT, t.ch)
	case ';':
		tok = newToken(token.SEMICOLON, t.ch)
	case ':':
		tok = newToken(token.COLON, t.ch)
	case '(':
		tok = newToken(token.LPAREN, t.ch)
	case ')':
		tok = newToken(token.RPAREN, t.ch)
	case '{':
		tok = newToken(token.LBRACE, t.ch)
	case '}':
		tok = newToken(token.RBRACE, t.ch)
	case '[':
		tok = newToken(token.LBRACKET, t.ch)
	case ']':
		tok = newToken(token.RBRACKET, t.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = t.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(t.ch) {
			tok.Literal = t.readIdentifier()
			tok.Type = token.IdentLookUp(tok.Literal)
			return tok
		} else if isDigit(t.ch) {
			tok.Literal = t.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, t.ch)
		}
	}

	t.readChar()
	return tok
}

// Peeks the next character in input without modifying indexes
func (t *Tokenizer) peekChar() byte {
	if t.readPosition > len(t.input) {
		return 0
	} else {
		return t.input[t.readPosition]
	}
}

// Fully reads Identifier
func (t *Tokenizer) readIdentifier() string {
	position := t.position
	for isLetter(t.ch) {
		t.readChar()
	}
	return t.input[position:t.position]
}

// Fully reads number
func (t *Tokenizer) readNumber() string {
	position := t.position
	for isDigit(t.ch) {
		t.readChar()
	}
	return t.input[position:t.position]
}

// Fully read String
func (t *Tokenizer) readString() string {
	t.readChar() // Skip opening quotes
	position := t.position
	for t.ch != 0 && t.ch != '"' {
		t.readChar()

	}
	return t.input[position:t.position]
}

// Reads next character of input
func (t *Tokenizer) readChar() {
	if t.readPosition >= len(t.input) {
		t.ch = 0 // ASCII code for NUL character (EOF)
	} else {
		t.ch = t.input[t.readPosition]
	}
	t.position = t.readPosition
	t.readPosition += 1
}

// Initializes new Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Checks if byte is letter or underscore
func isLetter(ch byte) bool {
	match, _ := regexp.Match("[A-Za-z_]", []byte{ch})
	return match
}

// Check if byte corresponds to number
func isDigit(ch byte) bool {
	match, _ := regexp.Match("[0-9]", []byte{ch})
	return match
}

// Skips alt whitespaces untit new char is read
func (t *Tokenizer) consumeWhitespace() {
	for t.ch == ' ' || t.ch == '\t' || t.ch == '\n' || t.ch == '\r' {
		t.readChar()
	}
}
