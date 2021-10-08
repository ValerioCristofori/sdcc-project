package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"syscall"
)

// server conf
var (
	path  		= "/"
	debugPath 	= "/debug"
	edgePort	= 12345
	port 		= 8080
)

type Listener int

type ReplyMessage struct {
	Ack bool
}

var(
	NodesAddress   []string
 	mutex 		    = sync.RWMutex{}
)


func (l *Listener)Register( address string, replyReg *ReplyMessage) error {

	mutex.Lock()
	log.Println("Adding address " + address )
	NodesAddress = append(NodesAddress, address)
	mutex.Unlock()

	replyReg.Ack = true
	return nil
}

func (l *Listener)GetAddresses(_ int, listAddresses *[]string) error {

	if len(NodesAddress) != 0 {
		*listAddresses = NodesAddress
	} else {
		return errors.New(fmt.Sprintf("no edge node running") )
	}
	return nil
}

func serveRequests()  {
	/**
	 * Start Listener
	 */
	listener := new(Listener)
	// Register a new RPC server and the struct we created above.
	server := rpc.NewServer()
	err := server.RegisterName("Listener", listener) // important for calling right func
	if err != nil {
		log.Fatal("Format of service listener is not correct: ", err)
	}
	// Register an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath
	server.HandleHTTP(path, debugPath)

	// Listen for incoming messages on port 8088
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


/*
 * Master node provide monitoring for edge node with heartbeats every 4 seconds
 * Also provide the authentication system for client nodes, listing all the addresses of the edge nodes
 */
func main()  {
	err := initDynamoDB("Sensors")
	if err != nil {
		syscall.Pause()
		log.Fatal("Error in Init DynamoDB: ", err)
	}
	serveRequests()
}



