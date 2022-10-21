package main

import (
	"context"
	"log"

	pb "lake.dev/filebird/client/grpc"
)

func CopyRemoteFile(remote ServerInfo, src string, dst string) {
	// Connect server
	conn := GenGRPCConn(remote.Addr, int(remote.Port))
	defer conn.Close()

	// Establish gRPC connection
	grpcClient := pb.NewFileOperateClient(conn)

	// call gRPC method
	req := pb.CopyFileReq{
		Src: src,
		Dst: dst,
	}
	rep, err := grpcClient.CopyFile(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	log.Println(rep)

}

// move remote file
func MoveRemoteFile(remote ServerInfo, src string, dst string) {
	// Connect server
	conn := GenGRPCConn(remote.Addr, int(remote.Port))
	defer conn.Close()

	// Establish gRPC connection
	grpcClient := pb.NewFileOperateClient(conn)

	// call gRPC method
	req := pb.MoveFileReq{
		Src: src,
		Dst: dst,
	}
	rep, err := grpcClient.MoveFile(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	log.Println(rep)

}

// delete remote file
func DeleteRemoteFile(remote ServerInfo, path string) {
	// Connect server
	conn := GenGRPCConn(remote.Addr, int(remote.Port))
	defer conn.Close()

	// Establish gRPC connection
	grpcClient := pb.NewFileOperateClient(conn)

	// call gRPC method
	req := pb.DeleteFileReq{
		Path: path,
	}
	rep, err := grpcClient.DeleteFile(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	log.Println(rep)

}
