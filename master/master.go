package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"syscall"
	"time"
)

// server conf
var (
	path  		= "/"
	debugPath 	= "/debug"
	port 		= 8080
	beatPort	= 8088
)

type Listener int

// simulate 3 edge nodes
// static name nodes
var edgeNodes = []string {"edge-node-1","edge-node-2","edge-node-3"}

// Clo check for running servers
var Clo chan bool

type ReplyMessage struct {
	Ack bool
	TaskStatus bool
}


func heartBeatAllAddresses() {
	for _, address := range edgeNodes {
		Clo = make(chan bool)
		client, err := jsonrpc.Dial("tcp", fmt.Sprintf("%s:%d", address, beatPort) )
		if err != nil {
			log.Fatal("dialhttp: ", err)
		}
		defer client.Close()
		go Handler(client)
		select {
		case <-Clo:
			fmt.Println("Shut down connection")
		}
	}
}

func Handler(client *rpc.Client) {
	go heartBeat(client)
}

func heartBeat(client *rpc.Client) {
	i :=0
	for {
		time.Sleep(4 * time.Second)

		var reply *ReplyMessage
		fmt.Println("send HeartBeat"+strconv.Itoa(i))

		err := client.Call("RpcHeartBeat.HeartBeat", 0, &reply)
		fmt.Println(*reply)
		if err != nil {
			log.Fatal("call heartBeat: ", err)
		}
		if reply.Ack{
			fmt.Println("heartBeat normal")
		}
		if reply.TaskStatus{
			fmt.Println("need to get task")
		}
		i++
		//go GetTask(client);


	}
}

func GetTask(client *rpc.Client) {

	//var reply *domain.DeployTaskJob
	//err := client.Call("MyRPC.DeployJob", "agent1", &reply)
	//if err != nil {
	//	log.Fatal("call DeployJob: ", err)
	//}
	//fmt.Println("getTask success")
	//fmt.Println(*reply)
}

func (l *Listener)GetAddresses(_ int, listAddresses *[]string) error {

	if len(edgeNodes) != 0 {
		*listAddresses = edgeNodes
	} else {
		return errors.New(fmt.Sprintf("no edge node running") )
	}
	return nil
}

func serveClientRequests()  {
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


/*
 * Master node provide monitoring for edge node with heartbeats every 4 seconds
 * Also provide the authentication system for client nodes, listing all the addresses of the edge nodes
 */
func main()  {

	go serveClientRequests()

	// Start a Heart Beat routine to check server running
	go heartBeatAllAddresses()

	syscall.Pause()

}

