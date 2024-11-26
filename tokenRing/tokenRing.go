package tokenring

import (
	"context"
	"distributed_p2p_network/calculator"
	"distributed_p2p_network/config"
	tpb "distributed_p2p_network/tokenRing/proto"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"golang.org/x/crypto/sha3"
)

type TokenRing struct {
	Neighbour Peer
	Token     Token
}
type Peer struct {
	Uid      string `json:"Uid"`
	Address  string `json:"Address"`
	Priority int    `json:"Priority"`
	Leader   bool   `json:"Leader"`
	TokenHolder bool `json:"TokenHolder"`
}

type Token struct {
	Value          string `json:"Value"`
	PreviousHolder string   `json:"PreviousHolder"`
}

var (
	LocalPeer Peer
	SharedToken Token
	PeerAddr = os.Getenv("NGH_ADDR")+":50002"
	LocalAddr = os.Getenv("LOCAL_NODE")+":50002" 
	Priority, _ = strconv.Atoi(os.Getenv("PRIORITY"))
	Uid = os.Getenv("PRIORITY")
)
func TokenGenerator(peerString string, mutex *sync.Mutex) Token{
	// send
	// Create a SHA3-256 hash
	hash := sha3.Sum256([]byte(peerString))
	log.Printf("SHA3-256 hash: %x\n", hash)
	mutex.Lock()
	SharedToken.Value = hex.EncodeToString(hash[:])
	SharedToken.PreviousHolder = LocalAddr
	mutex.Unlock()
	return SharedToken

}

func (s *TokenRingServer) LeaderElection(ctx context.Context, in *tpb.LeaderElectionRequest) (*tpb.LeaderElectionResult, error) {
	getPeerData := PeerJsonUnmarshal(in.GetLeader())
	clientIP, _, errGetPeer := calculator.GetPeerMetadata(ctx)
	if errGetPeer != nil {
		log.Fatalf("Error in GetPeer: %v", errGetPeer)
	}

	log.Printf("Leader election Request comming from peer %s with priority: %d\n", clientIP, getPeerData.Priority)

	if LocalPeer.Priority > getPeerData.Priority {
		config.Mutex.Lock()
		LocalPeer.Leader = true
		config.Mutex.Unlock()
		log.Printf("Leader election Result Leader: %s with priority: %d\n", LocalPeer.Address, LocalPeer.Priority)
		return &tpb.LeaderElectionResult{Leader: fmt.Sprintf("%t", false)}, nil
	} else {
		log.Printf("Leader election Result Leader: %s with priority: %d\n", getPeerData.Address, getPeerData.Priority)
		return &tpb.LeaderElectionResult{Leader: fmt.Sprintf("%t", true)}, nil
	}
}

func (s *TokenRingServer) TokenTransit(ctx context.Context, in *tpb.TokenRequest) (*tpb.TokenResponse, error) {

	getTokenData := TokenJsonUnmarshal(in.GetToken())	
	clientIP, _, errGetPeer := calculator.GetPeerMetadata(ctx)
	if errGetPeer != nil {
		log.Fatalf("Error in GetPeer: %v", errGetPeer)
	}

	SharedToken.Value = getTokenData.Value
	SharedToken.PreviousHolder = LocalPeer.Address

	log.Printf("Recieved token from peer %s with value: %s\n", clientIP, getTokenData.Value)
	log.Printf("Token Previous holder: %s\n", getTokenData.PreviousHolder)
	config.Mutex.Lock()
	LocalPeer.TokenHolder = true
	config.Mutex.Unlock()
	go CheckTokenHolder(&config.Mutex)

	return &tpb.TokenResponse{Token: fmt.Sprintf("Token recieved from peer %s", LocalPeer.Address)}, nil
}
