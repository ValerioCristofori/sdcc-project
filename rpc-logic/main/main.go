package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sdcc-project/rpc-logic/client"
	"strings"
	"time"
)



func main()  {

	// Provide user interface as also consumption site
	// BOUNDARY
	// run forever until user issue bye
	for {
		var command 	= ""
		var key			= ""
		var value 		= ""

		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Print(">")
		//********************************************************* Problem last char
		input, _ := consoleReader.ReadString('\n')
		input = strings.ToLower(input)
		if strings.HasPrefix(input, "bye") {
			fmt.Println("Good bye!")
			os.Exit(0)
		}

		// Parsing the user input
		stringTokens := strings.Split(input," ")
		if len(stringTokens) < 2 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}
		command = stringTokens[0]
		key = stringTokens[1]
		if len(stringTokens) > 2 {
			value = stringTokens[2]
		}
		timestamp := time.Now()

		// Controllo sintattico
		switch command {
		case "put": if len(stringTokens) < 3 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}else{
			break
		}
		case "append": if len(stringTokens) < 3 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}else{
			break
		}
		case "delete":
		case "get":
			break
		default:
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}

		//call RPC func
		client.RpcEdgeNode(command, key, value, timestamp )

	}




}