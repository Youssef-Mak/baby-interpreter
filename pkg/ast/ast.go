package ast

import (
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program will be the root node of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LET STATEMENT -> "let <identifier> = <expression>;"
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (lStatement *LetStatement) statementNode()       {}
func (lStatement *LetStatement) TokenLiteral() string { return lStatement.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }

// RETURN STATEMENT -> "return <expression>;"
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rS *ReturnStatement) statementNode()       {}
func (rS *ReturnStatement) TokenLiteral() string { return rS.Token.Literal }
