package main

import (
	"testing"
)

func TestGetRemoteDirAllFiles(t *testing.T) {
	files, err := GetRemoteDirAllFiles(&ServerInfo{
		Name: "server1",
		Addr: "127.0.0.1",
		Port: 2000,
	}, "/Users/nuc/Downloads")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		// fmt.Println(file.Path)
		t.Log(file.Path)
	}
}
