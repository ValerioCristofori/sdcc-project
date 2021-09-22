package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"testing"
)


func TestClient(t *testing.T)  {

	SetEdgeAddressTest()

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

	// Provide user interface as also consumption site
	// BOUNDARY
	// run forever until user issue bye
	for {
		var command 	= ""
		var key			= ""
		var value 		= ""

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">")
		scanner.Scan()
		in := scanner.Text()
		arguments := strings.Split(in, " ")
		if strings.HasPrefix(arguments[0], "bye") {
			fmt.Println("Good bye!")
			os.Exit(0)
		}

		if len(arguments) < 2 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}
		command = arguments[0]
		key = arguments[1]
		if len(arguments) > 2 {
			value = arguments[2]
		}
		timestamp := time.Now()

		// Controllo sintattico
		switch command {
		case "put": if len(arguments) < 3 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}else{
			break
		}
		case "append": if len(arguments) < 3 {
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}else{
			break
		}
		case "delete":
		case "get":
			break
		default:
			log.Fatal("Not valid args\nInsert args in the form: <get/put/delete/append> <key> {<value>}")
			os.Exit(1)
		}

		//call RPC func
		RpcEdgeNode(command, key, value, timestamp )

	}

}
