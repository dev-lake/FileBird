package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "lake.dev/filebird/client/grpc"
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

// check gRPC server status
func CheckGRPCServerStatus(addr string, port int) bool {
	// Connect server
	address := fmt.Sprintf("%s:%d", addr, port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return false
	}
	defer conn.Close()
	grpcClient := pb.NewUtilsClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = grpcClient.Ping(ctx, &pb.PingReq{})
	if err != nil {
		return false
	}
	return true
}
