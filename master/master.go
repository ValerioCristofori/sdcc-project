package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
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
	configuration   = Configuration{}
	running			= false
)

type Configuration struct {
	NumNodes    int
	AwsRegion	string
}

func readConfig()  {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
}


func (l *Listener)Register( address string, replyReg *ReplyMessage) error {

	mutex.Lock()
	log.Println("Adding address " + address )
	NodesAddress = append(NodesAddress, address)
	mutex.Unlock()

	for {
		if len(NodesAddress) == configuration.NumNodes{
			fmt.Println("All Nodes registered")
			running = true
			break
		}
	}

	replyReg.Ack = true
	return nil
}

func (l *Listener)GetAddresses(_ int, listAddresses *[]string) error {

	for{
		if running{
			break
		}
	}
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

func createDynamoDBTable( tableName string )  {
	err := initDynamoDB(tableName)
	if err != nil {
		log.Fatal("Error in Init DynamoDB: ", err)
	}
	fmt.Println("Created DynamoDB Table")
}


/*
 * Master node
 * provide the authentication system for client nodes, listing all the addresses of the edge nodes
 */
func main()  {
	readConfig()
	createDynamoDBTable("Sensors")
	serveRequests()

}



