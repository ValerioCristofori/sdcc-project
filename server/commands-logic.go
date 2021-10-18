package main

import (
	"fmt"
	"log"
)

func PutEntry(args *Args)  {
	// Build data struct
	data := Data{Value: args.Value, Counter: counter.inc()}
	// Save in the Datastore
	datastore.Store( args.Key, data)
	log.Println("PUT entry on datastore: {key: " + args.Key + "} {value: " + args.Value + "}" )
}

func DeleteEntry(args *Args)  {
	// Delete in the Datastore
	_, ok := datastore.Load(args.Key)
	if ok {
		datastore.Delete(args.Key)
	}else{
		//log.Printf(fmt.Sprintf("key %s not in datastore",args.Key))
		return
	}
	fmt.Println("DELETE entry on datastore: {key: " + args.Key + "}" )

}

func AppendEntry(args *Args)  {
	// Save in the Datastore
	data, ok := datastore.Load(args.Key)
	if ok {
		resultValue := data.(Data).Value + "\n" + args.Value
		datastore.Store(args.Key, Data{Value: resultValue, Counter: counter.inc()})
	}else{
		//log.Printf(fmt.Sprintf("key %s not in datastore",args.Key))
		//normal put
		datastore.Store(args.Key, Data{ Value: args.Value, Counter: counter.inc()})
	}
	log.Println("APPEND entry on datastore: {key: " + args.Key + "} {value: " + args.Value + "}" )

}


