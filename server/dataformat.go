package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
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
	Timestamp time.Time
}

type DataformatReply struct {
	DataResult  *Data
	Ack			bool
}

type Data struct {
	Value string
	Timestamp time.Time
}



// Map : K -> key, V -> data struct
var datastore map[string]Data
type Dataformat int //edge node

// mutex for sync
var mutex = sync.RWMutex{}

func InitMap() error {
	datastore = make(map[string]Data)
	return nil
}


func (t *Dataformat) Get(args Args, reply *DataformatReply) error {
	op := GET
	rfRPC.rf.Start(Command{Op: op,Key: args.Key,Timestamp: args.Timestamp})
	// Get from the datastore
	if d, found := datastore[args.Key]; found {
		reply.Ack = true
		*reply.DataResult = d
	} else {
		return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	}
	// timestamp of the PUT operation
	return nil
}


func (t *Dataformat) Put(args Args, reply *DataformatReply) error {
	op := PUT
	_,_,isLeader := rfRPC.rf.Start(Command{Op: op,Key: args.Key,Value: args.Value,Timestamp: args.Timestamp})
	if !isLeader {
		// the op called in a not leader edge node
		reply.Ack = false
		return nil
	}
	reply.Ack = true
	//// Build data struct
	//data := Data{args.Value, time.Now()}
	//// Save in the Datastore
	//mutex.Lock()
	//datastore[args.Key] = data
	//mutex.Unlock()
	//
	//// Return data to the caller
	//*dataResult = data

	return nil
}

func (t *Dataformat) Delete(args Args, reply *DataformatReply) error {
	op := DELETE
	_,_,isLeader := rfRPC.rf.Start(Command{Op: op,Key: args.Key,Timestamp: args.Timestamp})
	if !isLeader {
		// the op called in a not leader edge node
		reply.Ack = false
		return nil
	}
	reply.Ack = true
	//// Delete in the Datastore
	//if _, found := datastore[args.Key]; found {
	//	delete(datastore, args.Key)
	//}else {
	//	return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	//}
	return nil
}

func (t *Dataformat) Append(args Args, reply *DataformatReply) error {
	op := APPEND
	_,_,isLeader := rfRPC.rf.Start(Command{Op: op,Key: args.Key,Value: args.Value,Timestamp: args.Timestamp})
	if !isLeader {
		// the op called in a not leader edge node
		reply.Ack = false
		return nil
	}

	reply.Ack = true

	//// Build data struct
	//data := Data{args.Value, time.Now()}
	//// Save in the Datastore
	//mutex.Lock()
	//if d, found := datastore[args.Key]; found {
	//	d.Value = d.Value + "\n" + args.Value // dummy append
	//	d.Timestamp = data.Timestamp
	//	// update in memory
	//	datastore[args.Key] = d
	//	// update the result
	//	data = d
	//} else {
	//	// Normal Put func
	//	datastore[args.Key] = data
	//}
	//mutex.Unlock()
	//// Return data to the caller
	//*dataResult = data

	return nil
}