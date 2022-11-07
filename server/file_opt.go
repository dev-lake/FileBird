package main

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "lake.dev/filebird/server/grpc"
)

type FileOperateService struct {
	pb.FileOperateServer
}

func (s *FileOperateService) CopyFile(ctx context.Context, req *pb.CopyFileReq) (*pb.FileOptRep, error) {
	log.Println("CopyFile")
	if !CopyLocalFileRecursively(req.Src, req.Dst) {
		return nil, status.Error(codes.Unknown, "copy file failed")
	}
	return &pb.FileOptRep{
		Success: true,
	}, nil
}

func (s *FileOperateService) MoveFile(ctx context.Context, req *pb.MoveFileReq) (*pb.FileOptRep, error) {
	log.Println("MoveFile")
	if !MoveFile(req.Src, req.Dst) {
		return nil, status.Error(codes.Unknown, "move file failed")
	}
	return &pb.FileOptRep{
		Success: true,
	}, nil
}

func (s *FileOperateService) DeleteFile(ctx context.Context, req *pb.DeleteFileReq) (*pb.FileOptRep, error) {
	log.Println("DeleteFile")
	if !DeleteFile(req.Path) {
		return nil, status.Error(codes.Unknown, "delete file failed")
	}
	return &pb.FileOptRep{
		Success: true,
	}, nil
}

// mkdir
func (s *FileOperateService) MakeDir(ctx context.Context, req *pb.MakeDirReq) (*pb.FileOptRep, error) {
	log.Println("MakeDir")
	if !MakeLocalDir(req.Path) {
		return nil, status.Error(codes.Unknown, "make dir failed")
	}
	return &pb.FileOptRep{
		Success: true,
	}, nil
}
