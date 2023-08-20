#!/bin/bash

# mac
# CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/filebird-mac-amd64
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o bin/filebird-mac-arm64

# linux
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/filebird-linux-amd64
# CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o bin/filebird-linux-arm64

# windows
# CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o bin/filebird-windows-amd64.exe
# CGO_ENABLED=1 GOOS=windows GOARCH=arm64 go build -o bin/filebird-windows-arm64.exe