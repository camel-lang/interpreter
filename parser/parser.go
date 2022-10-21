package parser

import (
	"camel/ast"
	"camel/lexer"
	"camel/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,

	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.PLUS:     SUM,
	token.MINUS:    SUM,

	token.LT: LESSGREATER,
	token.GT: LESSGREATER,

	token.LPAREN: CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lex       *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
func New(lex *lexer.Lexer) *Parser {

	parser := &Parser{lex: lex}

	parser.nextToken()
	parser.nextToken()

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiterals)
	parser.registerPrefix(token.STRING , parser.parseStringLiteral)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral) 
	parser.registerPrefix(token.LBRACE, parser.parseHashLiteral) 

	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

	return parser
}

func (p *Parser) nextToken() {

	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {

	switch p.curToken.Type {

	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturn()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {

	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturn() *ast.ReturnStatement {

	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {

	prefixFunc := p.prefixParseFns[p.curToken.Type]
	if prefixFunc == nil {
		fmt.Printf("no prefix function registered for %v",
			p.curToken.Type)
		return nil
	}

	left := prefixFunc()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFunc := p.infixParseFns[p.peekToken.Type]
		if infixFunc == nil {
			return left
		}
		p.nextToken()
		left = infixFunc(left)
	}

	return left
}

func (p *Parser) parseIdentifier() ast.Expression {

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		fmt.Errorf(err.Error())
		return nil
	}
	lit.Value = value
	return lit
}
 
func (p *Parser) parseStringLiteral() ast.Expression { 
	return &ast.StringLiteral{ 
		Token: p.curToken, 
		Value: p.curToken.Literal,
	}
	
} 

func (p *Parser) parseArrayLiteral() ast.Expression { 

	array := &ast.ArrayLiteral{Token: p.curToken} 
	array.Elements = p.parseExpressionList(token.RBRACKET)
	
	return array
} 

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression { 
	
	list := []ast.Expression{}
	
	if p.peekTokenIs(end) { 
		p.nextToken() 
		return list
	} 		

	p.nextToken() 
	list = append(list, p.parseExpression(LOWEST))
	
	for p.peekTokenIs(token.COMMA) { 
		p.nextToken() 
		p.nextToken() 
		list = append(list, p.parseExpression(LOWEST)) 
	}

	if !p.expectPeek(end) { 
		return nil 
	}

	return list 
} 

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression { 
	
	exp := &ast.IndexExpression{Token: p.curToken, Left: left} 
	
	p.nextToken() 
	exp.Index = p.parseExpression(LOWEST) 

	if !p.expectPeek(token.RBRACKET) { 
		return nil 
	} 

	return exp 
} 

func (p *Parser) parseHashLiteral() ast.Expression {
	
	hash := &ast.HashLiteral{ Token: p.curToken } 
	hash.Pairs = make(map[ast.Expression]ast.Expression) 
	
	for !p.peekTokenIs(token.RBRACE) { 

		p.nextToken() 
		key := p.parseExpression(LOWEST) 
		if !p.expectPeek(token.COLON) { 
			return nil 
		} 
		p.nextToken()

		value := p.parseExpression(LOWEST) 
		hash.Pairs[key] = value  
	
		if !p.peekTokenIs(token.RBRACE) && 
		   !p.expectPeek(token.COMMA) { 
			return nil 
		} 

	} 

	p.nextToken() 
	return hash 
} 

func (p *Parser) parsePrefixExpression() ast.Expression {

	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {

	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {

	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiterals() ast.Expression {
	foo := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	foo.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	foo.Body = p.parseBlockStatement()
	return foo
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{}
	for p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

func (p *Parser) curTokenIs(tok token.TokenType) bool {
	return p.curToken.Type == tok
}

func (p *Parser) peekTokenIs(tok token.TokenType) bool {
	return p.peekToken.Type == tok
}

func (p *Parser) expectPeek(tok token.TokenType) bool {
	if p.peekTokenIs(tok) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) peekPrecedence() int {

	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {

	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
