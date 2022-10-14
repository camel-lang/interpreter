package parser 

import ( 
	"camel/ast" 
	"camel/token" 
	"camel/lexer"
	"strconv"
	"fmt"
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
) 

type ( 
	prefixParseFn func() ast.Expression 
	infixParseFn func(ast.Expression) ast.Expression 
)

type Parser struct { 
	lex *lexer.Lexer
	curToken token.Token 
	peekToken token.Token
	
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType , fn prefixParseFn) { 
	p.prefixParseFns[tokenType] = fn 
}

func (p *Parser) registerInfix(tokenType token.TokenType , fn infixParseFn) { 
	p.infixParseFns[tokenType] = fn
}
func New(lex *lexer.Lexer) *Parser { 

	parser := &Parser{lex: lex}
					 

	parser.nextToken() 
	parser.nextToken() 

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn) 	
	parser.registerPrefix(token.IDENT , parser.parseIdentifier) 
	parser.registerPrefix(token.INT , parser.parseIntegerLiteral) 

	return parser
}

func (p *Parser) nextToken() { 
	
	p.curToken = p.peekToken 
	p.peekToken = p.lex.NextToken() 
}

func (p *Parser) ParseProgram() (*ast.Program) { 
	
	program := &ast.Program{} 
	program.Statements = []ast.Statement{} 
	
	for p.curToken.Type != token.EOF { 

		stmt := p.parseStatement() 
		if stmt != nil { 
			program.Statements = append(program.Statements , stmt)
		}
		p.nextToken() 
	}
	return program 
}


func (p *Parser) parseStatement() ast.Statement { 

	switch p.curToken.Type { 

		case token.LET :
			return p.parseLetStatement() 
		case token.RETURN : 
			return p.parseReturn() 
		default : 
			return p.parseExpressionStatement() 
	}
}


func (p *Parser) parseLetStatement() *ast.LetStatement { 


	stmt := &ast.LetStatement{Token: p.curToken}
	 
	if !p.expectPeek(token.IDENT) {
		return nil 
	}

	stmt.Name = &ast.Identifier{Token: p.curToken , Value: p.curToken.Literal}
	
	if !p.expectPeek(token.ASSIGN) { 
		return nil 
	}

	for !p.curTokenIs(token.SEMICOLON) { 
		p.nextToken() 
	}

	return stmt 
}	

func (p *Parser) parseReturn() *ast.ReturnStatement {
	
	stmt := &ast.ReturnStatement{Token : p.curToken}

	p.nextToken() 
	
	for !p.curTokenIs(token.SEMICOLON) { 
		p.nextToken() 
	}
	
	return stmt 
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement { 
	
	stmt := &ast.ExpressionStatement{ Token: p.curToken } 
	
	stmt.Expression = p.parseExpression(LOWEST) 
	
	if p.peekTokenIs(token.SEMICOLON) { 
		p.nextToken() 
	}
		
	return stmt 
} 

func (p *Parser) parseExpression(precedence int) ast.Expression {
	
	prefix := p.prefixParseFns[p.curToken.Type] 
	if prefix == nil { 
		return nil 
	} 
	
	return prefix() 
} 

func (p *Parser) parseIdentifier() ast.Expression { 

	return &ast.Identifier{ Token: p.curToken , Value: p.curToken.Literal } 
}

func (p *Parser) parseIntegerLiteral() ast.Expression { 
	
	lit := &ast.IntegerLiteral{Token: p.curToken} 
	value , err := strconv.ParseInt(p.curToken.Literal , 0 , 64) 
	if err != nil { 
		fmt.Errorf(err.Error()) 
		return nil 
	}
	lit.Value = value 
	return lit
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


 
