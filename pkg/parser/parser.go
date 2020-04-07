package parser

import (
	"github.com/Youssef-Mak/baby-interpreter/pkg/ast"
	"github.com/Youssef-Mak/baby-interpreter/pkg/token"
	"github.com/Youssef-Mak/baby-interpreter/pkg/tokenizer"
)

type Parser struct {
	tokenizer *tokenizer.Tokenizer

	currentToken token.Token
	peekToken    token.Token
}

func New(tokenizer *tokenizer.Tokenizer) *Parser {
	p := &Parser{tokenizer: tokenizer}

	// Read tokens in pairs, so currToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.tokenizer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
