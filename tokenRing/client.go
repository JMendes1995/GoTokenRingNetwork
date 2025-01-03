package tokenring

import (
	"context"
	"encoding/json"
	"goTokenRingNetwork/calculator"
	tpb "goTokenRingNetwork/tokenRing/proto"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func LeaderElectionClient(mutex *sync.Mutex) {
	conn, err := grpc.NewClient(PeerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to server %s", err)
	}

	c := tpb.NewTokenRingClient(conn)
	defer conn.Close()

	peerJsonFormat, _ := json.Marshal(LocalPeer)
	log.Printf("Challeging leadership with peer %s\n", PeerAddr)
	reqErr := true
	for reqErr {
		// setup connection with timeout
		// if the peer is not ready try to establish connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		isLeader, err := c.LeaderElection(ctx, &tpb.LeaderElectionRequest{Leader: string(peerJsonFormat)})
		if err == nil {
			if isLeader.Leader == "true" {
				mutex.Lock()
				LocalPeer.Leader = true
				LocalPeer.TokenHolder = true
				mutex.Unlock()
				log.Printf("Peer %s is the Leader with priority of %d\n", LocalAddr, LocalPeer.Priority)
				TokenGenerator(string(peerJsonFormat), mutex)
				log.Printf("Generation New TOKEN: %s", SharedToken.Value)
				time.Sleep(time.Second * 10)
				CheckTokenHolder(mutex)
			} else {
				log.Printf("Peer %s is not the Leader with priority of %d\n", LocalAddr, LocalPeer.Priority)
			}
			reqErr = false
		} else {
			log.Printf("Peer not ready %s", PeerAddr)
		}
		time.Sleep(time.Second * 5)
	}
}
func CheckTokenHolder(mutex *sync.Mutex) {
	time.Sleep(time.Second * 10)
	log.Printf("Local Peer holding the token: %t\n", LocalPeer.TokenHolder)
	if LocalPeer.TokenHolder {
		calculator.Queue.Client(mutex)
		mutex.Lock()
		LocalPeer.TokenHolder = false
		mutex.Unlock()
		TokenTransitClient(mutex)
		log.Printf("Local Peer holding the tokan: %t\n", LocalPeer.TokenHolder)
	}

}
func TokenTransitClient(mutex *sync.Mutex) {
	conn, err := grpc.NewClient(PeerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to server %s", err)
	}
	defer conn.Close()

	c := tpb.NewTokenRingClient(conn)

	// setup connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tokenJsonFormat, _ := json.Marshal(SharedToken)

	log.Printf("Sending token to peer: %s \n", PeerAddr)
	response, errToken := c.TokenTransit(ctx, &tpb.TokenRequest{Token: string(tokenJsonFormat)})

	if errToken != nil {
		log.Printf("Error sending token to peer %s", errToken)
	} else {
		log.Print(response.Token)
	}

}
