package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

// operations
const (
	PUT int = iota
	GET
	APPEND
	DELETE
)

var DIMENSION int64 = 1000000//represent 8 bytes

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
	return nil
}

func PrintMap()  {
	// loop over elements of slice
	fmt.Println("Printing Map Datastore")
	for k, v := range datastore {
		fmt.Println(k, "value is", v)
	}
}


func (t *Dataformat) Get(args Args, dataResult *Data) error {
	// Get from the datastore
	mutex.Lock()
	defer mutex.Unlock()
	if d, found := datastore[args.Key]; found {
		*dataResult = d
		return nil
	}
	item := getItem(args.Key)
	if item.Value != "" {
		d := Data{item.Value}
		*dataResult = d
		return nil
	} else {
		return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	}
}

func checkDimension(args Args) bool{
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := f.Stat()

	if err!=nil{
		log.Fatal(err)
	}
	fileDim := fi.Size()

	if fileDim + int64(len(args.Key)) + int64(len(args.Value))>DIMENSION {
		return false
	}

	return true
}

func (t *Dataformat) Put(args Args, reply *DataformatReply) error {
	op := PUT

	isFree:=checkDimension(args)

	if !isFree{
		putItem(args)
		return nil
	}

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

	if getItem(args.Key).Value!="" {
		appendItem(args)
	}

	reply.Ack = true
	//if leader do immediately the op
	return nil
}
