package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"syscall"
)
// server conf
var(
	path  		= "/"
	debugPath 	= "/debug"
	port 		= 12345
	beatPort	= 8088
)

// log for fault recovery


// data conf
var df *Dataformat

// RpcHeartBeat type for heart beat routine
type RpcHeartBeat int

type ReplyMessage struct {
	Ack bool
	TaskStatus bool
}



func (r *RpcHeartBeat) HeartBeat(m int, reply *ReplyMessage) error {
	reply.Ack = true

	//Traverse the task queue here
	reply.TaskStatus = true
	return nil
}

func serveData(){
	//Create an instance of struct
	df = new(Dataformat)

	// Init
	InitMap()
	createLogFile()
	//InitDynamo()

	// Register a new RPC server and the struct we created above.
	server := rpc.NewServer()
	err := server.RegisterName("Dataformat", df) // important for calling right func
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

func serveHeartBeat() {
	// Create TCP Listener for heart beat
	rpcListener := new(RpcHeartBeat)
	rpc.Register(rpcListener)

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", beatPort) )
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
		fmt.Println("Start listening")
	}
}


func main()  {

	go serveData()

	go serveHeartBeat()

	syscall.Pause()
}


