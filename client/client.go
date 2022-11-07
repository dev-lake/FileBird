package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
	pb "lake.dev/filebird/client/grpc"
)

const (
	Version             string = "v0.1.0"
	ServerNameMaxLength int    = 20
)

var (
	IllegalServerNames = [...]string{"localhost", "local"}
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
	ls_file = ls.Arg("file", "File path").Default("").String()
	// pwd
	pwd             = app.Command("pwd", "Show current dir")
	pwd_server_name = pwd.Arg("name", "Server name").Default("").String()
	// cd dir
	cd             = app.Command("cd", "Change current dir")
	cd_server_path = cd.Arg("path", "Server and Path").Required().String() // format: server_name:path
	// mkdir
	mkdir      = app.Command("mkdir", "Create dir")
	mkdir_path = mkdir.Arg("path", "Server and Path").Required().String() // format: server_name:path
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
		file_info, _ := GetFileInfo(*dest_file)
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
			// Extract server name and path
			var server_name string = *pwd_server_name
			if strings.HasSuffix(server_name, ":") {
				server_name = (*pwd_server_name)[:len(*pwd_server_name)-1]
			}
			fmt.Println(GetServerPwd(InitDB(), server_name))
		}
	case cd.FullCommand():
		ChangeDir(*cd_server_path)
	case mkdir.FullCommand():
		// mkdir
		MakeDir(*mkdir_path)
	case cp.FullCommand():
		println(*cp_src_file, *cp_dst_file)
		CopyFile(*cp_src_file, *cp_dst_file)
	case mv.FullCommand():
		println(*mv_src_file, *mv_dst_file)
		MoveFile(*mv_src_file, *mv_dst_file)
	case rm.FullCommand():
		println(*rm_file)
		DeleteFile(*rm_file)
	}

}

// get file info from remote or local
func GetFileInfo(file_path string) (*FileInfo, error) {
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
		}, nil
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
			// log.Panicf("Call Route err: %v", err)
			return nil, err
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
		}, nil
	}
}

