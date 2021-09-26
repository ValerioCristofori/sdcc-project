package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"syscall"
	"time"
)
// server conf
var(
	dataPath	= "/data"
	dataDebugPath = "/debug-data"
	path  		= "/"
	debugPath 	= "/debug"
	port 		= 12345
	masterPort	= 8080
	masterAddr 	= fmt.Sprintf( "master:%d", masterPort)
    myAddress string

)

type ReplyMessage struct {
	Ack bool
}

type Cluster struct {
	Nodes    			[]string
	indexEdgeRequest 	int
}

// data conf
var (
	df 				*Dataformat
	cluster 			= new(Cluster)
	listEndPointsRPC 	= new([]*rpc.Client)
)


// interface to RaftRPC
var rfRPC *RaftRPC



func serveData(){
	//Create an instance of struct
	df = new(Dataformat)

	// Init
	InitMap()
	//InitDynamo()

	// Register a new RPC server and the struct we created above.
	server := rpc.NewServer()
	err := server.RegisterName("Dataformat", df) // important for calling right func
	if err != nil {
		log.Fatal("Format of service Datastore is not correct: ", err)
	}
	// Register an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath
	server.HandleHTTP(dataPath, dataDebugPath)

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


func register() {
	// RPC request to master:8080
	// retrieve list all edge node addresses
	// Try to connect to masterAddr using HTTP protocol
	var reply *ReplyMessage
	var client *rpc.Client

	// Try to connect to master
	client, err := rpc.DialHTTP("tcp", masterAddr)
	if err != nil {
		log.Println("Error in dialing: ", err)
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

func getMyAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}


func connectToAllNodes() {

	for _, nameAddress := range cluster.Nodes {
		if nameAddress != myAddress {
			address := fmt.Sprintf("%s:%d", nameAddress, port)
			var client *rpc.Client

			// Try to connect to master
			client, err := rpc.DialHTTP("tcp", address)
			if err != nil {
				log.Println("Error in dialing: ", err)
			}

			*listEndPointsRPC = append( *listEndPointsRPC, client )
		}
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

func main()  {

	// start configuration and initialization of raft cluster

	register()

	time.Sleep(3 * time.Second)
	getListEdgeNodes()
	//var rf = new(Raft)
	//go startListener(rf)
	serverRPC := rpc.NewServer()
	go startListener(serverRPC)
	time.Sleep(3 * time.Second)
	connectToAllNodes()

	persister := MakePersister()
	rfRPC = Make( *listEndPointsRPC, cluster.indexEdgeRequest, persister, nil)
	addHandlerRaft(serverRPC, rfRPC)
	addHandlerData(serverRPC, new(Dataformat))



	syscall.Pause()
}










