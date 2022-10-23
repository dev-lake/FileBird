package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
)

func ShowServerTable(servers []ServerInfo, filter string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Addr", "Port"})
	if filter == "" {
		for i, server := range servers {
			t.AppendRow(table.Row{i, server.Name, server.Addr, server.Port})
		}
	} else {
		for i, server := range servers {
			if server.Name == filter {
				t.AppendRow(table.Row{i, server.Name, server.Addr, server.Port})
			}
		}
	}
	t.Render()
	println("")
}

// extract server_name and file_path from server_name:file_path
func ExtractServerNameAndFilePath(serverNameAndFilePath string) (serverName string, filePath string) {
	for i, v := range serverNameAndFilePath {
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
	log.Println("trace CopyLocalFile")
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

// Delete file
func DeleteLocalFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}
