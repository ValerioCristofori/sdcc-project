package main


import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"sdcc-project/rpc-logic/dataformat"
)

var path  		= "/"
var debugPath 	= "/debug"
var port 		= 12345

func main()  {

	//Create an instance of struct which implements Arith interface
	dataformat := new(dataformat.Dataformat)

	// Register a new RPC server and the struct we created above.
	server := rpc.NewServer()
	err := server.RegisterName("Datastore", dataformat)
	if err != nil {
		log.Fatal("Format of service Datastore is not correct: ", err)
	}
	// Register an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath
	server.HandleHTTP(path, debugPath)

	// Listen for incoming messages on port 12345
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	log.Printf("RPC server on port %d", port)

	// Start go's http server on socket specified by listener
	err = http.Serve(lis, nil)
	if err != nil {
		log.Fatal("Serve error: ", err)
	}

}
