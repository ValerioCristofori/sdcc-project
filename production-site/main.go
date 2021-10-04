package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"syscall"
	"time"
)



func (s *Sensor) getMeasure() error {

	rand.Seed(time.Now().UnixNano())
	// Put the first measure
	measure := rand.Float64()*rangeFloats
	fmt.Printf("Measuring %f\n",measure)
	//timestamp := time.Now()
	//if len(leaderEdgeAddr) > 0 {
	//	RpcSingleEdgeNode("put", s.Id, fmt.Sprintf("%f", measure), timestamp, leaderEdgeAddr )
	//}else {
	//	RpcBroadcastEdgeNode("put", s.Id, fmt.Sprintf("%f", measure), timestamp)
	//}


	// Append every 5 seconds
	for range time.Tick(5 * time.Second){
		measure := rand.Float64()*rangeFloats
		if len(leaderEdgeAddr) > 0 {
			fmt.Println("RPC Leader address" + leaderEdgeAddr)
			RpcSingleEdgeNode("append", s.Id, fmt.Sprintf("%f", measure), leaderEdgeAddr )
		}else {
			fmt.Println("RPC broadcast")
			RpcBroadcastEdgeNode("append", s.Id, fmt.Sprintf("%f", measure) )
		}
	}

	return nil
}

func productionSite()  {
	// Init sensors
	fmt.Printf("Simulate %d sensors", NSensors)

	for i := 0; i < NSensors; i++ {
		id := fmt.Sprint("sensor-",i)
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
	syscall.Pause()
}



func main()  {
	time.Sleep(10 * time.Second)
	// Set right edge node address
	GetEdgeAddresses()

	productionSite()

}
