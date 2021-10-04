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

type Args struct {
	Key string
	Value string
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
	datastore = make(map[string]Data)
	//initDynamoDB("Sensors")
	////wait for table creation
	//for {
	//	tables := callTable()
	//	if tables == 0{
	//		print("Wait")
	//	}else {
	//		return nil
	//	}
	//}



	return nil
}


func (t *Dataformat) Get(args Args, dataResult *Data) error {
	// Get from the datastore
	mutex.Lock()
	defer mutex.Unlock()
	if d, found := datastore[args.Key]; found {
		*dataResult = d
	} else {
		return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	}
	return nil
}


func (t *Dataformat) Put(args Args, reply *DataformatReply) error {
	op := PUT
	_,_,isLeader := rfRPC.rf.Start(Command{Op: op,Key: args.Key,Value: args.Value})
	if !isLeader {
		// the op called in a not leader edge node
		reply.Ack = false
		return nil
	}
	reply.Ack = true
	//if leader do immediately the op

	//Communication with DynamoDB
	//se è troppo grande invia a dynamodb

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
	deleteItem(args)
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
	reply.Ack = true
	//if leader do immediately the op
	return nil
}
