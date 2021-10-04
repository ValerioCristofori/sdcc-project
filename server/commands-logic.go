package main

import (
	"fmt"
	"log"
	"os"
)

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
	data := Data{args.Value}
	// Save in the Datastore
	mutex.Lock()
	datastore[args.Key] = data
	mutex.Unlock()
	log.Println("PUT entry on datastore: {key: " + args.Key + "} {value: " + args.Value + "}" )
	appendOnLogFile("PUT{key: " + args.Key + "}{value: " + args.Value + "}\n")
}

func DeleteEntry(args *Args)  {
	// Delete in the Datastore
	mutex.Lock()
	defer mutex.Unlock()
	if _, found := datastore[args.Key]; found {
		mutex.Lock()
		delete(datastore, args.Key)
		mutex.Unlock()
	}else {
		log.Printf(fmt.Sprintf("key %s not in datastore",args.Key))
		return
	}
	log.Println("DELETE entry on datastore: {key: " + args.Key + "}" )

}

func AppendEntry(args *Args)  {
	// Build data struct
	data := Data{args.Value}
	// Save in the Datastore
	mutex.Lock()
	defer mutex.Unlock()
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


