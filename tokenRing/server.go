package tokenring

import (
	"encoding/json"
	"flag"
	tpb "goTokenRingNetwork/tokenRing/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type TokenRingServer struct {
	tpb.UnimplementedTokenRingServer
}

func PeerJsonUnmarshal(data string) Peer {
	var peer Peer
	err := json.Unmarshal([]byte(data), &peer)
	if err != nil {
		log.Fatalf("Error in Unmarshalling token: %s", err)
	}
	return peer
}
func TokenJsonUnmarshal(data string) Token {
	var token Token
	err := json.Unmarshal([]byte(data), &token)
	if err != nil {
		log.Fatalf("Error in Unmarshalling token: %s", err)
	}
	return token
}

func Server() {
	flag.Parse()
	lis, err := net.Listen("tcp", LocalAddr)
	if err != nil {
		log.Fatalf("failed to listen from tokenring server: %s", err)
	}
	log.Printf("Token Ring overlay network server listening at %s", LocalAddr)
	// start grpc server
	server := grpc.NewServer()
	// register server
	tpb.RegisterTokenRingServer(server, &TokenRingServer{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
