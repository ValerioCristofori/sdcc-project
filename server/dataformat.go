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
	// Get from the datastore
	mutex.Lock()
	if d, found := datastore[args.Key]; found {
		reply.Ack = true
		*reply.DataResult = d
	} else {
		return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	}
	mutex.Unlock()
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
	//if leader do immediately the op
	PutEntry(&args)

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
	//if leader do immediately the op
	DeleteEntry(&args)
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
	//if leader do immediately the op
	AppendEntry(&args)
	return nil
}