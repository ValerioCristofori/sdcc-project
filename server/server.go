package main

import (
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

func (c *Cluster) toString() string{
	return fmt.Sprintf("Cluster: %s, My position: %d", c.Nodes, c.indexEdgeRequest)
}

// data conf
var (
	cluster 			= new(Cluster)
	listEndPointsRPC 	= new([]*rpc.Client)
)


// interface to RaftRPC
var rfRPC *RaftRPC
// channel for newly committed messages
var applyCh chan ApplyMsg



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

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		fmt.Println("Saving persist log entries and Raft state...")
		if err := Save("./vol/backup", datastore); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Exiting...")
		os.Exit(2)

	}()

	go func() {
		for range time.Tick(5 * time.Second){
			fmt.Println("Saving persist log entries and Raft state...")
			if err := Save("./vol/backup", datastore); err != nil {
				log.Fatalln(err)
			}
		}
	}()
}

func main()  {

	// start configuration and initialization of raft cluster
	register()
	getListEdgeNodes()

	serverRPC := rpc.NewServer()
	go startListener(serverRPC)
	connectToAllNodes()

	err := InitMap()
	if err != nil {
		log.Fatal("Error in Init Map: ", err)
	}
	// wait for table creation
	//for {
	//	tables := callTable()
	//	if tables == 0 {
	//		print("Wait")
	//	} else {
	//		fmt.Println("Created dynamoDB table!")
	//		break
	//	}
	//}
	err = Load("./vol/backup", &datastore )
	if err != nil {
		log.Println("Not able to backup persistent state")
	}
	//PrintMap()
	shutdownHandler()
	// listen to messages from Raft indicating newly committed messages.
	applyCh = make(chan ApplyMsg)
	go applyChRoutine()

	rfRPC = Make( *listEndPointsRPC, cluster.indexEdgeRequest, applyCh)
	addHandlerRaft(serverRPC, rfRPC)
	addHandlerData(serverRPC, new(Dataformat))

	//time.Sleep(8 * time.Second)
	//if rfRPC.rf.state == LEADER {
	//	log.Println("EXITING FROM LEADER")
	//	os.Exit(1)
	//}
	//if cluster.indexEdgeRequest == 1 {
	//		log.Println("EXITING")
	//		os.Exit(1)
	//}
	syscall.Pause()

}










