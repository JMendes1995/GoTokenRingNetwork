package main

import (
	"flag"
	"fmt"
	"goTokenRingNetwork/calculator"
	"goTokenRingNetwork/config"
	tokenring "goTokenRingNetwork/tokenRing"
	"log"
	"os"
	"strconv"
)

func cli() (string, string, string, string) {
	// Define flags
	mode := flag.String("mode", "", " Options: [server, peer]. Specify the execution mode if the node is execution as a calculator server or a regular peer")
	priority := flag.String("priority", "", "Specify the Local Node priority")
	neighbourAddress := flag.String("neighbourAddress", "", "Specify the Neighbour address")
	serverAddress := flag.String("serverAddress", "", "Specify the Calculator server address (ignore this parameter when runing as calculator server mode)")

	help := flag.Bool("help", false, "Display help")

	// Parse flags
	flag.Parse()

	if *help {
		fmt.Println("Usage: cli [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	return *mode, *priority, *neighbourAddress, *serverAddress
}

func main() {
	localAddress, _ := tokenring.GetLocalAddr()

	mode, priority, neighbourAddress, serverAddress := cli()
	tokenring.PeerAddr = neighbourAddress + tokenring.Port
	calculator.ServerAddress = serverAddress + calculator.Port
	tokenring.Uid = priority
	tokenring.LocalAddr = localAddress + tokenring.Port
	intPriority, _ := strconv.Atoi(tokenring.Uid)
	tokenring.Priority = intPriority

	tokenring.LocalPeer = tokenring.Peer{
		Uid:         priority,
		Address:     localAddress,
		Priority:    intPriority,
		Leader:      false,
		TokenHolder: false,
	}
	// main function where is going to select if the node will be running as a calculator server or client
	if mode == "server" {
		go calculator.Server()
	} else if mode == "peer" {
		go tokenring.Server()
		go tokenring.LeaderElectionClient(&config.Mutex)
		go calculator.Queue.EventGenerator(&config.Mutex)
		//go calculator.Queue.Client(&config.Mutex, calcServerAddr)

	} else {
		log.Fatalf("node type requires to be peer or server")
	}
	select {}
}
