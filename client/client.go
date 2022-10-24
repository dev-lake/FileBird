package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	pb "lake.dev/filebird/client/grpc"
)

const (
	Version string = "v0.1.0"
)

type FileInfo struct {
	Path      string
	Name      string      // base name of the file
	Size      int64       // length in bytes for regular files; system-dependent for others
	Mode      fs.FileMode // file mode bits
	ModTime   string      // modification time
	IsDir     bool        // abbreviation for Mode().IsDir()
	UserName  string
	GroupName string
	// Sys     fs.any      // underlying data source (can return nil)
}

var (
	app = kingpin.New("filebird", "FileBird Client")
	// add server
	add_server      = app.Command("add_server", "Add a server to the list of servers")
	add_server_name = add_server.Flag("name", "Server name").Short('n').Required().String()
	add_server_addr = add_server.Flag("addr", "Server address").Short('a').Required().String()
	add_server_port = add_server.Flag("port", "Server port").Short('p').Default("2000").Int()
	// add_server_pass = add_server.Flag("pass", "Server Password").String()
	// add_server_desc = add_server.Flag("desc", "Server describe").String()
	// show server
	show_server      = app.Command("show_server", "Show the list of servers")
	show_server_name = show_server.Flag("name", "Server name").String()
	// delete server
	del_server      = app.Command("del_server", "Delete a server from the list of servers")
	del_server_name = del_server.Arg("name", "Server name").Required().String()
	// file info
	reset = app.Command("reset", "Reset the list of servers")
	// get file info
	info      = app.Command("info", "Get file info")
	dest_file = info.Arg("file", "File path").Required().String()
	// ls dir
	ls      = app.Command("ls", "Show files in dir")
	ls_file = ls.Arg("file", "File path").Required().String()
	// pwd
	pwd             = app.Command("pwd", "Show current dir")
	pwd_server_name = pwd.Arg("name", "Server name").Default("").String()
	// cd dir
	cd             = app.Command("cd", "Change current dir")
	cd_server_path = cd.Arg("path", "Server and Path").Required().String()
	// copy file
	cp          = app.Command("cp", "Copy file")
	cp_src_file = cp.Arg("src_file", "Source file path").Required().String()
	cp_dst_file = cp.Arg("dst_file", "Destination file path").Required().String()
	// move file
	mv          = app.Command("mv", "Move file")
	mv_src_file = mv.Arg("src_file", "Source file path").Required().String()
	mv_dst_file = mv.Arg("dst_file", "Destination file path").Required().String()
	// delete file
	rm      = app.Command("rm", "Remove file")
	rm_file = rm.Arg("file", "File path").Required().String()
)

func main() {
	// parse command line
	app.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Register user
	case add_server.FullCommand():
		println(*add_server_name, *add_server_addr, *add_server_port)
		// add server
		AddServer(InitDB(), *add_server_name, *add_server_addr, *add_server_port)
	case show_server.FullCommand():
		println(*show_server_name)
		// show server
		servers := ShowServer(InitDB())
		ShowServerTable(servers, *show_server_name)
	case del_server.FullCommand():
		println(*del_server_name)
		// delete server
		DeleteServer(InitDB(), *del_server_name)
	case reset.FullCommand():
		// reset server
		os.Remove(db_path)
	case info.FullCommand():
		// get file info
		file_info := GetFileInfo(*dest_file)
		ShowFileInfoTable(file_info)
	case ls.FullCommand():
		// ls dir or file
		info_list := GetDirFileInfoList(*ls_file)
		ShowFileInfoListTable(info_list)
	case pwd.FullCommand():
		// show current dir
		if *pwd_server_name == "localhost" || *pwd_server_name == "" {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(dir)
		} else {
			fmt.Println(GetServerPwd(InitDB(), *pwd_server_name))
		}
	case cp.FullCommand():
		println(*cp_src_file, *cp_dst_file)
		CopyFile(*cp_src_file, *cp_dst_file)
		// CopyLocalFile(*cp_src_file, *cp_dst_file)
	case mv.FullCommand():
		println(*mv_src_file, *mv_dst_file)
		MoveFile(*mv_src_file, *mv_dst_file)
	case rm.FullCommand():
		println(*rm_file)
		DeleteFile(*rm_file)
	}

}

// get file info from remote or local
func GetFileInfo(file_path string) (fileInfo *FileInfo) {
	// get server info
	server_name, path := ExtractServerNameAndFilePath(file_path)
	if server_name == "localhost" {
		file_info := GetLocalFileInfo(path)
		uname, gname := GetFileUserAndGroupName(file_info)
		return &FileInfo{
			Path:      file_path,
			Name:      file_info.Name(),
			Size:      file_info.Size(),
			ModTime:   file_info.ModTime().String(),
			IsDir:     file_info.IsDir(),
			Mode:      file_info.Mode(),
			UserName:  uname,
			GroupName: gname,
		}
	} else {
		server_info := GetServer(InitDB(), server_name)
		fmt.Println(path, server_info.Addr, server_info.Port)

		// Connect server
		conn := GenGRPCConn(server_info.Addr, int(server_info.Port))
		defer conn.Close()

		// Establish gRPC connection
		grpcClient := pb.NewFileInfoClient(conn)

		// Call gRPC method
		req := pb.FileReq{
			FilePath: path,
		}
		res, err := grpcClient.GetFileInfo(context.Background(), &req)
		if err != nil {
			log.Fatalf("Call Route err: %v", err)
		}
		// log.Println(res)
		return &FileInfo{
			Path:      res.Path,
			Name:      res.Name,
			Size:      int64(res.Size),
			ModTime:   res.ModifyTime,
			IsDir:     res.IsDir,
			Mode:      fs.FileMode(res.Mode),
			UserName:  res.UserName,
			GroupName: res.GroupName,
		}
	}
}

