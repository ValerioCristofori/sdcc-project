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
var toClean chan bool
type Dataformat int //edge node

// mutex for sync
var mutex = sync.RWMutex{}

func InitMap() error {
	datastore = make(map[string]Data)
	//toClean = make(chan bool, 1)
	//go cleanMap()
	return nil
}

func cleanMap()  {
	for{
		tooValues := <-toClean
		if tooValues {
			fmt.Println("Too Values on Local Map.\nSending to DynamoDB")
			putOnDynamoDB()
			toClean <- false
		}
	}


}

func PrintMap()  {
	// loop over elements of slice
	fmt.Println("Printing Map Datastore")
	for k, v := range datastore {
		fmt.Println(k, "value is", v)
	}
}

func checkDimension(args Args){

	memoryBytes := 0
	for k, v:= range datastore{
		memoryBytes += len(k) + len(v.Value) + 4
	}

	if memoryBytes >= (2/3)* DIMENSION {
		toClean <- true
	}




}

func putOnDynamoDB() {

	for len(datastore) >= DIMENSION*2/3 {

		count := 0
		var min int
		var key string
		//invia a dynamodb i valori con timestamp maggiore liberando spazio sull'edge node
		for k, v := range datastore {

			if count == 0 {

				min = v.Counter
				key = k
			} else {
				if min >= v.Counter {
					min = v.Counter
					key = k
				}
			}
			count++

		}
		item := Args{key, datastore[key].Value, datastore[key].Counter}
		putItem(item)
		DeleteEntry(&item)

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
	}else {
		return errors.New(fmt.Sprintf("key %s not in datastore and not in database",args.Key) )
	}/*
	item := getItem(args.Key)
	if item.Value != "" {
		d := Data{item.Value, item.Counter+1}
		*dataResult = d
		PutEntry(&item)
		return nil
	}*/
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
