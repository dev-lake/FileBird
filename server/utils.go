package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	copylib "github.com/otiai10/copy"
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

// Get File Info
func GetFileInfo(path string) fs.FileInfo {
	log.Println("Trace GetFileInfo.")
	fileStat, err := os.Stat(path)
	if err != nil {
		log.Panic(err)
	}
	// log
	log.Println(fileStat.Sys().(*syscall.Stat_t).Uid)
	return fileStat
}

// Read dir and return file list
func ReadDir(path string) []fs.FileInfo {
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

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// Judge file is dir or not
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Panic(err)
	}
	return fileInfo.IsDir()
}

// Move file
func MoveFile(src string, dst string) bool {
	err := os.Rename(src, dst)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}

// Copy file
func CopyFile(src string, dst string) bool {
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

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}

// Delete file
func DeleteFile(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
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

// Get linux current user
func GetLinuxCurrentUser() (*user.User, error) {
	user, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	return user, nil
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

// copy file recursively
func CopyLocalFileRecursively(src string, dst string) bool {
	err := copylib.Copy(src, dst)
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}
