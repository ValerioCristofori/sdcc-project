package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Args struct {
	Key string
	Value string
	Timestamp time.Time
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

// logFile all operations to datastore are registered here
var logFile *os.File

func InitMap() error {
	datastore = make(map[string]Data)
	return nil
}

func createLogFile() {
	file, err := os.OpenFile("operations.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logFile = file
	log.SetOutput(logFile)

}

func (t *Dataformat) Get(args Args, dataResult *Data) error {

	// Get from the datastore
	if d, found := datastore[args.Key]; found {
		*dataResult = d
	} else {
		return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	}
	// timestamp of the PUT operation
	log.Printf("GET: key:%s value:%s timestamp:%s \n", args.Key, dataResult.Value, dataResult.Timestamp.String()  )
	return nil
}


func (t *Dataformat) Put(args Args, dataResult *Data) error {


	// Build data struct
	data := Data{args.Value, time.Now()}
	// Save in the Datastore
	mutex.Lock()
	datastore[args.Key] = data
	mutex.Unlock()
	log.Printf("PUT: key:%s value:%s timestamp:%s \n", args.Key, data.Value, data.Timestamp.String()  )
	// Return data to the caller
	*dataResult = data

	return nil
}

func (t *Dataformat) Delete(args Args, dataResult *Data) error {

	// Delete in the Datastore
	if _, found := datastore[args.Key]; found {
		delete(datastore, args.Key)
	}else {
		return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	}
	log.Printf("DELETE: key:%s \n", args.Key )
	return nil
}

func (t *Dataformat) Append(args Args, dataResult *Data) error {

	// Build data struct
	data := Data{args.Value, time.Now()}
	// Save in the Datastore
	mutex.Lock()
	if d, found := datastore[args.Key]; found {
		d.Value = d.Value + "\n" + args.Value // dummy append
		d.Timestamp = data.Timestamp
		// update in memory
		datastore[args.Key] = d
		// update the result
		data = d
	} else {
		// Normal Put func
		datastore[args.Key] = data
	}
	mutex.Unlock()
	// Return data to the caller
	*dataResult = data
	log.Printf("APPEND: key:%s value:%s timestamp:%s \n", args.Key, data.Value, data.Timestamp.String()  )

	return nil
}