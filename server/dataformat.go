package main

import (
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



	return nil
}