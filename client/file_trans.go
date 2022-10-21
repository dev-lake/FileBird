package main

import (
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/grpc/metadata"
	pb "lake.dev/filebird/client/grpc"
)

func UploadFile(remote ServerInfo, local_path string, remote_path string) {
	// Connect to server
	conn := GenGRPCConn(remote.Addr, int(remote.Port))
	defer conn.Close()

	// Establish gRPC connection
	grpcClient := pb.NewFileTransClient(conn)

	// get local file meta info
	local_file_info := GetLocalFileInfo(local_path)
	file_meta := metadata.New(map[string]string{
		"name":        local_file_info.Name(),
		"origin_path": local_path,
		"remote_path": remote_path,
		// "size":        string(local_file_info.Size()),
		// "mode":        string(local_file_info.Mode()),
		"modify_time": local_file_info.ModTime().String(),
		"md5":         "",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), file_meta)

	// open file for read
	file, err := os.Open(local_path)
	if err != nil {
		log.Fatalf("open file err: %v", err)
	}
	defer file.Close()

	// Call gRPC method
	stream, err := grpcClient.Upload(ctx)
	if err != nil {
		log.Fatalf("call upload err: %v", err)
	}

	// read file and send to server
	buf := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				rep, err := stream.CloseAndRecv()
				if err != nil {
					log.Fatalf("close stream err: %v", err)
				}
				log.Println(rep.Msg)
			} else {
				log.Fatalf("read file err: %v", err)
			}
			break
		}
		err = stream.Send(&pb.UploadReq{Data: buf[:n]})
		if err != nil {
			log.Fatalf("send file err: %v", err)
		}
	}
}

func DownloadFile(remote ServerInfo, local_path string, remote_path string) {
	// Connect to server
	conn := GenGRPCConn(remote.Addr, int(remote.Port))
	defer conn.Close()

	// Establish gRPC connection
	grpcClient := pb.NewFileTransClient(conn)

	// Call gRPC method
	stream, err := grpcClient.Download(
		context.Background(),
		&pb.DownloadReq{
			Path: remote_path,
		},
	)
	if err != nil {
		log.Fatalf("call download err: %v", err)
	}

	// open file for write
	file, err := os.OpenFile(local_path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("open file err: %v", err)
	}
	defer file.Close()

	// receive file from server
	for {
		rep, err := stream.Recv()
		if err == io.EOF {
			log.Println("download file success")
			stream.CloseSend()
			break
		}
		if err != nil {
			log.Fatalf("receive file err: %v", err)
		}
		file.Write(rep.Data)
	}
}

// transmit file from remote to remote
func TransmitFile(src ServerInfo, dst ServerInfo, src_path string, dst_path string) {
	log.Println("transmit file from", src.Addr, "to", dst.Addr)
	// Connect to server
	src_conn := GenGRPCConn(src.Addr, int(src.Port))
	dst_conn := GenGRPCConn(dst.Addr, int(dst.Port))
	defer src_conn.Close()
	defer dst_conn.Close()

	// Establish gRPC connection
	src_client := pb.NewFileTransClient(src_conn)
	dst_client := pb.NewFileTransClient(dst_conn)

	// 下载流
	src_stream, err := src_client.Download(
		context.Background(),
		&pb.DownloadReq{
			Path: src_path,
		},
	)
	if err != nil {
		log.Fatalf("call download err: %v", err)
	}

	// 上传流
	file_meta := metadata.New(map[string]string{
		"name":        "",
		"origin_path": src_path,
		"remote_path": dst_path,
		// "size":        string(local_file_info.Size()),
		// "mode":        string(local_file_info.Mode()),
		"modify_time": "",
		"md5":         "",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), file_meta)
	dst_stream, err := dst_client.Upload(ctx)
	if err != nil {
		log.Fatalf("call upload err: %v", err)
	}

	for {
		rep, err := src_stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("download file success")
				src_stream.CloseSend()
				break
			} else {
				log.Fatalf("receive file err: %v", err)
			}
		}
		err = dst_stream.Send(&pb.UploadReq{Data: rep.Data})
		if err != nil {
			log.Fatalf("send file err: %v", err)
		}
	}
}