// get dir file info list from remote or local
func GetDirFileInfoList(dir_path string) (fileInfoList []*FileInfo) {
	// get server info
	server_name, path := ExtractServerNameAndFilePath(dir_path)
	if server_name == "localhost" {
		if path == "" {
			path = "."
		}
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

		// if path is empty, set it to "server_info.Pwd"
		if path == "" {
			path = server_info.Pwd
		}
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
func CopyFile(src string, dst string) (err error) {
	// get server info
	src_server_name, src_path := ExtractServerNameAndFilePath(src)
	dst_server_name, dst_path := ExtractServerNameAndFilePath(dst)

	// get server
	src_server := GetServer(InitDB(), src_server_name)
	dst_server := GetServer(InitDB(), dst_server_name)

	// final path to use
	var src_path_final, dst_path_final string

	if src_server_name == "localhost" && dst_server_name == "localhost" { // local to local
		// judge if src file exist
		if !IsLocalFileExist(src_path) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path)
		}
		src_path_final = src_path
		// judge if dst file exist
		if IsLocalFileExist(dst_path) {
			// if file is dir, append file name
			if LocalFileIsDir(dst_path) {
				dst_path_final = filepath.Join(dst_path, filepath.Base(src_path))
				// Judge if final dst path is exists, if exists, return err
				if IsLocalFileExist(dst_path_final) {
					fmt.Println("dst file exist")
					return errors.New("File Exist: " + dst_path_final)
				}
			} else {
				fmt.Println("dst file exist")
				return errors.New("File Exist: " + dst_path)
			}
		} else { // file not exist, use dst_path directly
			dst_path_final = dst_path
		}
	} else if src_server_name == "localhost" && dst_server_name != "localhost" { // local to remote
		// judge if src file exist
		if !IsLocalFileExist(src_path) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path)
		}
		src_path_final = src_path
		// convert dst_path to dst_path_abs
		var dst_path_abs string
		if !filepath.IsAbs(dst_path) {
			dst_path_abs = filepath.Join(dst_server.Pwd, dst_path)
		} else {
			dst_path_abs = dst_path
		}
		fmt.Println("dst_path_abs", dst_path_abs)
		// judge if dst file exist
		if IsRemoteFileExist(&dst_server, dst_path_abs) {
			// if file is dir, append file name
			if RemoteFileIsDir(&dst_server, dst_path_abs) {
				dst_path_final = filepath.Join(dst_path_abs, filepath.Base(src_path))
				fmt.Println("dst_path_final1", dst_path_final)
			} else {
				fmt.Println("dst file exist")
				return errors.New("File Exist: " + dst_path_abs)
			}
		} else {
			dst_path_final = dst_path_abs
		}
	} else if src_server_name != "localhost" && dst_server_name == "localhost" { // remote to remote or remote to local
		// convert src_path to src_path_abs
		var src_path_abs string
		if !filepath.IsAbs(src_path) {
			src_path_abs = filepath.Join(src_server.Pwd, src_path)
		} else {
			src_path_abs = src_path
		}
		src_path_final = src_path_abs
		fmt.Println("src_path_final", src_path_final)
		// judge if remote src file exist
		if !IsRemoteFileExist(&src_server, src_path_abs) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path_abs)
		} else { // src exists, judge dst
			// judge local dst file exists
			if IsLocalFileExist(dst_path) {
				// if file is dir, append file name
				if LocalFileIsDir(dst_path) {
					dst_path_final = filepath.Join(dst_path, filepath.Base(src_path))
					// Judge if final dst path is exists, if exists, return err
					if IsLocalFileExist(dst_path_final) {
						fmt.Println("dst file exist1")
						return errors.New("File Exist: " + dst_path_final)
					}
				} else { // file exists and NOT dir
					fmt.Println("dst file exist2")
					return errors.New("File Exist: " + dst_path)
				}
			} else { // file not exist, use dst_path directly
				dst_path_final = dst_path
			}
		}
	} else { // remote to remote
		// convert src_path to src_path_abs
		var src_path_abs, dst_path_abs string
		if !filepath.IsAbs(src_path) {
			src_path_abs = filepath.Join(src_server.Pwd, src_path)
		} else {
			src_path_abs = src_path
		}
		// convert dst_path to dst_path_abs
		if !filepath.IsAbs(dst_path) {
			dst_path_abs = filepath.Join(dst_server.Pwd, dst_path)
		} else {
			dst_path_abs = dst_path
		}
		src_path_final = src_path_abs
		dst_path_final = dst_path_abs
		// judge if remote src file exist
		if !IsRemoteFileExist(&src_server, src_path_abs) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path_abs)
		} else { // src exists, judge dst
			// judge remote dst file exists
			if IsRemoteFileExist(&dst_server, dst_path_abs) {
				// if file is dir, append file name
				if RemoteFileIsDir(&dst_server, dst_path_abs) {
					dst_path_final = filepath.Join(dst_path_abs, filepath.Base(src_path))
					// Judge if final dst path is exists, if exists, return err
					if IsRemoteFileExist(&dst_server, dst_path_final) {
						fmt.Println("dst file exist1")
						return errors.New("File Exist: " + dst_path_final)
					}
				} else { // file exists and NOT dir
					fmt.Println("dst file exist2")
					return errors.New("File Exist: " + dst_path_abs)
				}
			} else {
				dst_path_final = dst_path_abs
			}
		}
	}

	// 同一主机
	if src_server.Addr == dst_server.Addr {
		// 本地复制
		if src_server.Addr == "localhost" {
			// copy local file
			CopyLocalFileRecursively(src_path_final, dst_path_final)
			return nil
		} else {
			fmt.Println("final paht", src_path_final, dst_path_final)
			// copy remote file
			CopyRemoteFile(src_server, src_path_final, dst_path_final)
			return nil
		}
	} else if src_server.Addr == "localhost" || dst_server.Addr == "localhost" { // 一个本地一个远程
		if src_server.Addr == "localhost" {
			// upload local file to remote
			// UploadFile(dst_server, src_path_final, dst_path_final)
			UploadDir(dst_server, src_path_final, dst_path_final)
			return nil
		} else {
			// download remote file to local
			DownloadFile(src_server, dst_path_final, src_path_final)
			return nil
		}
	} else { // 两个远程
		TransmitFile(src_server, dst_server, src_path_final, dst_path_final)
	}
	return err
}

