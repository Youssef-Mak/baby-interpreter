package main

import (
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
	"github.com/Youssef-Mak/baby-interpreter/pkg/tokenizer"
	"testing"
)

func TestSymTokenizer(t *testing.T) {
	input := `,;(}{)+=`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.LPAREN, "("},
		{token.RBRACE, "}"},
		{token.LBRACE, "{"},
		{token.RPAREN, ")"},
		{token.PLUS, "+"},
		{token.ASSIGN, "="},
	}

	l := tokenizer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestSynTokenizer(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	let add = fn(x, y) {
	x + y;
	};
	let result = add(five, ten);
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIF, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.LET, "let"},
		{token.IDENTIF, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIF, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIF, "x"},
		{token.COMMA, ","},
		{token.IDENTIF, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIF, "x"},
		{token.PLUS, "+"},
		{token.IDENTIF, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIF, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIF, "add"},
		{token.LPAREN, "("},
		{token.IDENTIF, "five"},
		{token.COMMA, ","},
		{token.IDENTIF, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := tokenizer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}

}
