package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)
// server conf
var(

	path  		= "/"
	debugPath 	= "/debug"
	port 		= 12345
	masterPort	= 8080
	masterAddr 	= fmt.Sprintf( "master:%d", masterPort)
    myAddress string
)

type ReplyMessage struct {
	Ack bool
}

type Cluster struct {
	Nodes    			[]string
	indexEdgeRequest 	int
}

type Configuration struct {
	AwsRegion			string
	OptionCleaning		bool
}

func (c *Cluster) toString() string{
	return fmt.Sprintf("Cluster: %s, My position: %d", c.Nodes, c.indexEdgeRequest)
}

// data conf
var (
	cluster 			= new(Cluster)
	listEndPointsRPC 	= new([]*rpc.Client)
	configuration   	= Configuration{}
	// interface to RaftRPC
	rfRPC 			*RaftRPC
	// channel for newly committed messages
	applyCh 		chan ApplyMsg
)







func readConfig()  {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func getMyAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}


func applyChRoutine()  {
	for m := range applyCh {
		if m.UseSnapshot{
			//ignore snapshot
		}else{
			args := &Args{}
			args.Key = m.Command.Key
			args.Value = m.Command.Value
			switch m.Command.Op {
			case PUT: PutEntry(args)
			case APPEND: AppendEntry(args)
			case DELETE: DeleteEntry(args)
			default:
				continue
			}
		}

	}
}



func shutdownHandler() {
	m := map[string]string{}
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		datastore.Range(func(key, value interface{}) bool {
			m[key.(string)] = value.(Data).Value
			return true
		})
		fmt.Println(sig)
		fmt.Println("Saving persist log entries and Raft state...")
		if err := Save("./vol/backup", m); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Exiting...")
		os.Exit(2)

	}()

	go func() {
		for range time.Tick(10 * time.Second){
			datastore.Range(func(key, value interface{}) bool {
				m[key.(string)] = value.(Data).Value
				return true
			})
			fmt.Println("Saving persist log entries and Raft state...")
			if err := Save("./vol/backup", m); err != nil {
				log.Fatalln(err)
			}
		}
	}()
}

func main()  {

	// start configuration and initialization of raft cluster
	readConfig()
	register()
	getListEdgeNodes()

	serverRPC := rpc.NewServer()
	go startListener(serverRPC)
	connectToAllNodes()

	err := Load("./vol/backup", &datastore )
	if err != nil {
		log.Println("Not able to backup persistent state")
	}
	shutdownHandler()
	// listen to messages from Raft indicating newly committed messages.
	applyCh = make(chan ApplyMsg)
	go applyChRoutine()

	rfRPC = Make( *listEndPointsRPC, cluster.indexEdgeRequest, applyCh)
	addHandlerRaft(serverRPC, rfRPC)
	addHandlerData(serverRPC, new(Dataformat))
	if configuration.OptionCleaning{
		go cleanThread()
	}

	syscall.Pause()

}










