package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)


func register() {
	// RPC request to master:8080
	// retrieve list all edge node addresses
	// Try to connect to masterAddr using HTTP protocol
	var reply *ReplyMessage
	var client *rpc.Client

	retry:
	// Try to connect to master
	client, err := rpc.DialHTTP("tcp", masterAddr)
	if err != nil {
		log.Println("Error in dialing: ", err)
		// retry in 500 millisec
		time.Sleep(500*time.Millisecond)
		goto retry
	}
	defer client.Close()

	// Call remote procedure
	log.Printf("Synchronous call to RPC master for registration")
	myAddress = getMyAddress()
	err = client.Call("Listener.Register", myAddress, &reply)
	if err != nil {
		log.Fatal("Error in Listener.Register: ", err)
	}


}

func addHandlerData(server *rpc.Server, df *Dataformat) {
	// Register a new RPC server and the struct we created above.
	err := server.RegisterName("Dataformat", df) // important for calling right func
	if err != nil {
		log.Fatal("Format of service listener is not correct: ", err)
	}
}

func addHandlerRaft(server *rpc.Server, rfRPC *RaftRPC)  {
	// Register a new RPC server and the struct we created above.
	err := server.RegisterName("RaftRPC", rfRPC) // important for calling right func
	if err != nil {
		log.Fatal("Format of service listener is not correct: ", err)
	}
}

func startListener(server *rpc.Server) {
	// Start Listener for Raft

	// Register an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath
	server.HandleHTTP(path, debugPath)

	// Listen for incoming messages on port 8088
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	log.Printf("RPC server on port %d for Raft", port)

	// Start go's http server on socket specified by listener
	err = http.Serve(lis, nil)
	if err != nil {
		log.Fatal("Serve error: ", err)
	}

}



func connectToAllNodes() {

	for _, nameAddress := range cluster.Nodes {
		address := fmt.Sprintf("%s:%d", nameAddress, port)
		var client *rpc.Client

		retry:
		// Try to connect to master
		client, err := rpc.DialHTTP("tcp", address)
		if err != nil {
			log.Println("Error in dialing for raft: ", err)
			goto retry
		}

		*listEndPointsRPC = append( *listEndPointsRPC, client )

	}
}


func getListEdgeNodes() {

	var client *rpc.Client
	var me 	int
	listAddresses := new([]string)

	// Try to connect to master
	client, err := rpc.DialHTTP("tcp", masterAddr)
	if err != nil {
		log.Println("Error in dialing: ", err)
	}
	defer client.Close()

	// Call remote procedure
	log.Printf("Synchronous call to RPC master for list addresses")

	err = client.Call("Listener.GetAddresses", 0, listAddresses)
	if err != nil {
		log.Fatal("Error in Listener.GetAddresses: ", err)
	}
	for index, edgeNodeAddress := range *listAddresses {
		if myAddress == edgeNodeAddress {
			fmt.Printf("match in: %s in index %d\n" , myAddress , index )
			me = index
			break
		}
	}
	cluster.Nodes = *listAddresses
	cluster.indexEdgeRequest = me

}
