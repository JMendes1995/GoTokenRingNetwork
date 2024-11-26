package config

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc/peer"
)

var (
	Mutex sync.Mutex

)


func GetPeerMetadata(ctx context.Context) (string, string, error) {
	p, _ := peer.FromContext(ctx)

	addr := p.Addr
	ip, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", "", err
	}
	return ip, port, nil

}