// get dir file info list from remote or local
func GetDirFileInfoList(dir_path string) (fileInfoList []*FileInfo) {
	// get server info
	server_name, path := ExtractServerNameAndFilePath(dir_path)
	if server_name == "localhost" {
		file_info_list := ReadLocalDir(path)
		for _, file_info := range file_info_list {
			uname, gname := GetFileUserAndGroupName(file_info)
			fileInfoList = append(fileInfoList, &FileInfo{
				Path:      dir_path + "/" + file_info.Name(),
				Name:      file_info.Name(),
				Size:      file_info.Size(),
				ModTime:   file_info.ModTime().String(),
				IsDir:     file_info.IsDir(),
				Mode:      file_info.Mode(),
				UserName:  uname,
				GroupName: gname,
			})
		}
		return fileInfoList
	} else {
		server_info := GetServer(InitDB(), server_name)
		fmt.Println(path, server_info.Addr, server_info.Port)

		// Connect server
		conn := GenGRPCConn(server_info.Addr, int(server_info.Port))
		defer conn.Close()

		// Establish gRPC connection
		grpcClient := pb.NewFileInfoClient(conn)

		// Call gRPC method
		req := pb.FileReq{
			FilePath: path,
		}
		res, err := grpcClient.GetDirFileInfo(context.Background(), &req)
		if err != nil {
			log.Fatalf("Call Route err: %v", err)
		}
		// log.Println(res)
		for _, file_info := range res.FileInfoList {
			fileInfoList = append(fileInfoList, &FileInfo{
				Path:      file_info.Path,
				Name:      file_info.Name,
				Size:      int64(file_info.Size),
				ModTime:   file_info.ModifyTime,
				IsDir:     file_info.IsDir,
				Mode:      fs.FileMode(file_info.Mode),
				UserName:  file_info.UserName,
				GroupName: file_info.GroupName,
			})
		}
		return fileInfoList
	}
}

// copy file
func CopyFile(src string, dst string) (success bool) {
	// get server info
	src_server_name, src_path := ExtractServerNameAndFilePath(src)
	dst_server_name, dst_path := ExtractServerNameAndFilePath(dst)
	// get server
	src_server := GetServer(InitDB(), src_server_name)
	dst_server := GetServer(InitDB(), dst_server_name)

	// 同一主机
	if src_server.Addr == src_server.Addr {
		// 本地复制
		if src_server.Addr == "localhost" {
			// copy local file
			CopyLocalFile(src_path, dst_path)
			return true
		} else {
			// copy remote file
			CopyRemoteFile(src_server, src_path, dst_path)
			return true
		}
	} else if src_server.Addr == "localhost" || dst_server.Addr == "localhost" { // 一个本地一个远程
		if src_server.Addr == "localhost" {
			// upload local file to remote
			UploadFile(dst_server, src_path, dst_path)
			return true
		} else {
			// download remote file to local
			DownloadFile(src_server, src_path, dst_path)
			return true
		}
	} else { // 两个远程
		TransmitFile(src_server, dst_server, src_path, dst_path)
	}
	return false
}

// move file
func MoveFile(src string, dst string) (success bool) {
	// get server info
	src_server_name, src_path := ExtractServerNameAndFilePath(src)
	dst_server_name, dst_path := ExtractServerNameAndFilePath(dst)
	// get server
	src_server := GetServer(InitDB(), src_server_name)
	dst_server := GetServer(InitDB(), dst_server_name)

	// 同一主机
	if src_server.Addr == src_server.Addr {
		// 本地复制
		if src_server.Addr == "localhost" {
			// move local file
			MoveLocalFile(src_path, dst_path)
			return true
		} else {
			// move remote file
			MoveRemoteFile(src_server, src_path, dst_path)
			return true
		}
	} else if src_server.Addr == "localhost" || dst_server.Addr == "localhost" { // 一个本地一个远程
		if src_server.Addr == "localhost" {
			// upload local file to remote
			UploadFile(dst_server, src_path, dst_path)
			// delete local file
			DeleteLocalFile(src_path)
			return true
		} else {
			// download remote file to local
			DownloadFile(src_server, src_path, dst_path)
			// delete remote file
			DeleteRemoteFile(dst_server, dst_path)
			return true
		}
	} else { // 两个远程
		log.Println("Files will not be deleted for security reasons.")
		TransmitFile(src_server, dst_server, src_path, dst_path)
	}
	return false
}

// delete file
func DeleteFile(file_path string) (success bool) {
	// get server info
	server_name, path := ExtractServerNameAndFilePath(file_path)
	// get server
	server := GetServer(InitDB(), server_name)

	// 本地删除
	if server.Addr == "localhost" {
		// delete local file
		DeleteLocalFile(path)
		return true
	} else {
		// delete remote file
		DeleteRemoteFile(server, path)
		return true
	}
}
