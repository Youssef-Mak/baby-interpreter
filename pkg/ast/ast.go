package ast

import (
	"bytes"
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
	"strings"
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

type Boolean struct {
	Token token.Token // token.TRUE or token.FALSE
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type StringLiteral struct {
	Token token.Token // Token.STRING
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token // token.FUNCTION
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token // token.LBRACE
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// PREFIX EXPRESSION -> <prefix operator> <expression>
type PrefixExpression struct {
	Token    token.Token // Prefix Operator Token
	Operator string
	Right    Expression
}

func (prexp *PrefixExpression) expressionNode()      {}
func (prexp *PrefixExpression) TokenLiteral() string { return prexp.Token.Literal }
func (prexp *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(prexp.Operator)
	out.WriteString(prexp.Right.String())
	out.WriteString(")")
	return out.String()

}

// INFIX EXPRESSION -> <right expression> <operator> <left expression>
type InfixExpression struct {
	Token    token.Token // Infix Operator Token
	Operator string
	Right    Expression
	Left     Expression
}

func (inexp *InfixExpression) expressionNode()      {}
func (inexp *InfixExpression) TokenLiteral() string { return inexp.Token.Literal }
func (inexp *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(inexp.Left.String())
	out.WriteString(" " + inexp.Operator + " ")
	out.WriteString(inexp.Right.String())
	out.WriteString(")")
	return out.String()
}

// IF EXPRESSION -> "if (<condition>) <consequence> else <alternative>"
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifexp *IfExpression) expressionNode()      {}
func (ifexp *IfExpression) TokenLiteral() string { return ifexp.Token.Literal }
func (ifexp *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ifexp.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifexp.Consequence.String())
	out.WriteString(" ")
	if ifexp.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ifexp.Alternative.String())
	}
	return out.String()
}

// FUNCTION CALL EXPRESSION -> <expression>(<comma seperated expressions>)
type CallExpression struct {
	Token     token.Token // token.LPAREN
	Function  Expression  // Identifier of function OR FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// INDEX EXPRESSION -> <expression>[<expression>]
type IndexExpression struct {
	Token token.Token // token.LBRACE
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

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

// BLOCK STATEMENT: A SERIES OF STATEMENTS
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
