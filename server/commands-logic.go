package main

import (
	"fmt"
	"log"
	"time"
)

func PutEntry(args *Args)  {
	// Build data struct
	data := Data{args.Value, time.Now()}
	// Save in the Datastore
	datastore[args.Key] = data
	log.Println("PUT entry on datastore: {key: " + args.Key + "} {value: " + datastore[args.Key].Value + "}" )
}

func DeleteEntry(args *Args)  {
	// Delete in the Datastore
	if _, found := datastore[args.Key]; found {
		delete(datastore, args.Key)
	}else {
		log.Printf(fmt.Sprintf("key %s not in datastore",args.Key))
		return
	}
	log.Println("DELETE entry on datastore: {key: " + args.Key + "}" )

}

func AppendEntry(args *Args)  {
	// Build data struct
	data := Data{args.Value, time.Now()}
	// Save in the Datastore
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
	log.Println("APPEND entry on datastore: {key: " + args.Key + "} {value: " + datastore[args.Key].Value + "}" )


}


