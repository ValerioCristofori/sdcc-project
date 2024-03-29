package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var ErrorCmd = "Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}"

func consumptionSite()  {
	// Provide user interface as also consumption site
	// BOUNDARY
	// run forever until user issue bye
	fmt.Println("************************")
	fmt.Println("*     CLIENT SHELL     *")
	fmt.Println("************************")
	fmt.Println("Type 'exit' to exit the shell")
	for {
		var command 	= ""
		var key			= ""
		var value 		= ""

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>>")
		scanner.Scan()
		in := scanner.Text()
		arguments := strings.Split(in, " ")
		if strings.HasPrefix(arguments[0], "exit") {
			fmt.Println("Good bye!")
			os.Exit(0)
		}

		if len(arguments) < 2 {
			log.Println(ErrorCmd)
			continue
		}
		command = arguments[0]
		key = arguments[1]
		value = ""
		if len(arguments) > 2 {
			value = fmt.Sprintf("{%s %s}", time.Now().String(), arguments[2])
		}


		// Controllo sintattico
		switch command {
		case "put": if len(arguments) < 3 {
			log.Println(ErrorCmd)
			continue
		}else{
			break
		}
		case "append": if len(arguments) < 3 {
			log.Println(ErrorCmd)
			continue
		}else{
			break
		}
		case "delete":
		case "get":
			break
		default:
			log.Println(ErrorCmd)
			continue
		}

		//call RPC func
		RpcBroadcastEdgeNode(command, key, value)



	}
}

func main()  {
	// Set right edge node address
	GetEdgeAddresses()

	consumptionSite()
}