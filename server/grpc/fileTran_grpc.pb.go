// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: proto/fileTran.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FileTransClient is the client API for FileTrans service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileTransClient interface {
	Upload(ctx context.Context, opts ...grpc.CallOption) (FileTrans_UploadClient, error)
	Download(ctx context.Context, in *DownloadReq, opts ...grpc.CallOption) (FileTrans_DownloadClient, error)
}

type fileTransClient struct {
	cc grpc.ClientConnInterface
}

func NewFileTransClient(cc grpc.ClientConnInterface) FileTransClient {
	return &fileTransClient{cc}
}

func (c *fileTransClient) Upload(ctx context.Context, opts ...grpc.CallOption) (FileTrans_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileTrans_ServiceDesc.Streams[0], "/FileTrans/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileTransUploadClient{stream}
	return x, nil
}

type FileTrans_UploadClient interface {
	Send(*UploadReq) error
	CloseAndRecv() (*UploadRep, error)
	grpc.ClientStream
}

type fileTransUploadClient struct {
	grpc.ClientStream
}

func (x *fileTransUploadClient) Send(m *UploadReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileTransUploadClient) CloseAndRecv() (*UploadRep, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadRep)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileTransClient) Download(ctx context.Context, in *DownloadReq, opts ...grpc.CallOption) (FileTrans_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileTrans_ServiceDesc.Streams[1], "/FileTrans/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileTransDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileTrans_DownloadClient interface {
	Recv() (*DownloadRep, error)
	grpc.ClientStream
}

type fileTransDownloadClient struct {
	grpc.ClientStream
}

func (x *fileTransDownloadClient) Recv() (*DownloadRep, error) {
	m := new(DownloadRep)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileTransServer is the server API for FileTrans service.
// All implementations must embed UnimplementedFileTransServer
// for forward compatibility
type FileTransServer interface {
	Upload(FileTrans_UploadServer) error
	Download(*DownloadReq, FileTrans_DownloadServer) error
	mustEmbedUnimplementedFileTransServer()
}

// UnimplementedFileTransServer must be embedded to have forward compatible implementations.
type UnimplementedFileTransServer struct {
}

func (UnimplementedFileTransServer) Upload(FileTrans_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedFileTransServer) Download(*DownloadReq, FileTrans_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedFileTransServer) mustEmbedUnimplementedFileTransServer() {}

// UnsafeFileTransServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileTransServer will
// result in compilation errors.
type UnsafeFileTransServer interface {
	mustEmbedUnimplementedFileTransServer()
}

func RegisterFileTransServer(s grpc.ServiceRegistrar, srv FileTransServer) {
	s.RegisterService(&FileTrans_ServiceDesc, srv)
}

func _FileTrans_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileTransServer).Upload(&fileTransUploadServer{stream})
}

type FileTrans_UploadServer interface {
	SendAndClose(*UploadRep) error
	Recv() (*UploadReq, error)
	grpc.ServerStream
}

type fileTransUploadServer struct {
	grpc.ServerStream
}

func (x *fileTransUploadServer) SendAndClose(m *UploadRep) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileTransUploadServer) Recv() (*UploadReq, error) {
	m := new(UploadReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileTrans_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileTransServer).Download(m, &fileTransDownloadServer{stream})
}

type FileTrans_DownloadServer interface {
	Send(*DownloadRep) error
	grpc.ServerStream
}

type fileTransDownloadServer struct {
	grpc.ServerStream
}

func (x *fileTransDownloadServer) Send(m *DownloadRep) error {
	return x.ServerStream.SendMsg(m)
}

// FileTrans_ServiceDesc is the grpc.ServiceDesc for FileTrans service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileTrans_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FileTrans",
	HandlerType: (*FileTransServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _FileTrans_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _FileTrans_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/fileTran.proto",
}
