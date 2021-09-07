package dataformat

import (
	"fmt"
	"time"
)



type Data struct {
	Key string
	Value string
	Timestamp time.Time
}

// Dataformat can return Data ( if 'get')
type Dataformat int


func (d *Dataformat) Get(key string, dataResult *Data) error {
	fmt.Printf("Get\n")
	*dataResult = Data{key," ", time.Now()}
	return nil
}


func (d *Dataformat) Put(data Data, dataResult *Data) error {
	fmt.Printf("Put\n")
	return nil
}

func (d *Dataformat) Delete(data Data, dataResult *Data) error {
	fmt.Printf("Delete\n")
	return nil
}

func (d *Dataformat) Append(data Data, dataResult *Data) error {
	fmt.Printf("Append\n")
	return nil
}