package main

import "testing"

func TestUploadfile(t *testing.T) {
	UploadFile(ServerInfo{
		Name: "server1",
		Addr: "127.0.0.1",
		Port: 2000,
	}, "/Users/nuc/Downloads/copy.pdf", "/Users/nuc/Downloads/up.pdf")
}

func TestTransmitFile(t *testing.T) {
	TransmitFile(
		ServerInfo{
			Name: "server1",
			Addr: "127.0.0.1",
			Port: 2000,
		},
		ServerInfo{
			Name: "server1",
			Addr: "127.0.0.1",
			Port: 2000,
		},
		"/Users/nuc/Downloads/protoc111111",
		"/Users/nuc/Downloads/111111111111111/protoc111111",
	)
}
