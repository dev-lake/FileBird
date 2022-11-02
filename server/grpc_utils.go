package main

import (
	"context"
	"log"

	pb "lake.dev/filebird/server/grpc"
)

type UtilsService struct {
	pb.UtilsServer
}

func (s *UtilsService) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingRep, error) {
	log.Println("Trace Ping")
	return &pb.PingRep{}, nil
}
