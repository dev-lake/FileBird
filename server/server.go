package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"gopkg.in/ini.v1"
	pb "lake.dev/filebird/server/grpc"
)

var (
	BinPath    string
	BinHome    string
	ConfigPath string
	Config     *ini.File
	// listen on
	Addr string = ""
	Port string = "2000"
	// Proto Type
	Proto string = "tcp"
)

func init() {
	BinPath, _ = os.Executable()
	BinHome = filepath.Dir(BinPath)
	ConfigPath = filepath.Join(BinHome, "data", "config.ini")
	Config = LoadConfig()
	Addr = Config.Section("NORMAL").Key("ADDR").String()
	Port = Config.Section("NORMAL").Key("PORT").String()
}

func main() {
	fmt.Println("FileBird Server Starting...")
	// get listen on
	listener, err := net.Listen(Proto, Addr+":"+Port)
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
	fmt.Println("net listening", Addr+":"+Port, "...")
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
