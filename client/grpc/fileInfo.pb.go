// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: proto/fileInfo.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FileReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilePath string `protobuf:"bytes,1,opt,name=filePath,proto3" json:"filePath,omitempty"`
}

func (x *FileReq) Reset() {
	*x = FileReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_fileInfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileReq) ProtoMessage() {}

func (x *FileReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_fileInfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileReq.ProtoReflect.Descriptor instead.
func (*FileReq) Descriptor() ([]byte, []int) {
	return file_proto_fileInfo_proto_rawDescGZIP(), []int{0}
}

func (x *FileReq) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

type RegularFileInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path         string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	IsDir        bool   `protobuf:"varint,3,opt,name=isDir,proto3" json:"isDir,omitempty"`
	Size         uint64 `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"` // bytes
	Mode         uint32 `protobuf:"varint,5,opt,name=mode,proto3" json:"mode,omitempty"`
	UserName     string `protobuf:"bytes,6,opt,name=userName,proto3" json:"userName,omitempty"`
	GroupName    string `protobuf:"bytes,7,opt,name=groupName,proto3" json:"groupName,omitempty"`
	CreateTime   string `protobuf:"bytes,8,opt,name=createTime,proto3" json:"createTime,omitempty"`
	ModifyTime   string `protobuf:"bytes,9,opt,name=modifyTime,proto3" json:"modifyTime,omitempty"`
	LastOpenTime string `protobuf:"bytes,10,opt,name=lastOpenTime,proto3" json:"lastOpenTime,omitempty"`
	Exsit        bool   `protobuf:"varint,11,opt,name=exsit,proto3" json:"exsit,omitempty"`
}

func (x *RegularFileInfo) Reset() {
	*x = RegularFileInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_fileInfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegularFileInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegularFileInfo) ProtoMessage() {}

func (x *RegularFileInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_fileInfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegularFileInfo.ProtoReflect.Descriptor instead.
func (*RegularFileInfo) Descriptor() ([]byte, []int) {
	return file_proto_fileInfo_proto_rawDescGZIP(), []int{1}
}

func (x *RegularFileInfo) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *RegularFileInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RegularFileInfo) GetIsDir() bool {
	if x != nil {
		return x.IsDir
	}
	return false
}

func (x *RegularFileInfo) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *RegularFileInfo) GetMode() uint32 {
	if x != nil {
		return x.Mode
	}
	return 0
}

func (x *RegularFileInfo) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *RegularFileInfo) GetGroupName() string {
	if x != nil {
		return x.GroupName
	}
	return ""
}

func (x *RegularFileInfo) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

func (x *RegularFileInfo) GetModifyTime() string {
	if x != nil {
		return x.ModifyTime
	}
	return ""
}

func (x *RegularFileInfo) GetLastOpenTime() string {
	if x != nil {
		return x.LastOpenTime
	}
	return ""
}

func (x *RegularFileInfo) GetExsit() bool {
	if x != nil {
		return x.Exsit
	}
	return false
}

