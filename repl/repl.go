package repl

import (
	"bufio"
	"camel/lexer"
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
		fmt.Printf("%s\n", program.String())
	}
}
