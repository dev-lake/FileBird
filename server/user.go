package main

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "lake.dev/filebird/server/grpc"
)

type UserService struct {
	pb.UserServer
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRep, error) {
	log.Println("Trace GetUser")
	usr, err := GetLinuxCurrentUser()
	if err != nil {
		log.Panic(err)
		return nil, status.Error(codes.Internal, "get user info failed")
	}
	return &pb.GetUserRep{
		Uid:      usr.Uid,
		Gid:      usr.Gid,
		Username: usr.Username,
		Name:     usr.Name,
		HomeDir:  usr.HomeDir,
	}, nil
}
