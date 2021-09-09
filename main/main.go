package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sdcc-project/client"
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

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">")
		//********************************************************* Problem last char
		scanner.Scan()
		in := scanner.Text()
		arguments := strings.Split(in, " ")
		if strings.HasPrefix(arguments[0], "bye") {
			fmt.Println("Good bye!")
			os.Exit(0)
		}

		if len(arguments) < 2 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}
		command = arguments[0]
		key = arguments[1]
		if len(arguments) > 2 {
			value = arguments[2]
		}
		timestamp := time.Now()

		// Controllo sintattico
		switch command {
		case "put": if len(arguments) < 3 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}else{
			break
		}
		case "append": if len(arguments) < 3 {
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