// move file
func MoveFile(src string, dst string) (err error) {
	// get server info
	src_server_name, src_path := ExtractServerNameAndFilePath(src)
	dst_server_name, dst_path := ExtractServerNameAndFilePath(dst)
	// get server
	src_server := GetServer(InitDB(), src_server_name)
	dst_server := GetServer(InitDB(), dst_server_name)

	// final path to use
	var src_path_final, dst_path_final string

	if src_server_name == "localhost" && dst_server_name == "localhost" { // local to local
		// judge if src file exist
		if !IsLocalFileExist(src_path) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path)
		}
		src_path_final = src_path
		// judge if dst file exist
		if IsLocalFileExist(dst_path) {
			// if file is dir, append file name
			if LocalFileIsDir(dst_path) {
				dst_path_final = filepath.Join(dst_path, filepath.Base(src_path))
				// Judge if final dst path is exists, if exists, return err
				if IsLocalFileExist(dst_path_final) {
					fmt.Println("dst file exist")
					return errors.New("File Exist: " + dst_path_final)
				}
			} else {
				fmt.Println("dst file exist")
				return errors.New("File Exist: " + dst_path)
			}
		} else { // file not exist, use dst_path directly
			dst_path_final = dst_path
		}
	} else if src_server_name == "localhost" && dst_server_name != "localhost" { // local to remote
		// judge if src file exist
		if !IsLocalFileExist(src_path) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path)
		}
		src_path_final = src_path
		// convert dst_path to dst_path_abs
		var dst_path_abs string
		if !filepath.IsAbs(dst_path) {
			dst_path_abs = filepath.Join(dst_server.Pwd, dst_path)
		} else {
			dst_path_abs = dst_path
		}
		fmt.Println("dst_path_abs", dst_path_abs)
		// judge if dst file exist
		if IsRemoteFileExist(&dst_server, dst_path_abs) {
			// if file is dir, append file name
			if RemoteFileIsDir(&dst_server, dst_path_abs) {
				dst_path_final = filepath.Join(dst_path_abs, filepath.Base(src_path))
				fmt.Println("dst_path_final1", dst_path_final)
			} else {
				fmt.Println("dst file exist")
				return errors.New("File Exist: " + dst_path_abs)
			}
		} else {
			dst_path_final = dst_path_abs
		}
	} else if src_server_name != "localhost" && dst_server_name == "localhost" { // remote to remote or remote to local
		// convert src_path to src_path_abs
		var src_path_abs string
		if !filepath.IsAbs(src_path) {
			src_path_abs = filepath.Join(src_server.Pwd, src_path)
		} else {
			src_path_abs = src_path
		}
		src_path_final = src_path_abs
		fmt.Println("src_path_final", src_path_final)
		// judge if remote src file exist
		if !IsRemoteFileExist(&src_server, src_path_abs) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path_abs)
		} else { // src exists, judge dst
			// judge local dst file exists
			if IsLocalFileExist(dst_path) {
				// if file is dir, append file name
				if LocalFileIsDir(dst_path) {
					dst_path_final = filepath.Join(dst_path, filepath.Base(src_path))
					// Judge if final dst path is exists, if exists, return err
					if IsLocalFileExist(dst_path_final) {
						fmt.Println("dst file exist1")
						return errors.New("File Exist: " + dst_path_final)
					}
				} else { // file exists and NOT dir
					fmt.Println("dst file exist2")
					return errors.New("File Exist: " + dst_path)
				}
			} else { // file not exist, use dst_path directly
				dst_path_final = dst_path
			}
		}
	} else { // remote to remote
		// convert src_path to src_path_abs
		var src_path_abs, dst_path_abs string
		if !filepath.IsAbs(src_path) {
			src_path_abs = filepath.Join(src_server.Pwd, src_path)
		} else {
			src_path_abs = src_path
		}
		// convert dst_path to dst_path_abs
		if !filepath.IsAbs(dst_path) {
			dst_path_abs = filepath.Join(dst_server.Pwd, dst_path)
		} else {
			dst_path_abs = dst_path
		}
		src_path_final = src_path_abs
		dst_path_final = dst_path_abs
		// judge if remote src file exist
		if !IsRemoteFileExist(&src_server, src_path_abs) {
			fmt.Println("src file not exist")
			return errors.New("File NOT Found: " + src_path_abs)
		} else { // src exists, judge dst
			// judge remote dst file exists
			if IsRemoteFileExist(&dst_server, dst_path_abs) {
				// if file is dir, append file name
				if RemoteFileIsDir(&dst_server, dst_path_abs) {
					dst_path_final = filepath.Join(dst_path_abs, filepath.Base(src_path))
					// Judge if final dst path is exists, if exists, return err
					if IsRemoteFileExist(&dst_server, dst_path_final) {
						fmt.Println("dst file exist1")
						return errors.New("File Exist: " + dst_path_final)
					}
				} else { // file exists and NOT dir
					fmt.Println("dst file exist2")
					return errors.New("File Exist: " + dst_path_abs)
				}
			} else {
				dst_path_final = dst_path_abs
			}
		}
	}

	// 同一主机
	if src_server.Addr == src_server.Addr {
		// 本地复制
		if src_server.Addr == "localhost" {
			// move local file
			MoveLocalFile(src_path_final, dst_path_final)
			return nil
		} else {
			// move remote file
			MoveRemoteFile(src_server, src_path_final, dst_path_final)
			return nil
		}
	} else if src_server.Addr == "localhost" || dst_server.Addr == "localhost" { // 一个本地一个远程
		if src_server.Addr == "localhost" {
			// upload local file to remote
			UploadFile(dst_server, src_path_final, dst_path_final)
			// delete local file
			DeleteLocalFile(src_path_final)
			return nil
		} else {
			// download remote file to local
			DownloadFile(src_server, src_path_final, dst_path_final)
			// delete remote file
			DeleteRemoteFile(dst_server, dst_path_final)
			return nil
		}
	} else { // 两个远程
		log.Println("Files will not be deleted for security reasons.")
		TransmitFile(src_server, dst_server, src_path_final, dst_path_final)
	}
	return err
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
		// convert path to abs_path
		var abs_path string
		if !filepath.IsAbs(path) {
			abs_path = filepath.Join(server.Pwd, path)
		} else {
			abs_path = path
		}
		// delete remote file
		DeleteRemoteFile(server, abs_path)
		return true
	}
}

