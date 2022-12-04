package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// Connection create grpc client connection
func Connection(target string) *grpc.ClientConn {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	c, err := grpc.Dial(target, opts)
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	return c
}
