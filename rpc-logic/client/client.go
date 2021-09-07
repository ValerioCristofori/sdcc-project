package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sdcc-project/rpc-logic/dataformat"
	"strings"
	"time"
)

//address and port on which RPC server is listening
var port = 12345
var addr = fmt.Sprintf( "localhost:%d", port)


func main()  {

	// Try to connect to addr using HTTP protocol
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()
	numArgs := len(os.Args)
	errorArgs := false
	if numArgs < 3 {
		errorArgs = true
	}

	// Terminate program because of arguments error
	argumentsError:
	if errorArgs {
		log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
		os.Exit(1)
	}


	// Init variables
	command 	:= os.Args[1]
	key 		:= os.Args[2]
	timestamp 	:= time.Now() // current local timestamp
	value := ""
	if numArgs > 3 {
		value = os.Args[3]
	}

	// Init data input for RPC
	args := &dataformat.Data{Key: key, Value: value, Timestamp: timestamp}

	// Asynchronous call RPC
	if strings.EqualFold(command,"get") {

		// GET body
		divReply := new(dataformat.Get)
		log.Printf("Asynchronous call to RPC server")

		divCall := client.Go("Dataformat.Get", args, divReply, nil)
		divCall = <-divCall.Done
		if divCall.Error != nil {
			log.Fatal("Error in Dataformat.Get: ", divCall.Error.Error())
		}

		fmt.Printf("Dataformat.Get:\n Key:\t%s\nValue:\t%s\nTimestamp:\t%s\n", divReply.Key, divReply.Value, divReply.Timestamp.String() )

	} else if strings.EqualFold(command,"put") {

		// PUT body
		if numArgs < 4 {
			errorArgs = true
			goto argumentsError
		}
		divReply := new(dataformat.Put)
		log.Printf("Asynchronous call to RPC server")

		divCall := client.Go("Dataformat.Put", args, divReply, nil)
		divCall = <-divCall.Done
		if divCall.Error != nil {
			log.Fatal("Error in Dataformat.Put: ", divCall.Error.Error())
		}


	} else if strings.EqualFold(command,"delete") {

		// DELETE body
		divReply := new(dataformat.Delete)
		log.Printf("Asynchronous call to RPC server")

		divCall := client.Go("Dataformat.Delete", args, divReply, nil)
		divCall = <-divCall.Done
		if divCall.Error != nil {
			log.Fatal("Error in Dataformat.Delete: ", divCall.Error.Error())
		}

	} else if strings.EqualFold(command,"append") {

		// APPEND body
		if numArgs < 4 {
			errorArgs = true
			goto argumentsError
		}
		divReply := new(dataformat.Append)
		log.Printf("Asynchronous call to RPC server")

		divCall := client.Go("Dataformat.Append", args, divReply, nil)
		divCall = <-divCall.Done
		if divCall.Error != nil {
			log.Fatal("Error in Dataformat.Append: ", divCall.Error.Error())
		}


	}else {
		errorArgs = true
		goto argumentsError
	}



}