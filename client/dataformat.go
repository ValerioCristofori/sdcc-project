package main

import (
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

func (t *Dataformat) Get(args Args, dataResult *Data) error {
	return nil
}

func (t *Dataformat) Put(args Args, dataResult *Data) error {
	return nil
}

func (t *Dataformat) Delete(args Args, dataResult *Data) error {
	return nil
}

func (t *Dataformat) Append(args Args, dataResult *Data) error {
	return nil
}