package main

import "testing"

func TestUploadfile(t *testing.T) {
	UploadFile(ServerInfo{
		Name: "server1",
		Addr: "127.0.0.1",
		Port: 2000,
	}, "/Users/nuc/Downloads/copy.pdf", "/Users/nuc/Downloads/up.pdf")
}
