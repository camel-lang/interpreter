package object

import (
	"bytes"
	"camel/ast"
	"fmt"
	"strings"
	"hash/fnv"
)

type ObjectType string
type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ 	     = "STRING" 
	BOOLEAN_OBJ      = "BOOLEAN"
	ARRAY_OBJ 		 = "ARRAY"
	HASH_OBJ 		 = "HASH" 
	NULL             = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type String struct { 
	Value string 
}

func (s *String) Type() ObjectType { 
	return STRING_OBJ
} 
func (s *String) Inspect() string { 
	return s.Value
} 


type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "Error: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {

	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())
	out.WriteString("{\n")
	out.WriteString("\n}")

	return out.String()

}

type Builtin struct { 
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { 
	return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string { 
	return "builtin function" 
}


type Array struct { 
	Elements []Object
} 

func (a *Array) Type() ObjectType { 
	return ARRAY_OBJ
}
func (a *Array) Inspect() string { 
	
	var out bytes.Buffer 
	
	elems := []string{} 
	for _, e := range a.Elements { 
		elems = append(elems, e.Inspect())
	}
	out.WriteString("[") 
	out.WriteString(strings.Join(elems, ", ")) 
	out.WriteString("]") 
	
	return out.String() 
}

type Hashable interface { 
	HashKey() HashKey 
} 

type HashKey struct { 
	Type ObjectType
	Value uint64
}   


func (b *Boolean) HashKey() HashKey { 
	
	var value uint64 
		
	if b.Value { 
		value = 1 
	} else { 
		value = 0 
	} 

	return HashKey{Type: b.Type(), Value: value} 
} 
func (i *Integer) HashKey() HashKey { 
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
} 
func (s *String) HashKey() HashKey { 
	
	h := fnv.New64a() 
	h.Write([]byte(s.Value)) 

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct { 
	Key Object 
	Value Object 
} 

type Hash struct { 
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { 
	return HASH_OBJ
} 

func (h *Hash) Inspect() string { 
	
	var out bytes.Buffer 
	
	pairs := []string{}

	for _, p := range h.Pairs { 
		pairs = append(pairs, fmt.Sprintf("%s : %s",
		p.Key.Inspect(), p.Value.Inspect()))
	} 

	out.WriteString("{") 
	out.WriteString(strings.Join(pairs, ", ")) 
	out.WriteString("}") 
	
	return out.String() 
}  
