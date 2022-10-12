package parser 

import ( 
	"xxx/ast" 
	"xxx/token" 
	"xxx/lexer"
)

type Parser struct { 
	lex *lexer.Lexer
	curToken token.Token 
	peekToken token.Token
}

func New(lex *lexer.Lexer) *Parser { 

	parser := &Parser{lex: lex}
					 

	parser.nextToken() 
	parser.nextToken() 

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
			return nil 	
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

func (p *Parser) parseReturn() *ast.ReturnStatement {
	
	stmt := &ast.ReturnStatement{Token : p.curToken}

	p.nextToken() 
	
	for !p.curTokenIs(token.SEMICOLON) { 
		p.nextToken() 
	}
	
	return stmt 
} 
