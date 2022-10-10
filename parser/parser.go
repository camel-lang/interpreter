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

func (p *Parser) ParseProgram *ast.Program { 
	return nil
}
