package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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
	Value 		string
	Counter 	int64 // older if smaller
}


type GlobalCounter  int64
type Dataformat 	int //edge node


var (
	 DIMENSION 		= 1000
	 datastore 		sync.Map
	 counter 		GlobalCounter = 0
)

func (c *GlobalCounter) inc() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *GlobalCounter) get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func PrintMap()  {
	// loop over elements of slice
	fmt.Println("Printing Map Datastore")
	// with syncmap, looping over all keys is simple without locking the whole map for the entire loop
	datastore.Range(func(key, value interface{}) bool {
		// cast value to correct format
		val, ok := value.(string)
		if !ok {
			// this will break iteration
			return false
		}
		// do something with key/value
		fmt.Println(key, " value is ", val)

		// this will continue iterating
		return true
	})
}


func (t *Dataformat) Get(args Args, dataResult *Data) error {
	// Get from the datastore
	//if found in datastore return
	data, ok := datastore.Load(args.Key)
	if ok {
		*dataResult = Data{Value: data.(Data).Value, Counter: counter.inc()}
		datastore.Store( args.Key, *dataResult)
		return nil
	}
	//else search it in the cloud
	resp := GetLambda(args)
	if resp.Value != "" {
		d := Data{Value: resp.Value, Counter: counter.inc()}
		*dataResult = d
		datastore.Store( args.Key, *dataResult)
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

func Len( sm *sync.Map) int {
	length := 0
	sm.Range(func(key, value interface{}) bool{
		length = length + len(key.(string)) + len(value.(Data).Value) + 8
		return true
	})
	return length
}

func cleanThread()  {
	var min int64
	var keyToDelete string
	for {
		for Len(&datastore) >= 2*DIMENSION/3 {
			min = counter.get()
			//fmt.Printf("Clean Datastore until 2/3*MaxSize: counter %d, size %d\n", min, Len(&datastore) )
			datastore.Range(func(key, data interface{}) bool {
				if data.(Data).Counter < min {
					min = data.(Data).Counter
					keyToDelete = key.(string)
				}
				return true
			})
			// delete oldest entry
			datastore.Delete(keyToDelete)
			fmt.Println("CLEAN: Deleted key: ", keyToDelete)
		}
		time.Sleep(5*time.Second)
	}
}