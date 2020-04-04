package tokenizer

import (
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
	"regexp"
)

type Tokenizer struct {
	input        string
	position     int  // index of current character being processed
	readPosition int  // index of next character to be processed
	ch           byte // Current character being processed (ASCII)
}

func New(input string) *Tokenizer {
	l := &Tokenizer{input: input}
	l.readChar()
	return l
}

func (l *Tokenizer) NextToken() token.Token {
	var tok token.Token
	l.consumeWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERIX, l.ch)
	case '!':
		tok = newToken(token.NOT, l.ch)
	case '>':
		tok = newToken(token.GREATERTHAN, l.ch)
	case '<':
		tok = newToken(token.LESSTHAN, l.ch)
	case '&':
		tok = newToken(token.AND, l.ch)
	case '|':
		tok = newToken(token.OR, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IdentLookUp(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Fully reads Identifier
func (l *Tokenizer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Fully reads number
func (l *Tokenizer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Reads next character of input
func (l *Tokenizer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for NUL character (EOF)
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// Initializes new Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Checks if byte is letter or underscore
func isLetter(ch byte) bool {
	match, _ := regexp.Match("[A-Za-z|_]", []byte{ch})
	return match
}

// Check if byte corresponds to number
func isDigit(ch byte) bool {
	match, _ := regexp.Match("[0-9]", []byte{ch})
	return match
}

// Skips all whitespaces until new char is read
func (l *Tokenizer) consumeWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
