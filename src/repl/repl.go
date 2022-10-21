package repl

import (
	"bufio"
	"camel/eval"
	"camel/lexer"
	"camel/object"
	"camel/parser"
	"fmt"
	"io"
)

const Logo = `


     __,  .-.  .-.
    (__ _ |  \/  | _
       (_|||\__/||(/_|
          ||    ||   |_
         _||   _||
         ""    ""

------------------------------------------------`

/*
Thank you for visiting https://asciiart.website/
This ASCII pic can be found at
https://asciiart.website/index.php?art=animals/camels
*/

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {

		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		parser := parser.New(lex)
		program := parser.ParseProgram()
		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
