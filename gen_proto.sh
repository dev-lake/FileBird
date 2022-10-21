#!/bin/sh

protoc --go_out=./proto --go-grpc_out=./proto ./proto/*.proto