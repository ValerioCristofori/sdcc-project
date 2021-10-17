package main

import (
	"errors"
	"fmt"
	"sync"
)

// operations
const (
	PUT int = iota
	GET
	APPEND
	DELETE
)

var DIMENSION = 1000


type Args struct {
	Key string
	Value string
}

type Response struct {
	Ack 	bool
	Key 	string
	Value 	string
}

type DataformatReply struct {
	DataResult  *Data
	Ack			bool
}

type Data struct {
	Value string
}



// Map : K -> key, V -> data struct
var datastore map[string]Data
type Dataformat int //edge node

// mutex for sync
var mutex = sync.RWMutex{}

func InitMap() error {
	//create local datastore
	datastore = make(map[string]Data)
	return nil
}

func PrintMap()  {
	// loop over elements of slice
	fmt.Println("Printing Map Datastore")
	for k, v := range datastore {
		fmt.Println(k, "value is", v)
	}
}

//func checkDimension(args Args){
//
//	memoryBytes := 0
//
//	mutex.Lock()
//	//check how much storage is used
//	for k, v:= range datastore{
//		memoryBytes = memoryBytes + len(k) + len(v.Value) + 4
//	}
//	mutex.Unlock()
//	if (memoryBytes + len(args.Key) + len(args.Value) + 4) >= 2 * DIMENSION/3 {
//
//		fmt.Println("Too Values on Local Map.\nSending to DynamoDB")
//		go putOnDynamoDB()
//	}
//
//
//
//
//}

//func putOnDynamoDB() {
//
//	mutex.Lock()
//	defer mutex.Unlock()
//	//free up memory until it is half of the total
//	for len(datastore) >= DIMENSION / 2 {
//		count := 0
//		var max string
//
//
//		for k, v := range datastore {
//
//			if count == 0 {
//
//				max = k
//
//			} else {
//				if len(v.Value) >= len(datastore[max].Value){
//					max = k
//				}
//			}
//			count++
//
//		}
//
//		item := Args{max, datastore[max].Value}
//		//send to dynamodb the value with max dimension
//		go putItem(item)
//		//delete value from local storage
//		DeleteEntry(&item)
//		fmt.Println(item.Value)
//	}
//
//}


func (t *Dataformat) Get(args Args, dataResult *Data) error {
	// Get from the datastore
	mutex.Lock()
	defer mutex.Unlock()
	//if found in datastore return
	if d, found := datastore[args.Key]; found {
		*dataResult = d
		return nil
	}
	//else search it in the cloud
	resp := GetLambda(args)
	if resp.Value != "" {
		d := Data{resp.Value}
		*dataResult = d
		return nil
	}else {
		return errors.New(fmt.Sprintf("key %s not in datastore and not in database",args.Key) )
	}


}


func (t *Dataformat) Put(args Args, reply *DataformatReply) error {
	op := PUT

	//checkDimension(args)

	_,_,isLeader := rfRPC.rf.Start(Command{Op: op,Key: args.Key,Value: args.Value})
	if !isLeader {
		// the op called in a not leader edge node
		reply.Ack = false
		return nil
	}

	go PutLambda(args)

	reply.Ack = true
	//if leader do immediately the op

	return nil
}


func (t *Dataformat) Delete(args Args, reply *DataformatReply) error {
	op := DELETE
	_,_,isLeader := rfRPC.rf.Start(Command{Op: op,Key: args.Key})
	if !isLeader {
		// the op called in a not leader edge node
		reply.Ack = false
		return nil
	}
	reply.Ack = true
	go DeleteLambda(args)
	//if leader do immediately the op
	return nil
}

func (t *Dataformat) Append(args Args, reply *DataformatReply) error {
	op := APPEND
	_,_,isLeader := rfRPC.rf.Start(Command{Op: op,Key: args.Key,Value: args.Value})
	if !isLeader {
		// the op called in a not leader edge node
		reply.Ack = false
		return nil
	}

	//append also in dynamodb
	go AppendLambda(args)


	reply.Ack = true
	//if leader do immediately the op
	return nil
}
