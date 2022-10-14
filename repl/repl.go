package repl

import (
	"bufio"
	"camel/lexer"
	"camel/token"
	"fmt"
	"io"
)

const logo = `

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
	fmt.Println(logo)

	for {

		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
