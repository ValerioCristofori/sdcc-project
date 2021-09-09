package dataformat

import (
	"errors"
	"fmt"
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

func InitMap() error {
	datastore = make(map[string]Data)
	return nil
}

func (t *Dataformat) Get(args Args, dataResult *Data) error {

	// Get from the datastore
	if d, found := datastore[args.Key]; found {
		*dataResult = d
	} else {
		return errors.New(fmt.Sprintf("key %s not in datastore",args.Key) )
	}
	return nil
}


func (t *Dataformat) Put(args Args, dataResult *Data) error {

	// Build data struct
	data := Data{args.Value, time.Now()}
	// Save in the Datastore
	datastore[args.Key] = data
	//fmt.Printf("value: %s\ntimestamp: %s\n", data.Value, data.Timestamp.String()  )
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
	return nil
}

func (t *Dataformat) Append(args Args, dataResult *Data) error {

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

	// Return data to the caller
	*dataResult = data

	return nil
}