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
	INDEX       // array[2]
	DOT         // obj.prop
)

var precedences = map[token.TokenType]int{
	token.EQUALS:      EQUALS,
	token.NOTEQUALS:   EQUALS,
	token.LESSTHAN:    LESSGREATER,
	token.GREATERTHAN: LESSGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.SLASH:       PRODUCT,
	token.ASTERIX:     PRODUCT,
	token.LPAREN:      FUNCALL,
	token.LBRACKET:    INDEX,
	token.DOT:         DOT,
}

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
	p.addPrefix(token.STRING, p.parseStringLiteral)
	p.addPrefix(token.TRUE, p.parseBoolean)
	p.addPrefix(token.FALSE, p.parseBoolean)
	p.addPrefix(token.LPAREN, p.parseGroupExpression)
	p.addPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.addPrefix(token.LBRACE, p.parseHashLiteral)
	p.addPrefix(token.IF, p.parseIfExpression)
	p.addPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.addPrefix(token.NOT, p.parsePrefixOperationExpression)
	p.addPrefix(token.MINUS, p.parsePrefixOperationExpression)

	p.infixParseFuncs = make(map[token.TokenType]infixParseFunc)
	p.addInfix(token.PLUS, p.parseInfixExpression)
	p.addInfix(token.MINUS, p.parseInfixExpression)
	p.addInfix(token.SLASH, p.parseInfixExpression)
	p.addInfix(token.ASTERIX, p.parseInfixExpression)
	p.addInfix(token.EQUALS, p.parseInfixExpression)
	p.addInfix(token.NOTEQUALS, p.parseInfixExpression)
	p.addInfix(token.LESSTHAN, p.parseInfixExpression)
	p.addInfix(token.GREATERTHAN, p.parseInfixExpression)
	p.addInfix(token.LBRACKET, p.parseIndexExpression)
	p.addInfix(token.DOT, p.parseDotExpression)
	p.addInfix(token.LPAREN, p.parseCallExpression)

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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s was found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noInfixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no infix parse function for %s was found", t)
	p.errors = append(p.errors, msg)
}

/* PRECEDENCE MANAGEMENT */

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

/* TOKEN NAVIGATION */

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

// Assertion Function: enforce correctness of the order of tokens by checking type of next token
// If Matches toCheck, parser moves to next token
func (p *Parser) peekNextToken(toCheck token.TokenType, strict bool) bool {
	if p.checkIdNextToken(toCheck) {
		p.nextToken()
		return true
	} else {
		if strict {
			p.peekNextTokenError(toCheck)
		}
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
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.checkIdNextToken(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil {
			p.noInfixParseFnError(p.currentToken.Type)
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixOperationExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

/* SEMANTIC CODE FUNTIONS */

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.checkIdCurrentToken(token.TRUE)}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.peekNextToken(token.RPAREN, true) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.peekNextToken(token.LPAREN, true) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.peekNextToken(token.RPAREN, true) {
		return nil
	}

	if !p.peekNextToken(token.LBRACE, true) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if !p.peekNextToken(token.ELSE, false) {
		return expression
	}

	if !p.peekNextToken(token.LBRACE, true) {
		return nil
	}

	expression.Alternative = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.currentToken} // token.LBRACKET
	elems := []ast.Expression{}

	if p.peekNextToken(token.RBRACKET, false) {
		arr.Elements = elems
		return arr
	}

	p.nextToken()
	elems = append(elems, p.parseExpression(LOWEST))

	for p.peekNextToken(token.COMMA, false) {
		p.nextToken()
		elems = append(elems, p.parseExpression(LOWEST))
	}

	if !p.peekNextToken(token.RBRACKET, true) {
		return nil
	}

	arr.Elements = elems
	return arr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.currentToken} // token.LBRACE
	pairs := map[ast.Expression]ast.Expression{}

	if p.peekNextToken(token.RBRACE, false) {
		hash.Pairs = pairs
		return hash
	}

	for !p.checkIdNextToken(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.peekNextToken(token.COLON, true) {
			return nil
		}
		p.nextToken()
		val := p.parseExpression(LOWEST)
		pairs[key] = val

		if !p.checkIdNextToken(token.RBRACE) && !p.peekNextToken(token.COMMA, true) {
			return nil
		}
	}

	if !p.peekNextToken(token.RBRACE, true) {
		return nil
	}

	hash.Pairs = pairs
	return hash
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	funcExp := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.peekNextToken(token.LPAREN, true) {
		return nil
	}

	funcExp.Parameters = p.parseParameters()

	if !p.peekNextToken(token.LBRACE, true) {
		return nil
	}

	funcExp.Body = p.parseBlockStatement()

	return funcExp

}

func (p *Parser) parseParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}

	if p.peekNextToken(token.RPAREN, false) {
		return parameters
	}

	for !p.checkIdCurrentToken(token.RPAREN) {
		if p.peekNextToken(token.IDENTIF, false) {
			identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
			parameters = append(parameters, identifier)
			p.peekNextToken(token.COMMA, false)
		}
		p.peekNextToken(token.RPAREN, false)
	}

	return parameters
}

func (p *Parser) parseIndexExpression(array ast.Expression) ast.Expression {
	idxExp := &ast.IndexExpression{Token: p.currentToken, Left: array}

	p.nextToken()

	idxExp.Index = p.parseExpression(LOWEST)

	if !p.peekNextToken(token.RBRACKET, true) {
		return nil
	}

	return idxExp
}

func (p *Parser) parseDotExpression(hash ast.Expression) ast.Expression {
	dotExp := &ast.DotExpression{Token: p.currentToken, Left: hash}

	p.nextToken()
	dotExp.Attribute = p.parseExpression(LOWEST)

	return dotExp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExp := &ast.CallExpression{Token: p.currentToken, Function: function}
	callExp.Arguments = p.parseCallArguments()
	return callExp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekNextToken(token.RPAREN, false) {
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekNextToken(token.COMMA, false) {
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.peekNextToken(token.RPAREN, true) {
		return nil
	}

	return args
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

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
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

	if !p.peekNextToken(token.IDENTIF, true) {
		return nil
	}

	// There is an Identifier i.e is of form '''let <identifier> <...>'''
	letStatement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.peekNextToken(token.ASSIGN, true) {
		return nil
	}

	p.nextToken()

	letStatement.Value = p.parseExpression(LOWEST)

	for !p.checkIdCurrentToken(token.SEMICOLON) {
		p.nextToken()
	}

	// Its of form '''let <identifier> = <...>'''
	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	retStatement := &ast.ReturnStatement{Token: p.currentToken} // Going to be RETURN type

	p.nextToken()

	retStatement.ReturnValue = p.parseExpression(LOWEST)

	for !p.checkIdCurrentToken(token.SEMICOLON) {
		p.nextToken()
	}

	return retStatement
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStmt := &ast.BlockStatement{Token: p.currentToken}
	blockStmt.Statements = []ast.Statement{}

	p.nextToken()

	for !p.checkIdCurrentToken(token.RBRACE) && !p.checkIdCurrentToken(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			blockStmt.Statements = append(blockStmt.Statements, statement)
		}
		p.nextToken()
	}

	return blockStmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expStatement := &ast.ExpressionStatement{Token: p.currentToken}
	expStatement.Expression = p.parseExpression(LOWEST)

	if p.checkIdNextToken(token.SEMICOLON) {
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
