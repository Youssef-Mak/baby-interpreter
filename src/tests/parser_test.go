package tests

import (
	"github.com/Youssef-Mak/baby-interpreter/pkg/ast"
	"github.com/Youssef-Mak/baby-interpreter/pkg/parser"
	"github.com/Youssef-Mak/baby-interpreter/pkg/tokenizer"
	"testing"
)

func TestLets(t *testing.T) {
	input := `
	let two = 2;
	let three = 3;
	let fortyfive = 45;
	let thirtysix = 36;
	`

	tokenizer := tokenizer.New(input)
	parser := parser.New(tokenizer)

	program := parser.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 elements. Got %d elements",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"two"},
		{"three"},
		{"fortyfive"},
		{"thirtysix"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLet(t, statement, tt.expectedIdentifier) {
			return
		}
	}

}

func testLet(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true

}
