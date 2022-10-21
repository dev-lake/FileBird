package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
)

// Generate gRPC connection
func GenGRPCConn(addr string, port int) *grpc.ClientConn {
	// Connect server
	address := fmt.Sprintf("%s:%d", addr, port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
		return &grpc.ClientConn{}
	}
	return conn
}
