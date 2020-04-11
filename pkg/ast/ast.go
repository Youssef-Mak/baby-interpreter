package ast

import (
	"bytes"
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LET STATEMENT -> "let <identifier> = <expression>;"
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (lStatement *LetStatement) statementNode()       {}
func (lStatement *LetStatement) TokenLiteral() string { return lStatement.Token.Literal }
func (lStatement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(lStatement.TokenLiteral() + " ")
	out.WriteString(lStatement.Name.String())
	out.WriteString(" = ")

	if lStatement.Value != nil {
		out.WriteString(lStatement.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

/* EXPRESSIONS */

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string       { return ident.Value }

type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

/* STATEMENTS */

// RETURN STATEMENT -> "return <expression>;"
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rS *ReturnStatement) statementNode()       {}
func (rS *ReturnStatement) TokenLiteral() string { return rS.Token.Literal }
func (rS *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rS.TokenLiteral() + " ")
	if rS.ReturnValue != nil {
		out.WriteString(rS.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// EXPRESSION STATEMENT -> "<expression>;"
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (exp *ExpressionStatement) statementNode()       {}
func (exp *ExpressionStatement) TokenLiteral() string { return exp.Token.Literal }
func (exp *ExpressionStatement) String() string {
	if exp.Expression != nil {
		return exp.Expression.String()
	}
	return ""
}
