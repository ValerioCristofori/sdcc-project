package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sdcc-project/rpc-logic/dataformat"
	"strings"
)

//address and port on which RPC server is listening
var port 	= 12345
var addr 	= fmt.Sprintf( "localhost:%d", port)


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
	value 		:= ""
	if numArgs > 3 {
		value = os.Args[3]
	}

	// Init data input for RPC
	args := &dataformat.Args{Key: key, Value: value}

	// Asynchronous call RPC
	if strings.EqualFold(command,"get") {

		// GET body
		reply := new(dataformat.Data)
		log.Printf("Asynchronous call to RPC server")

		call := client.Go("Dataformat.Get", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Error in Dataformat.Get: ", call.Error.Error())
		}

		fmt.Printf("Dataformat.Get:\n Key:\t%s\nValue:\n%s\nTimestamp:\t%s\n", key, reply.Value, reply.Timestamp.String() )

	} else if strings.EqualFold(command,"put") {

		// PUT body
		if numArgs < 4 {
			errorArgs = true
			goto argumentsError
		}
		reply := new(dataformat.Data)
		log.Printf("Asynchronous call to RPC server")

		call := client.Go("Dataformat.Put", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Error in Dataformat.Put: ", call.Error.Error())
		}

		fmt.Printf("Dataformat.Put:\n Key:\t%s\nValue:\n%s\nTimestamp:\t%s\n", key, reply.Value, reply.Timestamp.String() )


	} else if strings.EqualFold(command,"delete") {

		// DELETE body
		reply := new(dataformat.Data)
		log.Printf("Asynchronous call to RPC server")

		call := client.Go("Dataformat.Delete", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Error in Dataformat.Delete: ", call.Error.Error())
		}

		fmt.Printf("Dataformat.Delete:\n Key:\t%s\nTimestamp:\t%s\n", key, reply.Timestamp.String() )


	} else if strings.EqualFold(command,"append") {

		// APPEND body
		if numArgs < 4 {
			errorArgs = true
			goto argumentsError
		}
		reply := new(dataformat.Data)
		log.Printf("Asynchronous call to RPC server")

		call := client.Go("Dataformat.Append", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Error in Dataformat.Append: ", call.Error.Error())
		}

		fmt.Printf("Dataformat.Append:\n Key:\t%s\nValue:\n%s\nTimestamp:\t%s\n", key, reply.Value, reply.Timestamp.String() )



	}else {
		errorArgs = true
		goto argumentsError
	}



}