package repl

import ( 

	"bufio" 
	"fmt" 
	"io"
	"xxx/token"
	"xxx/lexer"
)  

const PROMPT = ">> " 

func Start(in io.Reader , out io.Writer) { 

	scanner := bufio.NewScanner(in) 

	for { 

		fmt.Printf(PROMPT) 
		scanned := scanner.Scan() 
		if !scanned {
			return 
		} 
		 
		line := scanned.Text() 	
		lex := lexer.New(line) 
		for tok := lex.NextToken() ; tok.Type != token.EOF ; tok = lex.NextToken{ 
			fmt.Printf("%+v\n" , tok )
		} 
	} 
} 
