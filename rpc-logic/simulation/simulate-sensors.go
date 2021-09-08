package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sdcc-project/rpc-logic/client"
	"time"
)

type Sensor struct {
	Id string
}

var NSensors = 3
var rangeFloats float64 = 100.00
var sensors []Sensor



func (s *Sensor) getMeasure() error {

	// Put the first measure
	measure := rand.Float64()*rangeFloats
	timestamp := time.Now()
	client.RpcEdgeNode("put", s.Id, fmt.Sprintf("%f", measure), timestamp )

	// Append every 5 seconds
	for range time.Tick(5 * time.Second){
		measure := rand.Float64()*rangeFloats
		timestamp := time.Now()
		client.RpcEdgeNode("append", s.Id, fmt.Sprintf("%f", measure), timestamp )
	}

	return nil
}


func main()  {

	// Init sensors
	fmt.Printf("Simulate %d sensors", NSensors)

	for i := 0; i < NSensors; i++ {
		id := fmt.Sprint("id-sensor-",i)
		sensors = append( sensors, Sensor{id})
	}

	for sensorIndex := range sensors {
		go func(currentSensor *Sensor) {
			err := currentSensor.getMeasure()
			if err != nil {
				log.Fatal("error in simulating measure")
				os.Exit(1)
			}
		}(&sensors[sensorIndex])
	}

	time.Sleep(20 * time.Second)
	fmt.Println("Done!")
}
