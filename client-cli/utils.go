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
	port 		= 12345
	masterPort 	= 8080
	masterAddr 	= fmt.Sprintf( "master:%d", masterPort)
	leaderEdgeAddr 		string
	allNodesAddr		[]string
)
// random simulation rtt
var rangeRTT int64 = 20

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

func GetEdgeAddresses()  {
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
	log.Printf("Synchronous call to RPC server")
	err = client.Call("Listener.GetAddresses", 0, &allNodesAddr)
	if err != nil {
		log.Fatal("Error in Listener.GetAddresses: ", err)
	}
}

func RpcBroadcastEdgeNode(command string, key string, value string, timestamp time.Time)  {

	for _, address := range allNodesAddr {
		RpcSingleEdgeNode(command, key, value, timestamp, address)
	}
}

func RpcSingleEdgeNode(command string, key string, value string, timestamp time.Time, edgeAddr string )  {

	var client *rpc.Client

	// Try to connect to edgeAddr using HTTP protocol
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf( "%s:%d", edgeAddr, port))
	if err != nil{
		log.Fatal("Error in dialing: ", err)
	}

	// Init data input for RPC
	args := &Args{Key: key, Value: value, Timestamp: timestamp}

	// Asynchronous call RPC
	if strings.EqualFold(command,"get") {

		// GET body
		reply := DataformatReply{}
		reply.DataResult = &Data{}

		call := client.Go("Dataformat.Get", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Println("Error in Dataformat.Get: ", call.Error.Error())
			return
		}

		fmt.Printf("Dataformat.Get:\n Key:\t%s\nValue:\n%s\nTimestamp:\t%s\n", key, reply.DataResult.Value, reply.DataResult.Timestamp.String() )

	} else if strings.EqualFold(command,"put") {

		// PUT body
		reply := &DataformatReply{}

		call := client.Go("Dataformat.Put", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Error in Dataformat.Put: ", call.Error.Error())
		}

		//fmt.Printf("Dataformat.Put:\n Key:\t%s\nValue:\n%s\nTimestamp:\t%s\n", key, reply.Value, reply.Timestamp.String() )


	} else if strings.EqualFold(command,"delete") {
		// DELETE body
		reply := &DataformatReply{}


		call := client.Go("Dataformat.Delete", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Error in Dataformat.Delete: ", call.Error.Error())
		}

		//fmt.Printf("Dataformat.Delete:\n Key:\t%s\nTimestamp:\t%s\n", key, reply.Timestamp.String() )


	} else if strings.EqualFold(command,"append") {

		// APPEND body
		reply := &DataformatReply{}


		call := client.Go("Dataformat.Append", args, reply, nil)
		call = <-call.Done
		if call.Error != nil {
			log.Fatal("Error in Dataformat.Append: ", call.Error.Error())
		}

		//fmt.Printf("Dataformat.Append:\n Key:\t%s\nValue:\n%s\nTimestamp:\t%s\n", key, reply.Value, reply.Timestamp.String() )



	}



}