type DirFileInfoList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path         string             `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	FileInfoList []*RegularFileInfo `protobuf:"bytes,2,rep,name=fileInfoList,proto3" json:"fileInfoList,omitempty"`
}

func (x *DirFileInfoList) Reset() {
	*x = DirFileInfoList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_fileInfo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirFileInfoList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirFileInfoList) ProtoMessage() {}

func (x *DirFileInfoList) ProtoReflect() protoreflect.Message {
	mi := &file_proto_fileInfo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirFileInfoList.ProtoReflect.Descriptor instead.
func (*DirFileInfoList) Descriptor() ([]byte, []int) {
	return file_proto_fileInfo_proto_rawDescGZIP(), []int{2}
}

func (x *DirFileInfoList) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *DirFileInfoList) GetFileInfoList() []*RegularFileInfo {
	if x != nil {
		return x.FileInfoList
	}
	return nil
}

type FileExistRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exsit bool `protobuf:"varint,1,opt,name=exsit,proto3" json:"exsit,omitempty"`
}

func (x *FileExistRep) Reset() {
	*x = FileExistRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_fileInfo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileExistRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileExistRep) ProtoMessage() {}

func (x *FileExistRep) ProtoReflect() protoreflect.Message {
	mi := &file_proto_fileInfo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileExistRep.ProtoReflect.Descriptor instead.
func (*FileExistRep) Descriptor() ([]byte, []int) {
	return file_proto_fileInfo_proto_rawDescGZIP(), []int{3}
}

func (x *FileExistRep) GetExsit() bool {
	if x != nil {
		return x.Exsit
	}
	return false
}

type FileIsDirRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsDir bool `protobuf:"varint,1,opt,name=isDir,proto3" json:"isDir,omitempty"`
}

func (x *FileIsDirRep) Reset() {
	*x = FileIsDirRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_fileInfo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileIsDirRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileIsDirRep) ProtoMessage() {}

func (x *FileIsDirRep) ProtoReflect() protoreflect.Message {
	mi := &file_proto_fileInfo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileIsDirRep.ProtoReflect.Descriptor instead.
func (*FileIsDirRep) Descriptor() ([]byte, []int) {
	return file_proto_fileInfo_proto_rawDescGZIP(), []int{4}
}

func (x *FileIsDirRep) GetIsDir() bool {
	if x != nil {
		return x.IsDir
	}
	return false
}

var File_proto_fileInfo_proto protoreflect.FileDescriptor

var file_proto_fileInfo_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x07, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0xab, 0x02,
	0x0a, 0x0f, 0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x73, 0x44,
	0x69, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x73, 0x44, 0x69, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4f, 0x70, 0x65, 0x6e, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4f, 0x70, 0x65,
	0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x73, 0x69, 0x74, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x78, 0x73, 0x69, 0x74, 0x22, 0x5b, 0x0a, 0x0f, 0x44,
	0x69, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x12, 0x34, 0x0a, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69,
	0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x52, 0x65, 0x67, 0x75, 0x6c,
	0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0c, 0x66, 0x69, 0x6c, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x24, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65,
	0x45, 0x78, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x73, 0x69,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x78, 0x73, 0x69, 0x74, 0x22, 0x24,
	0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x73, 0x44, 0x69, 0x72, 0x52, 0x65, 0x70, 0x12, 0x14,
	0x0a, 0x05, 0x69, 0x73, 0x44, 0x69, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69,
	0x73, 0x44, 0x69, 0x72, 0x32, 0xe9, 0x01, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x2b, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x08, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x52, 0x65, 0x67,
	0x75, 0x6c, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00, 0x12, 0x2e,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x44, 0x69, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x08, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x44, 0x69, 0x72,
	0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x2e,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x44, 0x69, 0x72, 0x41, 0x6c, 0x6c, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x12, 0x08, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x44, 0x69, 0x72,
	0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x28,
	0x0a, 0x0b, 0x49, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74, 0x12, 0x08, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x45, 0x78,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x22, 0x00, 0x12, 0x26, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65,
	0x49, 0x73, 0x44, 0x69, 0x72, 0x12, 0x08, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x1a,
	0x0d, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x73, 0x44, 0x69, 0x72, 0x52, 0x65, 0x70, 0x22, 0x00,
	0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x2e, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_fileInfo_proto_rawDescOnce sync.Once
	file_proto_fileInfo_proto_rawDescData = file_proto_fileInfo_proto_rawDesc
)

func file_proto_fileInfo_proto_rawDescGZIP() []byte {
	file_proto_fileInfo_proto_rawDescOnce.Do(func() {
		file_proto_fileInfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_fileInfo_proto_rawDescData)
	})
	return file_proto_fileInfo_proto_rawDescData
}

var file_proto_fileInfo_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_fileInfo_proto_goTypes = []interface{}{
	(*FileReq)(nil),         // 0: FileReq
	(*RegularFileInfo)(nil), // 1: RegularFileInfo
	(*DirFileInfoList)(nil), // 2: DirFileInfoList
	(*FileExistRep)(nil),    // 3: FileExistRep
	(*FileIsDirRep)(nil),    // 4: FileIsDirRep
}
var file_proto_fileInfo_proto_depIdxs = []int32{
	1, // 0: DirFileInfoList.fileInfoList:type_name -> RegularFileInfo
	0, // 1: FileInfo.GetFileInfo:input_type -> FileReq
	0, // 2: FileInfo.GetDirFileInfo:input_type -> FileReq
	0, // 3: FileInfo.GetDirAllFiles:input_type -> FileReq
	0, // 4: FileInfo.IsFileExist:input_type -> FileReq
	0, // 5: FileInfo.FileIsDir:input_type -> FileReq
	1, // 6: FileInfo.GetFileInfo:output_type -> RegularFileInfo
	2, // 7: FileInfo.GetDirFileInfo:output_type -> DirFileInfoList
	2, // 8: FileInfo.GetDirAllFiles:output_type -> DirFileInfoList
	3, // 9: FileInfo.IsFileExist:output_type -> FileExistRep
	4, // 10: FileInfo.FileIsDir:output_type -> FileIsDirRep
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_fileInfo_proto_init() }
func file_proto_fileInfo_proto_init() {
	if File_proto_fileInfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_fileInfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_fileInfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegularFileInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_fileInfo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirFileInfoList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_fileInfo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileExistRep); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_fileInfo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileIsDirRep); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_fileInfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_fileInfo_proto_goTypes,
		DependencyIndexes: file_proto_fileInfo_proto_depIdxs,
		MessageInfos:      file_proto_fileInfo_proto_msgTypes,
	}.Build()
	File_proto_fileInfo_proto = out.File
	file_proto_fileInfo_proto_rawDesc = nil
	file_proto_fileInfo_proto_goTypes = nil
	file_proto_fileInfo_proto_depIdxs = nil
}
