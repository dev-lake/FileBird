package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/jedib0t/go-pretty/v6/table"
	copylib "github.com/otiai10/copy"
	"github.com/schollz/progressbar/v3"
)

func ShowServerTable(servers []ServerInfo, filter string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Addr", "Port", "User", "Password", "HomeDir", "Pwd", "Description"})
	if filter == "" {
		for i, server := range servers {
			t.AppendRow(table.Row{i, GetServerStatus(&server) + server.Name, server.Addr, server.Port, server.User, server.Pass, server.Home, server.Pwd, server.Desc})
		}
	} else {
		for i, server := range servers {
			if server.Name == filter {
				t.AppendRow(table.Row{i, GetServerStatus(&server) + server.Name, server.Addr, server.Port, server.User, server.Pass, server.Home, server.Pwd, server.Desc})
			}
		}
	}
	t.Render()
	println("")
}

// extract server_name and file_path from server_name:file_path
func ExtractServerNameAndFilePath(serverNameAndFilePath string) (serverName string, filePath string) {
	if serverNameAndFilePath == "" {
		return "localhost", ""
	}
	for i, v := range serverNameAndFilePath {
		if v == '/' || v == '\\' || i > ServerNameMaxLength { // no server name, only file path
			return "localhost", serverNameAndFilePath
		}
		if v == ':' {
			serverName = serverNameAndFilePath[:i]
			filePath = serverNameAndFilePath[i+1:]
			break
		}
	}
	if serverName == "" {
		serverName = "localhost"
		filePath = serverNameAndFilePath
	}
	return serverName, filePath
}

// show file info
func ShowFileInfoTable(file_info *FileInfo) {
	println("\nFile Info")
	println("---------")
	println("Name:", file_info.Name)
	println("Size:", file_info.Size)
	println("Owner:", file_info.UserName, file_info.GroupName)
	// println("Mode:", file_info.Mode)
	println("ModTime:", file_info.ModTime)
	println("IsDir:", file_info.IsDir)
	println("Mode:", file_info.Mode.String())
	// println("Sys:", file_info.Sys)
	println("Path:", file_info.Path)
	println("\n")
}

// show file list
func ShowFileInfoListTable(info_list []*FileInfo) {
	for _, info := range info_list {
		fmt.Println("FileName:", info.Name)
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.Style().Options.SeparateRows = true
	t.Style().Options.SeparateColumns = false
	t.Style().Options.DrawBorder = false
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Mode", "Name", "Owner/Group", "Size", "ModTime"})
	for _, file := range info_list {
		short_time := strings.Split(file.ModTime, ".")[0]
		t.AppendRow(table.Row{file.Mode, file.Name, file.UserName + "/" + file.GroupName, file.Size, short_time})
	}
	t.Render()
	println("")
}

// Get File Info
func GetLocalFileInfo(path string) fs.FileInfo {
	// log.Println("Trace GetFileInfo.")
	fileStat, err := os.Stat(path)
	if err != nil {
		log.Panic(err)
	}
	// log
	// b, err := json.Marshal(fileStat)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Println(string(b))
	return fileStat
}

// Read dir and return file list
func ReadLocalDir(path string) []fs.FileInfo {
	log.Println("Trace ReadDir.")
	dir, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer dir.Close()

	fileInfoList, err := dir.Readdir(-1)
	if err != nil {
		log.Panic(err)
	}
	return fileInfoList
}

// read dir all files recursively
func ReadLocalDirAll(path string) []FileInfo {
	log.Println("Trace ReadDirAll.", path)
	var file_list []FileInfo
	fileInfoList, err := os.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}
	// convert []fs.DirEntry to FileInfo struct
	for _, file := range fileInfoList {
		if file.Name() == ".DS_Store" {
			continue
		}
		info, _ := file.Info()
		file_list = append(file_list, FileInfo{
			Name: file.Name(),
			Size: info.Size(),
			// UserName: info.Sys().(*syscall.Stat_t).Uid,
			// GroupName: file.Sys().(*syscall.Stat_t).Gid,
			Mode:    info.Mode(),
			ModTime: info.ModTime().String(),
			IsDir:   file.IsDir(),
			Path:    filepath.Join(path, file.Name()),
		})
	}
	// reverse next layer file list
	for _, fileInfo := range file_list {
		if fileInfo.IsDir {
			file_list = append(file_list, ReadLocalDirAll(fileInfo.Path)...)
		}
	}

	return file_list
}

func GetFileUserAndGroupName(fStat fs.FileInfo) (username string, groupname string) {
	// get file user id and group id
	uid := fStat.Sys().(*syscall.Stat_t).Uid
	gid := fStat.Sys().(*syscall.Stat_t).Gid
	u := strconv.FormatUint(uint64(uid), 10)
	g := strconv.FormatUint(uint64(gid), 10)
	usr, err := user.LookupId(u)
	if err != nil {
		log.Panic(err)
	}
	group, err := user.LookupGroupId(g)
	if err != nil {
		log.Panic(err)
	}
	return usr.Username, group.Name
}

// Move file
func MoveLocalFile(src string, dst string) bool {
	err := os.Rename(src, dst)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}

// Copy file
func CopyLocalFile(src string, dst string) bool {
	log.Println("trace CopyLocalFile", src, dst)
	// init progress bar
	progressbar := progressbar.DefaultBytes(
		GetLocalFileInfo(src).Size(),
		"copying",
	)

	srcFile, err := os.Open(src)
	if err != nil {
		log.Panic(err)
		return false
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		log.Panic(err)
		return false
	}
	defer dstFile.Close()

	_, err = io.Copy(io.MultiWriter(dstFile, progressbar), srcFile)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}

// copy file recursively
func CopyLocalFileRecursively(src string, dst string) bool {
	err := copylib.Copy(src, dst)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}

// Delete file
func DeleteLocalFile(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}

// Get linux current user
func GetLinuxCurrentUser() (*user.User, error) {
	user, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	return user, nil
}

// Check server name
func CheckServerName(serverName string) error {
	if serverName == "" {
		return errors.New("server name can NOT be empty")
	}
	for _, name := range IllegalServerNames {
		if name == serverName {
			return errors.New("server name can NOT be \"" + name + "\"")
		}
	}
	for i, v := range serverName {
		if v == '/' || v == '\\' || i > ServerNameMaxLength || v == ':' {
			return errors.New("server name can NOT contain '/' '\\' : and length can NOT be more than " + strconv.Itoa(ServerNameMaxLength))
		}
	}

	return nil
}

func IsLocalFileExist(dst_path string) bool {
	_, err := os.Stat(dst_path)
	if err != nil {
		return false
	}
	return true
}

func LocalFileIsDir(dst_path string) bool {
	file_info := GetLocalFileInfo(dst_path)
	return file_info.IsDir()
}

// make dir
func MakeLocalDir(path string) bool {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}
