package main

import (
	"context"
	"io/fs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "lake.dev/filebird/server/grpc"
)

type FileInfoService struct {
	pb.FileInfoServer
}

func (s *FileInfoService) GetFileInfo(ctx context.Context, req *pb.FileReq) (*pb.RegularFileInfo, error) {
	// log.Println("========", req.FilePath)
	exists, _ := PathExists(req.FilePath)
	if !exists {
		return nil, status.Error(codes.NotFound, "File Not Exists")
	}
	fStat := GetFileInfo(req.FilePath)
	// get file user id and group id
	username, groupname := GetFileUserAndGroupName(fStat)
	reply := pb.RegularFileInfo{
		Path:      req.FilePath,
		Name:      fStat.Name(),
		IsDir:     fStat.IsDir(),
		Size:      uint64(fStat.Size()),
		UserName:  username,
		GroupName: groupname,
		// CreateTime:   "",
		Mode:       uint32(fStat.Mode()),
		ModifyTime: fStat.ModTime().String(),
		// LastOpenTime: "",
	}
	return &reply, nil
}

func (s *FileInfoService) GetDirFileInfo(ctx context.Context, req *pb.FileReq) (*pb.DirFileInfoList, error) {
	var fileInfoList []*pb.RegularFileInfo
	var info_list []fs.FileInfo

	// judge file exists and is not dir
	exists, _ := PathExists(req.FilePath)
	if !exists {
		return nil, status.Error(codes.NotFound, "File Not Exists")
	}
	isDir := IsDir(req.FilePath)
	if !isDir {
		info := GetFileInfo(req.FilePath)
		info_list = append(info_list, info)
	} else {
		info_list = ReadDir(req.FilePath)
	}
	// fill reply
	for _, info := range info_list {
		println(info.Name())
		username, groupname := GetFileUserAndGroupName(info)
		fileInfo := pb.RegularFileInfo{
			Path:         req.FilePath,
			Name:         info.Name(),
			IsDir:        info.IsDir(),
			Size:         uint64(info.Size()),
			UserName:     username,
			GroupName:    groupname,
			Mode:         uint32(info.Mode()),
			ModifyTime:   info.ModTime().String(),
			LastOpenTime: "",
		}
		fileInfoList = append(fileInfoList, &fileInfo)
	}
	// return
	return &pb.DirFileInfoList{
		Path:         req.FilePath,
		FileInfoList: fileInfoList,
	}, nil
}

// get dir all files
func (s *FileInfoService) GetDirAllFiles(ctx context.Context, req *pb.FileReq) (*pb.DirFileInfoList, error) {
	var fileInfoList []*pb.RegularFileInfo
	var info_list []FileInfo

	// judge file exists and is not dir
	exists, _ := PathExists(req.FilePath)
	if !exists {
		return nil, status.Error(codes.NotFound, "File Not Exists")
	}
	isDir := IsDir(req.FilePath)
	if !isDir {
		info := GetFileInfo(req.FilePath)
		username, groupname := GetFileUserAndGroupName(info)
		info_list = append(info_list, FileInfo{
			Path:      req.FilePath,
			Name:      info.Name(),
			Size:      info.Size(),
			IsDir:     info.IsDir(),
			Mode:      info.Mode(),
			ModTime:   info.ModTime().String(),
			UserName:  username,
			GroupName: groupname,
		})
	} else {
		info_list = ReadLocalDirAll(req.FilePath)
	}
	// fill reply
	for _, info := range info_list {
		username, groupname := GetFileUserAndGroupName(GetFileInfo(info.Path))
		fileInfo := pb.RegularFileInfo{
			Path:         info.Path,
			Name:         info.Name,
			IsDir:        info.IsDir,
			Size:         uint64(info.Size),
			UserName:     username,
			GroupName:    groupname,
			Mode:         uint32(info.Mode),
			ModifyTime:   info.ModTime,
			LastOpenTime: "",
		}
		fileInfoList = append(fileInfoList, &fileInfo)
	}
	// return
	return &pb.DirFileInfoList{
		Path:         req.FilePath,
		FileInfoList: fileInfoList,
	}, nil
}
