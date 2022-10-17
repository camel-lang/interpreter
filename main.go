package main

import (
	"camel/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Println(repl.Logo) 
	fmt.Printf("Hello! %s, This is xxx programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
