package main

import (
	"distributed_p2p_network/calculator"
	"distributed_p2p_network/config"
	tokenring "distributed_p2p_network/tokenRing"
	"log"
	"os"
)

func main() {
	tokenring.LocalPeer = tokenring.Peer{
		Uid:    tokenring.Uid,
		Address:  tokenring.LocalAddr,
		Priority: tokenring.Priority,
		Leader:   false,
		TokenHolder: false,
	}
	// main function where is going to select if the node will be running as a calculator server or client
	if os.Args[1] == "server" {
		go calculator.Server()		
		//go tokenring.Server()
	} else if os.Args[1] == "peer" {
		go calculator.Queue.EventGenerator(&config.Mutex)
		go tokenring.Server()
		go tokenring.LeaderElectionClient(&config.Mutex)
		//go calculator.Queue.Client(&config.Mutex, calcServerAddr)

	} else {
		log.Fatalf("node type requires to be peer or server")
	}
	select {}
}
