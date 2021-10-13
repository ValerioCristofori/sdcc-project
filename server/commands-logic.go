package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var mutexD = sync.Mutex{}

func appendOnLogFile( entry string )  {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte(entry)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func PutEntry(args *Args)  {
	// Build data struct
	data := Data{args.Value, args.Counter}
	// Save in the Datastore
	mutexD.Lock()
	defer mutexD.Unlock()
	datastore[args.Key] = data
	log.Println("PUT entry on datastore: {key: " + args.Key + "} {value: " + args.Value + "}" )
	appendOnLogFile("PUT{key: " + args.Key + "}{value: " + args.Value + "}\n")
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
	fmt.Println("DELETE entry on datastore: {key: " + args.Key + "}" )

}

func AppendEntry(args *Args)  {
	// Build data struct
	data := Data{args.Value, args.Counter}
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
		// Normal Put func
		datastore[args.Key] = data
	}
	log.Println("APPEND entry on datastore: {key: " + args.Key + "} {value: " + args.Value + "}" )
	appendOnLogFile("APPEND{key: " + args.Key + "}{value: " + args.Value + "}\n")

}


