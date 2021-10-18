package main

import (
	"fmt"
	"log"
	"sync"
)

var mutexD = sync.Mutex{}


func PutEntry(args *Args)  {
	// Build data struct
	data := Data{args.Value}
	// Save in the Datastore
	mutexD.Lock()
	defer mutexD.Unlock()
	datastore[args.Key] = data
	//log.Println("PUT entry on datastore: {key: " + args.Key + "} {value: " + args.Value + "}" )
}

func DeleteEntry(args *Args)  {
	// Delete in the Datastore
	mutexD.Lock()
	defer mutexD.Unlock()
	if _, found := datastore[args.Key]; found {
		delete(datastore, args.Key)
	}else {
		log.Printf(fmt.Sprintf("key %s not in datastore",args.Key))
		return
	}
	//fmt.Println("DELETE entry on datastore: {key: " + args.Key + "}" )

}

func AppendEntry(args *Args)  {
	// Build data struct
	data := Data{args.Value}
	// Save in the Datastore
	mutexD.Lock()
	defer mutexD.Unlock()
	if d, found := datastore[args.Key]; found {
		d.Value = d.Value + "\n" + args.Value // dummy append
		// update in memory
		datastore[args.Key] = d
		// update the result
		data = d
	} else {
		// Normal Put func if key is not in datastore
		datastore[args.Key] = data
	}
	//log.Println("APPEND entry on datastore: {key: " + args.Key + "} {value: " + args.Value + "}" )

}


