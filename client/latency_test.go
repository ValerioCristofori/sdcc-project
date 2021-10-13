package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var(
	TotQuery = 10000
	rangeKeys = 100
)

func appendOnLogFile( entry string )  {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("results/test-latency.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func getCalls(wg *sync.WaitGroup, numGet int)  {
	defer wg.Done()
	for i:=0; i<numGet; i++{
		// build random key
		var key 	= rand.Intn(rangeKeys)
		// send rpc get to a random address
		var address = allNodesAddr[rand.Intn(len(allNodesAddr))]
		RpcSingleEdgeNode("get", fmt.Sprintf("%d", key), "", address)
	}
}

func putCalls(wg *sync.WaitGroup, numPut int)  {
	defer wg.Done()
	for i:=0; i<numPut; i++{
		// build random key-value
		var key 	= rand.Intn(rangeKeys)
		var value 	= rand.Float64()*rangeFloats
		timestamp := time.Now()
		// send rpc put
		if len(leaderEdgeAddr) > 0 {
			RpcSingleEdgeNode("put", fmt.Sprintf("%d", key), fmt.Sprintf("{%s %f}", timestamp.Format(timestampFormat), value), leaderEdgeAddr )
		}else {
			RpcBroadcastEdgeNode("put", fmt.Sprintf("%d", key), fmt.Sprintf("{%s %f}", timestamp.Format(timestampFormat), value) )
		}
	}
}

func appendCalls(wg *sync.WaitGroup, numAppend int) {
	defer wg.Done()
	for i:=0; i<numAppend; i++{
		// build random key-value
		var key 	= rand.Intn(rangeKeys)
		var value 	= rand.Float64()*rangeFloats
		timestamp := time.Now()

		// send rpc put
		if len(leaderEdgeAddr) > 0 {
			RpcSingleEdgeNode("append", fmt.Sprintf("%d", key), fmt.Sprintf("{%s %f}", timestamp.Format(timestampFormat), value), leaderEdgeAddr )
		}else {
			RpcBroadcastEdgeNode("append", fmt.Sprintf("%d", key), fmt.Sprintf("{%s %f}", timestamp.Format(timestampFormat), value) )
		}
	}

}

func test1()  {
	// test for 85% GET and 15% PUT
	rand.Seed(time.Now().UnixNano())
	var (
		numGet 		= TotQuery*85/100
		numPut 		= TotQuery*15/100
		wg			= new(sync.WaitGroup)
	)
	wg.Add(2)
	go getCalls(wg, numGet)
	go putCalls(wg, numPut)

	wg.Wait()

}

func test2()  {
	// test for 40% PUT, 20% APPEND and 40% GET
	var (
		numGet 		 = TotQuery*40/100
		numPut 		 = TotQuery*40/100
		numAppend 	 = TotQuery*20/100
		wg			 = new(sync.WaitGroup)
	)
	wg.Add(3)
	go getCalls(wg, numGet)
	go putCalls(wg, numPut)
	go appendCalls(wg, numAppend)

	wg.Wait()

}

func Test(t *testing.T)  {
	// Set right edge node address
	GetEdgeAddresses()

	fmt.Println("STARTING TESTS..")

	start := time.Now()
	test1()
	timeTest1 := time.Since(start).Milliseconds()
	appendOnLogFile(fmt.Sprintf("test1,%d,%d\n", TotQuery,timeTest1))
	fmt.Printf("Test 1 finished in %d milliseconds\n", timeTest1)

	start = time.Now()
	test2()
	timeTest2 := time.Since(start).Milliseconds()
	appendOnLogFile(fmt.Sprintf("test2,%d,%d\n", TotQuery, timeTest2))
	fmt.Printf("Test 2 finished in %d milliseconds\n", timeTest2)

}
