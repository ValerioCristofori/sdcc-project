package main

import (
	"fmt"
	"log"
	"net/rpc"
	"strings"
	"time"
)

//address and port on which RPC server is listening
var(
	port 				= 12345
	masterPort 			= 8080
	masterAddr 			= fmt.Sprintf( "master:%d", masterPort)
	leaderEdgeAddr 		string
	allNodesAddr		[]string
	errorLeader			= false
)

type DataformatReply struct {
	DataResult  *Data
	Ack			bool
}

type Args struct {
	Key string
	Value string
	Timestamp time.Time
}

type Data struct {
	Value string
	Timestamp time.Time
}

func SetEdgeAddressTest(){
	leaderEdgeAddr = fmt.Sprintf("localhost:%d", port)
}

func GetEdgeAddresses() {
	// RPC request to master:8080
	// retrieve list all edge node addresses
	// Try to connect to masterAddr using HTTP protocol
	var client *rpc.Client

	// Try to connect to master
	client, err := rpc.DialHTTP("tcp", masterAddr)
	if err != nil {
		log.Println("Error in dialing: ", err)
	}
	defer client.Close()


	// Call remote procedure
	err = client.Call("Listener.GetAddresses", 0, &allNodesAddr)
	if err != nil {
		log.Fatal("Error in Listener.GetAddresses: ", err)
	}
}

func RpcBroadcastEdgeNode(command string, key string, value string, timestamp time.Time)  {
restart:
	for i:=0; i<len(allNodesAddr); i++ {
		if !RpcSingleEdgeNode(command, key, value, timestamp, allNodesAddr[i]) {
			errorLeader = true
			break
		}
	}
	if errorLeader {
		fmt.Println("Retrying send command...")
		time.Sleep(1 * time.Second)
		goto restart
	} else {
		// command executed, reset var to false
		errorLeader = false
	}
}

func RpcSingleEdgeNode(command string, key string, value string, timestamp time.Time, edgeAddr string ) bool {

	var client *rpc.Client

	// Try to connect to edgeAddr using HTTP protocol
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf( "%s:%d", edgeAddr, port))
	if err != nil{
		fmt.Println("Error in dialing: ", err)
		if strings.Contains(err.Error(), "connect") {
			// no connection to the host -> edge node down
			if edgeAddr == leaderEdgeAddr {
				fmt.Println("Leader is down. Founding new leader...")
				// leader is down
				// next call is broadcast rpc to all edge
				leaderEdgeAddr = ""
				return false
			} else {
				return true
			}

		}
	}

	// Init data input for RPC
	args := &Args{Key: key, Value: value, Timestamp: timestamp}

	// Asynchronous call RPC
	if strings.EqualFold(command,"get") {

		// GET body
		reply := &DataformatReply{}
		reply.DataResult = &Data{}

		call := client.Go("Dataformat.Get", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Println("Error in Dataformat.Get: ", call.Error.Error())
		}

	} else if strings.EqualFold(command,"put") {

		// PUT body
		reply := &DataformatReply{}

		call := client.Go("Dataformat.Put", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Println("Error in Dataformat.Put: ", call.Error.Error())
		}
		// check if i call the leader or not
		if reply.Ack{
			leaderEdgeAddr = edgeAddr
		}


	} else if strings.EqualFold(command,"delete") {
		// DELETE body
		reply := &DataformatReply{}


		call := client.Go("Dataformat.Delete", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Println("Error in Dataformat.Delete: ", call.Error.Error())
		}
		// check if i call the leader or not
		if reply.Ack{
			leaderEdgeAddr = edgeAddr
		}

	} else if strings.EqualFold(command,"append") {

		// APPEND body
		reply := &DataformatReply{}


		call := client.Go("Dataformat.Append", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Println("Error in Dataformat.Append: ", call.Error.Error())
		}
		// check if i call the leader or not
		if reply.Ack{
			leaderEdgeAddr = edgeAddr
		}


	}


	return true



}