package ast 

import ( 
	"xxx/token"
) 


type Node interface{ 
	TokenLiteral() string
} 

type Statement interface{ 
	Node
	statementNode()
}

type Expression interface{ 
	Node 
	expressionNode() 
}

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


type LetStatement struct { 
	Token token.Token
	Name *Identifier 
	Value Expression 
} 

func (ls *LetStatement) TokenLiteral() string { 
	return ls.Token.Literal
}
func (ls *LetStatement) statementNode () {} 

type Identifier struct {
	Token token.Token 
	Value string 
}

func (i *Identifier) expressionNode() {} 
func(i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type ReturnStatement struct { 
	Token token.Token 
	ReturnValue Expression 
}

func (rs *ReturnStatement) TokenLiteral() string { 
	return rs.Token.Literal 
} 

func (rs *ReturnStatement) statementNode () {} 
 
