package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc/metadata"
	pb "lake.dev/filebird/client/grpc"
)

func UploadFile(remote ServerInfo, local_path string, remote_path string) {
	log.Println("trace UploadFile", local_path, "to", remote_path)
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

	// init progress
	progress := progressbar.DefaultBytes(local_file_info.Size(), "uploading")

	// read file and send to server
	buf := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("upload file success")
				break
			} else {
				log.Fatalf("read file err: %v", err)
			}
			break
		}
		err = stream.Send(&pb.UploadReq{Data: buf[:n]})
		if err != nil {
			if err == io.EOF {
				log.Println("upload file success")
				break
			} else {
				log.Panicln("upload file", file.Name(), "failed, err:", err)
			}
		}
		progress.Add(n) // show progress
	}
}

// upload file recursively
func UploadDir(remote ServerInfo, local_path string, remote_path string) error {
	log.Println("trace UploadDir", local_path, "to", remote_path)
	// judge file exist
	if !IsLocalFileExist(local_path) {
		log.Panic("file not exist")
		return fmt.Errorf("file not exist")
	}
	// judge file is dir
	if !LocalFileIsDir(local_path) {
		UploadFile(remote, local_path, remote_path)
	} else {
		// get local dir files
		files := ReadLocalDirAll(local_path)
		for _, file := range files {
			fmt.Println("---", file.Path)
		}
		// upload files
		for _, file := range files {
			if file.Name == ".DS_Store" { // skip .DS_Store. macos system file
				continue
			}
			// get remote path
			relative_path := file.Path[len(local_path):]
			remote_file_path := remote_path + relative_path
			if file.IsDir {
				fmt.Println("mkdir", file.Name)
				MakeRemoteDir(remote, remote_file_path)
			} else {
				fmt.Println("upload", file.Name)
				UploadFile(remote, file.Path, remote_file_path)
			}
			// UploadDir(remote, local_path+"/"+file.Name(), remote_path+"/"+file.Name())
		}
	}
	return nil
}

func DownloadFile(remote ServerInfo, local_path string, remote_path string) {
	fmt.Println("trace DownloadFile", remote_path, "to", local_path)
	fmt.Println(remote.Addr, remote.Port)
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

	// get remote file meta info
	infoClient := pb.NewFileInfoClient(conn)
	info_res, err := infoClient.GetFileInfo(
		context.Background(), &pb.FileReq{FilePath: remote_path},
	)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}

	// init progress
	progress := progressbar.DefaultBytes(int64(info_res.Size), "downloading")

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
		// show progress
		progress.Add(len(rep.Data))
	}
}

func DownloadDir(remote ServerInfo, local_path string, remote_path string) error {
	log.Println("trace DownloadDir", remote_path, "to", local_path)
	// get remote dir all files
	files, err := GetRemoteDirAllFiles(&remote, remote_path)
	if err != nil {
		return err
	}
	// download files
	for _, file := range files {
		// get local path
		relative_path := file.Path[len(remote_path):]
		local_file_path := local_path + relative_path
		if file.IsDir {
			fmt.Println("mkdir", file.Name)
			ok := MakeLocalDir(local_file_path)
			if !ok {
				return fmt.Errorf("mkdir %s failed", local_file_path)
			}
		} else {
			fmt.Println("download", file.Name)
			DownloadFile(remote, local_file_path, file.Path)
		}
		// DownloadDir(remote, local_path+"/"+file.Name(), remote_path+"/"+file.Name())
	}
	return nil
}

// transmit file from remote to remote
func TransmitFile(src ServerInfo, dst ServerInfo, src_path string, dst_path string) {
	log.Println("trace transmit file from", src.Addr, "to", dst.Addr)

	// get remote file meta info
	infoClient := pb.NewFileInfoClient(GenGRPCConn(src.Addr, int(src.Port)))
	info_res, err := infoClient.GetFileInfo(
		context.Background(), &pb.FileReq{FilePath: src_path},
	)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}

	// init progress
	progress := progressbar.DefaultBytes(int64(info_res.Size), "transmitting")

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
				log.Println("Transmit success")
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
		// show progress
		progress.Add(len(rep.Data))
	}
}
