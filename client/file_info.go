package main

import (
	"context"

	pb "lake.dev/filebird/client/grpc"
)

func IsRemoteFileExist(remote *ServerInfo, path_abs string) bool {
	_, err := GetFileInfo(remote.Name + ":" + path_abs)
	if err != nil {
		return false
	}
	return true
}

func RemoteFileIsDir(remote *ServerInfo, path_abs string) bool {
	file_info, err := GetFileInfo(remote.Name + ":" + path_abs)
	if err != nil {
		return false
	}
	if file_info.IsDir {
		return true
	}
	return false
}

func GetRemoteDirAllFiles(remote *ServerInfo, path_abs string) ([]*FileInfo, error) {
	// connect server
	conn := GenGRPCConn(remote.Addr, int(remote.Port))
	defer conn.Close()
	// establish gRPC connection
	grpcClient := pb.NewFileInfoClient(conn)
	// call gRPC method
	req := pb.FileReq{
		FilePath: path_abs,
	}
	rep, err := grpcClient.GetDirAllFiles(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	// convert rep.FileInfoList to []*FileInfo
	var file_info_list []*FileInfo
	for _, file_info := range rep.FileInfoList {
		file_info_list = append(file_info_list, &FileInfo{
			Path: file_info.Path,
			Name: file_info.Name,
			Size: int64(file_info.Size),
			// Mode:      file_info.Mode.FileMode(),
			ModTime:   file_info.ModifyTime,
			IsDir:     file_info.IsDir,
			UserName:  file_info.UserName,
			GroupName: file_info.GroupName,
		})
	}
	return file_info_list, nil
}
