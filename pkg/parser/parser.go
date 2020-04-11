package parser

import (
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/ast"
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
	"github.com/Youssef-Mak/baby-interpreter/pkg/tokenizer"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // !x or -x
	FUNCALL     // func(x)
)

type Parser struct {
	tokenizer *tokenizer.Tokenizer
	errors    []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFuncs map[token.TokenType]prefixParseFunc
	infixParseFuncs  map[token.TokenType]infixParseFunc
}

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

func New(tokenizer *tokenizer.Tokenizer) *Parser {
	p := &Parser{tokenizer: tokenizer, errors: []string{}}

	p.prefixParseFuncs = make(map[token.TokenType]prefixParseFunc)
	p.addPrefix(token.IDENTIF, p.parseIdentifier)
	p.addPrefix(token.INT, p.parseIntegerLiteral)

	p.infixParseFuncs = make(map[token.TokenType]infixParseFunc)

	// Read tokens in pairs, so currToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) addPrefix(tokenType token.TokenType, fn prefixParseFunc) {
	p.prefixParseFuncs[tokenType] = fn
}

func (p *Parser) addInfix(tokenType token.TokenType, fn infixParseFunc) {
	p.infixParseFuncs[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.tokenizer.NextToken()
}

func (p *Parser) checkIdCurrentToken(toCheck token.TokenType) bool {
	return p.currentToken.Type == toCheck
}

func (p *Parser) checkIdNextToken(toCheck token.TokenType) bool {
	return p.peekToken.Type == toCheck
}

// Assertion function: enforce correctness of the order of tokens by checking type of next token
func (p *Parser) peekNextToken(toCheck token.TokenType) bool {
	if p.checkIdNextToken(toCheck) {
		p.nextToken()
		return true
	} else {
		p.peekNextTokenError(toCheck)
		return false
	}
}

func (p *Parser) peekNextTokenError(t token.TokenType) {
	errMsg := fmt.Sprintf("expected token %s, but got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, errMsg)
}

/* EXPRESSION PARSING */

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFuncs[p.currentToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

/* SEMANTIC CODE FUNTIONS */

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as Integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

/* STATEMENT PARSING */

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}

}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{Token: p.currentToken} // Going to be LET type

	if !p.peekNextToken(token.IDENTIF) {
		return nil
	}

	// There is an Identifier i.e is of form '''let <identifier> <...>'''
	letStatement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.peekNextToken(token.ASSIGN) {
		return nil
	}

	// TODO: handle expressions
	for !p.checkIdCurrentToken(token.SEMICOLON) {
		p.nextToken()
	}

	// Its of form '''let <identifier> = <...>'''
	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	retStatement := &ast.ReturnStatement{Token: p.currentToken} // Going to be LET type

	p.nextToken()

	// TODO: handle expression
	for !p.checkIdCurrentToken(token.SEMICOLON) {
		p.nextToken()
	}

	return retStatement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expStatement := &ast.ExpressionStatement{Token: p.currentToken}

	expStatement.Expression = p.parseExpression(LOWEST)

	if p.peekNextToken(token.SEMICOLON) {
		p.nextToken()
	}

	return expStatement
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.checkIdCurrentToken(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) GetErrors() []string {
	return p.errors
}
