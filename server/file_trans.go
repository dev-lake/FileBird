package main

import (
	"io"
	"log"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "lake.dev/filebird/server/grpc"
)

type FileTransService struct {
	pb.FileTransServer
}

func (s *FileTransService) Upload(stream pb.FileTrans_UploadServer) error {
	log.Println("Upload")
	// get stream metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "missing context metadata")
	}
	// get file name and open for write
	filePath := md.Get("remote_path")[0]
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot open file:")
	}
	defer file.Close()
	// write file
	for {
		uploadReq, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&pb.UploadRep{
					Ok:  true,
					Msg: "upload success",
				})
			} else {
				return status.Error(codes.Unknown, "cannot receive file data")
			}
		}
		file.Write(uploadReq.Data)
	}
}

func (s *FileTransService) Download(req *pb.DownloadReq, stream pb.FileTrans_DownloadServer) error {
	log.Println("Download")
	// get file name and open for read
	filePath := req.Path
	file, err := os.Open(filePath)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot open file:")
	}
	defer file.Close()
	// read file
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return status.Error(codes.Unknown, err.Error())
			}
		}
		if err := stream.Send(&pb.DownloadRep{
			Data: buf[:n],
		}); err != nil {
			return status.Error(codes.Unknown, "cannot send file data")
		}
	}
}
