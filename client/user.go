package main

import (
	"context"
	"log"
	"os/user"

	pb "lake.dev/filebird/client/grpc"
)

// Get remote user info
func GetRemoteUserInfo(remote ServerInfo) (*user.User, error) {
	// Connect server
	conn := GenGRPCConn(remote.Addr, int(remote.Port))
	defer conn.Close()

	// Establish gRPC connection
	grpcClient := pb.NewUserClient(conn)

	// call gRPC method
	req := pb.GetUserReq{}
	rep, err := grpcClient.GetUser(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	log.Println(rep)
	return &user.User{
		Uid:      rep.Uid,
		Gid:      rep.Gid,
		Username: rep.Username,
		Name:     rep.Name,
		HomeDir:  rep.HomeDir,
	}, err
}
