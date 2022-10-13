package ast 

import ( 
	"xxx/token"
	"bytes"
) 


type Node interface{ 
	TokenLiteral() string
	String string
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

func (p *Program) String() string { 
	var out bytes.Buffer 
	
	for _ , s := range p.Statements { 
		out.WriteString(s.String()) 	
	}	
	return out.String() 
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

func (ls *LetStatement) String() string { 
	
	var out bytes.Buffer
	
	out.WriteString(ls.TokenLiteral() + " ") 
	out.WriteString(ls.Name.Value) 
	out.WriteString(" = ")
	
	if ls.Value != nil { 
		out.WriteString(ls.Value.String()) 
	} 
	out.WriteString(";") 
	return out.String()
} 

type Identifier struct {
	Token token.Token 
	Value string 
}

func (i *Identifier) expressionNode() {} 
func(i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string { 
	return i.Value 
} 

type ReturnStatement struct { 
	Token token.Token 
	ReturnValue Expression 
}

func (rs *ReturnStatement) TokenLiteral() string { 
	return rs.Token.Literal 
} 

func (rs *ReturnStatement) statementNode () {} 
 
func (rs *ReturnStatement) String() string {
	
	var out bytes.Buffer 
	
	out.WriteString(rs.TokenLiteral() + " ") 
	if rs.ReturnValue != nil { 
		out.WriteString(ReturnValue.String()) 
	}

	out.WriteString(";") 
	return out.String() 
} 

type ExpressionStatement struct { 
	Token token.Token 
	Expression Expression 	
}

func (exs *ExpressionStatement) TokenLiteral() string { 
	return exs.Token.Literal 
} 

func (exs *ExpressionStatement) statementNode() {} 

func (exs *ExpressionStatement) String() string { 
	if exs.Expression != nil {
		return exs.Expression.String() 
	} 
	return "" 
} 