// change dir
// server_name format: server_name:dir_path
func ChangeDir(server_path string) error {
	// get server info
	server_name, path := ExtractServerNameAndFilePath(server_path)
	if server_name == "localhost" || server_name == "" { // local
		fmt.Println("To Change Local Position, Please Use 'cd' Command.")
		return nil
	}
	// get server info
	server := GetServer(InitDB(), server_name)
	// if path is empty, set it to "server_info.Home"
	if path == "" {
		path = server.Home
	}
	// Judage weather path is absolute path(windows and linux), if not, add prefix
	var absPath string
	if !filepath.IsAbs(path) {
		fmt.Println("path", path)
		absPath = filepath.Join(server.Pwd, path)
	} else {
		absPath = path
	}
	log.Println("Change Dir:", absPath)
	// get fileinfo and check path exists or not
	file_info, err := GetFileInfo(server_name + ":" + absPath)
	if err != nil {
		fmt.Println("Fail to Change Dir:", err)
		return err
	}
	if !file_info.IsDir {
		fmt.Println(file_info.Name, "is not a directory.")
		return errors.New(file_info.Name + " is not a directory.")
	}

	// change dir
	ModifyServerPwd(InitDB(), server_name, absPath)
	return nil
}

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

// mkdir
func MakeDir(server_path string) error {
	// get server info
	server_name, path := ExtractServerNameAndFilePath(server_path)
	// get server
	server := GetServer(InitDB(), server_name)

	// 本地创建
	if server.Addr == "localhost" {
		// create local dir
		MakeLocalDir(path)
		return nil
	} else {
		// convert path to abs_path
		var abs_path string
		if !filepath.IsAbs(path) {
			abs_path = filepath.Join(server.Pwd, path)
		} else {
			abs_path = path
		}
		// make remote dir
		MakeRemoteDir(server, abs_path)

		return nil
	}
}
