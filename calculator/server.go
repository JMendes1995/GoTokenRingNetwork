package calculator

import (
	"context"
	cpb "distributed_p2p_network/calculator/proto"
	"distributed_p2p_network/config"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type CalculatorServer struct {
	cpb.UnimplementedCalculatorServer
}

func jsonUnmarshal(data string) Calc {
	var reqCalc Calc
	err := json.Unmarshal([]byte(data), &reqCalc)
	if err != nil {
		fmt.Println("Error in Unmarshalling calc:", err)
	}
	return reqCalc
}

func GetPeerMetadata(ctx context.Context) (string, string, error) {
  p, _ := peer.FromContext(ctx)

  addr := p.Addr
  ip, port, err := net.SplitHostPort(addr.String())
  if err != nil {
      return "", "", err
  }
  return ip, port, nil

}

func (s *CalculatorServer) Calculate(ctx context.Context, in *cpb.CalculateRequest) (*cpb.CalculateResponse, error) {
	clientIP, _, errGetPeer := config.GetPeerMetadata(ctx)
	if errGetPeer != nil {
		log.Fatalf("Error in GetPeer: %v", errGetPeer)
	}

	reqCalc := jsonUnmarshal(in.GetCalc())
	calcResult, err := CalcOperation(reqCalc)
	if err != nil {
		log.Fatalf("Error in Calculation: %v", err)
	}
	log.Printf("Request from %s: Calculation %s%s%s=%s", clientIP, fmt.Sprintf("%d", reqCalc.Operand1), reqCalc.Operator,
		fmt.Sprintf("%d", reqCalc.Operand2), fmt.Sprintf("%f", calcResult))

	// sending result to client
	return &cpb.CalculateResponse{Result: fmt.Sprintf("%f", calcResult)}, nil
}

func Server() {
	flag.Parse()
	lis, err := net.Listen("tcp", calcServerAddr)
	if err != nil {
		log.Fatalf("failed to listen from calc server: %v", err)
	}
	//start grpc server
	server := grpc.NewServer()
	//register server
	cpb.RegisterCalculatorServer(server, &CalculatorServer{})
	log.Printf("Calculator server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
