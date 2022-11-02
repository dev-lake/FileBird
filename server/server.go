package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "lake.dev/filebird/server/grpc"
)

const (
	// listen on
	Addr string = ":2000"
	// Proto Type
	Proto string = "tcp"
)

func main() {
	fmt.Println("FileBird Server Starting...")
	listener, err := net.Listen(Proto, Addr)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
		return
	}
	server := grpc.NewServer()
	pb.RegisterFileInfoServer(server, &FileInfoService{})
	pb.RegisterFileOperateServer(server, &FileOperateService{})
	pb.RegisterFileTransServer(server, &FileTransService{})
	pb.RegisterUserServer(server, &UserService{})
	pb.RegisterUtilsServer(server, &UtilsService{})

	// Start listen
	fmt.Println("net listening", Addr, "...")
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
