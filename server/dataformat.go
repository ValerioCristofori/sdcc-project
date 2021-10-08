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

var DIMENSION int64 = 100000//represent 8 bytes

type Args struct {
	Key string
	Value string
	Counter int
}

type DataformatReply struct {
	DataResult  *Data
	Ack			bool
}

type Data struct {
	Value string
	Counter int
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

func checkDimension(args Args){
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
		putSomeItemsOnDynamoDB()
	}



}

func putSomeItemsOnDynamoDB() {

	for int64(len(datastore)) >= DIMENSION*2/3{

		//invia a dynamodb i valori con timestamp maggiore liberando spazio sull'edge node

	}


}


func (t *Dataformat) Get(args Args, dataResult *Data) error {
	// Get from the datastore
	mutex.Lock()
	defer mutex.Unlock()
	if d, found := datastore[args.Key]; found {
		*dataResult = d
		d.Counter = d.Counter + 1
		return nil
	}
	item := getItem(args.Key)
	if item.Value != "" {
		d := Data{item.Value, item.Counter+1}
		*dataResult = d
		return nil
	} else {
		return errors.New(fmt.Sprintf("key %s not in datastore and not in database",args.Key) )
	}
}


func (t *Dataformat) Put(args Args, reply *DataformatReply) error {
	op := PUT

	/*isFree:=checkDimension(args)

	if !isFree{
		fmt.Println("PUT ON DYNAMODB")
		putItem(args)
		return nil
	}*/

	checkDimension(args)

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
